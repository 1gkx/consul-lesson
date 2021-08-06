FROM rabbitmq:3.8.9-alpine

COPY certificates /etc/rabbitmq/certificates
# Specific for Erlang. Container used user rabbitmq
# but if copy or mount files in containerm his owner is root
# and Erlang don`t get cert, because permission denied
RUN chown -R rabbitmq:rabbitmq /etc/rabbitmq/certificates
COPY config/rabbit.conf /etc/rabbitmq/rabbitmq.conf

RUN rabbitmq-plugins enable rabbitmq_management
