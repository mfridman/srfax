package srfax

import (
	"github.com/pkg/errors"
)

// StopFaxResp is the response from a StopFax operation.
type StopFaxResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// stopFaxReq defines the POST variables for a StopFax request
type stopFaxReq struct {
	Action string `json:"action"`
	Client
	FaxDetailsID int `json:"sFaxDetailsID"`
}

// StopFax deletes a specified queued fax which has not yet been processed.
// FaxDetailsID returned from Queue_Fax
func (c *Client) StopFax(id int) (*StopFaxResp, error) {

	if id <= 0 {
		return nil, errors.New("FaxDetailsID cannot be zero or negative number")
	}

	req := stopFaxReq{
		Action:       actionStopFax,
		Client:       *c,
		FaxDetailsID: id,
	}

	var resp StopFaxResp
	if err := run(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
