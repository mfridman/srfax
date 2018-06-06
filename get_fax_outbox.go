package srfax

import "github.com/pkg/errors"

// OutboxOptions specify optional arguments when retrieving outbox items.
type OutboxOptions struct {
	// ALL or RANGE if not provided defaults to ALL
	Period string `json:"sPeriod,omitempty"`

	// Only required if RANGE is specified in sPeriod – date format must be YYYYMMDD
	StartDate string `json:"sStartDate,omitempty"`
	EndDate   string `json:"sEndDate,omitempty"`

	// Set to Y to include faxes received by a sub user of the account as well
	IncludeSubUsers string `json:"sIncludeSubUsers,omitempty"`
}

func (o *OutboxOptions) validate() error {
	if o.Period != "" {
		switch o.Period {
		case "RANGE":
			if ok := validDateOrTime("20060102", o.StartDate, o.EndDate); !ok {
				return errors.New("when Period set to RANGE must supply StartDate and EndDate; format must be YYYYMMDD")
			}
		case "ALL":
			if o.StartDate != "" || o.EndDate != "" {
				return errors.New("StartDate and/or EndDate not required when setting Period to ALL")
			}
		default:
			return errors.New("Period can be blank or one of ALL, RANGE. If not provided defaults to ALL")
		}
	}
	if o.IncludeSubUsers != "" && o.IncludeSubUsers != yes {
		return errors.Errorf("IncludeSubUsers must be omitted or set to %q", yes)
	}
	return nil
}

// Outbox represents fax outbox information. More information can be found on the official docs:
// https://www.srfax.com/api-page/get_fax_outbox/, look for JSON Returned Variables.
type Outbox struct {
	Status string
	Result []struct {
		FileName      string
		SentStatus    string
		DateQueued    string
		DateSent      string
		EpochTime     string
		ToFaxNumber   string
		RemoteID      string
		ErrorCode     string
		AccountCode   string
		Subject       string
		UserID        string
		UserFaxNumber string
		Pages         int
		Duration      int
		Size          int
	}
}

type mappedOutbox struct {
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
		UserID        string `mapstructure:"User_ID,omitempty"`
		UserFaxNumber string `mapstructure:"User_FaxNumber,omitempty"`
		Pages         int    `mapstructure:"Pages"`
		Duration      int    `mapstructure:"Duration"`
		Size          int    `mapstructure:"Size"`
	} `mapstructure:"Result"`
}

// outboxOperation defines the POST variables for a GetFaxOutbox request
type outboxOperation struct {
	Action string `json:"action"`
	Client
	OutboxOptions
}

func newOutboxOperation(c *Client, o *OutboxOptions) *outboxOperation {
	return &outboxOperation{Action: actionGetFaxOutbox, Client: *c, OutboxOptions: *o}
}

func newOutboxOptions(options ...OutboxOptions) (*OutboxOptions, error) {
	opts := OutboxOptions{}
	if len(options) > 0 {
		if err := options[0].validate(); err != nil {
			return nil, err
		}
		opts = options[0]
	}
	return &opts, nil
}

// GetFaxOutbox retrieves a list of faxes sent for a specified period of time.
func (c *Client) GetFaxOutbox(options ...OutboxOptions) (*Outbox, error) {
	opts, err := newOutboxOptions(options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed options")
	}

	operation, err := constructReader(newOutboxOperation(c, opts))
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct a reader from newOutboxOperation")
	}

	result := mappedOutbox{}
	if err := run(operation, &result, c.url); err != nil {
		return nil, err
	}

	out := Outbox(result)
	return &out, nil
}
