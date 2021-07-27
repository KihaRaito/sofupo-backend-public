provider "github" {
  owner = "KihaRaito"
}

resource "github_repository_webhook" "github_repository_webhook" {
  repository = "sofupo-backend"

  configuration {
    url          = aws_codepipeline_webhook.webhook.url
    secret       = "VeryRandomStringMoreThan20Byte!"
    content_type = "json"
    insecure_ssl = false
  }

  events = ["push"]
}
