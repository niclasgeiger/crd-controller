build:
	@./cli/build-image.sh
start:
	@kubectl create -f k8s/resources.yml
stop:
	@kubectl delete -f k8s/resources.yml