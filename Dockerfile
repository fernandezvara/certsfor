ARG ARCH

FROM multiarch/alpine:${ARCH}-latest-stable

ARG USER=cfd
ENV HOME /home/${USER}

RUN apk add --no-cache bash
RUN adduser -D -s /bin/bash ${USER}

USER ${USER}
WORKDIR ${HOME}

# binary
COPY cfd /usr/local/bin/

# cfd variables
ENV CFD_CONFIG ${HOME}/.cfd/config.yaml
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

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/cfd"]
CMD [ "--help" ]
