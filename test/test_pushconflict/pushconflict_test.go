package pushconflict

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-technologies/terraform-provider-gitfile/test/helpers"
)

func setup() {
	cleanup()
	example_dir := "example.git"

	os.MkdirAll(example_dir, 0755)
	helpers.GitCommand(example_dir, "init")
	os.Create(fmt.Sprintf("%s/.exists", example_dir))
	helpers.GitCommand(example_dir, "add", ".exists")
	helpers.GitCommand(example_dir, "commit", "-m", "Initial Commit")
	helpers.GitCommand(example_dir, "checkout", "-b", "move_HEAD")
}

func createConflict(git_folder string, file_no int, wg *sync.WaitGroup) {
	defer wg.Done()
	helpers.GitCommand(git_folder, "checkout", "master")

	message := []byte(fmt.Sprintf("conflict changes %d", file_no))
	err := ioutil.WriteFile(fmt.Sprintf("%s/terraform", git_folder), message, 0644)
	if err != nil {
		log.Fatal(err)
	}
	helpers.GitCommand(git_folder, "add", "terraform")
	helpers.GitCommand(git_folder, "commit", "-m", fmt.Sprintf("test conflict %d", file_no))
	helpers.GitCommand(git_folder, "checkout", "move_HEAD")
}

func spinWorkers(w *sync.WaitGroup) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go createConflict("example.git", i, &wg)
		time.Sleep(time.Second)
	}
	wg.Wait()
	w.Done()
}

func TestPushConflict(t *testing.T) {
	setup()

	o := &helpers.TerratestDefaultOptions
	terraform.Init(t, o)
	helpers.GeneratePlan(t, o, "plan.out")
	defer os.RemoveAll("plan.out")

	var wg sync.WaitGroup
	wg.Add(1)
	go spinWorkers(&wg)
	helpers.ApplyWithPlanFile(t, o, "plan.out")
	wg.Wait()

	expected_commit_msg := "Created by terraform gitfile_commit"

	tests := []struct {
		output   string
		expected string
	}{
		{"gitfile_checkout_path", "checkout"},
		{"gitfile_commit_commit_message", expected_commit_msg},
	}

	for _, test := range tests {
		actual, _ := terraform.OutputE(t, o, test.output)
		assert.Equal(t, test.expected, actual, "terraform output '%s'", test.output)
	}
	found, err := helpers.GitLogContains("checkout", expected_commit_msg)
	assert.NoError(t, err)
	assert.True(t, found, fmt.Sprintf("checkout should have commit message '%s'", expected_commit_msg))

	helpers.GitCommand("checkout", "fetch")
	foundLogMaster, err := helpers.GitLogContains("checkout", expected_commit_msg, "origin/master")
	assert.NoError(t, err)
	assert.True(t, foundLogMaster, fmt.Sprintf("checkout should have commit message '%s'", expected_commit_msg))
	assert.FileExists(t, "checkout/terraform", "terraform file does not exist in checkout")

	time.Sleep(2 * time.Second)

	helpers.GenerateDestroyPlan(t, o, "destroy.out")
	defer os.RemoveAll("destroy.out")
	helpers.ApplyWithPlanFile(t, o, "destroy.out")

	assert.NoDirExists(t, "checkout", "checkout dir should be deleted on destroy")
}

func cleanup() {
	os.RemoveAll("example.git")
	os.RemoveAll("checkout")
	os.RemoveAll("terraform.tfstate")
	os.RemoveAll("terraform.tfstate.backup")
}
