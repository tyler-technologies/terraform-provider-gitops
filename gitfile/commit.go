package gitfile

import (
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const CommitBodyHeader string = "The following files are managed by terraform:"

func commitResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"commit_message": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Created by terraform gitfile_commit",
				ForceNew: true,
			},
			"handles": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      hashString,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
		Create: CommitCreate,
		Read:   CommitRead,
		Delete: CommitDelete,
		Exists: CommitExists,
	}
}

func CommitCreate(d *schema.ResourceData, m interface{}) error {
	c := getConfig(m)
	checkout_dir := m.(*GitFileConfig).Path
	retry_count := d.Get("retry_count").(int)
	retry_interval := d.Get("retry_interval").(int)
	lockCheckout(checkout_dir)
	defer unlockCheckout(checkout_dir)
	if err := cloneIfNotExist(c); err != nil {
		return err
	}

	handles := d.Get("handles").(*schema.Set)
	filepaths := []string{}
	for _, handle := range handles.List() {
		filepaths = append(filepaths, parseHandle(handle.(string)).path)
	}
	commit_message := d.Get("commit_message").(string)
	commit_body := fmt.Sprintf("%s\n%s", CommitBodyHeader, strings.Join(filepaths, "\n"))

	if err := commit(checkout_dir, commit_message, commit_body); err != nil {
		return errwrap.Wrapf("push error: {{err}}", err)
	}

	if err := push(checkout_dir, 0, retry_count, retry_interval); err != nil {
		return err
	}

	out, err := gitCommand(checkout_dir, "rev-parse", "HEAD")
	if err != nil {
		return err
	}

	sha := strings.TrimRight(string(out), "\n")

	d.SetId(fmt.Sprintf("%s %s", sha, checkout_dir))
	return nil
}

func CommitRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func CommitExists(d *schema.ResourceData, m interface{}) (bool, error) {
	checkout_dir := m.(*GitFileConfig).Path
	lockCheckout(checkout_dir)
	defer unlockCheckout(checkout_dir)
	commitID := strings.Split(d.Id(), " ")[0]

	_, err := gitCommand(checkout_dir, flatten("show", commitID)...)

	if err != nil {
		return false, nil
	} else {
		return true, nil
	}

}

func CommitDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
