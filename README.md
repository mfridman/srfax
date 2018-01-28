[![](https://godoc.org/github.com/mfridman/srfax?status.svg)](http://godoc.org/github.com/mfridman/srfax)
[![Go Report Card](https://goreportcard.com/badge/github.com/mfridman/srfax)](https://goreportcard.com/report/github.com/mfridman/srfax)
# srfax

`srfax` is a Go client library for interacting with the SRFax API service. The client supports all official operations and provides convenience functions for ease-of-use.

The official SRFax API documentation can be found [here](https://www.srfax.com/api-page/getting-started/)

**Important**, this library will handle Result errors in a Go like manner. Instead of mixing error messages with real results like so:

```json
{
    "Status": "",
    "Result": "",
}
```

This client will always return an operation-specific response (suffixed with Resp) composed of the following:

```js
{
    "Status": "",       /* Success or Failed */
    "Result": "",       /* result will be specific to the operation performed */
    "ResultError": ""   /* error will be empty unless there was an error */
}
```

It is the caller's responsibility to check the Error value prior to using Result. If Error is not an empty string assume something has gone wrong and Result should not be used.

If SRFax had publicly available errors then could compose error types, but for now it'll be a string.

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

// for testing purposes use convenience func PP to pretty print the response to terminal
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
    "ResultError": ""
}
```