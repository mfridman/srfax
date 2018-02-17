package srfax

import (
	"github.com/pkg/errors"
)

// InboxOptions contains optional arguments when retrieving inbox items.
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
				return errors.New("StartDate or EndDate only required when Period set to RANGE")
			}
		default:
			return errors.New("Period must be ALL|RANGE")
		}

	}

	if o.IncludeSubUsers != "" && o.IncludeSubUsers != "Y" {
		return errors.New(`IncludeSubUsers must be blank or set to "Y"`)
	}

	if o.ViewedStatus != "" {
		switch o.ViewedStatus {
		case "UNREAD":
			break
		case "READ":
			break
		case "ALL":
			break
		default:
			return errors.New("ViewedStatus must be blank or one of READ|UNREAD|ALL")
		}
	}

	return nil
}

// Inbox represents fax inbox information.
type Inbox struct {
	Status string `mapstructure:"Status"`
	Result []struct {
		FileName      string `mapstructure:"FileName"`
		ReceiveStatus string `mapstructure:"ReceiveStatus"`
		Date          string `mapstructure:"Date"`
		CallerID      string `mapstructure:"CallerID"`
		RemoteID      string `mapstructure:"RemoteID"`
		ViewedStatus  string `mapstructure:"ViewedStatus"`
		UserID        string `mapstructure:"User_ID" json:",omitempty"`
		UserFaxNumber string `mapstructure:"User_FaxNumber" json:",omitempty"`
		EpochTime     int    `mapstructure:"EpochTime"`
		Pages         int    `mapstructure:"Pages"`
		Size          int    `mapstructure:"Size"`
	} `mapstructure:"Result"`
}

// inboxRequest defines the POST variables for a GetFaxInbox request
type inboxRequest struct {
	Action string `json:"action"`
	Client
	InboxOptions
}

// GetFaxInbox retrieves a list of faxes received for a specified period of time.
func (c *Client) GetFaxInbox(options ...InboxOptions) (*Inbox, error) {
	opts := InboxOptions{}
	if len(options) >= 1 {
		opts = options[0]
		if err := opts.validate(); err != nil {
			return nil, err
		}
	}

	req := inboxRequest{
		Action:       actionGetFaxInbox,
		Client:       *c,
		InboxOptions: opts,
	}

	var resp Inbox
	if err := run(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
