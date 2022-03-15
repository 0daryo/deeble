# deeble
deeble is kafka client worker to translate debezium messages

#  The db history topic is missing
% did=$(docker container ls -a | grep debezium-server | awk '{print $1}') && docker container rm $did

# mongo
https://github.com/debezium/debezium-examples/tree/main/debezium-server-mongo-pubsub#how-to-run

```
docker exec mongodb bash -c '/usr/local/bin/init-inventory.sh'
```
