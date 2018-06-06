package srfax

import (
	"strings"

	"github.com/pkg/errors"
)

// MulFaxStatus represents the status of multiple sent faxes.
type MulFaxStatus struct {
	Status string
	Result []struct {
		Pages       string
		EpochTime   string
		Duration    string
		Size        string
		FileName    string
		SentStatus  string
		DateQueued  string
		DateSent    string
		ToFaxNumber string
		RemoteID    string
		ErrorCode   string
		AccountCode string
	}
}

type mappedMulFaxStatus struct {
	Status string `mapstructure:"Status"`
	Result []struct {
		Pages       string `mapstructure:"Pages"`
		EpochTime   string `mapstructure:"EpochTime"`
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
	return &mulFaxStatusOperation{Action: actionGetMulFaxStatus, Client: *c, IDs: strings.Join(ids, "|")}
}

// GetMulFaxStatus retrieves status of multiple sent faxes. Works only with outbound faxes.
// Accepts a multiple id, i.e., FaxDetailsID, which is the result value from QueueFax or ForwardFax.
// Note, this method will take care of formatting ids accordingly with pipe(s).
func (c *Client) GetMulFaxStatus(ids []string) (*MulFaxStatus, error) {
	if len(ids) == 0 {
		return nil, errors.New("must supply one or more identifiers")
	}

	operation, err := constructReader(newMulFaxUsageOperation(c, ids))
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct a reader from newMulFaxUsageOperation")
	}

	result := mappedMulFaxStatus{}
	if err := run(operation, &result, c.url); err != nil {
		return nil, err
	}

	out := MulFaxStatus(result)
	return &out, nil
}
