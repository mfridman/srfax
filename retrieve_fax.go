package srfax

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// RetrieveFaxOpts contains optional arguments when retriving faxes.
type RetrieveFaxOpts struct {
	SubUserID    string `json:"sSubUserID,omitempty"`
	FaxFormat    string `json:"sFaxFormat,omitempty"`
	MarkAsViewed string `json:"sMarkasViewed,omitempty"`
}

// RetrieveFaxResp is the response from retrieving a fax.
type RetrieveFaxResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// DecodeResult decodes a base64-encoded Result string and returns the raw bytes.
func (r *RetrieveFaxResp) DecodeResult() ([]byte, error) {
	if !(strings.ToLower(r.Status) == "success") {
		return nil, errors.Errorf("cannot call decode on a [%v] Status", r.Status)
	}
	b, err := base64.StdEncoding.DecodeString(r.Result)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode base64 Result")
	}
	return b, nil
}

// retrieveFaxReq defines the POST variables for a RetrieveFax request
type retrieveFaxReq struct {
	Action string `json:"action"`
	Client
	FaxDetailsID int    `json:"sFaxDetailsID,omitempty"` // Either the FaxFileName or the FaxDetailsID must be supplied
	FaxFileName  string `json:"sFaxFileName,omitempty"`  // Either the FaxFileName or the FaxDetailsID must be supplied
	Direction    string `json:"sDirection"`              // "IN" or "OUT" for inbound or outbound fax
	RetrieveFaxOpts
}

// RetrieveFax returns a sent or received fax file in PDF or TIFF format.
//
// ident is the sFaxDetailsID & sFaxFileName returned from GetFaxInbox or GetFaxOutbox,
// only one of sFaxDetailsID or sFaxFileName must be supplied as a string.
//
// If operation succeeds the Result value contain a base64-encoded string.
// The file format will be "PDF" or "TIF" – defaults to account settings if FaxFormat not supplied in optional args.
func (c *Client) RetrieveFax(ident, dir string, options ...RetrieveFaxOpts) (*RetrieveFaxResp, error) {
	opts := RetrieveFaxOpts{}
	if len(options) >= 1 {
		opts = options[0]
	}

	if !(dir == inbound || dir == outbound) {
		return nil, errors.Errorf("Direction must be either: %s or %s", inbound, outbound)
	}

	req := retrieveFaxReq{
		Action:          actionRetrieveFax,
		Client:          *c,
		Direction:       dir,
		RetrieveFaxOpts: opts,
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

	msg, err := sendPost(req)
	if err != nil {
		return nil, errors.Wrap(err, "RetrieveFaxResp SendPost error")
	}

	var resp RetrieveFaxResp
	if err := decodeResp(msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
