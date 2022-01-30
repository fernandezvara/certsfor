
# Quick Example

We have 3 files prepared to create CA, server and client certificates.

**ca.yaml**
```yaml
dn:
  cn: ca.certsfor.dev
  c: ""
  l: ""
  o: ""
  ou: ""
  p: ""
  pc: ""
  st: ""
san: []
key: rsa:4096
exp: 730
client: false
```

**server.yaml**
```yaml
dn:
  cn: server
  c: ""
  l: ""
  o: ""
  ou: ""
  p: ""
  pc: ""
  st: ""
san: 
  - "server.certsfor.dev"
  - "192.168.192.168"
key: rsa:4096
exp: 365
client: false
```

**client.yaml**
```yaml
dn:
  cn: client
  c: ""
  l: ""
  o: ""
  ou: ""
  p: ""
  pc: ""
  st: ""
san: []
key: rsa:4096
exp: 365
client: true
```

# Steps to create all the certificates

```bash
>
> # If you need an empty template
> cfd create template
> # this will create a YAML template for fill the information needed in
>
> 
> # in this example there are already 3 templates filled (see above its contents):
> # ca.yaml
> # server.yaml
> # client.yaml

>cfd create ca -c ./ca.cer.pem -k ./ca.key.pem -f ./ca.yaml

CA Created. ID: '6b834a85-0ad2-4eeb-a148-e8d2eda4d8aa'

> # we just created a CA. Every CA is identified internally by an ID. 
> # So, we need to ensure we save the ID for future reference.
> # for convenience, we can set it in an environment variable
>
> # if you are on linux / macos
> export CFD_CA_ID="6b834a85-0ad2-4eeb-a148-e8d2eda4d8aa"
>
> # on windows
> set CFD_CA_ID=6b834a85-0ad2-4eeb-a148-e8d2eda4d8aa
>
> # Create server and client certificates:
> cfd create certificate -c ./server.cer.pem -k ./server.key.pem -b server.bundle.pem --pfx ./server.pfx -f ./server.yaml
> cfd create certificate -c ./client.cer.pem -k ./client.key.pem --pfx ./client.pfx -f ./client.yaml
>
> # Now all your certificates are ready to use.
>
> dir
 Volume in drive C has no label.
 Volume Serial Number is 52D6-BE7A

 Directory of C:\certificates

01/28/2022  08:06 AM    <DIR>          .
01/28/2022  08:06 AM    <DIR>          ..
01/28/2022  07:57 AM             1,830 ca.cer.pem
01/28/2022  07:57 AM             3,272 ca.key.pem
01/28/2022  07:54 AM               130 ca.yaml
01/28/2022  08:06 AM             1,781 client.cer.pem
01/28/2022  08:06 AM             3,272 client.key.pem
01/28/2022  08:06 AM             5,375 client.pfx
01/28/2022  07:54 AM               120 client.yaml
01/28/2022  07:55 AM             2,048 readme.yaml
01/28/2022  08:04 AM             3,583 server.bundle.pem
01/28/2022  08:04 AM             1,753 server.cer.pem
01/28/2022  08:04 AM             3,272 server.key.pem
01/28/2022  08:04 AM             5,359 server.pfx
01/28/2022  07:54 AM               121 server.yaml
01/28/2022  07:55 AM               125 template.yaml

> # You can list the certificates created
>cfd list cert

┌─────────────────┬────────────────────┬─────────────────────────────┐
│ Common Name     │ Distinguished Name │ Expires In                  │
├─────────────────┼────────────────────┼─────────────────────────────┤
│ ca.certsfor.dev │ CN=ca.certsfor.dev │ 729 days (28/01/2024 06:57) │
│ client          │ CN=client          │ 364 days (28/01/2023 07:06) │
│ server          │ CN=server          │ 364 days (28/01/2023 07:04) │
└─────────────────┴────────────────────┴─────────────────────────────┘

``` 

# Fully automated

You can automate all the workflow of certifications to be usable in any pipeline. This example automates the creation of a new CA, server and client certificates. This allows recreating an environment from scratch using fresh certificates every time.


```bash
>
> # if you are on linux / macos
> export CFD_CA_ID=`cfd create ca -c ./ca.cer.pem -k ./ca.key.pem -f ./ca.yaml -q`
>
> # on windows (in this example we use a temporal file)
> cfd create ca -c ./ca.cer.pem -k ./ca.key.pem -f ./ca.yaml -q>tmp.txt
> set /P CFD_CA_ID=<tmp.txt
> del tmp.txt

> # now create server and client certificates:
> cfd create certificate -c ./server.cer.pem -k ./server.key.pem -b server.bundle.pem -f ./server.yaml
> cfd create certificate -c ./client.cer.pem -k ./client.key.pem -f ./client.yaml
>
```