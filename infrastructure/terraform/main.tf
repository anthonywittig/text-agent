terraform {
  backend "s3" {
    key    = "project/text-agent/terraform.tfstate"
    region = "us-west-2"
  }
}

provider "aws" {
  region = "us-west-2"
}

locals {
  region = "us-west-2"
}

variable "aws_profile" {
  type = string
}

variable "git_sha" {
  type = string
}
