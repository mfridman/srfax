/*
Package srfax is a Go client for the SRFax API service, supporting all
POST operations. The official SRFax API docs can be found here:

https://www.srfax.com/developers/our-fax-api/

There are a few notable differences between the official docs and this
client as a result of variable types coming back from the server. See README file.

The response format is JSON-only. Despite the API returning both XML
and JSON, it's a design decision (for now) to support JSON only respones.
*/
package srfax
