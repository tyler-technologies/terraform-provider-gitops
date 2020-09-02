# <resource name> File

Creates a symlink within a git repository from terraform

## Example Usage

```hcl
resource "gitops_symlink" "test_symlink" {
  checkout = gitops_checkout.test.id
  path = "terraform"
  target = "/etc/passwd"
}
```

## Argument Reference

* `checkout` - (Required) The ID of the checkout resource the files are associated with
* `path` - (Required) The path within the checkout to create the symlink at
* `target` - (Required) The place the symlink should point to. Can be an absolute or relative path

## Attribute Reference

* `id` - The id of the created symlink. This is usually passed to gitops_commit