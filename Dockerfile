FROM scratch

COPY bin /usr/local/bin 

ENTRYPOINT ["/usr/local/bin/powerdns-zone-provisioner"]
