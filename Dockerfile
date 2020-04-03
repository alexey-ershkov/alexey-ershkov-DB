FROM ubuntu:18.04

FROM ubuntu:18.04

MAINTAINER Howle

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
# Обвновление списка пакетов
RUN apt-get -y update



ENV PGVER 10
RUN apt-get install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER farcoad WITH SUPERUSER PASSWORD 'postgres';" &&\
    createdb -O farcoad forum &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root


RUN apt-get install -y curl
RUN curl —silent —location https://deb.nodesource.com/setup_13.x | bash -
RUN apt-get install -y nodejs
RUN apt-get install -y build-essential

# создание директории приложения
WORKDIR /usr/src/app

# установка зависимостей
# символ астериск ("*") используется для того чтобы по возможности
# скопировать оба файла: package.json и package-lock.json
COPY package*.json ./

RUN npm install
# Если вы создаете сборку для продакшн
# RUN npm ci --only=production

# копируем исходный код
COPY . .

EXPOSE 5000
ENV PGPASSWORD postgres
CMD service postgresql start && psql -h localhost -d forum -U farcoad -p 5432 -a -q -f ./init/init.sql && node index.js
