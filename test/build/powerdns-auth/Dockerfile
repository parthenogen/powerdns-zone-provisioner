FROM ubuntu

RUN apt update

RUN apt install -y pdns-server pdns-backend-lmdb

COPY test/build/powerdns-auth/pdns.conf /etc/powerdns/pdns.conf

CMD /usr/sbin/pdns_server --guardian=no --daemon=no --disable-syslog --log-timestamp=no --write-pid=no
