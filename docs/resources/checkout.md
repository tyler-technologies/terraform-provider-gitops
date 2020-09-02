# <resource name> gitops_checkout

Checks out a git repository onto your local filesystem from within a terraform provider.

This is mostly used to ensure that a checkout is present, before using the _gitops_commit_
resource to commit some Terraform generated data.

## Example Usage

```hcl
resource "gitops_checkout" "test_checkout" {}
```

## Argument Reference

* `retry_count` - (Optional) The number of git checkout retries (default: `10`)
* `retry_interval` - (Optional) The number of seconds between git checkout retries (default: `5`)

## Attribute Reference

* `path` -  The file path on filesystem where the repository has been checked out
* `repo` - The repository url that was checked out
* `branch` - The branch being checked out
* `head` - The git head value