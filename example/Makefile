.PHONY: mysql-cli
mysql-cli:
	mysql -h"127.0.0.1" -P"3306" -uroot -pdebezium

setup:
	docker-compose -f docker-compose.db.yaml up --build --force-recreate -d
	docker exec debezium-mongodb bash -c '/usr/local/bin/init-inventory.sh'
	docker-compose -f docker-compose.debezium.yaml up --build --force-recreate -d

# https://github.com/debezium/docker-images/blob/main/examples/mongodb/0.10/init-inventory.sh

logdb:
	docker-compose -f docker-compose.db.yaml logs -f

logdeb:
	docker-compose -f docker-compose.debezium.yaml logs -f
