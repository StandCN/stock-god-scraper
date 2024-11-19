IMAGE_NAME=stock-god-scraper
TAG=0.0.1

build:
	# 注意这里的缩进，必须使用Tab键，而不是空格
	docker build -t $(IMAGE_NAME):$(TAG) .
