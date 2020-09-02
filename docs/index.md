# <provider> Provider

Summary of what the provider is for, including use cases and links to
app/service documentation.

## Example Usage

```hcl
provider "gitops" {
  repo_url = "https://myverisoncontrolprovider.com/my/repo"
  branch = "master"
  path = "tmp.mycheckoutdestination"
}
```
## Argument Reference

* `repo_url` - (Required) The url for the git repository
* `branch` - (Required) The git branch to operate on for the gitops project
* `path` - (Required) The local path where the repository will be checked out temporarily during terraform actions