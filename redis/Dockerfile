FROM debian:latest
LABEL authors="Dmitry"


RUN apt-get update && apt-get install -y lsb-release && apt-get clean all  -y
RUN apt-get install curl gpg -y
RUN curl -fsSL https://packages.redis.io/gpg | gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg

RUN    echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" |  tee /etc/apt/sources.list.d/redis.list

RUN    apt-get update
RUN    apt-get install redis -y

EXPOSE 6379
CMD redis-server /usr/local/etc/redis/redis.conf &&  tail -f /dev/null
