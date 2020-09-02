package helpers

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

var TerratestDefaultOptions = terraform.Options{
	Logger: logger.Discard,
}

func GeneratePlan(t *testing.T, o *terraform.Options, f string) error {
	outArg := fmt.Sprintf("-out=%s", f)
	if _, err := terraform.RunTerraformCommandE(t, o, terraform.FormatArgs(o, "plan", outArg)...); err != nil {
		err_message := fmt.Sprintf("Failed to generate plan with error: %s", err.Error())
		errors.New(err_message)
	}
	return nil
}

func GenerateDestroyPlan(t *testing.T, o *terraform.Options, f string) error {
	outArg := fmt.Sprintf("-out=%s", f)
	args := terraform.FormatArgs(o, "plan", "-destroy")
	args = append(args, outArg)
	if _, err := terraform.RunTerraformCommandE(t, o, args...); err != nil {
		err_message := fmt.Sprintf("Failed to generate plan with error: %s", err.Error())
		errors.New(err_message)
	}
	return nil
}

func ApplyWithPlanFile(t *testing.T, o *terraform.Options, f string) (string, error) {
	args := terraform.FormatArgs(o, "apply")
	args = append(args, f)
	return terraform.RunTerraformCommandE(t, o, args...)
}

func GitLogContains(checkout_dir, expected string, args ...string) (bool, error) {
	out, err := GitCommand(checkout_dir, "log")
	if err != nil {
		return false, err
	}
	return strings.Contains(string(out), expected), nil
}

func GitCommand(checkout_dir string, args ...string) ([]byte, error) {
	command := exec.Command("git", args...)
	command.Dir = checkout_dir
	out, err := command.CombinedOutput()
	if err != nil {
		log.Fatalf("error running git command: %s, %v", strings.Join(args, " "), err)

	}
	return out, err
}
