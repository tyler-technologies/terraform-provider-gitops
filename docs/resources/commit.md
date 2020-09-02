# <resource name> gitfile_commit

Makes a git commit of a set of gitfile_commit and gitfile_file resources in a git
repository, and pushes it to origin.

Note that even if the a file with the same contents Terraform creates already exists,
Terraform will create an empty commit with the specified commit message.

## Example Usage

```hcl
resource "gitfile_commit" "test_commit" {
  commit_message = "Created by terraform gitfile_commit"
  handles = [gitfile_file.test_file.id, gitfile_file.test_symlink.id]
}
```

## Argument Reference
* `commit_message` - (Required) The commit message to use for the commit
* `handles` - (Required) An array of ids from gitfile_file or gitfile_symlink resources which should be included in this commit
* `retry_count` - (Optional) The number of git commit retries (default: `10`)
* `retry_interval` - (Optional) The number of seconds between git commit retries (default: `5`)


## Attribute Reference

* `commit_message` - The commit message for the commit that will be made