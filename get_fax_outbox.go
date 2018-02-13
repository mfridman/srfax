package srfax

import "github.com/pkg/errors"

// GetFaxOutboxOpts contains optional arguments when retrieving outbox items.
type GetFaxOutboxOpts struct {
	Period          string `json:"sPeriod,omitempty"`
	StartDate       string `json:"sStartDate,omitempty"`
	EndDate         string `json:"sEndDate,omitempty"`
	IncludeSubUsers string `json:"sIncludeSubUsers,omitempty"`
}

// GetFaxOutboxResp represents fax outbox information.
type GetFaxOutboxResp struct {
	Status string `mapstructure:"Status"`
	Result []struct {
		FileName      string `mapstructure:"FileName"`
		SentStatus    string `mapstructure:"SentStatus"`
		DateQueued    string `mapstructure:"DateQueued"`
		DateSent      string `mapstructure:"DateSent"`
		EpochTime     string `mapstructure:"EpochTime"`
		ToFaxNumber   string `mapstructure:"ToFaxNumber"`
		RemoteID      string `mapstructure:"RemoteID"`
		ErrorCode     string `mapstructure:"ErrorCode"`
		AccountCode   string `mapstructure:"AccountCode"`
		Subject       string `mapstructure:"Subject"`
		UserID        string `mapstructure:"User_ID" json:",omitempty"`
		UserFaxNumber string `mapstructure:"User_FaxNumber" json:",omitempty"`
		Pages         int    `mapstructure:"Pages"`
		Duration      int    `mapstructure:"Duration"`
		Size          int    `mapstructure:"Size"`
	} `mapstructure:"Result"`
}

// getFaxOutboxReq defines the POST variables for a GetFaxOutbox request
type getFaxOutboxReq struct {
	Action string `json:"action"`
	Client
	GetFaxOutboxOpts
}

// GetFaxOutbox retrieves a list of faxes sent for a specified period of time.
func (c *Client) GetFaxOutbox(optArgs ...GetFaxOutboxOpts) (*GetFaxOutboxResp, error) {
	opts := GetFaxOutboxOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	req := getFaxOutboxReq{
		Action:           actionGetFaxOutbox,
		Client:           *c,
		GetFaxOutboxOpts: opts,
	}

	msg, err := sendPost(req)
	if err != nil {
		return nil, errors.Wrap(err, "GetFaxOutboxResp SendPost error")
	}

	var resp GetFaxOutboxResp
	if err := decodeResp(msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
