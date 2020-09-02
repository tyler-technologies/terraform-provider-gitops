provider "gitfile" {
    repo_url = "../example.git"
    branch = "master"
    path = "checkout"
}

resource "gitfile_checkout" "test" {}

resource "gitfile_file" "test" {
    checkout = gitfile_checkout.test.id
    path = "terraform"
    contents = "Terraform making commits"
}

resource "gitfile_commit" "test" {
    commit_message = "Created by terraform gitfile_commit"
    handles = [gitfile_file.test.id]
}

output "gitfile_checkout_path" {
    value = gitfile_checkout.test.path
}

output "gitfile_commit_commit_message" {
    value = gitfile_commit.test.commit_message
}
