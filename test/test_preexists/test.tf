terraform {
  required_providers {
    gitops = {
      source  = "tyler-technologies/gitops"
      version = "0.0.2-rc"
    }
  }
}

provider "gitops" {
  repo_url = "../example.git"
  branch   = "master"
  path     = "checkout"
}

resource "gitops_checkout" "test" {}

resource "gitops_file" "test" {
  checkout = gitops_checkout.test.id
  path     = "terraform"
  contents = "preexisting_commits\n"
}
resource "gitops_commit" "test" {
  commit_message = "Created by terraform gitops_commit"
  handles        = ["${gitops_file.test.id}"]
}

output "gitops_checkout_path" {
  value = gitops_checkout.test.path
}

output "gitops_commit_commit_message" {
  value = gitops_commit.test.commit_message
}
