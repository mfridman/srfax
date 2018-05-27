package srfax

import (
	"github.com/pkg/errors"
)

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
	if o.IncludeSubUsers != "" && o.IncludeSubUsers != yes {
		return errors.Errorf(`IncludeSubUsers must be blank or set to "%s"`, yes)
	}
	return nil
}

// FaxUsage is the response from a GetFaxUsage operation.
type FaxUsage struct {
	Status string
	Result []struct {
		Period        string
		ClientName    string
		BillingNumber string
		UserID        int
		SubUserID     int
		NumberOfFaxes int
		NumberOfPages int
	}
}

type mappedFaxUsage struct {
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

// faxUsageOperation defines the POST variables for a GetFaxUsage request
type faxUsageOperation struct {
	Action string `json:"action"`
	Client
	FaxUsageOptions
}

func newFaxUsageOperation(c *Client, opts *FaxUsageOptions) *faxUsageOperation {
	return &faxUsageOperation{Action: actionGetFaxUsage, Client: *c, FaxUsageOptions: *opts}
}

// GetFaxUsage reports usage for a specified user and period.
func (c *Client) GetFaxUsage(options ...FaxUsageOptions) (*FaxUsage, error) {
	opts := FaxUsageOptions{}
	if len(options) >= 1 {
		if err := options[0].validate(); err != nil {
			return nil, err
		}
		opts = options[0]
	}

	operation, err := constructFromStruct(newFaxUsageOperation(c, &opts))
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct a reader from newFaxUsageOperation")
	}

	result := mappedFaxUsage{}
	if err := run(operation, &result); err != nil {
		return nil, err
	}

	out := FaxUsage(result)
	return &out, nil
}
