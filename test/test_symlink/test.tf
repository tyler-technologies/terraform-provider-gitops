provider "gitfile" {
    repo_url = "../example.git"
    branch = "master"
    path = "checkout"
}

resource "gitfile_checkout" "test" {}

resource "gitfile_symlink" "test" {
    checkout = gitfile_checkout.test.id
    path = "terraform"
    target = "/etc/passwd"
}

resource "gitfile_commit" "test" {
    commit_message = "Created by terraform gitfile_commit"
    handles = [gitfile_symlink.test.id]
}

output "gitfile_checkout_path" {
    value = gitfile_checkout.test.path
}

output "gitfile_commit_commit_message" {
    value = gitfile_commit.test.commit_message
}