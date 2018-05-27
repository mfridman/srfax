package srfax

import (
	"github.com/pkg/errors"
)

// InboxOptions specify optional arguments when retrieving inbox items.
type InboxOptions struct {
	Period          string `json:"sPeriod,omitempty"`
	StartDate       string `json:"sStartDate,omitempty"`
	EndDate         string `json:"sEndDate,omitempty"`
	ViewedStatus    string `json:"sViewedStatus,omitempty"`
	IncludeSubUsers string `json:"sIncludeSubUsers,omitempty"`
}

func (o *InboxOptions) validate() error {

	if o.Period != "" {
		switch o.Period {
		case "RANGE":
			if ok := validDateOrTime("20060102", o.StartDate, o.EndDate); !ok {
				return errors.New("when Period set to RANGE must supply StartDate and EndDate; format must be YYYYMMDD")
			}
		case "ALL":
			if o.StartDate != "" || o.EndDate != "" {
				return errors.New("StartDate and EndDate are not required when Period set to ALL")
			}
		default:
			return errors.New("Period must be ALL|RANGE")
		}
	}

	if o.IncludeSubUsers != "" && o.IncludeSubUsers != yes {
		return errors.Errorf("IncludeSubUsers must be omitted or set to %q", yes)
	}

	if o.ViewedStatus != "" {
		switch o.ViewedStatus {
		case "UNREAD", "READ", "ALL":
			break
		default:
			return errors.New("ViewedStatus must be blank or one of READ, UNREAD or ALL")
		}
	}

	return nil
}

// mappedInbox represents an internal mapstructure inbox
type mappedInbox struct {
	Status string `mapstructure:"Status"`
	Result []struct {
		FileName      string `mapstructure:"FileName"`
		ReceiveStatus string `mapstructure:"ReceiveStatus"`
		Date          string `mapstructure:"Date"`
		CallerID      string `mapstructure:"CallerID"`
		RemoteID      string `mapstructure:"RemoteID"`
		ViewedStatus  string `mapstructure:"ViewedStatus"`
		UserID        string `mapstructure:"User_ID,omitempty"`
		UserFaxNumber string `mapstructure:"User_FaxNumber,omitempty"`
		EpochTime     int    `mapstructure:"EpochTime"`
		Pages         int    `mapstructure:"Pages"`
		Size          int    `mapstructure:"Size"`
	} `mapstructure:"Result"`
}

// Inbox represents fax inbox information.
type Inbox struct {
	Status string
	Result []struct {
		FileName      string
		ReceiveStatus string
		Date          string
		CallerID      string
		RemoteID      string
		ViewedStatus  string
		UserID        string
		UserFaxNumber string
		EpochTime     int
		Pages         int
		Size          int
	}
}

// inboxOperation defines the POST variables for a GetFaxInbox operation.
type inboxOperation struct {
	Action string `json:"action"`
	Client
	InboxOptions
}

func newInboxOperation(c *Client, o *InboxOptions) *inboxOperation {
	return &inboxOperation{Action: actionGetFaxInbox, Client: *c, InboxOptions: *o}
}

func newInboxOptions(options ...InboxOptions) (*InboxOptions, error) {
	opts := InboxOptions{}
	if len(options) >= 1 {
		if err := options[0].validate(); err != nil {
			return nil, err
		}
		opts = options[0]
	}
	return &opts, nil
}

// GetFaxInbox retrieves a list of faxes received for a specified period of time.
func (c *Client) GetFaxInbox(options ...InboxOptions) (*Inbox, error) {
	opts, err := newInboxOptions(options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed options")
	}

	operation, err := constructFromStruct(newInboxOperation(c, opts))
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct a reader from newInboxOperation")
	}

	result := mappedInbox{}
	if err := run(operation, &result); err != nil {
		return nil, err
	}

	out := Inbox(result)
	return &out, nil
}
