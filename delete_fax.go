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

// DeleteFax deletes one or more received or sent faxes.
//
// dir is the direction of fax: "IN" or "OUT" for inbound or outbound fax
func (c *Client) DeleteFax(ids []string, dir string) (*DeleteResp, error) {
	if !(dir == inbound || dir == outbound) {
		return nil, errors.New(`dir (direction) must be one of either "IN" or "OUT"`)
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

	if len(ids) <= 0 {
		return nil, errors.New("must supply one or more identifiers when deleting faxes. Accepts multiple fax file names or multiple IDs")
	}

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
