VER=0.11.3

build-server:	
	cd server; CGO_ENABLED=0 GOOS=linux go build -o ../.

build-ui:
	cd ui; npm ci; npm run build

build-docker:
	docker build . --build-arg VER=${VER} -t robrotheram/watch2gether:${VER}
	docker build . --build-arg VER=${VER} -t robrotheram/watch2gether:latest
	
publish:
	docker push robrotheram/watch2gether:${VER}
	docker push robrotheram/watch2gether:latest

build: build-docker

run: 
	docker-compose up -d
