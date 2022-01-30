# Usage

`cfd` comes with integrated help that you can refer, but here you can find examples just in case command help is not enough.

Global Flags:

| Flag | Explanation | Environment Var | Required |
| ---- | ----------- | --------------- | :------: |
| `--config` | Config file location. (Default: `$HOME/.cfg/config.yaml`) | CFD_CONFIG | |
| `-q`, `--quiet` | Supress the command output (Default: `false`)) | CFD_QUIET | |

## configfile

**Usage:** `cfd configfile`

By default, the tool comes configured with sane defaults to start working. If you are a solo developer probably, creating a configuration file won't be needed.

Default configuration file is `$HOME/.cfd/config.yaml` and default database will be created on `$HOME/.cfd/db`. You can use other configuration file by using the global flag `--config`.

## create ca

Creates a new Certification Authority based on the answers received, interactively or by using a template file. You can customize where to store the certificate file (`-c`, `--cert`) and its key (`-k`, `--key`).

By using template files, you can automate all certificate creation.

**Usage:** `cfd create ca [flags]`

| Flag | Explanation | Environment Var | Required |
| ---- | ----------- | --------------- | -------- |
| `-c`, `--cert` | Where to store the CA Certificate after its creation. | | |
| `-k`, `--key` | Where to store the key file. *Normally, you don't need this file*. | | |
| `-f`, `--file` | File with the answers in YAML format. | | |

> [!TIP]
> The CA ID identifies the group of certificates to use. It needs to passed to the tool to make the operations with the right CA. The ID can be passed by:
>
>- Flag: `--ca-id="your uuid"`
>- Environment variable (`$CFD_CA_ID`)
>- Configuration switch on the config file. (`ca-id`)

> [!ATTENTION]
> **Do not share the CA key file. It can be used to hijack TLS connections and sniff (read) the traffic.**

## create certificate

Certificate creation expects the same answers as the CA.

**Usage:** `cfd create certificate [flags]`

Aliases: `create certificate`, `create cert`

| Flag | Explanation | Environment Var | Required |
| ---- | ----------- | --------------- | :------: |
| `-b`, `--bundle` | Bundle file location. Some services like [NGINX](https://www.nginx.org) uses this kind of file. | | |
| `--ca-cert` | Where to store the CA Certificate. | | |
| `--ca-id` | ID of the CA to interact to. | CFD_CA_ID | :heavy_check_mark: |
| `-c`, `--cert` | Where to store the Certificate after its creation. | | |
| `-k`, `--key` | Where to store the key file. | | |
| `--pfx` | Where to store the Certificate in pkcs12 format. | | |
| `--pfx-password` | PFX file password (Default: `changeit`) | | |
| `-f`, `--file` | File with the answers in YAML format. | | |

>[!NOTE|label:Output to console]
>Certificate, key, CA certificate, bundle and PFX files can be echoed to console by using `out` or `stdout` as value.
>
> Ex: `cdf create cert -f ./template.yaml -c stdout`

## create template

Writes an empty certificate creation template as YAML file.

**Usage:** `cfd create template [flags]`

| Flag | Explanation | Environment Var | Required |
| ---- | ----------- | --------------- | :------: |
| `--ca` | Template is for a CA | | |
| `-f`, `--file` | Where to store the template file in YAML format. If not provided it will be requested interactively. | | |

## get certificate / get cert

Retrieve any certificate using its Common Name as Identifier. This command will get the certificate stored on the database if valid or will get a new updated one.

By default, when a certificate is retrieved using the CLI, it will ask the CA to renew it if the time remaining for its expiration is less than the desired percent.

**Usage:** `cfd get certificate [flags]`

| Flag | Explanation | Environment Var | Required |
| ---- | ----------- | --------------- | :------: |
| `--ca-id` | ID of the CA to interact to. | CFD_CA_ID | :heavy_check_mark: |
| `--cn` | Common name of the Certificate to retrieve. | | :heavy_check_mark: |
| `-c`, `--cert` | Where to store the Certificate after its creation. | | |
| `-k`, `--key` | Where to store the key file. | | |
| `-b`, `--bundle` | Bundle file location. | | |
| `--ca-cert` | Where to store the CA Certificate. | | |
| `--renew` | Time (expresed as percent) to determine if the certificate must be renewed **(defaults to 20%)**. | | |
| `--pfx` | Where to store the Certificate in pkcs12 format. | | |
| `--pfx-password` | PFX file password (Default: `changeit`) | | |

>[!NOTE|label:Output to console]
>Certificate, key, CA certificate, bundle and PFX files can be echoed to console by using `out` or `stdout` as value.
>
> Ex: `cdf get cert --ca-id <uuid> --cn <common-name> -c stdout`

## info certificate / info cert

**Usage:** `cfd info certificate [flags]`

| Flag | Explanation | Environment Var | Required |
| ---- | ----------- | --------------- | :------: |
| `--ca-id` | ID of the CA to interact to. | CFD_CA_ID | if `--cn` |
| `--cn` | Common name of the Certificate to query its information. | | if `--ca-id` |
| `-f`, `--file` | Certificate file to get its information. | | |
| `-u`, `--url` | URL to get its information | | |
| `--insecure` | Do no make validations on the connection. (Only used if `--url`). | | |
| `--timeout` | Timeout making the request to the URL. (Only used if `--url`). | | |
| `--markdown` | Return data in `1markdown` format. | | |
| `--csv` | Return data as CSV. | | |

```bash
> # example
> cfd info cert --url www.google.com      
┌────────────────┬───────────────────────────────────────────────────────────────────┬──────────────────┬─────────────────────────────┬────────────────┬───────┐
│ Common Name    │ Distinguished Name                                                │ Not Before       │ Expires                     │ SANs           │ CA?   │
├────────────────┼───────────────────────────────────────────────────────────────────┼──────────────────┼─────────────────────────────┼────────────────┼───────┤
│ www.google.com │ CN=www.google.com,O=Google LLC,L=Mountain View,ST=California,C=US │ 19/01/2021 08:04 │ 63 days (13/04/2021 08:04)  │ www.google.com │ false │
│ GTS CA 1O1     │ CN=GTS CA 1O1,O=Google Trust Services,C=US                        │ 15/06/2017 00:00 │ 309 days (15/12/2021 00:00) │                │ true  │
└────────────────┴───────────────────────────────────────────────────────────────────┴──────────────────┴─────────────────────────────┴────────────────┴───────┘
```

## list certificate / list certificates / list cert

Return a list with all the certificates on the CA.

**Usage:** `cfd list certificate [flags]`

| Flag | Explanation | Environment Var | Required |
| ---- | ----------- | --------------- | :------: |
| `--ca-id` | ID of the CA to interact to. | CFD_CA_ID | :heavy_check_mark: |
| `--csv` | Output as CSV. | |  |
| `--md` | Output as Markdown. | | |

```bash
>
> # example
> cdf get list cert
┌─────────────┬──────────────────────────────────────┬─────────────────────────────┐
│ Common Name │ Distinguished Name                   │ Expires In                  │
├─────────────┼──────────────────────────────────────┼─────────────────────────────┤
│ cert3       │ CN=cert3,OU=dev,O=CFD,ST=Madrid,C=ES │ 86 days (04/05/2021 00:08)  │
│ ca          │ CN=ca                                │ 361 days (02/02/2022 22:55) │
│ cert1       │ CN=cert1                             │ 86 days (04/05/2021 00:06)  │
│ cert2       │ CN=cert2                             │ 86 days (04/05/2021 00:07)  │
└─────────────┴──────────────────────────────────────┴─────────────────────────────┘


```

## start api

Starts cfd in daemon-mode. This mode allows remote cfd clients or simple call (like curl) usage.

**Usage:** `cfd start api`

> [!ATTENTION]
> By default API **does not have any security applied**, so its recommended to create certificates to secure the communication on transit, before its use.

Refer to the [API endpoints documentation](api.md) for its usage.

>If you are using a HA capable data store, you run many instances that will behave as one (normally behind a load balancer)
>
>When using custom certificates on the API servers you can proxy TCP traffic directly from the load balancer to ensure point-to-point in transit data encryption.

## start webserver

Starts a simple webserver that serves files from the selected directory using a certificate selected by its common name. It's required to select the CA ID and Common Name to execute this mode.

By default webserver serves the files in the current directory, listen in all IPs and TPC port 8443. If certificate needs to be renewed it defaults to a 20% lifetime.

**Usage:** `cfd start webserver [flags]`

| Flag | Explanation | Environment Var | Required |
| ---- | ----------- | --------------- | :------: |
| `--ca-id` | ID of the CA to interact to. | CFD_CA_ID | :heavy_check_mark: |
| `--cn` | Common name of the Certificate to use. | | :heavy_check_mark: |
| `--root` | Directory with the files to serve. (Default: `.`) | | |
| `--listen` | IP:TCP Port where the content will be serve. (Default: `0.0.0.0:8443`) | | |
| `--renew` | Time (expresed as percent) to determine if the certificate must be renewed **(defaults to 20%)**. | | |

## status

**Usage:** `cfd status`

Checks if service is usable. If it's operating in a local mode it will open the database and make a simple test to ensure it's ok.

On remote mode, as API client, it will make a request to the API and will show versions on both sides.
