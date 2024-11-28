IMAGE_NAME=stock-god-scraper
TAG=latest

build:
	@if [ $$(docker ps -a -q -f name=$(IMAGE_NAME)) ]; then \
		docker stop $(IMAGE_NAME); \
		docker rm $(IMAGE_NAME); \
	fi
	@if [ $$(docker images -q $(IMAGE_NAME)) ]; then \
		docker rmi $(IMAGE_NAME); \
	fi
	# 注意这里的缩进，必须使用Tab键，而不是空格
	docker build -t $(IMAGE_NAME):$(TAG) .
	docker run -d --name $(IMAGE_NAME) $(IMAGE_NAME):$(TAG)
