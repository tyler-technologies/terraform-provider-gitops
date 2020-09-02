package gitfile

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

func checkoutResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"head": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retry_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     10,
				Description: "Number of git commit retries",
			},
			"retry_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     5,
				Description: "Number of seconds between git commit retries",
			},
		},
		Create: CheckoutCreate,
		Read:   CheckoutRead,
		Update: nil,
		Delete: CheckoutDelete,
	}
}

func read(d *schema.ResourceData) error {
	checkout_dir := d.Id()
	var repo string
	var branch string
	var head string

	if out, err := gitCommand(checkout_dir, "config", "--get", "remote.origin.url"); err != nil {
		return err
	} else {
		repo = strings.TrimRight(string(out), "\n")
	}
	if out, err := gitCommand(checkout_dir, "rev-parse", "--abbrev-ref", "HEAD"); err != nil {
		return err
	} else {
		branch = strings.TrimRight(string(out), "\n")
	}

	if _, err := gitCommand(checkout_dir, "pull", "--ff-only", "origin"); err != nil {
		return err
	}

	if out, err := gitCommand(checkout_dir, "rev-parse", "HEAD"); err != nil {
		return err
	} else {
		head = strings.TrimRight(string(out), "\n")
	}
	_ = d.Set("path", checkout_dir)
	_ = d.Set("repo", repo)
	_ = d.Set("branch", branch)
	_ = d.Set("head", head)
	return nil
}

func CheckoutCreate(d *schema.ResourceData, m interface{}) error {
	c := getConfig(m)
	lockCheckout(c.Path)
	defer unlockCheckout(c.Path)
	if err := cloneIfNotExist(c); err != nil {
		return err
	}
	d.SetId(c.Path)
	return read(d)
}

func CheckoutRead(d *schema.ResourceData, m interface{}) error {
	checkout_id := d.Id()
	c := getConfig(m)

	if c.Path != checkout_id {
		err_message := fmt.Sprintf("[ERROR] Checkout directory state mismatch. Checkout Directory is: %s. Expected: %s", c.Path, checkout_id)
		return errors.New(err_message)
	}

	lockCheckout(c.Path)
	defer unlockCheckout(c.Path)
	if err := cloneIfNotExist(c); err != nil {
		return err
	}

	return read(d)
}

func CheckoutDelete(d *schema.ResourceData, m interface{}) error {
	checkout_id := d.Id()
	retry_count := d.Get("retry_count").(int)
	retry_interval := d.Get("retry_interval").(int)

	var repo string
	var branch string
	expected_repo := d.Get("repo").(string)
	expected_branch := d.Get("branch").(string)
	expected_head := d.Get("head").(string)
	c := getConfig(m)

	if c.Path != checkout_id {
		err_message := fmt.Sprintf("[ERROR] Checkout directory state mismatch. Checkout Directory is: %s. Expected: %s", c.Path, checkout_id)
		return errors.New(err_message)
	}

	lockCheckout(c.Path)
	defer unlockCheckout(c.Path)
	if err := cloneIfNotExist(c); err != nil {
		return err
	}

	// sanity check
	var head string

	if out, err := gitCommand(c.Path, "config", "--get", "remote.origin.url"); err != nil {
		return err
	} else {
		repo = strings.TrimRight(string(out), "\n")
	}
	if out, err := gitCommand(c.Path, "rev-parse", "--abbrev-ref", "HEAD"); err != nil {
		return err
	} else {
		branch = strings.TrimRight(string(out), "\n")
	}

	if _, err := gitCommand(c.Path, "pull", "--ff-only", "origin"); err != nil {
		return err
	}

	if out, err := gitCommand(c.Path, "rev-parse", "HEAD"); err != nil {
		return err
	} else {
		head = strings.TrimRight(string(out), "\n")
	}

	if expected_repo != repo {
		return fmt.Errorf("expected repo to be %s, was %s", expected_repo, repo)
	}
	if expected_branch != branch {
		return fmt.Errorf("expected branch to be %s, was %s", expected_branch, branch)
	}
	if expected_head != head {
		return fmt.Errorf("expected head to be %s, was %s", expected_head, head)
	}

	if err := commit(c.Path, "Removed by Terraform", ""); err != nil {
		return errwrap.Wrapf("push error: {{err}}", err)
	}

	if err := push(c.Path, 0, retry_count, retry_interval); err != nil {
		return err
	}

	// actually delete
	if err := os.RemoveAll(c.Path); err != nil {
		return err
	}

	return nil
}
