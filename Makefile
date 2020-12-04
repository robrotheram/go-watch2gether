VER=v0.1.2

build-server:
	CGO_ENABLED=0 GOOS=linux go build

build-ui:
	cd ui; yarn build
	sed -i 's/{WATCH2GETHER_VERSION}/$(VER)/g' ui/build/index.html

build-docker:
	docker build . -t robrotheram/watch2gether:${VER}

publish:
	docker push robrotheram/watch2gether:${VER}

build: build-server build-ui build-docker
