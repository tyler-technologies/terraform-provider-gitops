<p align="center">
  <img src="git.png" alt="gitfile provider" width="140"/>

  <h3 align="center">Terraform Gitfile Provider</h3>

  <p align="center">
    <a href="https://github.com/tyler-technologies/terraform-provider-gitfile/actions?query=workflow%3Abuild"><img alt="Build" src="https://github.com/tyler-technologies/terraform-provider-gitfile/workflows/build/badge.svg"></a>
    <a href="https://github.com/tyler-technologies/terraform-provider-gitfile/actions?query=workflow%3Atest"><img alt="Test" src="https://github.com/tyler-technologies/terraform-provider-gitfile/workflows/test/badge.svg"></a>
    <a href="https://github.com/tyler-technologies/terraform-provider-gitfile/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/tyler-technologies/terraform-provider-gitfile"></a>
    <a href="https://github.com/tyler-technologies/terraform-provider-gitfile/releases/latest"><img alt="Downloads" src="https://img.shields.io/github/downloads/tyler-technologies/terraform-provider-gitfile/total?color=orange"></a>
    <a href="https://github.com/tyler-technologies/terraform-provider-gitfile/tree/checkout-tfe-support"><img alt="Latest Commit" src="https://img.shields.io/github/last-commit/tyler-technologies/terraform-provider-gitfile?color=ff69b4"></a>
  </p>
</p>

---

## Synopsis

A [Terraform](http://terraform.io) cloud optimized plugin to manage files in git repositories.
> Note: While this provider is optimized for terraform cloud/enterprise, it will still work excellently for terraform open source. 

This allows you to export terraform managed state into other systems which are controlled
by git repositories - for example commit server IPs to DNS config repositories,
or write out hiera data into your puppet configuration.

## Example:
```hcl
  provider "gitfile" {
    repo_url = "https://myverisoncontrolprovider.com/my/repo"
    branch = "master"
    path = "tmp.mycheckoutdestination"
  }

  resource "gitfile_checkout" "test_checkout" {}

  resource "gitfile_file" "test_file" {
    checkout = gitfile_checkout.test_checkout.id
    path = "terraform"
    contents = "Terraform making commits"
  }

  resource "gitfile_symlink" "test_symlink" {
    checkout = gitfile_checkout.test.id
    path = "terraform"
    target = "/etc/passwd"
  }

  resource "gitfile_commit" "test_commit" {
    commit_message = "Created by terraform gitfile_commit"
    handles = [gitfile_file.test_file.id, gitfile_file.test_symlink.id]
  }
```

## Resources

### gitfile_checkout

Checks out a git repository onto your local filesystem from within a terraform provider.

This is mostly used to ensure that a checkout is present, before using the _gitfile_commit_
resource to commit some Terraform generated data.

Inputs:
  - retry_count - The number of git checkout retries
  - retry_interval - The number of seconds between git checkout retries
  
Outputs:
  - path - The file path on filesystem where the repository has been checked out
  - repo - The repository url that was checked out
	- branch - The branch being checked out
	- head - The git head value

### gitfile_file

Creates a file within a git repository with some content from terraform

Inputs:
  - checkout - The ID of the checkout resource the files are associated with
  - path - The path within the checkout to create the file at
  - contents - The contents of the file

Outputs:
  - id - The id of the created file. This is usually passed to _gitfile_commit_

### gitfile_symlink

Creates a symlink within a git repository from terraform

Inputs:
  - checkout - The ID of the checkout resource the files are associated with
  - path - The path within the checkout to create the symlink at
  - target - The place the symlink should point to. Can be an absolute or relative path

Outputs:
  - id - The id of the created symlink. This is usually passed to _gitfile_commit_

### gitfile_commit

Makes a git commit of a set of _gitfile_commit_ and _gitfile_file_ resources in a git
repository, and pushes it to origin.

Note that even if the a file with the same contents Terraform creates already exists,
Terraform will create an empty commit with the specified commit message.

Inputs:
  - commit_message - The commit message to use for the commit
  - handles - An array of ids from _gitfile_file_ or _gitfile_symlink_ resources which should be included in this commit
  - retry_count - The number of git commit retries
  - retry_interval - The number of seconds between git commit retries

Outputs:
  - commit_message - The commit message for the commit that will be made

# License

Apache2 - See the included LICENSE file for more details.

