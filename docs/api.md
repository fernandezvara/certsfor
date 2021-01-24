# API

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

#### **Response**

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
    }
}
```

#### **Curl**

```bash
>>curl -X PUT -d '{ 
    "dn": {
        "cn": "mycert1",
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
    "key": "rsa:4096",
    "exp": 90,
    "client": false
}' http://localhost:8080/v1/ca/0d22ee90-aed7-455e-b983-1631c91295ab/certificates/mycert1
{"key":"BASE64","certificate":"BASE64","ca_certificate":"BASE64",
"request":{"dn":{"cn":"mycert1","c":"ES","l":"MyLocality","o":"MyOrganization","ou":"MyOU",
"p":"MyProvince","pc":"00000","st":"MyStreet"},"san":["service1.example.com","192.168.1.2"],
"key":"rsa:4096","exp":90,"client":false}}
```

#### **Go**

```go
package main

import (
	"fmt"

	"github.com/fernandezvara/certsfor/pkg/client"
)

func main() {

	cli, err := client.New("api.certsfor.dev:443", "", "", "")
	if err != nil {
		panic(err)
	}

	request := client.APICertificateRequest{
		DN: client.APIDN{
			CN: "mycert1",
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
		Key:            "rsa:4096",
		ExpirationDays: 90,
		Client:         false,
	}

	caCert, err := cli.CACreate(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(caCert)
}


```

<!-- tabs:end -->

``` 
		"GET": {
			"/status": {
				Handler: a.getStatus,
				Matcher: []string{""},
			},
			"/v1/ca/:caid/certificates/:cn": {
				Handler: a.getCertificate,
				Matcher: []string{"", "", "", "", "[a-zA-Z0-9.-_]+"},
			},
		},
		"POST": {
			"/v1/ca": {
				Handler: a.postCA,
				Matcher: []string{"", ""},
			},
		},
		"PUT": {
			"/v1/ca/:caid/certificates/:cn": {
				Handler: a.putCertificate,
				Matcher: []string{"", "", "", "", "[a-zA-Z0-9.-_]+"},
			},
		},

```
