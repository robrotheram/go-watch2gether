VER=0.11.0

build-server:
	sed -i '/ /s/".*"/"${VER}"/' server/pkg/datastore/version.go
	cd server; CGO_ENABLED=0 GOOS=linux go build -o ../.

build-ui:
	cd ui; npm ci; npm run build
	sed -i 's/{WATCH2GETHER_VERSION}/$(VER)/g' ui/build/index.html

build-docker:
	docker build . --build-arg VER=${VER} -t robrotheram/watch2gether:${VER}
	docker build . --build-arg VER=${VER} -t robrotheram/watch2gether:latest
	
publish:
	docker push robrotheram/watch2gether:${VER}
	docker push robrotheram/watch2gether:latest

build: build-docker

run: 
	docker-compose up -d
