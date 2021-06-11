VER=v0.4.5

build-server:
	cd server; CGO_ENABLED=0 GOOS=linux go build -o ../.

build-ui:
	cd ui; yarn; yarn build
	sed -i 's/{WATCH2GETHER_VERSION}/$(VER)/g' ui/build/index.html

build-docker:
	docker build . -t robrotheram/watch2gether:${VER}
	docker build . -t robrotheram/watch2gether:latest
	
publish:
	docker push robrotheram/watch2gether:${VER}
	docker push robrotheram/watch2gether:latest

build: build-server build-ui build-docker

run: 
	docker-compose up -d
