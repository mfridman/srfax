package srfax

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// RetrieveOptions specify optional arguments when retriving faxes.
type RetrieveOptions struct {
	// The account number of a sub account,
	// if you want to use a master account to download a sub account’s fax
	SubUserID string `json:"sSubUserID,omitempty"`

	// "PDF" or "TIFF", defaults to account settings if not supplied
	FaxFormat string `json:"sFaxFormat,omitempty"`

	// "Y" mark fax as viewed once method completes successfully.
	// "N" leave viewed status as is (default)
	MarkAsViewed string `json:"sMarkasViewed,omitempty"`
}

// RetrieveResp is the response from retrieving a fax.
type RetrieveResp struct {
	Status string

	// If successful the Result field will contain Base64 encoded fax file contents
	Result string
}

type mappedRetrieveResp struct {
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

// retrieveOperation defines the POST variables for a RetrieveFax request
type retrieveOperation struct {
	Action string `json:"action"`
	Client

	// Either the FaxFileName or the FaxDetailsID must be supplied
	FaxDetailsID int    `json:"sFaxDetailsID,omitempty"`
	FaxFileName  string `json:"sFaxFileName,omitempty"`

	// "IN" or "OUT" for inbound or outbound fax
	Direction string `json:"sDirection"`

	RetrieveOptions
}

func newRetrieveOperation(c *Client, ident, direction string, o *RetrieveOptions) (*retrieveOperation, error) {
	op := &retrieveOperation{Action: actionRetrieveFax, Client: *c, Direction: direction, RetrieveOptions: *o}
	if strings.Contains(ident, "|") {
		op.FaxFileName = ident
	} else {
		n, err := strconv.Atoi(ident)
		if err != nil {
			return nil, errors.New("failed ident string to int conversion")
		}
		op.FaxDetailsID = n
	}
	return op, nil
}

// RetrieveFax returns a sent or received fax file in PDF or TIFF format.
//
// ident will be either an sFaxDetailsID or sFaxFileName, returned from GetFaxInbox or GetFaxOutbox operation.
// Note, only one of sFaxDetailsID or sFaxFileName must be supplied.
//
// If operation succeeds the Result value contains a base64-encoded string.
// The file format will be "PDF" or "TIF" – defaults to account settings if FaxFormat not supplied in optional args.
func (c *Client) RetrieveFax(ident, direction string, options ...RetrieveOptions) (*RetrieveResp, error) {
	opts := RetrieveOptions{}
	if len(options) >= 1 {
		opts = options[0]
	}
	if !(direction == inbound || direction == outbound) {
		return nil, errors.Errorf("Direction must be one of: %s or %s", inbound, outbound)
	}

	opr, err := newRetrieveOperation(c, ident, direction, &opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build a newRetrieveOperation")
	}

	operation, err := constructFromStruct(opr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct a reader for newRetrieveOperation")
	}

	result := mappedRetrieveResp{}
	if err := run(operation, &result); err != nil {
		return nil, err
	}

	out := RetrieveResp(result)
	return &out, nil
}
