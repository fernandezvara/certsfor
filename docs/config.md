# config.yaml

The configuration file allows to configure several parts of the tool/service. Normally for a solo developer default configuration is enough.

By default, its stored in `$HOME/.cfd/config.yaml`.

**Default file:**

```yaml
api:
  addr: 127.0.0.1:8080
  enabled: true
  log:
    access:
    - stdout
    debug: false
    error:
    - stderr
ca-id: ""
db:
  connection: /home/afv/.cfd/db
  type: badger
tls:
  ca: ""
  certificate: ""
  force: false
  key: ""
```

| Key | Description | Default value |
| --- | ----------- | ------- |
| api.addr | *(string)* IP:PORT where the *client will connect* (if enabled) or the *API will listen* | `127.0.0.1:8080` |
| api.enabled | *(boolean)* Indicates when the client will connect to the cfd API | `false` |
| api.log.access | *(array<string>)* Only applies to the API. Where to store the access log. | `stdout` |
| api.log.error | *(array<string>)* Only applies to the API. Where to store the error log. | `stderr` |
| api.log.debug | *(boolean)* Only applies to the API. Write debug log. | `false` |
| ca-id | *(string)* If you will use just one CA from the service in client mode, write the UUID of the CA to use. This setting will be overwritten with `--ca-id` flag and `$CFD_CA_ID` environment variable if set. | "" |
| db.connection | *(string)* Connection string for the database store. More information on [data stores](./data-stores.md) | `$HOME/.cfd/db` |
| db.type | *(string)* Data store driver to use. | `badger` |
| tls.ca | *(string)* CA Certificate file to use for connect to the API (if client mode) or serve the API. | "" |
| tls.certificate | *(string)* Certificate file to use for connect to the API (if client mode) or serve the API. | "" |
| tls.key | *(string)* Key file to use for connect to the API (if client mode) or serve the API. | "" |
| tls.force | *(boolean)* If no certificates are set client will configure as `http`. If the API is served by trusted certificates by the system, this setting will try to connect as `https` instead. | `false` |
