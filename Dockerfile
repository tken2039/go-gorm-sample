FROM --platform=linux/amd64 mysql:8.0.28

COPY ./init.sql /docker-entrypoint-initdb.d/

VOLUME [ "/var/lib/mysql" ]
