package srfax

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

// StopFaxResp is the response from a StopFax operation.
type StopFaxResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// StopFax deletes a specified queued fax which has not yet been processed.
//
// Must supply a valid FaxDetailsID, which is a return value when calling QueueFax.
func (c *Client) StopFax(id int) (io.Reader, error) {

	if id <= 0 {
		return nil, errors.New("id (sFaxDetailsID) cannot be zero or negative number")
	}

	msg := struct {
		Action string `json:"action"`
		Client
		FaxDetailsID int `json:"sFaxDetailsID"`
	}{
		Action:       actionStopFax,
		Client:       *c,
		FaxDetailsID: id,
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
