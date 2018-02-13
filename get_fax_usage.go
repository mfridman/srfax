package srfax

import "github.com/pkg/errors"

// GetFaxUsageOpts contains optional arguments to modify fax usage report.
type GetFaxUsageOpts struct {
	Period          string `json:"sPeriod,omitempty"`
	StartDate       string `json:"sStartDate,omitempty"`
	EndDate         string `json:"sEndDate,omitempty"`
	IncludeSubUsers string `json:"sIncludeSubUsers,omitempty"`
}

// GetFaxUsageResp is the response from a GetFaxUsage operation.
type GetFaxUsageResp struct {
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

// getFaxUsageReq defines the POST variables for a GetFaxUsage request
type getFaxUsageReq struct {
	Action string `json:"action"`
	Client
	GetFaxUsageOpts
}

// GetFaxUsage reports usage for a specified user and period.
func (c *Client) GetFaxUsage(optArgs ...GetFaxUsageOpts) (*GetFaxUsageResp, error) {
	opts := GetFaxUsageOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	req := getFaxUsageReq{
		Action:          actionGetFaxUsage,
		Client:          *c,
		GetFaxUsageOpts: opts,
	}

	msg, err := sendPost(req)
	if err != nil {
		return nil, errors.Wrap(err, "GetFaxUsageResp SendPost error")
	}

	var resp GetFaxUsageResp
	if err := decodeResp(msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
