IMAGE_NAME=stock-god-scraper
TAG=latest

build:
	docker stop $(IMAGE_NAME)
	docker rm $(IMAGE_NAME)
	docker rmi $(IMAGE_NAME)
	# 注意这里的缩进，必须使用Tab键，而不是空格
	docker build -t $(IMAGE_NAME):$(TAG) .
	docker run -d --name $(IMAGE_NAME) $(IMAGE_NAME):$(TAG)
