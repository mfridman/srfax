package srfax

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// RetrieveFaxOpts contains optional arguments when retriving faxes.
type RetrieveFaxOpts struct {
	SubUserID      string `json:"sSubUserID,omitempty"`
	ResponseFormat string `json:"sResponseFormat,omitempty"`
	FaxFormat      string `json:"sFaxFormat,omitempty"`
	MarkAsViewed   string `json:"sMarkasViewed,omitempty"`
}

// RetrieveFaxResp is the response from retrieving a fax.
// Note, error message from Result, if any, will be stored in ResultError.
type RetrieveFaxResp struct {
	Status      string `mapstructure:"Status"`
	Result      string `mapstructure:"Result"`
	ResultError string
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

// RetrieveFax returns a sent or received fax file in PDF or TIFF format.
//
// ident is the sFaxDetailsID or sFaxFileName returned from GetFaxInbox or GetFaxOutbox
// dir is the direction; "IN" or "OUT" for inbound or outbound fax
//
// If operation succeeds the Result value contain a base64-encoded string.
// The file format will be "PDF" or "TIF" – defaults to account settings if FaxFormat not supplied in optional args.
func (c *Client) RetrieveFax(ident, dir string, optArgs ...RetrieveFaxOpts) (*RetrieveFaxResp, error) {
	if !(dir == "IN" || dir == "OUT") {
		return nil, errors.New(`dir (direction) must be one of either "IN" or "OUT"`)
	}

	opts := RetrieveFaxOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	msg := struct {
		Action string `json:"action"`
		Client
		FaxDetailsID int    `json:"sFaxDetailsID,omitempty"`
		FaxFileName  string `json:"sFaxFileName,omitempty"`
		Direction    string `json:"sDirection"`
		RetrieveFaxOpts
	}{
		Action:          actionRetrieveFax,
		Client:          *c,
		Direction:       dir,
		RetrieveFaxOpts: opts,
	}

	if strings.Contains(ident, "|") {
		msg.FaxFileName = ident
	} else {
		n, err := strconv.Atoi(ident)
		if err != nil {
			return nil, errors.New("failed id string to int conversion")
		}
		msg.FaxDetailsID = n
	}

	resp, err := sendPost(msg, c.url)
	if err != nil {
		return nil, errors.Wrap(err, "sendPost failed")
	}

	if st, err := checkStatus(resp); err != nil {
		return &RetrieveFaxResp{Status: st, ResultError: fmt.Sprint(err)}, nil
	}

	var result RetrieveFaxResp
	var md mapstructure.Metadata
	cfg := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Metadata:         &md,
		Result:           &result,
	}

	if err := decodeResp(resp, cfg); err != nil {
		return nil, err
	}

	return &result, nil
}
