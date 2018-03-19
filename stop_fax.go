package srfax

import (
	"github.com/pkg/errors"
)

// StopFaxResp is the response from a StopFax operation.
type StopFaxResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// stopFaxOperation defines the POST variables for a StopFax request
type stopFaxOperation struct {
	Action string `json:"action"`
	Client
	FaxDetailsID int `json:"sFaxDetailsID"`
}

func newStopFaxOperation(c *Client, id int) *stopFaxOperation {
	return &stopFaxOperation{Action: actionStopFax, Client: *c, FaxDetailsID: id}
}

// StopFax deletes a specified queued fax which has not yet been processed.
// FaxDetailsID returned from Queue_Fax
func (c *Client) StopFax(id int) (*StopFaxResp, error) {
	if id <= 0 {
		return nil, errors.New("id cannot be zero or negative number")
	}
	resp := StopFaxResp{}
	op := newStopFaxOperation(c, id)
	if err := run(op, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
