package srfax

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// DeleteResp is the response from a DeleteFax operation.
type DeleteResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// DeleteFax deletes one or more received or sent faxes for a given direction.
//
// dir (direction) must be one of "IN" or "OUT" for inbound or outbound.
//
// ids is a slice of faxes to delete based on FaxFileName or FaxDetailsID.
// These are unique identifiers returned from a GetFaxOutbox or GetFaxInbox operation.
// Note, this method will take care of formatting the request accordingly, so it is
// safe to mix filenames with IDs: []string{"20170721124555-1213-4_0|272568938", "172568938"}
func (c *Client) DeleteFax(ids []string, dir string) (*DeleteResp, error) {
	if !(dir == inbound || dir == outbound) {
		return nil, errors.New(`dir (direction) must be one of either "IN" or "OUT"`)
	}

	if len(ids) <= 0 {
		return nil, errors.New("must supply one or more identifiers when deleting faxes")
	}

	req := map[string]interface{}{
		"action":     actionDeleteFax,
		"access_id":  c.AccessID,
		"access_pwd": c.AccessPwd,
		"sDirection": dir,
	}

	const (
		prefixName = "sFaxFileName_"
		prefixID   = "sFaxDetailsID_"
	)

	for i, j := range ids {
		if strings.Contains(j, "|") {
			req[prefixName+strconv.Itoa(i)] = j
		} else {
			req[prefixID+strconv.Itoa(i)] = j
		}
	}

	var resp DeleteResp
	if err := run(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
