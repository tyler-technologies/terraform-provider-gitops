provider "gitops" {
    repo_url = "../example.git"
    branch = "master"
    path = "checkout"
}

resource "gitops_checkout" "test" {}

output "gitops_checkout_path" {
    value = gitops_checkout.test.path
}

resource "gitops_file" "test" {
    checkout = gitops_checkout.test.id
    path = "terraform"
    contents = "Terraform making commits"
}
resource "gitops_file" "test1" {
    checkout = gitops_checkout.test.id
    path = "terraform1"
    contents = "Terraform making commits"
}
resource "gitops_file" "test2" {
    checkout = gitops_checkout.test.id
    path = "terraform2"
    contents = "Terraform making commits"
}
resource "gitops_file" "test3" {
    checkout = gitops_checkout.test.id
    path = "terraform3"
    contents = "Terraform making commits"
}
resource "gitops_file" "test4" {
    checkout = gitops_checkout.test.id
    path = "terraform4"
    contents = "Terraform making commits"
}

resource "gitops_commit" "test" {
    commit_message = "Created by terraform gitops_commit"
    handles = [gitops_file.test.id, gitops_file.test1.id,gitops_file.test2.id,gitops_file.test3.id,gitops_file.test4.id]
}

output "gitops_commit_commit_message" {
    value = gitops_commit.test.commit_message
}

output "checkout" {
    value = gitops_checkout.test
}
