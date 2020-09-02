provider "gitfile" {
    repo_url = "../example.git"
    branch = "master"
    path = "checkout"
}

resource "gitfile_checkout" "checkout" {}

resource "gitfile_file" "file1" {
    checkout = gitfile_checkout.checkout.id
    path = "terraform"
    contents = "Terraform making commits"
}

resource "gitfile_file" "file2" {
    checkout = gitfile_checkout.checkout.id
    path = "myfile"
    contents = "Terraform shizz"
}

resource "gitfile_commit" "commit" {
    commit_message = "Created by terraform gitfile_commit"
    # handles = ["${gitfile_file.testfile.id}"]
    handles = ["${gitfile_file.file1.id}", "${gitfile_file.file2.id}"]
}

output "gitfile_commit_commit_message" {
    value = gitfile_commit.commit.commit_message
}

output "gitfile_checkout_path" {
    value = gitfile_checkout.checkout.path
}
