package srfax

import "github.com/pkg/errors"

// OutboxOptions specify optional arguments when retrieving outbox items.
type OutboxOptions struct {
	Period          string `json:"sPeriod,omitempty"`
	StartDate       string `json:"sStartDate,omitempty"`
	EndDate         string `json:"sEndDate,omitempty"`
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

// Outbox represents fax outbox information.
type Outbox struct {
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

// outboxOperation defines the POST variables for a GetFaxOutbox request
type outboxOperation struct {
	Action string `json:"action"`
	Client
	OutboxOptions
}

func newOutboxOperation(c *Client, o *OutboxOptions) *outboxOperation {
	return &outboxOperation{Action: actionGetFaxOutbox, Client: *c, OutboxOptions: *o}
}

// GetFaxOutbox retrieves a list of faxes sent for a specified period of time.
func (c *Client) GetFaxOutbox(options ...OutboxOptions) (*Outbox, error) {
	opts := OutboxOptions{}
	if len(options) >= 1 {
		if err := options[0].validate(); err != nil {
			return nil, err
		}
		opts = options[0]
	}
	resp := Outbox{}
	opr := newOutboxOperation(c, &opts)
	if err := run(opr, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
