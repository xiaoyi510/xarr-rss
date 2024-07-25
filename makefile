r:
	docker build --no-cache -t xiaoyi510/xarr-rss . && docker push xiaoyi510/xarr-rss
	docker build  -t xiaoyi510/xarr-rss:v1.4 . && docker push xiaoyi510/xarr-rss:v1.4

push_prd:
	docker buildx build --no-cache --platform linux/arm64,linux/amd64,linux/arm -t xiaoyi510/xarr-rss:v2.0.3.4 -t xiaoyi510/xarr-rss . --push
