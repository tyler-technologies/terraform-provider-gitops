provider "gitops" {
    repo_url = "../example.git"
    branch = "master"
    path = "checkout"
}

resource "gitops_checkout" "checkout" {}

resource "gitops_file" "file1" {
    checkout = gitops_checkout.checkout.id
    path = "terraform"
    contents = "Terraform making commits"
}

resource "gitops_file" "file2" {
    checkout = gitops_checkout.checkout.id
    path = "myfile"
    contents = "Terraform shizz"
}

resource "gitops_commit" "commit" {
    commit_message = "Created by terraform gitops_commit"
    # handles = ["${gitops_file.testfile.id}"]
    handles = ["${gitops_file.file1.id}", "${gitops_file.file2.id}"]
}

output "gitops_commit_commit_message" {
    value = gitops_commit.commit.commit_message
}

output "gitops_checkout_path" {
    value = gitops_checkout.checkout.path
}
