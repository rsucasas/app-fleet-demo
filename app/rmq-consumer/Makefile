build:
	docker build -t rmq-consumer . --target production

tag:
	docker tag rmq-consumer debbabi/rmq-consumer:0.1.0

push:
	docker push debbabi/rmq-consumer:0.1.0

start:
	docker run --rm -p 4001:4001 --name rmq-consumer -e RABBITMQ_HOST=some-rabbit rmq-consumer