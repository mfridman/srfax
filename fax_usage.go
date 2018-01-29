package srfax

import (
	"bytes"
	"encoding/json"
	"io"
)

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

// GetFaxUsage reports usage for a specified user and period.
func (c *Client) GetFaxUsage(optArgs ...FaxUsageOpts) (io.Reader, error) {
	opts := FaxUsageOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	msg := struct {
		Action string `json:"action"`
		Client
		FaxUsageOpts
	}{
		Action:       actionGetFaxUsage,
		Client:       *c,
		FaxUsageOpts: opts,
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
