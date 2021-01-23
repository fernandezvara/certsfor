# Quick Start

Easy certificate management tool for development environments for Linux, Windows and macOS. While you are a solo developer in your workstation or distributed team of developers with many stations and servers.

Manage multiple CA without hassle. Automate different environments or servers without manual steps or magical flags that can make you lose too much time if some detail is missing.

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
