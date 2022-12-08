rabbitmq-dev:
	docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management
	
build:
	docker build -t rmq-publisher . --target production

tag:
	docker tag rmq-publisher debbabi/rmq-publisher:0.1.2

tag-latest:
	docker tag rmq-publisher debbabi/rmq-publisher:latest

push:
	docker push debbabi/rmq-publisher:0.1.2

push-latest:
	docker push debbabi/rmq-publisher:latest

start:
	docker run --rm -p 4001:4001 --name rmq-publisher -e RABBITMQ_HOST=some-rabbit rmq-publisher