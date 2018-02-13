package srfax

import (
	"github.com/pkg/errors"
)

// GetFaxStatusResp represents the status of a single sent fax.
type GetFaxStatusResp struct {
	Status string `mapstructure:"Status"`
	Result *struct {
		FileName    string `mapstructure:"FileName"`
		SentStatus  string `mapstructure:"SentStatus"`
		DateQueued  string `mapstructure:"DateQueued"`
		DateSent    string `mapstructure:"DateSent"`
		ToFaxNumber string `mapstructure:"ToFaxNumber"`
		RemoteID    string `mapstructure:"RemoteID"` // API docs incorrect, this is a string, not an "integer"
		ErrorCode   string `mapstructure:"ErrorCode"`
		AccountCode string `mapstructure:"AccountCode"`
		Pages       int    `mapstructure:"Pages"`
		EpochTime   string `mapstructure:"EpochTime"` // FFS, API docs say this is a string, comes across as a number
		Duration    int    `mapstructure:"Duration"`
		Size        int    `mapstructure:"Size"`
	} `mapstructure:"Result"`
}

// getFaxStatusReq defines the POST variables for a GetFaxStatus request
type getFaxStatusReq struct {
	Action string `json:"action"`
	Client
	ID int `json:"sFaxDetailsID"`
}

// GetFaxStatus retrieves the status of a single sent fax. Works only with outbound faxes.
// Accepts a single id, i.e., FaxDetailsID, which is the return value from a QueueFax operation.
func (c *Client) GetFaxStatus(id int) (*GetFaxStatusResp, error) {

	if id <= 0 {
		return nil, errors.New("id cannot be zero or negative number")
	}

	req := getFaxStatusReq{
		Action: actionGetFaxStatus,
		Client: *c,
		ID:     id,
	}

	msg, err := sendPost(req)
	if err != nil {
		return nil, errors.Wrap(err, "GetFaxStatusResp SendPost error")
	}

	var resp GetFaxStatusResp
	if err := decodeResp(msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
