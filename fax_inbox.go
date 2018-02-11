package srfax

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

// FaxInboxReq defines the POST variables for a GetFaxInbox request
type FaxInboxReq struct {
	Action string `json:"action"`
	Client
	FaxInboxOpts
}

// GetFaxInbox retrieves a list of faxes received for a specified period of time.
func (c *Client) GetFaxInbox(optArgs ...FaxInboxOpts) (*FaxInboxReq, error) {
	opts := FaxInboxOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	req := FaxInboxReq{
		Action:       actionGetFaxInbox,
		Client:       *c,
		FaxInboxOpts: opts,
	}

	return &req, nil
}
