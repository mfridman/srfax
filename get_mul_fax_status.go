package srfax

import "github.com/pkg/errors"

// GetMulFaxStatusResp represents the status of multiple sent faxes.
type GetMulFaxStatusResp struct {
	Status string `mapstructure:"Status"`
	Result []struct {
		Pages string `mapstructure:"Pages"` // API returns empty string when there is an error
		// TODO this comes across as an error regardless
		// decoding resp body: mapstructure: cannot unmarshal string into Go struct field .EpochTime of type int
		// decoding resp body: mapstructure: cannot unmarshal number into Go struct field .EpochTime of type string
		EpochTime   string `mapstructure:"EpochTime"` // API returns empty string when there is an error
		Duration    string `mapstructure:"Duration"`
		Size        string `mapstructure:"Size"`
		FileName    string `mapstructure:"FileName"`
		SentStatus  string `mapstructure:"SentStatus"`
		DateQueued  string `mapstructure:"DateQueued"`
		DateSent    string `mapstructure:"DateSent"`
		ToFaxNumber string `mapstructure:"ToFaxNumber"`
		RemoteID    string `mapstructure:"RemoteID"` // API docs incorrect, this is a string, not an "integer"
		ErrorCode   string `mapstructure:"ErrorCode"`
		AccountCode string `mapstructure:"AccountCode"`
	} `mapstructure:"Result"`
}

// getMulFaxStatusReq defines the POST variables for a GetMulFaxStatus request
type getMulFaxStatusReq struct {
	Action string `json:"action"`
	Client
	IDs string `json:"sFaxDetailsID"`
}

// GetMulFaxStatus retrieves status of multiple sent faxes. Works only with outbound faxes.
// Multiple FaxDetailIDs (ids) can be requested by separating each FaxDetailsID with a "|" (pipe) character.
// Where FaxDetailsID returned from a QueueFax operation.
func (c *Client) GetMulFaxStatus(ids string) (*GetMulFaxStatusResp, error) {

	// if len(ids) == 0 {
	// 	return nil, errors.New("when getting multiple fax status, must supply one or more FaxDetailsIDs (ids)")
	// }

	req := getMulFaxStatusReq{
		Action: actionGetMulFaxStatus,
		IDs:    ids,
		Client: *c,
	}

	msg, err := sendPost(req)
	if err != nil {
		return nil, errors.Wrap(err, "GetMulFaxStatusResp SendPost error")
	}

	var resp GetMulFaxStatusResp
	if err := decodeResp(msg, &resp); err != nil {
		return nil, errors.Wrap(err, "GetMulFaxStatusResp decodeResp error")
	}

	return &resp, nil
}
