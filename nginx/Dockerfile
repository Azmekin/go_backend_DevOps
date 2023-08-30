FROM debian:latest
LABEL authors="Dmitry"


RUN  apt-get update
RUN  apt-get install nginx -y
RUN rm /etc/nginx/nginx.conf

COPY nginx.conf /etc/nginx/

CMD nginx -c /etc/nginx/nginx.conf && nginx -s reload && tail -f /dev/null

