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
terraform destroy

rm ~/.kube/config
aws eks --region ${AWS_REGION} update-kubeconfig --name ${EKS_CLUSTER_NAME} --profile ${AWS_PROFILE}

terraform state list
terraform state rm module.eks.kubernetes_config_map.aws_auth
```
