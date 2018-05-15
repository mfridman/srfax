[![](https://godoc.org/github.com/mfridman/srfax?status.svg)](http://godoc.org/github.com/mfridman/srfax)
[![Go Report Card](https://goreportcard.com/badge/github.com/mfridman/srfax)](https://goreportcard.com/report/github.com/mfridman/srfax) [![Build Status](https://travis-ci.com/mfridman/srfax.svg?branch=master)](https://travis-ci.com/mfridman/srfax)
# srfax

`srfax` is a Go client library for interacting with the SRFax API service. The client supports all SRFax operations and provides convenience functions for ease-of-use. Examples for all methods can be found in the [wiki](https://github.com/mfridman/srfax/wiki).

The official SRFax API documentation can be found [here](https://www.srfax.com/api-page/getting-started/)

---

### A note about error handling

This client will return an error if the Status value is "Failed". Do not attempt to access error messages via the Result field, use standard Go error handling instead.

It is the caller's responsibility to check for errors prior to accessing a response struct. If error is not `nil` assume something has gone wrong.

Caller can check for the `*ResultError` error type and retrieve the original Status and Result message. Example:

```go
if err != nil {
	switch e := err.(type) {
	case *srfax.ResultError:
		fmt.Println(e.Status) // Failed
		fmt.Println(e.Raw)    // Invalid Access Code / Password
	}
}

fmt.Println(err) // Failed: Invalid Access Code / Password
```

If SRFax had publicly available errors then could compose more specific error types, but for now `ResultError` will do.

## Installation

    go get -u github.com/mfridman/srfax

## Usage

Import the library. `"github.com/mfridman/srfax"`

Start by initializing a `ClientCfg` and pass it to `NewClient`, where ID (account number) and Pwd (password) are unique to your SRFax Account.

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

There is a convenience method to check authentication:

```go
ok, err := client.CheckAuth()
if err != nil {
    fmt.Printf("ok: %t, err: %s\n", ok, err) // ok: false, err: Failed: Invalid Access Code / Password
    os.Exit(1)
}
fmt.Println(ok) // true
```

With a `*Client` one runs all the supported SRFax operations. For each method this client will construct a request operation, send it via POST and decode the response into a corresponding type.

Some methods require many arguments so they'll typically be wrapped in a struct with a Cfg suffix (E.g., ForwardFax, QueueFax and UpdateViewedStatus). 

Some methods accept optional arguments and will be wrapped in a struct with an Options suffix.

Examples for all methods will be found in the [wiki](https://github.com/mfridman/srfax/wiki). The following is a quick example to get you started:

#### Example:

```go
resp, err := client.GetFaxUsage() // resp is of type *FaxUsage
if err != nil {
    // handle error
}

// use convenience func PP to pretty print response to terminal.
srfax.PP(resp)
```

Output:

```json
{
  "Status": "Success",
  "Result": [
    {
      "Period": "ALL",
      "ClientName": "Michael Fridman",
      "BillingNumber": "mf192@icloud.com",
      "UserID": 92126,
      "SubUserID": 0,
      "NumberOfFaxes": 11,
      "NumberOfPages": 11
    },
    {
      "Period": "ALL",
      "ClientName": "Michael Fridman",
      "BillingNumber": "6476892276",
      "UserID": 92126,
      "SubUserID": 0,
      "NumberOfFaxes": 11,
      "NumberOfPages": 11
    }
  ]
}
```