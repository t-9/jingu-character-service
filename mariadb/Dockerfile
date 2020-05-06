FROM mariadb:10.5.2-bionic

COPY ./charset.cnf /etc/mysql/conf.d/charset.cnf
RUN chmod 644 /etc/mysql/conf.d/*

RUN apt-get update && \
  apt-get install -y locales && \
  rm -rf /var/lib/apt/lists/* && \
  echo "ja_JP.UTF-8 UTF-8" > /etc/locale.gen && \
  locale-gen ja_JP.UTF-8
ENV LC_ALL ja_JP.UTF-8