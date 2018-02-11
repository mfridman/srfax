package srfax

import (
	"github.com/pkg/errors"
)

// ViewedStatusResp is the response from a UpdateViewedStatus operation.
type ViewedStatusResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// ViewedStatusCfg contains mandatory arguments when updating the viewed status of a fax.
type ViewedStatusCfg struct {
	FaxDetailsID int    `json:"sFaxDetailsID,omitempty"` // Either the FaxFileName or the FaxDetailsID must be supplied
	FaxFileName  string `json:"sFaxFileName,omitempty"`  // Either the FaxFileName or the FaxDetailsID must be supplied
	Direction    string `json:"sDirection"`              // "IN" or "OUT" for inbound or outbound fax
	MarkAsViewed string `json:"sMarkasViewed"`           // "Y" marks fax READ, "N" marks fax UNREAD
}

// ViewedStatusReq defines the POST variables for a UpdateViewedStatus request
type ViewedStatusReq struct {
	Action string `json:"action"`
	Client
	ViewedStatusCfg
}

// UpdateViewedStatus marks an inbound or outbound fax as read or unread.
//
// A note about ident:
// when passing FaxFileName, the entire name (including pipe and ID) must be supplied.
// E.g., 20170101230101-8812-34_0|29124120
func (c *Client) UpdateViewedStatus(cfg ViewedStatusCfg) (*ViewedStatusReq, error) {

	if cfg.FaxDetailsID > 0 && cfg.FaxFileName != "" {
		return nil, errors.New("Either FaxFileName or FaxDetailsID must be supplied, not both")
	}

	if !(cfg.Direction == inbound || cfg.Direction == outbound) {
		return nil, errors.Errorf("Direction must be either: %s or %s", inbound, outbound)
	}

	if !(cfg.MarkAsViewed == yes || cfg.MarkAsViewed == no) {
		return nil, errors.Errorf("MarkAsViewed must be either: %s or %s", yes, no)
	}

	req := ViewedStatusReq{
		Action:          actionUpdateViewedStatus,
		Client:          *c,
		ViewedStatusCfg: cfg,
	}

	return &req, nil
}
