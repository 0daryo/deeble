.PHONY: mysql-cli
mysql-cli:
	mysql -h"127.0.0.1" -P"3306" -uroot -pdebezium
setup-connector:
	curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d '{ "name": "inventory-connector", "config": { "connector.class": "io.debezium.connector.mysql.MySqlConnector", "tasks.max": "1", "database.hostname": "mysql", "database.port": "3306", "database.user": "debezium", "database.password": "dbz", "database.server.id": "184054", "database.server.name": "dbserver1", "database.include.list": "inventory", "database.history.kafka.bootstrap.servers": "kafka:9092", "database.history.kafka.topic": "dbhistory.inventory" } }'
	curl -H "Accept:application/json" localhost:8083/connectors/
	curl -H "Accept:application/json" localhost:8083/connectors/inventory-connector/config | jq .
