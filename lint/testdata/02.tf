provider "google" {
  project = "my-project-id"
  region  = "us-central1"
}
provider "aws" {
  version = "~> 2.0"
  region  = "us-east-1"
}
provider "google" {
  project = "my-project-id"
  region  = "us-west1"
  alias   = "west"
}
