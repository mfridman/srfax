package srfax

import (
	"github.com/pkg/errors"
)

// GetFaxInboxOpts contains optional arguments when retrieving inbox items.
type GetFaxInboxOpts struct {
	Period          string `json:"sPeriod,omitempty"`
	StartDate       string `json:"sStartDate,omitempty"`
	EndDate         string `json:"sEndDate,omitempty"`
	ViewedStatus    string `json:"sViewedStatus,omitempty"`
	IncludeSubUsers string `json:"sIncludeSubUsers,omitempty"`
}

// GetFaxInboxResp represents fax inbox information.
type GetFaxInboxResp struct {
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

// getFaxInboxReq defines the POST variables for a GetFaxInbox request
type getFaxInboxReq struct {
	Action string `json:"action"`
	Client
	GetFaxInboxOpts
}

// GetFaxInbox retrieves a list of faxes received for a specified period of time.
func (c *Client) GetFaxInbox(options ...GetFaxInboxOpts) (*GetFaxInboxResp, error) {
	opts := GetFaxInboxOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	if len(options) >= 1 {
		opts = options[0]
	}

	req := getFaxInboxReq{
		Action:          actionGetFaxInbox,
		Client:          *c,
		GetFaxInboxOpts: opts,
	}

	msg, err := sendPost(req)
	if err != nil {
		return nil, errors.Wrap(err, "GetFaxInbox SendPost error")
	}

	var resp GetFaxInboxResp
	if err := decodeResp(msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
