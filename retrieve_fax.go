package srfax

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// RetrieveOptions specify optional arguments when retriving faxes.
type RetrieveOptions struct {
	SubUserID    string `json:"sSubUserID,omitempty"`
	FaxFormat    string `json:"sFaxFormat,omitempty"`
	MarkAsViewed string `json:"sMarkasViewed,omitempty"`
}

// RetrieveResp is the response from retrieving a fax.
type RetrieveResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// DecodeResult decodes a base64-encoded Result string and returns the raw bytes.
func (r *RetrieveResp) DecodeResult() ([]byte, error) {
	if !(strings.ToLower(r.Status) == "success") {
		return nil, errors.Errorf("cannot call decode on a [%v] Status", r.Status)
	}
	b, err := base64.StdEncoding.DecodeString(r.Result)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode base64 Result")
	}
	return b, nil
}

// retrieveRequest defines the POST variables for a RetrieveFax request
type retrieveRequest struct {
	Action string `json:"action"`
	Client
	FaxDetailsID int    `json:"sFaxDetailsID,omitempty"` // Either the FaxFileName or the FaxDetailsID must be supplied
	FaxFileName  string `json:"sFaxFileName,omitempty"`  // Either the FaxFileName or the FaxDetailsID must be supplied
	Direction    string `json:"sDirection"`              // "IN" or "OUT" for inbound or outbound fax
	RetrieveOptions
}

// RetrieveFax returns a sent or received fax file in PDF or TIFF format.
//
// ident is the sFaxDetailsID & sFaxFileName returned from GetFaxInbox or GetFaxOutbox,
// only one of sFaxDetailsID or sFaxFileName must be supplied as a string.
//
// If operation succeeds the Result value contain a base64-encoded string.
// The file format will be "PDF" or "TIF" – defaults to account settings if FaxFormat not supplied in optional args.
func (c *Client) RetrieveFax(ident, dir string, options ...RetrieveOptions) (*RetrieveResp, error) {
	opts := RetrieveOptions{}
	if len(options) >= 1 {
		opts = options[0]
	}

	if !(dir == inbound || dir == outbound) {
		return nil, errors.Errorf("Direction must be either: %s or %s", inbound, outbound)
	}

	req := retrieveRequest{
		Action:          actionRetrieveFax,
		Client:          *c,
		Direction:       dir,
		RetrieveOptions: opts,
	}

	if strings.Contains(ident, "|") {
		req.FaxFileName = ident
	} else {
		n, err := strconv.Atoi(ident)
		if err != nil {
			return nil, errors.New("failed ident string to int conversion")
		}
		req.FaxDetailsID = n
	}

	var resp RetrieveResp
	if err := run(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
