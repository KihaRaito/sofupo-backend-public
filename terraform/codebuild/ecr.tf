resource "aws_ecr_repository" "sofupo-backend" {
  name = "sofupo-backend"
}

resource "aws_ecr_lifecycle_policy" "sofupo-backend" {
  repository = aws_ecr_repository.sofupo-backend.name

  policy = <<EOF
  {
    "rules": [
      {
        "rulePriority": 1,
        "description": "Keep last 30 release tagged images",
        "selection": {
          "tagStatus": "tagged",
          "tagPrefixList": ["release"],
          "countType": "imageCountMoreThan",
          "countNumber": 30
        },
        "action": {
          "type": "expire"
        }
      }
    ]
  }
EOF
}

resource "aws_ecr_repository" "sofupo-db" {
  name = "sofupo-db"
}

resource "aws_ecr_lifecycle_policy" "sofupo-db" {
  repository = aws_ecr_repository.sofupo-db.name

  policy = <<EOF
  {
    "rules": [
      {
        "rulePriority": 1,
        "description": "Keep last 30 release tagged images",
        "selection": {
          "tagStatus": "tagged",
          "tagPrefixList": ["release"],
          "countType": "imageCountMoreThan",
          "countNumber": 30
        },
        "action": {
          "type": "expire"
        }
      }
    ]
  }
EOF
}
