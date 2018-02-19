/*
Package srfax is a Go client library for interacting with the SRFax API service,
supporting all POST operations. The official SRFax API docs can be found here:

https://www.srfax.com/developers/our-fax-api/

There are a few notable differences between the official API docs and this Go client:

1.	As a result of variable types being returned in the Result field, see note about error handling: https://github.com/mfridman/srfax#a-note-about-error-handling

2.	Response format supports JSON-only. Despite the API being able to return both XML
and JSON, it is a design decision to support JSON only. As a result,
sResponseFormat is not available as an optional parameter.
*/
package srfax
