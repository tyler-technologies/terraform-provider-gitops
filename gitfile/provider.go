package gitfile

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"repo_url": {
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"branch": {
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					i := strings.IndexRune(value, '/')
					if i == 0 {
						es = append(es, fmt.Errorf("Paths which begin with / not allowed in %q", k))
					}
					return
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gitfile_checkout": checkoutResource(),
			"gitfile_file":     fileResource(),
			"gitfile_symlink":  symlinkResource(),
			"gitfile_commit":   commitResource(),
		},
		ConfigureFunc: gitfileConfigure,
	}
}

func gitfileConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &GitFileConfig{
		RepoUrl: d.Get("repo_url").(string),
		Path:    d.Get("path").(string),
		Branch:  d.Get("branch").(string),
	}
	return config, nil
}

type GitFileConfig struct {
	RepoUrl string
	Path    string
	Branch  string
}
