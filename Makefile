build:
	cd ui; yarn build
	cd ../
	CGO_ENABLED=0 GOOS=linux go build
docker:
	docker build . -t robrotheram/watch2gether


