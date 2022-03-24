# deeble
deeble is message client worker to translate debezium messages

# Trouble shooting

##  The db history topic is missing
The reason may be orphan container
```
% did=$(docker container ls -a | grep debezium-server | awk '{print $1}') && docker container rm $did
```
# mongo
```
% pwd
deeble/example
% make setup
```

also refer.
https://github.com/debezium/debezium-examples/tree/main/debezium-server-mongo-pubsub#how-to-run

