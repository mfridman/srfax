package srfax

import (
	"github.com/pkg/errors"
)

// FaxStatusResp represents the status of a single sent fax.
type FaxStatusResp struct {
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

// FaxStatusReq defines the POST variables for a GetFaxStatus request
type FaxStatusReq struct {
	Action string `json:"action"`
	Client
	ID int `json:"sFaxDetailsID"`
}

// GetFaxStatus retrieves the status of a single sent fax. Works only with outbound faxes.
// Accepts a single id, i.e., FaxDetailsID, which is the return value from a QueueFax operation.
func (c *Client) GetFaxStatus(id int) (*FaxStatusReq, error) {

	if id <= 0 {
		return nil, errors.New("id cannot be zero or negative number")
	}

	req := FaxStatusReq{
		Action: actionGetFaxStatus,
		Client: *c,
		ID:     id,
	}

	return &req, nil
}
