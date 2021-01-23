# Usage

`cfd` comes with integrated help that you can refer, but here you can find examples just in case command help is not enough.

## configfile

By default, the tool comes configured with sane defaults to start working. If you are a solo developer probably you won't need to change anything.

Configuration directory is `$HOME/.cfd` where the `config.yaml` lives if you want to create it and default database.

## create ca

Creates a new Certification Authority based on the answers received, interactively or by using a template file. You can customize where to store the certificate file (-c, --cert) and its key (-k, --key).

By using template files you can automate all certificate creation.

```bash
Usage:
  cfd create ca [flags]

Flags:
  -c, --cert string   Certificate file location.
  -f, --file string   File with the answers in YAML format.
  -k, --key string    Key file location. NOTE: Do not share this file.
```

?> The CA ID is the identification that needs to be passed to the tool to make the operations with the right CA. The ID can be passed by:<br>- Flag: --ca-id="your uuid"<br>- Environment variable ($CFD_CA_ID)<br>- Configuration switch on the config file. (ca-id)


!> **Do not share the CA key file. It can be used to hijack TLS connections and sniff (read) the traffic.**

## create certificate

Certificate creation expects the same answers than the CA.

Flags allows to set where to store a bundle certificate file, used in some services like [NGINX](https://www.nginx.org/), certificate and key files.

```bash
Usage:
  cfd create certificate [flags]

Aliases:
  certificate, cert

Flags:
  -b, --bundle string   Bundle file location.
      --ca-id string    CA Identifier. (required). [$CFD_CA_ID]
  -c, --cert string     Certificate file location.
  -f, --file string     File with the answers in YAML format.
  -k, --key string      Key file location. NOTE: Do not share this file.
```

