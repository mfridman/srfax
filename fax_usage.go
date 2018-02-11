package srfax

// FaxUsageOpts contains optional arguments to modify fax usage report.
type FaxUsageOpts struct {
	Period          string `json:"sPeriod,omitempty"`
	StartDate       string `json:"sStartDate,omitempty"`
	EndDate         string `json:"sEndDate,omitempty"`
	IncludeSubUsers string `json:"sIncludeSubUsers,omitempty"`
}

// FaxUsageResp is the response from a GetFaxUsage operation.
type FaxUsageResp struct {
	Status string `mapstructure:"Status"`
	Result []struct {
		Period        string `mapstructure:"Period"`
		ClientName    string `mapstructure:"ClientName"`
		BillingNumber string `mapstructure:"BillingNumber"`
		UserID        int    `mapstructure:"UserID"`
		SubUserID     int    `mapstructure:"SubUserID"`
		NumberOfFaxes int    `mapstructure:"NumberOfFaxes"`
		NumberOfPages int    `mapstructure:"NumberOfPages"`
	} `mapstructure:"Result"`
}

// FaxUsageReq defines the POST variables for a GetFaxUsage request
type FaxUsageReq struct {
	Action string `json:"action"`
	Client
	FaxUsageOpts
}

// GetFaxUsage reports usage for a specified user and period.
func (c *Client) GetFaxUsage(optArgs ...FaxUsageOpts) (*FaxUsageReq, error) {
	opts := FaxUsageOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	req := FaxUsageReq{
		Action:       actionGetFaxUsage,
		Client:       *c,
		FaxUsageOpts: opts,
	}

	return &req, nil
}
