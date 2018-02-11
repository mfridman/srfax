package srfax

// FaxOutboxOpts contains optional arguments when retrieving outbox items.
type FaxOutboxOpts struct {
	Period          string `json:"sPeriod,omitempty"`
	StartDate       string `json:"sStartDate,omitempty"`
	EndDate         string `json:"sEndDate,omitempty"`
	IncludeSubUsers string `json:"sIncludeSubUsers,omitempty"`
}

// FaxOutboxResp represents fax outbox information.
type FaxOutboxResp struct {
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

// FaxOutboxReq defines the POST variables for a GetFaxOutbox request
type FaxOutboxReq struct {
	Action string `json:"action"`
	Client
	FaxOutboxOpts
}

// GetFaxOutbox retrieves a list of faxes sent for a specified period of time.
func (c *Client) GetFaxOutbox(optArgs ...FaxOutboxOpts) (*FaxOutboxReq, error) {
	opts := FaxOutboxOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	req := FaxOutboxReq{
		Action:        actionGetFaxOutbox,
		Client:        *c,
		FaxOutboxOpts: opts,
	}

	return &req, nil
}
