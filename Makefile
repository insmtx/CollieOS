PROJECT := SingerOS
REGISTRY ?= registry.yygu.cn/insmtx/

.PHONY: docker-build docker-build-singer docker-build-skill-proxy docker-push docker-release docker-run

docker-build:
	docker build -t $(REGISTRY)$(PROJECT):latest -f deployments/build/Dockerfile .

docker-build-singer:
	docker build --target singer -t $(REGISTRY)$(PROJECT)-singer:latest -f deployments/build/Dockerfile .

docker-build-skill-proxy:
	docker build --target skill-proxy -t $(REGISTRY)$(PROJECT)-skill-proxy:latest -f deployments/build/Dockerfile .

docker-build-all: docker-build-singer docker-build-skill-proxy

docker-push:
	docker push $(REGISTRY)$(PROJECT):latest

docker-push-singer:
	docker push $(REGISTRY)$(PROJECT)-singer:latest

docker-push-skill-proxy:
	docker push $(REGISTRY)$(PROJECT)-skill-proxy:latest

docker-push-all: docker-push-singer docker-push-skill-proxy

docker-release: docker-build docker-push

docker-release-all: docker-build-all docker-push-all

docker-run:
	-docker rm -f $(PROJECT)-dev
	docker run -d --name $(PROJECT)-dev -p 8080:8080 $(REGISTRY)$(PROJECT):latest

docker-run-singer:
	-docker rm -f $(PROJECT)-singer-dev
	docker run -d --name $(PROJECT)-singer-dev -p 8080:8080 $(REGISTRY)$(PROJECT)-singer:latest

docker-run-skill-proxy:
	-docker rm -f $(PROJECT)-skill-proxy-dev
	docker run -d --name $(PROJECT)-skill-proxy-dev -p 8081:8080 $(REGISTRY)$(PROJECT)-skill-proxy:latest

install-protoc:
	which protoc || (echo "protoc not found, please install it first" && exit 1)
	wget -O /tmp/protoc-34.0-linux-x86_64.zip https://github.com/protocolbuffers/protobuf/releases/download/v34.0/protoc-34.0-linux-x86_64.zip
	unzip /tmp/protoc-34.0-linux-x86_64.zip -d protoc
	sudo mv protoc/bin/* /usr/local/bin/
	sudo mv protoc/include/* /usr/local/include/

generate-proto-go:
	which protoc-gen-go || go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
	which protoc-gen-go-grpc || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	mkdir -p gen
	protoc --go_out=gen --go-grpc_out=gen --proto_path=proto --proto_path=third_party proto/**/*.proto

generate-proto-node:
	which protoc-gen-es || npm install -g @bufbuild/protoc-gen-es
	protoc --plugin=protoc-gen-es --es_out=. proto/*.proto
	