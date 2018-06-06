package srfax

import (
	"github.com/pkg/errors"
)

// ViewedStatus is the response from a UpdateViewedStatus operation.
type ViewedStatus struct {
	Status string
	Result string
}

type mappedViewedStatus struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// ViewedStatusCfg specifies mandatory arguments when updating the Viewed status of a fax.
type ViewedStatusCfg struct {
	// Either the FaxFileName or the FaxDetailsID must be supplied
	//
	// When passing FaxFileName, the entire name (including pipe and ID)
	// must be supplied. E.g., 20170101230101-1212-21_7|12124720
	FaxDetailsID int    `json:"sFaxDetailsID,omitempty"`
	FaxFileName  string `json:"sFaxFileName,omitempty"`

	// IN or OUT for inbound or outbound fax
	Direction string `json:"sDirection"`

	// Y marks fax READ, N marks fax UNREAD
	MarkAsViewed string `json:"sMarkasViewed"`
}

func (c *ViewedStatusCfg) validate() error {
	if c.FaxDetailsID > 0 && c.FaxFileName != "" {
		return errors.New("Either FaxFileName or FaxDetailsID must be supplied, not both")
	}
	if !(c.Direction == inbound || c.Direction == outbound) {
		return errors.Errorf("Direction must be either: %s or %s", inbound, outbound)
	}
	if !(c.MarkAsViewed == yes || c.MarkAsViewed == no) {
		return errors.Errorf("MarkAsViewed must be either: %s or %s", yes, no)
	}
	return nil
}

// viewedStatusOperation defines the POST variables for a UpdateViewedStatus request
type viewedStatusOperation struct {
	Action string `json:"action"`
	Client
	ViewedStatusCfg
}

func newViewedStatusOperation(c *Client, cfg *ViewedStatusCfg) *viewedStatusOperation {
	return &viewedStatusOperation{Action: actionUpdateViewedStatus, Client: *c, ViewedStatusCfg: *cfg}
}

// UpdateViewedStatus marks an inbound or outbound fax as read or unread.
func (c *Client) UpdateViewedStatus(cfg ViewedStatusCfg) (*ViewedStatus, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	operation, err := constructReader(newViewedStatusOperation(c, &cfg))
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct a reader for newViewedStatusOperation")
	}

	result := mappedViewedStatus{}
	if err := run(operation, &result, c.url); err != nil {
		return nil, err
	}

	out := ViewedStatus(result)
	return &out, nil
}
