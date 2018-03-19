package srfax

import (
	"strings"

	"github.com/pkg/errors"
)

// MulFaxStatus represents the status of multiple sent faxes.
type MulFaxStatus struct {
	Status string `mapstructure:"Status"`
	Result []struct {
		Pages string `mapstructure:"Pages"` // API returns empty string when there is an error
		// TODO this comes across as an error regardless
		// decoding resp body: mapstructure: cannot unmarshal string into Go struct field .EpochTime of type int
		// decoding resp body: mapstructure: cannot unmarshal number into Go struct field .EpochTime of type string
		EpochTime   string `mapstructure:"EpochTime"` // API returns empty string when there is an error
		Duration    string `mapstructure:"Duration"`
		Size        string `mapstructure:"Size"`
		FileName    string `mapstructure:"FileName"`
		SentStatus  string `mapstructure:"SentStatus"`
		DateQueued  string `mapstructure:"DateQueued"`
		DateSent    string `mapstructure:"DateSent"`
		ToFaxNumber string `mapstructure:"ToFaxNumber"`
		RemoteID    string `mapstructure:"RemoteID"`
		ErrorCode   string `mapstructure:"ErrorCode"`
		AccountCode string `mapstructure:"AccountCode"`
	} `mapstructure:"Result"`
}

// mulFaxStatusOperation defines the POST variables for a GetMulFaxStatus request
type mulFaxStatusOperation struct {
	Action string `json:"action"`
	Client
	IDs string `json:"sFaxDetailsID"`
}

func newMulFaxUsageOperation(c *Client, ids []string) *mulFaxStatusOperation {
	s := strings.Join(ids, "|")
	return &mulFaxStatusOperation{Action: actionGetMulFaxStatus, Client: *c, IDs: s}
}

// GetMulFaxStatus retrieves status of multiple sent faxes. Works only with outbound faxes.
// Accepts a multiple id, i.e., FaxDetailsID, which is the result value from QueueFax or ForwardFax.
// Note, this method will take care of formatting ids accordingly with pipe(s).
func (c *Client) GetMulFaxStatus(ids []string) (*MulFaxStatus, error) {
	if len(ids) == 0 {
		return nil, errors.New("must supply one or more identifiers")
	}
	resp := MulFaxStatus{}
	opr := newMulFaxUsageOperation(c, ids)
	if err := run(opr, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
