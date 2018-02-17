package srfax

// MulFaxStatus represents the status of multiple sent faxes.
type MulFaxStatus struct {
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
		RemoteID    string `mapstructure:"RemoteID"`
		ErrorCode   string `mapstructure:"ErrorCode"`
		AccountCode string `mapstructure:"AccountCode"`
	} `mapstructure:"Result"`
}

// mulFaxStatusRequest defines the POST variables for a GetMulFaxStatus request
type mulFaxStatusRequest struct {
	Action string `json:"action"`
	Client
	IDs string `json:"sFaxDetailsID"`
}

// GetMulFaxStatus retrieves status of multiple sent faxes. Works only with outbound faxes.
// Multiple FaxDetailIDs (ids) can be requested by separating each FaxDetailsID with a "|" (pipe) character.
// Where FaxDetailsID returned from a QueueFax operation.
func (c *Client) GetMulFaxStatus(ids string) (*MulFaxStatus, error) {

	// if len(ids) == 0 {
	// 	return nil, errors.New("when getting multiple fax status, must supply one or more FaxDetailsIDs (ids)")
	// }

	req := mulFaxStatusRequest{
		Action: actionGetMulFaxStatus,
		IDs:    ids,
		Client: *c,
	}

	var resp MulFaxStatus
	if err := run(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
