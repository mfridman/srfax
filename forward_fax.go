package srfax

import (
	"strings"

	"github.com/pkg/errors"
)

// ForwardFaxOpts contains optional arguments when forwarding a fax.
type ForwardFaxOpts struct {
	SubUserID     int    `json:"sSubUserID,omitempty"`
	AccountCode   string `json:"sAccountCode,omitempty"`
	Retries       int    `json:"sRetries,omitempty"`
	FaxFromHeader string `json:"sFaxFromHeader,omitempty"`
	NotifyURL     string `json:"sNotifyURL,omitempty"`
	QueueFaxDate  string `json:"sQueueFaxDate,omitempty"` // YYYY-MM-DD
	QueueFaxTime  string `json:"sQueueFaxTime,omitempty"` // HH:MM, using 24 hour time
}

// ForwardFaxCfg contains mandatory arguments when forwarding a fax.
type ForwardFaxCfg struct {
	FaxDetailsID string `json:"sFaxDetailsID,omitempty"` // Either FaxFileName or FaxDetailsID must be supplied
	FaxFileName  string `json:"sFaxFileName,omitempty"`  // Either FaxFileName or FaxDetailsID must be supplied
	Direction    string `json:"sDirection"`              // "IN" or "OUT" for inbound or outbound fax
	CallerID     int    `json:"sCallerID"`               // sender's fax number (must be 10 digits)
	SenderEmail  string `json:"sSenderEmail"`            // sender's email address
	FaxType      string `json:"sFaxType"`                // "SINGLE" or "BROADCAST"
	ToFaxNumber  string `json:"sToFaxNumber"`            // 11 digit number or up to 50 11 digit fax numbers separated by a “|” (pipe)
}

// ForwardFaxResp represents information about a forwarded fax.
type ForwardFaxResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// forwardFaxReq defines the POST variables for a ForwardFax request.
type forwardFaxReq struct {
	Action string `json:"action"`
	Client
	ForwardFaxCfg
	ForwardFaxOpts
}

// ForwardFax forwards a fax to other fax numbers.
func (c *Client) ForwardFax(cfg ForwardFaxCfg, optArgs ...ForwardFaxOpts) (*ForwardFaxResp, error) {

	opts := ForwardFaxOpts{}
	if len(optArgs) >= 1 {
		opts = optArgs[0]
	}

	if !(cfg.Direction == inbound || cfg.Direction == outbound) {
		return nil, errors.Errorf("Direction must be either: %s or %s", inbound, outbound)
	}

	ss := strings.Split(cfg.ToFaxNumber, "|")

	if len(ss) > 1 && cfg.FaxType != broadcast {
		return nil, errors.New("when supplying more than one fax number in ToFaxNumber, the FaxType must be set to BROADCAST")
	}

	if len(ss) == 1 && cfg.FaxType != single {
		return nil, errors.New("when supplying one fax number in ToFaxNumber, the FaxType must be set to SINGLE")
	}

	if cfg.FaxDetailsID == "" || cfg.FaxFileName == "" {
		return nil, errors.New("must supply either FaxDetailsID or FaxFileName")
	}

	req := forwardFaxReq{
		Action:         actionForwardFax,
		Client:         *c,
		ForwardFaxCfg:  cfg,
		ForwardFaxOpts: opts,
	}

	msg, err := sendPost(req)
	if err != nil {
		return nil, errors.Wrap(err, "ForwardFaxResp SendPost error")
	}

	var resp ForwardFaxResp
	if err := decodeResp(msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
