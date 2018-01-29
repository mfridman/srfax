package srfax

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ViewedStatusResp is the response from a UpdateViewedStatus operation.
type ViewedStatusResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// UpdateViewedStatus marks an inbound or outbound fax as read or unread.
//
// dir (direction) will be "IN" or "OUT" for inbound or outbound fax.
// view will be "Y" – mark fax as READ, or "N" – mark fax as UNREAD
// A note about ident:
//
// when passing sFaxFileName, the entire name (including pipe and ID) must be supplied.
// E.g., 20180101230101-8812-34_0|31524120
//
// If updating a fax based on sFaxDetailsID, pass in the number as a string.
// Formatting handled automatically.
func (c *Client) UpdateViewedStatus(ident, dir, view string) (io.Reader, error) {
	// TODO consider wrapping the string params "ident, dir, view" into a struct

	msg := struct {
		Action string `json:"action"`
		Client
		FaxDetailsID int    `json:"sFaxDetailsID,omitempty"` // mutually exclusive
		FaxFileName  string `json:"sFaxFileName,omitempty"`  // mutually exclusive
		Direction    string `json:"sDirection"`
		MarkAsViewed string `json:"sMarkasViewed"`
	}{
		Action:       actionUpdateViewedStatus,
		Client:       *c,
		Direction:    dir,
		MarkAsViewed: view,
	}

	if strings.Contains(ident, "|") {
		msg.FaxFileName = ident
	} else {
		n, err := strconv.Atoi(ident)
		if err != nil {
			return nil, errors.Errorf("failed updating viewed status. sFaxDetailsID (id) string to int conversion, got [%[1]v] of type [%[1]T].", ident)
		}
		msg.FaxDetailsID = n
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
