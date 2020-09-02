# <resource name> gitfile_file

Use this resource to create, read, update, and delete files in the underlying gitops repository.

## Example Usage

```hcl
resource "gitfile_file" "test_file" {
  checkout = gitfile_checkout.test_checkout.id
  path = "terraform"
  contents = "Terraform making commits"
}
```

## Argument Reference

* `checkout` - (Required) The ID of the checkout resource (This makes it possible to track changes accross many files in a single commit)
* `path` - (Required) The filepath inside the checked out repository.
* `contents` - (Required) The files contents.

## Attribute Reference

* `id` - The id of the created file. This is usually passed to gitfile_commit