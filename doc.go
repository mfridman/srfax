/*
Package srfax is a Go client for the SRFax API service, supporting all
POST operations. The official SRFax API docs can be found here:

https://www.srfax.com/developers/our-fax-api/

There are a few notable differences between the official docs and this
client as a result of variable types being returned.

See note about error handling: https://github.com/mfridman/srfax#srfax

The response format is JSON-only. Despite the API returning both XML
and JSON, it's a design decision (for now) to support JSON only.
*/
package srfax
