# Deploy
```
tfenv list-remote
tfenv install 0.14.0
tfenv use 0.14.0

terraform login
terraform init

. .env

terraform plan
terraform apply

aws s3 rm s3://sofupo-backend-artifact/ --recursive --profile ${AWS_PROFILE}
terraform destroy
```
