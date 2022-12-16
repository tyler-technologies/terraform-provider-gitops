package preexists

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-technologies/terraform-provider-gitops/test/helpers"
)

func setup() {
	cleanup()
	example_dir := "example.git"

	os.MkdirAll(example_dir, 0755)
	helpers.GitCommand(example_dir, "init")
	os.Create(fmt.Sprintf("%s/.exists", example_dir))
	helpers.GitCommand(example_dir, "add", ".exists")
	helpers.GitCommand(example_dir, "commit", "-m", "Initial Commit")

	os.Create(fmt.Sprintf("%s/terraform", example_dir))
	message := []byte("preexisting_commits\n")
	err := ioutil.WriteFile(fmt.Sprintf("%s/terraform", example_dir), message, 0644)
	if err != nil {
		log.Fatal(err)
	}
	helpers.GitCommand(example_dir, "add", "terraform")
	helpers.GitCommand(example_dir, "commit", "-m", "PRE")
	helpers.GitCommand(example_dir, "checkout", "-b", "move_HEAD")
}

func TestPreexists(t *testing.T) {
	setup()

	o := &helpers.TerratestDefaultOptions
	terraform.InitAndApply(t, o)
	expected_commit_msg := "Created by terraform gitops_commit"

	tests := []struct {
		output   string
		expected string
	}{
		{"gitops_checkout_path", "\"checkout\""},
		{"gitops_commit_commit_message", "\"" + expected_commit_msg + "\""},
	}

	for _, test := range tests {
		actual, _ := terraform.OutputE(t, o, test.output)
		assert.Equal(t, test.expected, actual, "terraform output '%s'", test.output)
	}
	found, err := helpers.GitLogContains("checkout", expected_commit_msg)
	assert.NoError(t, err)
	assert.True(t, found, fmt.Sprintf("checkout should have commit message '%s'", expected_commit_msg))

	out, err := helpers.GitCommand("checkout", "diff", "HEAD~1..HEAD")
	assert.NoError(t, err)
	assert.Empty(t, string(out), "commit should show no diff")
	assert.FileExists(t, "checkout/terraform")
}

func cleanup() {
	os.RemoveAll("example.git")
	os.RemoveAll("checkout")
	os.RemoveAll("terraform.tfstate")
	os.RemoveAll("terraform.tfstate.backup")
}
