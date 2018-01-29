package srfax

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/pkg/errors"
)

// MulFaxStatusResp represents the status of multiple sent faxes.
type MulFaxStatusResp struct {
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

// GetMulFaxStatus retrieves status of multiple sent faxes. Works only with outbound faxes.
// Accepts a slice of ids, i.e., FaxDetailIDs. Formatting handled automatically.
// FaxDetailsID returned from a QueueFax operation.
func (c *Client) GetMulFaxStatus(ids []string) (io.Reader, error) {

	if len(ids) == 0 {
		return nil, errors.New("when getting multiple fax status, must supply one or more FaxDetailsIDs (ids)")
	}

	msg := struct {
		Action string `json:"action"`
		ID     string `json:"sFaxDetailsID"`
		Client
	}{
		Action: actionGetMulFaxStatus,
		ID:     strings.Join(ids, "|"),
		Client: *c,
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
