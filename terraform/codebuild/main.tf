terraform {
  required_version = "0.14.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "2.20.0"
    }
    github = {
      source  = "integrations/github"
      version = "4.12.1"
    }
  }
}

variable "aws_region" {
    default = "ap-northeast-1"
}
variable "aws_profile" {}

provider "aws" {
    region  = var.aws_region
    profile = var.aws_profile
}
