package destroy

import (
	"fmt"
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
	helpers.GitCommand(example_dir, "checkout", "-b", "move_HEAD")
}

func TestDestroy(t *testing.T) {
	setup()

	o := &helpers.TerratestDefaultOptions
	terraform.Init(t, o)
	helpers.GeneratePlan(t, o, "plan.out")
	os.RemoveAll("checkout")
	helpers.ApplyWithPlanFile(t, o, "plan.out")
	os.RemoveAll("plan.out")

	expected_commit_msg := "Created by terraform gitops_commit"

	tests := []struct {
		output   string
		expected string
	}{
		{"gitops_checkout_path", "checkout"},
		{"gitops_commit_commit_message", expected_commit_msg},
	}

	for _, test := range tests {
		actual, _ := terraform.OutputE(t, o, test.output)
		assert.Equal(t, test.expected, actual, "terraform output '%s'", test.output)
	}

	// Check git log
	foundLog, err := helpers.GitLogContains("checkout", expected_commit_msg)
	assert.NoError(t, err)
	assert.True(t, foundLog, fmt.Sprintf("checkout should have commit message '%s'", expected_commit_msg))

	// Check contents of checkout folder / git repo
	helpers.GitCommand("checkout", "fetch")
	foundLogMaster, err := helpers.GitLogContains("checkout", expected_commit_msg, "origin/master")
	assert.NoError(t, err)
	assert.True(t, foundLogMaster, fmt.Sprintf("checkout should have commit message '%s'", expected_commit_msg))
	assert.FileExists(t, "checkout/terraform", "terraform file does not exist in checkout")

	// Check contents of example.git folder / git repo
	helpers.GitCommand("example.git", "checkout", "master")
	assert.FileExists(t, "example.git/terraform", "terraform file does not exist in example.git")

	helpers.GitCommand("example.git", "checkout", "move_HEAD")
	os.RemoveAll("checkout")

	helpers.GenerateDestroyPlan(t, o, "destroy.out")
	os.RemoveAll("checkout")
	helpers.ApplyWithPlanFile(t, o, "destroy.out")
	os.RemoveAll("destroy.out")

	helpers.GitCommand("example.git", "checkout", "master")
	assert.NoFileExists(t, "example.git/terraform", "terraform file does not exist in example.git")
}

func cleanup() {
	os.RemoveAll("example.git")
	os.RemoveAll("checkout")
	os.RemoveAll("terraform.tfstate")
	os.RemoveAll("terraform.tfstate.backup")
	os.RemoveAll("plan.out")
}
