# deeble
deeble is kafka client worker to translate debezium messages

#  The db history topic is missing
% did=$(docker container ls -a | grep debezium-server | awk '{print $1}') && docker container rm $did
