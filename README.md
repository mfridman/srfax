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

Import the library.

To begin using the client initialize a `ClientCfg` and pass it to `NewClient`.

```go
cfg := srfax.ClientCfg{
    ID:  00001,
    PWD: "password",
    // URL: optional, defaults to https://www.srfax.com/SRF_SecWebSvc.php
}

client, err := srfax.NewClient(cfg)
if err != nil {
    // check errors
}
```

With a `*Client` one runs all SRFax operations.

#### Example:

```go
resp, err := client.GetFaxUsage() // resp is a pointer to a FaxUsageResp.
if err != nil {
    // check errors
}

// for testing purposes use convenience func PP to pretty print response to terminal. PP uses MarshalIndent.
PP(resp) 
```
Output:
```json
{
    "Status": "Success",
    "Result": [{
        "Period": "ALL",
        "ClientName": "mike",
        "BillingNumber": "m@frid.io",
        "UserID": 00001,
        "SubUserID": 0,
        "NumberOfFaxes": 140,
        "NumberOfPages": 40
    }],
}
```