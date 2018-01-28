package srfax

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// FaxInboxOpts contains optional arguments when retrieving inbox items.
type FaxInboxOpts struct {
	Period          string `json:"sPeriod,omitempty"`
	StartDate       string `json:"sStartDate,omitempty"`
	EndDate         string `json:"sEndDate,omitempty"`
	ViewedStatus    string `json:"sViewedStatus,omitempty"`
	IncludeSubUsers string `json:"sIncludeSubUsers,omitempty"`
}

// FaxInboxResp represents fax inbox information.
type FaxInboxResp struct {
	Status string `mapstructure:"Status"`
	Result []struct {
		FileName      string `mapstructure:"FileName"`
		ReceiveStatus string `mapstructure:"ReceiveStatus"`
		Date          string `mapstructure:"Date"`
		CallerID      string `mapstructure:"CallerID"`
		RemoteID      string `mapstructure:"RemoteID"`
		ViewedStatus  string `mapstructure:"ViewedStatus"`
		UserID        string `mapstructure:"User_ID" json:",omitempty"`        // only if sIncludeSubUsers is set to “Y”
		UserFaxNumber string `mapstructure:"User_FaxNumber" json:",omitempty"` // only if sIncludeSubUsers is set to “Y”
		EpochTime     int    `mapstructure:"EpochTime"`
		Pages         int    `mapstructure:"Pages"`
		Size          int    `mapstructure:"Size"`
	} `mapstructure:"Result"`
}

// GetFaxInbox retrieves a list of faxes received for a specified period of time.
func (c *Client) GetFaxInbox(optArgs ...FaxInboxOpts) (*FaxInboxResp, error) {
	opts := FaxInboxOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	msg := struct {
		Action string `json:"action"`
		Client
		FaxInboxOpts
	}{
		Action:       actionGetFaxInbox,
		Client:       *c,
		FaxInboxOpts: opts,
	}

	resp, err := sendPost(msg, c.url)
	if err != nil {
		return nil, errors.Wrap(err, "sendPost failed")
	}

	if st, err := checkStatus(resp); err != nil {
		return nil, &ResultError{Status: st, Raw: fmt.Sprint(err)}
	}

	var result FaxInboxResp
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
