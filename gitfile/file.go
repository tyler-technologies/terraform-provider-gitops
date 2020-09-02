package gitfile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/hashicorp/terraform/helper/schema"
)

func fileResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"checkout": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contents": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Create: fileCreateUpdate,
		Read:   fileRead,
		Delete: fileDelete,
		Exists: fileExists,
	}
}

func fileCreateUpdate(d *schema.ResourceData, m interface{}) error {
	c := getConfig(m)

	lockCheckout(c.Path)
	defer unlockCheckout(c.Path)
	if err := cloneIfNotExist(c); err != nil {
		return err
	}

	filepath := d.Get("path").(string)
	contents := d.Get("contents").(string)

	filename := path.Join(c.Path, filepath)
	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, []byte(contents), 0666); err != nil {
		return err
	}

	if _, err := gitCommand(c.Path, "add", "--", filepath); err != nil {
		err_mess := fmt.Sprintf("SORRY BOUTCHA: %s", err.Error())
		return errors.New(err_mess)
	}

	hand := handle{
		kind: "file",
		hash: hashString(contents),
		path: filepath,
	}

	d.SetId(hand.String())
	return nil
}

func fileRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func fileExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := getConfig(m)

	lockCheckout(c.Path)
	defer unlockCheckout(c.Path)

	if err := cloneIfNotExist(c); err != nil {
		return false, err
	}
	filepath := d.Get("path").(string)

	out, err := ioutil.ReadFile(path.Join(c.Path, filepath))

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if string(out) == d.Get("contents").(string) {
		return true, nil
	} else {
		return false, nil
	}
}

func fileDelete(d *schema.ResourceData, m interface{}) error {
	c := getConfig(m)

	lockCheckout(c.Path)
	defer unlockCheckout(c.Path)
	if err := cloneIfNotExist(c); err != nil {
		return err
	}

	filepath := d.Get("path").(string)
	file := path.Join(c.Path, filepath)
	if err := os.Remove(file); err != nil {
		return err
	}
	if _, err := gitCommand(c.Path, "rm", "--", filepath); err != nil {
		return err
	}
	return nil
}
