FROM golang:latest

ADD setup-mariadb.sh /setup-mariadb.sh

ENV MARIADB_VERSION 10.1
ENV MYSQL_ROOT_PASSWORD root

# Update and Fix Language
RUN \
 apt-get update && apt-get -y upgrade &&\
 apt-get install -y --no-install-recommends locales &&\
 echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen &&\
 locale-gen en_US.UTF-8 &&\
 /usr/sbin/update-locale LANG=en_US.UTF-8

# Install mariadb
RUN ["/bin/bash", "-c", "bash /setup-mariadb.sh"]

RUN go get github.com/gorilla/mux
RUN go get github.com/ziutek/mymysql/godrv

# Clean up
#RUN \
# apt-get autoclean && apt-get clean && apt-get autoremove


WORKDIR /app

COPY . .

RUN go build -o main . 

EXPOSE 8000

RUN mysqld_safe & until mysqladmin ping; do sleep 1; done && \
    mysql -uroot -e "CREATE DATABASE team_members;" && \
    mysql -uroot -e "CREATE USER 'tm'@'localhost' IDENTIFIED BY 'tm';" && \
    mysql -uroot -e "GRANT ALL PRIVILEGES ON team_members.* TO 'tm'@'localhost';" && \
    mysql -u root team_members < db_import.sql

CMD /etc/init.d/mysql start && ./main

