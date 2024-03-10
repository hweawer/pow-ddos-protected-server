dockerbuild:
	docker build --tag=pow_server:latest .


dockerbuildclient:
	docker build -t pow_client -f Dockerfile.client .
