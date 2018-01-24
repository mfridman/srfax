package srfax

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// StopFaxOpts contains optional arguments when stopping a fax.
type StopFaxOpts struct {
	ResponseFormat string `json:"sResponseFormat,omitempty"`
}

// StopFaxResp is the response from a StopFax operation.
// Note, error message from Result, if any, will be stored in ResultError.
type StopFaxResp struct {
	Status      string `mapstructure:"Status"`
	Result      string `mapstructure:"Result"`
	ResultError string
}

// StopFax deletes a specified queued fax which has not yet been processed.
//
// Must supply a valid FaxDetailsID, which is a return value when calling QueueFax.
func (c *Client) StopFax(id int, optArgs ...StopFaxOpts) (*StopFaxResp, error) {
	opts := StopFaxOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	if id <= 0 {
		return nil, errors.New("id (sFaxDetailsID) cannot be zero or negative number")
	}

	msg := struct {
		Action string `json:"action"`
		Client
		FaxDetailsID int `json:"sFaxDetailsID"`
		StopFaxOpts
	}{
		Action:       actionStopFax,
		Client:       *c,
		FaxDetailsID: id,
		StopFaxOpts:  opts,
	}

	resp, err := sendPost(msg, c.url)
	if err != nil {
		return nil, errors.Wrap(err, "sendPost failed")
	}

	if st, err := checkStatus(resp); err != nil {
		return &StopFaxResp{Status: st, ResultError: fmt.Sprint(err)}, nil
	}

	var result StopFaxResp
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
