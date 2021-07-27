# Directory Structure
```
.
├── backend  # Goのhandler/utils/init関数
├── dockerfiles  # Dockerfile
├── kubernetes
│   └── kustomize
│       ├── base  # Kustomizeのbase
│       └── overlays
│           ├── develop  # 開発用のmanifestファイル
│           └── production  # 本番用のmanifestファイル
├── model  # DBモデル
├── repository  # CRUDで使用するメソッド
└── terraform
    ├── codebuild  # CodeBuildのIaC
    └── eks  # EKSのIaC
```

# 実行方法
```
$ export $(cat .env | xargs)
$ go run main.go
```

# UnitTest
```
$ go test ./... -v -cover
```

# Vet
```
$ go vet ./...
```

# Docker
```
# Build
$ docker build -t db:${DB_IMAGE_TAG} -f dockerfiles/Dockerfile.db .
$ docker build -t server:${BACKEND_IMAGE_TAG} -f dockerfiles/Dockerfile.backend .

# Run
$ docker run -d --name db --hostname db -p ${DB_PORT}:${DB_PORT} --env-file .env -it db:${DB_IMAGE_TAG}
$ docker run -d --name server --hostname server -p ${BACKEND_PORT}:${BACKEND_PORT} --env-file .env -e DB_HOST=$(docker exec -it db hostname -i) -it server:${BACKEND_IMAGE_TAG}
```

# E2ETest
```
$ newman run SofupoE2ETest.postman_collection.json
```

# ECR
```
$ aws ecr get-login-password --region ap-northeast-1 --profile ${AWS_PROFILE} | docker login --username AWS --password-stdin ${REGISTRY_ID}.dkr.ecr.ap-northeast-1.amazonaws.com

$ docker build -t sofupo-backend:${BACKEND_IMAGE_TAG} -f dockerfiles/Dockerfile.backend .
$ REGISTRY_ID=`aws ecr create-repository --repository-name sofupo-backend --region ap-northeast-1 --profile ${AWS_PROFILE} |jq -r .repository.registryId`
$ docker tag sofupo-backend:${BACKEND_IMAGE_TAG} ${REGISTRY_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/sofupo-backend:${BACKEND_IMAGE_TAG}
$ docker push ${REGISTRY_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/sofupo-backend:${BACKEND_IMAGE_TAG}

$ docker build -t sofupo-db:${DB_IMAGE_TAG} -f dockerfiles/Dockerfile.db .
$ REGISTRY_ID=`aws ecr create-repository --repository-name sofupo-db --region ap-northeast-1 --profile ${AWS_PROFILE} |jq -r .repository.registryId`
$ docker tag sofupo-db:${DB_IMAGE_TAG} ${REGISTRY_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/sofupo-db:${DB_IMAGE_TAG}
$ docker push ${REGISTRY_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/sofupo-db:${DB_IMAGE_TAG}
```
