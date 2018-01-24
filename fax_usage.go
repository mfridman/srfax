package srfax

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// FaxUsageOpts contains optional arguments to modify fax usage report.
type FaxUsageOpts struct {
	ResponseFormat  string `json:"sResponseFormat,omitempty"`
	Period          string `json:"sPeriod,omitempty"`
	StartDate       string `json:"sStartDate,omitempty"`
	EndDate         string `json:"sEndDate,omitempty"`
	IncludeSubUsers string `json:"sIncludeSubUsers,omitempty"`
}

// FaxUsageResp is the response from a GetFaxUsage operation.
// Note, error message from Result, if any, will be stored in ResultError.
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
	ResultError string
}

// GetFaxUsage reports usage for a specified user and period.
func (c *Client) GetFaxUsage(optArgs ...FaxUsageOpts) (*FaxUsageResp, error) {
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

	resp, err := sendPost(msg, c.url)
	if err != nil {
		return nil, errors.Wrap(err, "sendPost failed")
	}

	if st, err := checkStatus(resp); err != nil {
		return &FaxUsageResp{Status: st, ResultError: fmt.Sprint(err)}, nil
	}

	var result FaxUsageResp
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
