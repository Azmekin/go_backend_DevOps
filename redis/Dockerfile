FROM debian:latest
LABEL authors="Dmitry"


RUN apt-get update && apt-get install -y lsb-release && apt-get clean all  -y
RUN apt-get install curl gpg -y
RUN curl -fsSL https://packages.redis.io/gpg | gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg

RUN    echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" |  tee /etc/apt/sources.list.d/redis.list

RUN    apt-get update
RUN    apt-get install redis -y

COPY redis.conf /usr/local/etc/redis/redis.conf
COPY users.acl /etc/redis/users.acl
COPY start_redis.sh /etc/redis/start_redis.sh
COPY ca.crt /etc/redis/tls/ca.crt
COPY redis.key /etc/redis/tls/redis.key
COPY redis.crt /etc/redis/tls/redis.crt
COPY redis.dh /etc/redis/tls/redis.dh


EXPOSE 6379
CMD ["/etc/redis/start_redis.sh"]