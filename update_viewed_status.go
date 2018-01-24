package srfax

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// ViewedStatusOpts contains optional arguments when updating fax status.
type ViewedStatusOpts struct {
	ResponseFormat string `json:"sResponseFormat,omitempty"`
}

// ViewedStatusResp is the response from a UpdateViewedStatus operation.
// Note, error message from Result, if any, will be stored in ResultError.
type ViewedStatusResp struct {
	Status      string `mapstructure:"Status"`
	Result      string `mapstructure:"Result"`
	ResultError string
}

// UpdateViewedStatus marks an inbound or outbound fax as read or unread.
//
// dir (direction) will be "IN" or "OUT" for inbound or outbound fax.
// view will be "Y" – mark fax as READ, or "N" – mark fax as UNREAD
// A note about ident:
//
// when passing sFaxFileName, the entire name (including pipe and ID) must be supplied.
// E.g., 20180101230101-8812-34_0|31524120
//
// If updating a fax based on sFaxDetailsID, pass in the number as a string.
// Formatting handled automatically.
func (c *Client) UpdateViewedStatus(ident, dir, view string, optArgs ...ViewedStatusOpts) (*ViewedStatusResp, error) {
	// TODO consider wrapping the string params "ident, dir, view" into a struct
	opts := ViewedStatusOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	msg := struct {
		Action string `json:"action"`
		Client
		FaxDetailsID int    `json:"sFaxDetailsID,omitempty"` // mutually exclusive
		FaxFileName  string `json:"sFaxFileName,omitempty"`  // mutually exclusive
		Direction    string `json:"sDirection"`
		Viewed       string `json:"sMarkasViewed"`
		ViewedStatusOpts
	}{
		Action:           actionUpdateViewedStatus,
		Client:           *c,
		Direction:        dir,
		Viewed:           view,
		ViewedStatusOpts: opts,
	}

	if strings.Contains(ident, "|") {
		msg.FaxFileName = ident
	} else {
		n, err := strconv.Atoi(ident)
		if err != nil {
			return nil, errors.Errorf("failed updating viewed status. sFaxDetailsID (id) string to int conversion, got [%[1]v] of type [%[1]T].", ident)
		}
		msg.FaxDetailsID = n
	}

	resp, err := sendPost(msg, c.url)
	if err != nil {
		return nil, errors.Wrap(err, "sendPost failed")
	}

	if st, err := checkStatus(resp); err != nil {
		return &ViewedStatusResp{Status: st, ResultError: fmt.Sprint(err)}, nil
	}

	var result ViewedStatusResp
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
