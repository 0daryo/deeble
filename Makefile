.PHONY: mysql-cli
mysql-cli:
	mysql -h"127.0.0.1" -P"3306" -uroot -pdebezium
setup-connector:
	curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d '{ "name": "inventory-connector", "config": { "connector.class": "io.debezium.connector.mysql.MySqlConnector", "tasks.max": "1", "database.hostname": "mysql", "database.port": "3306", "database.user": "debezium", "database.password": "dbz", "database.server.id": "184054", "database.server.name": "dbserver1", "database.include.list": "inventory", "database.history.kafka.bootstrap.servers": "kafka:9092", "database.history.kafka.topic": "dbhistory.inventory" } }'
	curl -H "Accept:application/json" localhost:8083/connectors/
	curl -i -X GET -H "Accept:application/json" localhost:8083/connectors/inventory-connector

setup:
	docker-compose -f docker-compose.db.yaml up --build --force-recreate -d
	docker exec debezium-mongodb bash -c '/usr/local/bin/init-inventory.sh'
	docker-compose -f docker-compose.debezium.yaml up --build --force-recreate -d

# https://github.com/debezium/docker-images/blob/main/examples/mongodb/0.10/init-inventory.sh

logdb:
	docker-compose -f docker-compose.db.yaml logs -f

logdeb:
	docker-compose -f docker-compose.debezium.yaml logs -f
