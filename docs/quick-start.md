

[![Release](https://img.shields.io/github/release/fernandezvara/certsfor.svg?style=for-the-badge)](https://github.com/fernandezvara/certsfor/releases/latest)
![GitHub all releases](https://img.shields.io/github/downloads/fernandezvara/certsfor/total?style=for-the-badge)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](/LICENSE)
[![Build status](https://img.shields.io/github/workflow/status/fernandezvara/certsfor/goreleaser?style=for-the-badge)](https://github.com/fernandezvara/certsfor/actions?workflow=goreleaser)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/fernandezvara/certsfor)

Easy certificate management tool for development environments for Linux, Windows and macOS. While you are a solo developer in your workstation or distributed team of developers with many stations and servers.

Manage multiple CA without hassle. Automate different environments or servers without manual steps or magical flags that can make you lose too much time if some detail is missing.

# Quick Start

#### Using docker

>If you have `docker` installed you can just copy and paste this snippet to follow the guide. This will open an interactive console where operate the command.

```bash
# prepare local directory # this will allow file and directory creation
mkdir -p $HOME/.cfd && chmod 777 $HOME/.cfd

# run an interactive console to try
docker run -v $HOME/.cfd:/home/cfd/.cfd --entrypoint /bin/bash -it ghcr.io/fernandezvara/cfd:latest
```

>[!NOTE|label:Persistence with docker]
>A directory will be created on the home directory (`$HOME/.cfd`) where configuration and default database will be stored. If you change to binary, configuration will beheave the same.

#### Using a binary

>Download a binary from the [releases page](https://github.com/fernandezvara/certsfor/releases), for full information go to the [installation guide](./installation).

## Create your first CA interactively

```bash
> cfd create ca
✔ Common Name: myca
✔ Country (optional): 
✔ Province (optional): 
✔ Locality (optional): 
✔ Postal Code (optional): 
✔ Street (optional): 
✔ Organization (optional): MyOrg
✔ Organizational Unit (optional): Dev
✔ Hosts and IPs. (blank if finish): localhost
✔ Hosts and IPs. (blank if finish): 127.0.0.1
✔ Hosts and IPs. (blank if finish): ca.test.cfd.local
✔ Hosts and IPs. (blank if finish): 172.16.1.2
✔ Hosts and IPs. (blank if finish): 
✔ Expires in (days): 365
Use the arrow keys to navigate: ↓ ↑ → ← 
? Key Algorithm: 
↑   RSA (4096 bytes)
    ECDSA (EC-224)
    ECDSA (EC-256)
    ECDSA (EC-384)
  ▸ ECDSA (EC-521)

CA Created. ID: '2f36810e-5ef3-4bc0-9d19-cdf1d944e38d'
```

## Create your first certificate

```bash
> CFD_CA_ID="2f36810e-5ef3-4bc0-9d19-cdf1d944e38d" cfd create certificate -c ./my-cert.crt -k ./my-cert-key.crt
✔ Common Name: cert1
✔ Country (optional): 
✔ Province (optional): 
✔ Locality (optional): 
✔ Postal Code (optional): 
✔ Street (optional): 
✔ Organization (optional): MyOrg
✔ Organizational Unit (optional): Dev
✔ Hosts and IPs. (blank if finish): cert1.test.cfd.local
✔ Hosts and IPs. (blank if finish): 172.16.1.3
✔ Hosts and IPs. (blank if finish): 
✔ Expires in (days): 90
✔ RSA (4096 bytes)


Certificate Created.

> cat ./my-cert.crt 
-----BEGIN CERTIFICATE-----
MIID0DCCxxxxxxxxXXXXXXXXxxxxxxxxXXXXXXXXxxxxxxxxXXXXXXXXxxxxxxxx
................................................................
S4mpl3S4mpl3S4mpl3S4mpl3S4mpl3=
-----END CERTIFICATE-----
```
