provider "gitfile" {
    repo_url = "../example.git"
    branch = "master"
    path = "checkout"
}

resource "gitfile_checkout" "test" {}

output "gitfile_checkout_path" {
    value = gitfile_checkout.test.path
}

resource "gitfile_file" "test" {
    checkout = gitfile_checkout.test.id
    path = "terraform"
    contents = "Terraform making commits"
}
resource "gitfile_file" "test1" {
    checkout = gitfile_checkout.test.id
    path = "terraform1"
    contents = "Terraform making commits"
}
resource "gitfile_file" "test2" {
    checkout = gitfile_checkout.test.id
    path = "terraform2"
    contents = "Terraform making commits"
}
resource "gitfile_file" "test3" {
    checkout = gitfile_checkout.test.id
    path = "terraform3"
    contents = "Terraform making commits"
}
resource "gitfile_file" "test4" {
    checkout = gitfile_checkout.test.id
    path = "terraform4"
    contents = "Terraform making commits"
}

resource "gitfile_commit" "test" {
    commit_message = "Created by terraform gitfile_commit"
    handles = [gitfile_file.test.id, gitfile_file.test1.id,gitfile_file.test2.id,gitfile_file.test3.id,gitfile_file.test4.id]
}

output "gitfile_commit_commit_message" {
    value = gitfile_commit.test.commit_message
}

output "checkout" {
    value = gitfile_checkout.test
}
