variable "aws_region" {
    default = "ap-northeast-1"
}
variable "aws_profile" {}

locals {
    cluster_name    = "sofupo-eks"
    cluster_version = "1.19"
}

provider "aws" {
    region  = var.aws_region
    profile = var.aws_profile
}
