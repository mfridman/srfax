package srfax

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ForwardFaxOpts contains optional arguments when forwarding a fax.
type ForwardFaxOpts struct {
	SubUserID     int    `json:"sSubUserID"`
	AccountCode   string `json:"sAccountCode"`
	Retries       int    `json:"sRetries"`
	FaxFromHeader string `json:"sFaxFromHeader"`
	NotifyURL     string `json:"sNotifyURL"`
	QueueFaxDate  string `json:"sQueueFaxDate"` // YYYY-MM-DD
	QueueFaxTime  string `json:"sQueueFaxTime"` // HH:MM, using 24 hour time
}

// ForwardFaxCfg contains mandatory arguments when forwarding a fax.
type ForwardFaxCfg struct {
	CallerID    int      // sender's fax number (must be 10 digits)
	SenderEmail string   // sender's email address
	FaxType     string   // "SINGLE" or "BROADCAST"
	ToFaxNumber []string // each number must be 11 digits represented as a String
}

// ForwardFaxResp represents information about a forwarded fax.
type ForwardFaxResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// ForwardFax forwards a fax to other fax numbers
//
// ident is the sFaxDetailsID or sFaxFileName returned from GetFaxInbox or GetFaxOutbox
// dir is the direction; "IN" or "OUT" for inbound or outbound fax
func (c *Client) ForwardFax(dir, ident string, fc ForwardFaxCfg, optArgs ...ForwardFaxOpts) (io.Reader, error) {

	if !(dir == "IN" || dir == "OUT") {
		return nil, errors.New(`dir (direction) must be one of either "IN" or "OUT"`)
	}

	l := len(fc.ToFaxNumber)
	if l > 1 && fc.FaxType != "BROADCAST" {
		return nil, errors.New("when supplying many fax numbers in ToFaxNumber, the type must be BROADCAST")
	}

	opts := ForwardFaxOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	msg := struct {
		Action string `json:"action"`
		Client
		FaxDetailsID int    `json:"sFaxDetailsID,omitempty"`
		FaxFileName  string `json:"sFaxFileName,omitempty"`
		Direction    string `json:"sDirection"`
		CallerID     int    `json:"sCallerID"`
		SenderEmail  string `json:"sSenderEmail"`
		FaxType      string `json:"sFaxType"`
		ToFaxNumber  string `json:"sToFaxNumber"`
		ForwardFaxOpts
	}{
		Action:         actionForwardFax,
		Client:         *c,
		Direction:      dir,
		CallerID:       fc.CallerID,
		SenderEmail:    fc.SenderEmail,
		FaxType:        fc.FaxType,
		ToFaxNumber:    strings.Join(fc.ToFaxNumber, "|"),
		ForwardFaxOpts: opts,
	}

	if strings.Contains(ident, "|") {
		msg.FaxFileName = ident
	} else {
		n, err := strconv.Atoi(ident)
		if err != nil {
			return nil, errors.New("failed id string to int conversion")
		}
		msg.FaxDetailsID = n
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
