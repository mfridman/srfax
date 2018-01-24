package srfax

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// DeleteFaxOpts contains optional arguments when deleting faxes.
type DeleteFaxOpts struct {
	ResponseFormat string
}

// DeleteFaxResp is the response from a DeleteFax operation.
// Note, error message from Result, if any, will be stored in ResultError.
type DeleteFaxResp struct {
	Status      string `mapstructure:"Status"`
	Result      string `mapstructure:"Result"`
	ResultError string
}

// DeleteFax deletes either, one ore more, received or sent faxes.
//
// dir is the direction of fax: "IN" or "OUT" for inbound or outbound fax
//
// TODO (MF): Status always seems to be "Success", even after item(s) deleted. Even if ID is a
// fake, still "Success"
// Also, passing in an incorrect "Name" with a valid ID will delete item.
// backend SRFax system uses ID or pipe+ID to delete
// OUT || IN:
// wrong name, correct id = deletion ... foobar|31524120baz will trigger a deletion
// correct name, wrong id = nothing
// wrong name, wronag id = nothing (just in case)
func (c *Client) DeleteFax(ids []string, dir string, optArgs ...DeleteFaxOpts) (*DeleteFaxResp, error) {
	if !(dir == "IN" || dir == "OUT") {
		return nil, errors.New(`dir (direction) must be one of either "IN" or "OUT"`)
	}

	msg := map[string]interface{}{
		"action":     actionDeleteFax,
		"access_id":  c.AccessID,
		"access_pwd": c.AccessPwd,
		"sDirection": dir,
	}

	if len(optArgs) >= 1 {
		msg["sResponseFormat"] = optArgs[0].ResponseFormat
	}

	const (
		prefixName = "sFaxFileName_"
		prefixID   = "sFaxDetailsID_"
	)

	if len(ids) <= 0 {
		return nil, errors.New("must supply one or more identifiers when deleting faxes. Accepts multiple fax file names or multiple IDs")
	}

	// no strict checking because SRFax API will only delete valid IDs
	for i, j := range ids {
		if strings.Contains(j, "|") {
			// this is useless as the backend SRFax system uses ID or pipe+ID to delete
			// foobar|31524120baz will trigger a deletion
			msg[prefixName+strconv.Itoa(i)] = j
		} else {
			msg[prefixID+strconv.Itoa(i)] = j
		}
	}

	resp, err := sendPost(msg, c.url)
	if err != nil {
		return nil, errors.Wrap(err, "sendPost failed")
	}

	if st, err := checkStatus(resp); err != nil {
		return &DeleteFaxResp{Status: st, ResultError: fmt.Sprint(err)}, nil
	}

	var result DeleteFaxResp
	var md mapstructure.Metadata
	cfg := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Metadata:         &md,
		Result:           &result,
	}

	if err := decodeResp(resp, cfg); err != nil {
		return nil, err
	}

	return &result, nil
}
