# Deploy
```
kubectl create ns sofupo

# Develop
kubectl create configmap config --from-env-file=.env.dev -n sofupo
kubectl create secret generic secret --from-env-file=.aws -n sofupo
kustomize build kustomize/overlays/develop |REGISTRY_ID=$REGISTRY_ID BACKEND_IMAGE_TAG=$BACKEND_IMAGE_TAG BACKEND_PORT=$BACKEND_PORT NODE_PORT=$NODE_PORT DB_IMAGE_TAG=$DB_IMAGE_TAG envsubst |kubectl apply -f -

# Production
kubectl create configmap config --from-env-file=.env -n sofupo
kubectl create secret generic secret --from-env-file=.aws -n sofupo
kustomize build kustomize/overlays/production |REGISTRY_ID=$REGISTRY_ID BACKEND_IMAGE_TAG=$BACKEND_IMAGE_TAG BACKEND_PORT=$BACKEND_PORT NODE_PORT=$NODE_PORT envsubst |kubectl apply -f -

curl <DNS>:80
curl https://<DNS> --insecure
curl https://${DOMAIN}

kustomize build kustomize/overlays/develop |kubectl delete -f -
kustomize build kustomize/overlays/production |kubectl delete -f -
```
# Setting Docker Registry
```
TOKEN=`aws ecr get-authorization-token --region=${REGISTRY_REGION} --profile ${AWS_PROFILE} --registry-ids ${REGISTRY_ID} --output text --query authorizationData[].authorizationToken | base64 -D | cut -d ':' -f 2`
kubectl delete secret --ignore-not-found registry-auth
kubectl create secret docker-registry registry-auth \
 --docker-server=https://${REGISTRY_ID}.dkr.ecr.${REGISTRY_REGION}.amazonaws.com \
 --docker-username=AWS \
 --docker-password="${TOKEN}"
kubectl patch serviceaccount default -p '{"imagePullSecrets":[{"name":"registry-auth"}]}'
```
# ArgoCD
```
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl port-forward svc/argocd-server -n argocd ${BACKEND_PORT}:443

kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d '/' -f 2

argocd login 127.0.0.1:${BACKEND_PORT} --insecure
argocd account update-password

argocd repo add https://github.com/KihaRaito/sofupo-backend --username <username> --password <token>

open https://127.0.0.1:${BACKEND_PORT}

aws ecr describe-images --profile ${AWS_PROFILE} --repository-name sofupo-backend

kubectl apply -n argocd -f argocd-app.yaml

kubectl delete -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl delete -n argocd -f argocd-app.yaml
```
# Grafana
```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add stable https://charts.helm.sh/stable
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack -f values.yaml -n monitoring --create-namespace

kubectl port-forward svc/prometheus-grafana -n monitoring 3000:80
open http://127.0.0.1:3000

「3119」のDashboardをImportしてLoad
```
# Logging
```
kubectl create ns amazon-cloudwatch
kubectl create configmap fluent-bit-cluster-info \
    --from-literal=cluster.name=sofupo-eks \
    --from-literal=http.server=On \
    --from-literal=http.port=2020 \
    --from-literal=read.head=Off \
    --from-literal=read.tail=On \
    --from-literal=logs.region=ap-northeast-1 -n amazon-cloudwatch

kubectl apply -f https://raw.githubusercontent.com/aws-samples/amazon-cloudwatch-container-insights/latest/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/fluent-bit/fluent-bit.yaml
```
# References
* https://blog.star-flare.com/2020/12/24/start-gitops-using-argo-cd/
