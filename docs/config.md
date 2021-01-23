# config.yaml

Default path: `$HOME/.cfd/config.yaml`

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
  tls:
    ca: ""
    certificate: ""
    key: ""
  web: false
ca-id: ""
db:
  connection: /home/afv/.cfd/db
  type: badger
```
