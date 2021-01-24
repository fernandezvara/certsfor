# Usage

`cfd` comes with integrated help that you can refer, but here you can find examples just in case command help is not enough.

## configfile

By default, the tool comes configured with sane defaults to start working, so there is not configuration file needed. If you are a solo developer probably you won't need to change anything.

Default configuration file is `$HOME/.cfd/config.yaml` and default database will be created on `$HOME/.cfd/db`. You can use other configuration file by using the global flag `--config

## create ca

Creates a new Certification Authority based on the answers received, interactively or by using a template file. You can customize where to store the certificate file (`-c`, `--cert`) and its key (`-k`, `--key`).

By using template files you can automate all certificate creation.

```bash
Usage:
  cfd create ca [flags]

Flags:
  -c, --cert string   Certificate file location.
  -f, --file string   File with the answers in YAML format.
  -k, --key string    Key file location. NOTE: Do not share this file.
```

> [!TIP]
> The CA ID is the identification that needs to be passed to the tool to make the operations with the right CA. The ID can be passed by:
>
>- Flag: `--ca-id="your uuid"`
>- Environment variable (`$CFD_CA_ID`)
>- Configuration switch on the config file. (`ca-id`)

> [!ATTENTION|style:flat]
> **Do not share the CA key file. It can be used to hijack TLS connections and sniff (read) the traffic.**

## create certificate

Certificate creation expects the same answers than the CA.

Flags allows to set where to store a bundle certificate file, used in some services like [NGINX](https://www.nginx.org/), certificate and/or key files.

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

## get

Retrieve any certificate using its Common Name as Identifier. This command will get the certificate stored on the database if valid or will get a new updated one.

By default, when a certificate is retrieved using the cli, it will ask the CA to renew it if the time remaining for its expiration is less than the desired percent.

```bash
Flags:
  -b, --bundle string   Bundle file location.
      --ca-id string    CA Identifier. (required). [$CFD_CA_ID]
  -c, --cert string     Certificate file location.
      --cn string       Common Name. (required).
  -h, --help            help for get
  -k, --key string      Key file location.
      --renew int       Time (expresed as percent) to be used to determine if the certificate 
      must be renewed (defaults to 20%). Key remains the same. (default 20)
```

## start api

Starts cfd in daemon-mode. This mode allows remote cfd clients or simple call (like curl) usage.

!> By default API does not have any security applied, so its recommended to create certificates to secure the communication on transit.

Refer to the [API documentation](api.md) for full information.

## start webserver

Starts a simple webserver that serves files from the selected directory using a certificate selected by its common name. It's required to select the CA ID and Common Name to execute this mode.

By default webserver serves the files in the current directory, listen in all IPs and TPC port 8443. If certificate needs to be renewed it defaults to a 20% lifetime.

```bash
Flags:
      --ca-id string    CA Identifier. (required). [$CFD_CA_ID]
      --cn string       Common Name. (required).
  -h, --help            help for webserver
      --listen string   IP:TCP Port where the server will be served. Defaults to all network 
      interfaces and port 8443. (default "0.0.0.0:8443")
      --renew int       Time (expresed as percent) to be used to determine if the certificate
      must be renewed (defaults to 20 %). Key remains the same. (default 20)
      --root string     Directory where the files reside, defaults to current (.). (default ".")
```

## status

Checks if service is usable. If it's operating in a local mode it will open the database and make a simple test to ensure it's ok.

On remote mode, as API client, it will make a request to the API and will show versions on both sides.

There is no required flags since it reads the configuration file before execution.
