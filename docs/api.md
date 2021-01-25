# API

`cdf` comes with a simple REST API to allow any others interact.

## Constants

### Certificate Types

`cdf` supports the creation of certificates using ESCDA and RSA algorithms.

RSA (Rivest Shamir Adleman) asymmetric encryption algorithm. It was invented by Ron Rivest, Adi Shamir and Leonard Adleman in 1977. In this method two titanic-sized random prime numbers are multiplied to create another gigantic number. Multiplying both numbers is a simple task, but determining the original prime numbers is 'virtually' an impossible task, at least for now.

ECDSA (elliptic curve digital signature algorithm), is the successor of the digital signature algorithm (DSA). ECDSA was born when two mathematicians named Neal Koblitz and Victor S. Miller proposed the use of elliptical curves in cryptography. ECDSA is an assymmetric encryption algorithm that0s contructed around elliptical curves and a function known as 'trapdoor function'. An elliptic curve represents the set of points that satify a mathematical equation (y² = x³ +ax +b).

The following table enumerates which combinations are supported:

| Short Code | Algorithm |
| ---------- | --------- |
| rsa:2048   | RSA       |
| rsa:3072   | RSA       |
| rsa:4096   | RSA       |
| ecdsa:224  | ECDSA     |
| ecdsa:256  | ECDSA     |
| ecdsa:384  | ECDSA     |
| ecdsa:521  | ECDSA     |

>[!INFO]
>RSA 1024 is not allowed since is not considered secure and can be cracked (CVE-2017-7526).

>[!TIP]
>On Certificate request you must set `key` to the shortcode you need for the certificate.

## Create CA

```
POST /v1/ca
```

<!-- tabs:start -->

#### **Request**

**Body**

```json
{
    "dn": {
        "cn": "myca",
        "c": "ES",
        "l": "MyLocality",
        "o": "MyOrganization",
        "ou": "MyOU",
        "p": "MyProvince",
        "pc": "00000",
        "st": "MyStreet"
    },
    "san": [
        "www.example.com",
        "192.168.1.1"
    ],
    "key": "rsa:4096",
    "exp": 90,
    "client": false
}
```

#### **Responses**

| Code | Description |
| ---- | ----------- |
| 201  | CA created successfully |
| 400  | Request does not meet the requirements |


**Body**

```json
{
    "key": "BASE64 string",
    "certificate": "BASE64 string",
    "request": {
        "dn": {
            "cn": "mycert",
            "c": "ES",
            "l": "MyLocality",
            "o": "MyOrganization",
            "ou": "MyOU",
            "p": "MyProvince",
            "pc": "00000",
            "st": "MyStreet"
        },
        "san": [
            "www.example.com",
            "192.168.1.1"
        ],
        "key": "rsa:4096",
        "exp": 90,
        "client": false
    },
    "ca_id":"a600097f-d860-4f53-9269-28f1b8bd15b8"
}
```

#### **Curl**

```bash
>>curl -X POST -d '{ 
    "dn": {
        "cn": "myca",
        "c": "ES",
        "l": "MyLocality",
        "o": "MyOrganization",
        "ou": "MyOU",
        "p": "MyProvince",
        "pc": "00000",
        "st": "MyStreet"
    },
    "san": [
        "ca.example.com",
        "192.168.1.1"
    ],
    "key": "rsa:4096",
    "exp": 90,
    "client": false
}' https://api.certsfor.dev:8443/v1/ca
{"key":"BASE64","certificate":"BASE64","ca_certificate":"BASE64",
"request":{"dn":{"cn":"myca","c":"ES","l":"MyLocality","o":"MyOrganization","ou":"MyOU",
"p":"MyProvince","pc":"00000","st":"MyStreet"},"san":["ca.example.com","192.168.1.1"],
"key":"rsa:4096","exp":90,"client":false},"ca_id":"a600097f-d860-4f53-9269-28f1b8bd15b8"}
```

#### **Go**

```go
package main

import (
	"fmt"

	"github.com/fernandezvara/certsfor/pkg/client"
)

func main() {

	cli, err := client.New("api.certsfor.dev:8443", "", "", "", true)
	if err != nil {
		panic(err)
	}

	request := client.APICertificateRequest{
		DN: client.APIDN{
			CN: "myca",
			C:  "ES",
			L:  "MyLocality",
			O:  "MyOrganization",
			OU: "MyOU",
			P:  "MyProvince",
			PC: "00000",
			ST: "MyStreet",
		},
		SAN: []string{
			"ca.example.com",
			"192.168.1.1",
		},
		Key:            "rsa:4096",
		ExpirationDays: 90,
		Client:         false,
	}

	caCert, err := cli.CACreate(request)
	if err != nil {
		panic(err)
	}

    fmt.Println(caCert.CAID)
	fmt.Println(string(caCert.Certificate))
}
```

<!-- tabs:end -->

## Create/Update Certificate

```
PUT /v1/ca/:caid:/certificates/:common-name:
```

>[!ATTENTION]
>Updating a certificate does not renew it, **it makes a new pair combination from scratch**. Be sure to don't overwrite certificates. Common Name is used as ID.
>
>Renewal is done retriving the certificate simplifying the workflow.
>
>**CA Certificate cannot be updated**, it will fail with Conflict HTTP error code. CA certificates are always stored as 'ca'.

<!-- tabs:start -->

#### **Request**

**Body**

```json
{
    "dn": {
        "cn": "service1",
        "c": "ES",
        "l": "MyLocality",
        "o": "MyOrganization",
        "ou": "MyOU",
        "p": "MyProvince",
        "pc": "00000",
        "st": "MyStreet"
    },
    "san": [
        "service1.example.com",
        "192.168.1.2"
    ],
    "key": "ecdsa:521",
    "exp": 30,
    "client": false
}
```

#### **Responses**

| Code | Description |
| ---- | ----------- |
| 200  | Certificate created / updated successfully |
| 400  | Request does not meet the requirements |
| 409  | Updating the certificate will overwrite the CA certificate, so it's not permitted |

**Body**

```json
{
    "key": "BASE64 string",
    "certificate": "BASE64 string",
    "request": {
        "dn": {
            "cn": "service1",
            "c": "ES",
            "l": "MyLocality",
            "o": "MyOrganization",
            "ou": "MyOU",
            "p": "MyProvince",
            "pc": "00000",
            "st": "MyStreet"
        },
        "san": [
            "service1.example.com",
            "192.168.1.2"
        ],
        "key": "ecdsa:521",
        "exp": 90,
        "client": false
    }
}
```

#### **Curl**

```bash
>>curl -X PUT -d '{ 
    "dn": {
        "cn": "service1",
        "c": "ES",
        "l": "MyLocality",
        "o": "MyOrganization",
        "ou": "MyOU",
        "p": "MyProvince",
        "pc": "00000",
        "st": "MyStreet"
    },
    "san": [
        "service1.example.com",
        "192.168.1.2"
    ],
    "key": "ecdsa:521"
    "exp": 90,
    "client": false
}' https://api.certsfor.dev:8443/v1/ca/a600097f-d860-4f53-9269-28f1b8bd15b8/certificates/service1
{"key":"BASE64","certificate":"BASE64","ca_certificate":"BASE64",
"request":{"dn":{"cn":"service1","c":"ES","l":"MyLocality","o":"MyOrganization","ou":"MyOU",
"p":"MyProvince","pc":"00000","st":"MyStreet"},"san":["service1.example.com","192.168.1.2"],
"key":"ecdsa:521","exp":90,"client":false}}
```

#### **Go**

```go
package main

import (
	"fmt"

	"github.com/fernandezvara/certsfor/pkg/client"
)

func main() {

	cli, err := client.New("api.certsfor.dev:8443", "", "", "", true)
	if err != nil {
		panic(err)
	}

	requestService1 := client.APICertificateRequest{
		DN: client.APIDN{
			CN: "service1",
			C:  "ES",
			L:  "MyLocality",
			O:  "MyOrganization",
			OU: "MyOU",
			P:  "MyProvince",
			PC: "00000",
			ST: "MyStreet",
		},
		SAN: []string{
			"service1.example.com",
			"192.168.1.2",
		},
		Key:            "ecdsa:521",
		ExpirationDays: 90,
		Client:         false,
	}

	cert, err := cli.CertificateCreate("a600097f-d860-4f53-9269-28f1b8bd15b8", "service1", requestService1)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(cert.Certificate))
    
}
```

<!-- tabs:end -->

## Get Certificate

```
GET /v1/ca/:caid:/certificates/:common-name:?renew=XX
```

>[!NOTE]
>Renewal is done automatically on retrieving to simplify the workflow. If not specified it will be renewed if the time to expire is (20% or less than certificate lifetime).

>[!NOTE]
>To retrieve the CA certificate use *ca* as common name when calling the API

<!-- tabs:start -->

#### **Request**

**Parameters**

| Parameter | Description |
| --------- | ----------- |
| renew  | *(optional)* Percent of time used to calculate if the certificate needs to be renewed. If the threshold is met, the certificate will be auto-renewed and returned on the response. **(default: 20)** |


#### **Responses**

| Code | Description |
| ---- | ----------- |
| 200  | Certificate retrieved successfully |
| 404  | Certificate not found |

**Body**

```json
{
    "key": "BASE64 string",
    "certificate": "BASE64 string",
    "request": {
        "dn": {
            "cn": "service1",
            "c": "ES",
            "l": "MyLocality",
            "o": "MyOrganization",
            "ou": "MyOU",
            "p": "MyProvince",
            "pc": "00000",
            "st": "MyStreet"
        },
        "san": [
            "service1.example.com",
            "192.168.1.2"
        ],
        "key": "ecdsa:521",
        "exp": 90,
        "client": false
    }
}
```

#### **Curl**

##### Basic Usage

```bash
>>curl https://api.certsfor.dev:8443/v1/ca/a600097f-d860-4f53-9269-28f1b8bd15b8/certificates/service1
{"key":"BASE64","certificate":"BASE64","ca_certificate":"BASE64",
"request":{"dn":{"cn":"service1","c":"ES","l":"MyLocality","o":"MyOrganization","ou":"MyOU",
"p":"MyProvince","pc":"00000","st":"MyStreet"},"san":["service1.example.com","192.168.1.2"],
"key":"ecdsa:521","exp":90,"client":false}}
```

##### Save files

```bash
# certificate
curl -X GET https://api.certsfor.dev:8443/v1/ca/a600097f-d860-4f53-9269-28f1b8bd15b8/certificates/service1 
| jq '.certificate' -r | base64 -d > cert.crt


# key
curl -X GET https://api.certsfor.dev:8443/v1/ca/a600097f-d860-4f53-9269-28f1b8bd15b8/certificates/service1 
| jq '.key' -r | base64 -d > key.crt
```

#### **Go**

```go
package main

import (
	"fmt"

	"github.com/fernandezvara/certsfor/pkg/client"
)

func main() {

	cli, err := client.New("api.certsfor.dev:8443", "", "", "", true)
	if err != nil {
		panic(err)
	}

	var cert client.Certificate

	cert, err = cli.CertificateGet("a600097f-d860-4f53-9269-28f1b8bd15b8", "service1", 20)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(cert.CACertificate))
	fmt.Println(string(cert.Certificate))
	fmt.Println(string(cert.Key))

}
```

<!-- tabs:end -->

## Status

```
GET /status
```

<!-- tabs:start -->

#### **Request**

**Parameters**

There are not parameters for this endpoint.

#### **Responses**

| Code | Description |
| ---- | ----------- |
| 200  | API is ok |
| 500  | There is an error on the API |

**Body**

```json
{
    "version":"0.1"
}
```

#### **Curl**

##### Basic Usage

```bash
>>curl https://api.certsfor.dev:8443/status
{"version":"0.1"}
```

#### **Go**

```go
package main

import (
	"fmt"

	"github.com/fernandezvara/certsfor/pkg/client"
)

func main() {

	cli, err := client.New("api.certsfor.dev:8443", "", "", "", true)
	if err != nil {
		panic(err)
	}

    status, err := cli.Status()
	if err != nil {
		panic(err)
	}

	fmt.Println(status.Version)

}
```

<!-- tabs:end -->

