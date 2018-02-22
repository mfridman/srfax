package srfax

import "github.com/pkg/errors"

// FaxUsageOptions specify optional arguments to modify fax usage report.
type FaxUsageOptions struct {
	Period          string `json:"sPeriod,omitempty"`
	StartDate       string `json:"sStartDate,omitempty"`
	EndDate         string `json:"sEndDate,omitempty"`
	IncludeSubUsers string `json:"sIncludeSubUsers,omitempty"`
}

func (o *FaxUsageOptions) validate() error {
	if o.Period != "" {
		switch o.Period {
		case "RANGE":
			if ok := validDateOrTime("20060102", o.StartDate, o.EndDate); !ok {
				return errors.New("when Period set to RANGE must supply StartDate and EndDate; format must be YYYYMMDD")
			}
		case "ALL":
			if o.StartDate != "" || o.EndDate != "" {
				return errors.New("StartDate or EndDate only required when Period set to RANGE")
			}
		default:
			return errors.New("Period must be ALL|RANGE")
		}

	}

	if o.IncludeSubUsers != "" && o.IncludeSubUsers != "Y" {
		return errors.New(`IncludeSubUsers must be blank or set to "Y"`)
	}

	return nil
}

// FaxUsage is the response from a GetFaxUsage operation.
type FaxUsage struct {
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

// faxUsageRequest defines the POST variables for a GetFaxUsage request
type faxUsageRequest struct {
	Action string `json:"action"`
	Client
	FaxUsageOptions
}

// GetFaxUsage reports usage for a specified user and period.
func (c *Client) GetFaxUsage(options ...FaxUsageOptions) (*FaxUsage, error) {
	opts := FaxUsageOptions{}
	if len(options) >= 1 {
		opts = options[0]
		if err := opts.validate(); err != nil {
			return nil, err
		}
	}

	req := faxUsageRequest{
		Action:          actionGetFaxUsage,
		Client:          *c,
		FaxUsageOptions: opts,
	}

	var resp FaxUsage
	if err := run(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
