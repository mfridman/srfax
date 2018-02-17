package srfax

import (
	"github.com/pkg/errors"
)

// FaxStatus represents the status of a single sent fax.
type FaxStatus struct {
	Status string `mapstructure:"Status"`
	Result *struct {
		FileName    string `mapstructure:"FileName"`
		SentStatus  string `mapstructure:"SentStatus"`
		DateQueued  string `mapstructure:"DateQueued"`
		DateSent    string `mapstructure:"DateSent"`
		ToFaxNumber string `mapstructure:"ToFaxNumber"`
		RemoteID    string `mapstructure:"RemoteID"`
		ErrorCode   string `mapstructure:"ErrorCode"`
		AccountCode string `mapstructure:"AccountCode"`
		Pages       int    `mapstructure:"Pages"`
		EpochTime   string `mapstructure:"EpochTime"`
		Duration    int    `mapstructure:"Duration"`
		Size        int    `mapstructure:"Size"`
	} `mapstructure:"Result"`
}

// faxStatusRequest defines the POST variables for a GetFaxStatus request
type faxStatusRequest struct {
	Action string `json:"action"`
	Client
	ID int `json:"sFaxDetailsID"`
}

// GetFaxStatus retrieves the status of a single sent fax. Works only with outbound faxes.
// Accepts a single id, i.e., FaxDetailsID, which is the result value from QueueFax or ForwardFax.
func (c *Client) GetFaxStatus(id int) (*FaxStatus, error) {

	if id <= 0 {
		return nil, errors.New("id cannot be zero or negative number")
	}

	req := faxStatusRequest{
		Action: actionGetFaxStatus,
		Client: *c,
		ID:     id,
	}

	var resp FaxStatus
	if err := run(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
