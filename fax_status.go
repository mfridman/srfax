package srfax

import (
	"bytes"
	"encoding/json"
	"io"

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

// GetFaxStatus retrieves the status of a single sent fax. Works only with outbound faxes.
// Accepts single id, i.e., FaxDetailID. Where FaxDetailsID returned from a QueueFax operation.
func (c *Client) GetFaxStatus(id int) (io.Reader, error) {

	if id <= 0 {
		return nil, errors.New("id (sFaxDetailsID) cannot be zero or negative number")
	}

	msg := struct {
		Action string `json:"action"`
		ID     int    `json:"sFaxDetailsID"`
		Client
	}{
		Action: actionGetFaxStatus,
		ID:     id,
		Client: *c,
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
