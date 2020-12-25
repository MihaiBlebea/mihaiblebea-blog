IMAGE := mihaiblebea-blog
TAG := 0.1

setup: build run

bundle:
	go-bindata -o=assets/bindata.go --pkg=assets static/templates/... static/markdown/...

build:
	docker build -t ${IMAGE}:${TAG} .

run:
	docker run -d --rm --name ${IMAGE} --env-file=.env -p 8099:8099 ${IMAGE}:${TAG}

remove:
	docker stop ${IMAGE} && docker rm ${IMAGE}