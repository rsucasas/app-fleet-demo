build:
	docker build -t rmq-consumer . --target production

tag:
	docker tag rmq-consumer debbabi/rmq-consumer:0.2.2

tag-latest:
	docker tag rmq-consumer debbabi/rmq-consumer:latest

push:
	docker push debbabi/rmq-consumer:0.2.2

push-latest:
	docker push debbabi/rmq-consumer:latest

start:
	docker run --rm -p 4001:4001 --name rmq-consumer -e RABBITMQ_HOST=some-rabbit rmq-consumer