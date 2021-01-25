FROM alpine:latest

EXPOSE 8080

ENV CFD_CONFIG /etc/cfd/config.yaml
ENV CFD_DB_TYPE ""
ENV CFD_DB_CONNECTION ""
ENV CFD_API_ENABLED ""
ENV CFD_API_ADDR ""
ENV CFD_LOG_ACCESS ""
ENV CFD_LOG_ERROR ""
ENV CFD_LOG_DEBUG ""
ENV CFD_TLS_CA ""
ENV CFD_TLS_CERT ""
ENV CFD_TLS_KEY ""
ENV CFD_TLS_FORCE ""
ENV CFD_CA_ID ""

COPY cfd /usr/local/bin/
# RUN chmod +x /usr/local/bin/cfd

RUN /usr/local/bin/cfd --config="/etc/cfd/config.yaml" configfile

ENTRYPOINT ["/usr/local/bin/cfd"]
CMD [ "--help" ]
