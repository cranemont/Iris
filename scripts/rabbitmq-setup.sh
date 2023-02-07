#!/bin/sh

while ! nc -z "$RABBITMQ_HOST" "$RABBITMQ_PORT"; do sleep 3; done
>&2 echo "rabbitmq is up - server running..."

# install rabbitmqadmin
wget http://rabbitmq:15672/cli/rabbitmqadmin
chmod +x rabbitmqadmin

# Make an Exchange
./rabbitmqadmin -H rabbitmq -u skku -p 1234 declare exchange name="$JUDGE_EXCHANGE_NAME" type=direct

# Make queues
./rabbitmqadmin -H rabbitmq -u skku -p 1234 declare queue name="$JUDGE_RESULT_QUEUE_NAME" durable=true
./rabbitmqadmin -H rabbitmq -u skku -p 1234 declare queue name="$JUDGE_SUBMISSION_QUEUE_NAME" durable=true

# Make bindings
./rabbitmqadmin -H rabbitmq -u skku -p 1234 declare binding source="$JUDGE_EXCHANGE_NAME"\
                                destination_type=queue destination="$JUDGE_RESULT_QUEUE_NAME" routing_key="$JUDGE_RESULT_ROUTING_KEY"
./rabbitmqadmin -H rabbitmq -u skku -p 1234 declare binding source="$JUDGE_EXCHANGE_NAME"\
                                destination_type=queue destination="$JUDGE_SUBMISSION_QUEUE_NAME" routing_key="$JUDGE_SUBMISSION_ROUTING_KEY"


rm rabbitmqadmin

