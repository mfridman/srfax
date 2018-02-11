package srfax

import (
	"github.com/pkg/errors"
)

// UpdateViewedStatusResp is the response from a UpdateViewedStatus operation.
type UpdateViewedStatusResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// UpdateViewedStatusCfg contains mandatory arguments when updating the Viewed status of a fax.
type UpdateViewedStatusCfg struct {
	FaxDetailsID int    `json:"sFaxDetailsID,omitempty"` // Either the FaxFileName or the FaxDetailsID must be supplied
	FaxFileName  string `json:"sFaxFileName,omitempty"`  // Either the FaxFileName or the FaxDetailsID must be supplied
	Direction    string `json:"sDirection"`              // "IN" or "OUT" for inbound or outbound fax
	MarkAsViewed string `json:"sMarkasViewed"`           // "Y" marks fax READ, "N" marks fax UNREAD
}

// updatedViewedStatusReq defines the POST variables for a UpdateViewedStatus request
type updatedViewedStatusReq struct {
	Action string `json:"action"`
	Client
	UpdateViewedStatusCfg
}

// UpdateViewedStatus marks an inbound or outbound fax as read or unread.
//
// A note about ident:
// when passing FaxFileName, the entire name (including pipe and ID) must be supplied.
// E.g., 20170101230101-8812-34_0|29124120
func (c *Client) UpdateViewedStatus(cfg UpdateViewedStatusCfg) (*UpdateViewedStatusResp, error) {

	if cfg.FaxDetailsID > 0 && cfg.FaxFileName != "" {
		return nil, errors.New("Either FaxFileName or FaxDetailsID must be supplied, not both")
	}

	if !(cfg.Direction == inbound || cfg.Direction == outbound) {
		return nil, errors.Errorf("Direction must be either: %s or %s", inbound, outbound)
	}

	if !(cfg.MarkAsViewed == yes || cfg.MarkAsViewed == no) {
		return nil, errors.Errorf("MarkAsViewed must be either: %s or %s", yes, no)
	}

	req := updatedViewedStatusReq{
		Action:                actionUpdateViewedStatus,
		Client:                *c,
		UpdateViewedStatusCfg: cfg,
	}

	msg, err := sendPost(req)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateViewedStatusResp SendPost error")
	}

	var resp UpdateViewedStatusResp
	if err := decodeResp(msg, &resp); err != nil {
		return nil, errors.Wrap(err, "UpdateViewedStatusResp decodeResp error")
	}

	return &resp, nil
}
