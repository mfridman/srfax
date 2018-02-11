[![](https://godoc.org/github.com/mfridman/srfax?status.svg)](http://godoc.org/github.com/mfridman/srfax)
[![Go Report Card](https://goreportcard.com/badge/github.com/mfridman/srfax)](https://goreportcard.com/report/github.com/mfridman/srfax)
# srfax

`srfax` is a Go client library for interacting with the SRFax API service. The client supports all SRFax operations and provides convenience functions for ease-of-use. Examples for all methods can be found in the [wiki](https://github.com/mfridman/srfax/wiki).

The official SRFax API documentation can be found [here](https://www.srfax.com/api-page/getting-started/)

**The current client is under development and the API may change. There is no guarantee of backwards compatibility at this time.**

---

A note about error handling:

This client will return an error if the Status value is "Failed". Do not attempt to access error messages via the Result field, use standard Go error handling instead.

It is the caller's responsibility to check for errors prior to accessing a response struct. If error is not `nil` assume something has gone wrong.

Caller can check for `*ResultError` error type and retrieve the original Status and Result message. Example:

```go
if err != nil {
	switch e := err.(type) {
	case *srfax.ResultError:
		fmt.Println(e.Status)
		fmt.Println(e.Raw)
	}
}
// Failed
// Invalid Access Code / Password

fmt.Println(err)
// Failed: Invalid Access Code / Password
```

If SRFax had publicly available errors then could compose error types, but for now `ResultError` will do.

## Installation

    go get -u github.com/mfridman/srfax

## Usage

Import the library. `"github.com/mfridman/srfax"`

Start by initializing a `ClientCfg` and pass it to `NewClient`.

```go
cfg := srfax.ClientCfg{
    ID:  00001,
    Pwd: "password",
}

client, err := srfax.NewClient(cfg)
if err != nil {
    // check errors
}
```

With a `*Client` one runs an SRFax operation to obtain a request struct (or map), sends the request via POST and decodes the response into the corresponding struct.

User's of this library have the flexibility to implement their own POST and pass a `map[string]interface{}` directly to `DecodeResp` along with the corresponding response type.

#### Example:

```go
req, err := client.GetFaxUsage()
if err != nil {
	// handle error
}

msg, err := srfax.SendPost(req)
if err != nil {
	// handle error
}

var resp srfax.GetFaxUsageResp
if err := srfax.DecodeResp(msg, &resp); err != nil {
	// handle error
}

// use convenience func PP to pretty print response to terminal.
srfax.PP(resp)
```

Output:

```json
{
    "Status": "Success",
    "Result": [{
        "Period": "ALL",
        "ClientName": "Michael Fridman",
        "BillingNumber": "mf192@icloud.com",
        "UserID": 00001,
        "SubUserID": 0,
        "NumberOfFaxes": 140,
        "NumberOfPages": 240
    }]
}
```