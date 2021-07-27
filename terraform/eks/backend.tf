terraform {
  backend "remote" {
    hostname     = "app.terraform.io" // Terraform Cloudのホスト上に保存
    organization = "kiharaito"
    workspaces {
      name = "sofupo-eks"
    }
  }
}
