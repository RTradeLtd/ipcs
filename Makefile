GOOS?=linux
GOARCH?=amd64

convert:
	@GO111MODULE=on IPFS_PATH=./tmp/ipfs go run ./cmd/convert docker.io/library/alpine:latest localhost:5000/library/alpine:p2p

compare:
	@GO111MODULE=on IPFS_PATH=./tmp/ipfs go run ./cmd/compare docker.io/library/ubuntu:xenial docker.io/titusoss/ubuntu:latest

ipcs:
	@mkdir -p ./tmp/containerd/root/plugins
	@GO111MODULE=on go build -buildmode=plugin -o ./tmp/containerd/root/plugins/ipcs-$(GOOS)-$(GOARCH).so cmd/ipcs/main.go

containerd-binary:
	@mkdir -p ./bin
	@GO111MODULE=on go build -o ./bin/containerd ./cmd/containerd

containerd: containerd-binary ipcs
	@mkdir -p ./tmp
	@IPFS_PATH=./tmp/ipfs rootlesskit --copy-up=/etc \
	  --state-dir=./tmp/rootlesskit-containerd \
	    ./bin/containerd -l debug --config ./cmd/containerd/config.toml
	    
ipfs:
	@mkdir -p ./tmp
	@IPFS_PATH=./tmp/ipfs ipfs daemon --init

clean:
	@rm -rf ./tmp ./bin

.PHONY: convert registry ipcs containerd-binary containerd
