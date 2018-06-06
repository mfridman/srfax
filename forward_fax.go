package srfax

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ForwardOptions specify optional arguments when forwarding a fax.
type ForwardOptions struct {
	// The account number of a sub account, if you want to use a master account to
	// download a sub account’s fax
	SubUserID int `json:"sSubUserID,omitempty"`

	// Internal Reference Number (Maximum of 20 characters)
	AccountCode string `json:"sAccountCode,omitempty"`

	// Number of times the system is to retry a number if busy or an error is
	// encountered – number from 0 to 6
	Retries int `json:"sRetries,omitempty"`

	// From: On the Fax Header Line(Maximum of 30 characters)
	FaxFromHeader string `json:"sFaxFromHeader,omitempty"`

	// Provide an absolute URL (beginning with http:// or https://) and the SRFax
	// system will POST back the fax status record when the fax completes. See docs
	// for more details: https://www.srfax.com/api-page/forward_fax/
	NotifyURL string `json:"sNotifyURL,omitempty"`

	// The date you want to schedule a future fax for.
	// Must be in the format YYYY-MM-DD. Required if using QueueFaxTime
	QueueFaxDate string `json:"sQueueFaxDate,omitempty"`
	// The time you want to schedule a future fax for. Must be in the format HH:MM,
	// using 24 hour time (ie, 00:00 – 23:59). Required if using QueueFaxDate.
	// The timezone set on the account will be used when scheduling.
	QueueFaxTime string `json:"sQueueFaxTime,omitempty"`
}

func (o *ForwardOptions) validate() error {
	if len(o.AccountCode) > 20 {
		return errors.New("AccountCode must be a maximum of 20 characters")
	}
	if len(o.FaxFromHeader) > 30 {
		return errors.New("FaxFromHeader must be a maximum of 30 characters")
	}
	if o.NotifyURL != "" {
		http, https := "http://", "https://"
		if !strings.HasPrefix(o.NotifyURL, http) && !strings.HasPrefix(o.NotifyURL, https) {
			return errors.Errorf("NotifyURL must have prefix %q or %q", http, https)
		}
	}
	if o.Retries > 6 || o.Retries < 0 {
		return errors.New("Retries must be a number between 0-6")
	}
	if o.QueueFaxDate != "" && o.QueueFaxTime == "" {
		return errors.New("QueueFaxTime cannot be blank when supplying QueueFaxDate")
	}
	if o.QueueFaxTime != "" && o.QueueFaxDate == "" {
		return errors.New("QueueFaxDate cannot be blank when supplying QueueFaxTime")
	}
	if ok := validDateOrTime("2006-01-02", o.QueueFaxDate); !ok {
		return errors.New("QueueFaxDate must have format: YYYY-MM-DD")
	}
	if ok := validDateOrTime("15:04", o.QueueFaxTime); !ok {
		return errors.New("QueueFaxTime must have format: HH:MM, using 24 hour time")
	}
	return nil
}

// ForwardCfg specifies mandatory arguments when forwarding a fax.
type ForwardCfg struct {
	// Either FaxFileName or FaxDetailsID must be supplied
	//
	// FaxDetailsID of the fax – the ID is located after the "|" (pipe) character
	// of the FaxFileName
	FaxDetailsID string `json:"sFaxDetailsID,omitempty"`
	// FaxFileName returned from Get_Fax_Inbox or Get_Fax_Outbox
	FaxFileName string `json:"sFaxFileName,omitempty"`

	// IN or OUT for inbound or outbound
	Direction string `json:"sDirection"`

	// Sender fax number (must be 10 digits)
	CallerID int `json:"sCallerID"`

	// Sender email address
	SenderEmail string `json:"sSenderEmail"`

	// SINGLE when sending to one number; BROADCAST when sending to multiple numbers
	FaxType string `json:"sFaxType"`
	// Slice of string representing an 11 digit fax number
	ToFaxNumber []string `json:"-"`
}

// TODO: validate email with regex
func (c *ForwardCfg) validate() error {
	if err := hasEmpty(*c); err != nil {
		return errors.Wrap(err, "all fields in ForwardCfg are mandatory")
	}
	if c.FaxDetailsID != "" && c.FaxFileName != "" {
		return errors.New("either FaxFileName or FaxDetailsID must be supplied, not both")
	}
	if c.FaxDetailsID == "" && c.FaxFileName == "" {
		return errors.New("must supply either FaxDetailsID or FaxFileName")
	}
	if !(c.Direction == inbound || c.Direction == outbound) {
		return errors.Errorf("direction must be either: %s or %s", inbound, outbound)
	}
	invalid := make([]string, 0)
	for i, n := range c.ToFaxNumber {
		if ok := isNChars(n, 11); !ok {
			invalid = append(invalid, c.ToFaxNumber[i])
		}
	}
	if len(invalid) > 0 {
		return errors.Errorf("to fax number(s) must be 11 digits, found errors: %s", strings.Join(invalid, ", "))
	}
	if len(c.ToFaxNumber) > 1 {
		if c.FaxType != broadcast {
			return errors.Errorf("when supplying more than one fax number in ToFaxNumber, the FaxType must be set to %s", broadcast)
		}
	}
	if len(c.ToFaxNumber) == 1 && c.FaxType != single {
		return errors.Errorf("when supplying one fax number in ToFaxNumber, the FaxType must be %s", single)
	}
	if ok := isNChars(strconv.Itoa(c.CallerID), 10); !ok {
		return errors.Errorf("CallerID must be 10 digits: %d", c.CallerID)
	}
	return nil
}

// ForwardResp represents information about a forwarded fax.
type ForwardResp struct {
	Status string
	Result string
}

type mappedForwardResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// forwardOperation defines the POST variables for a ForwardFax request.
type forwardOperation struct {
	Action string `json:"action"`
	Client
	ForwardCfg
	ToFaxNumbers string `json:"sToFaxNumber"`
	ForwardOptions
}

func newForwardOperation(c *Client, cfg *ForwardCfg, opts *ForwardOptions) *forwardOperation {
	op := forwardOperation{Action: actionForwardFax, Client: *c, ForwardCfg: *cfg, ForwardOptions: *opts}
	op.ToFaxNumbers = strings.Join(cfg.ToFaxNumber, "|")
	return &op
}

// ForwardFax forwards a fax to other fax numbers.
func (c *Client) ForwardFax(cfg ForwardCfg, options ...ForwardOptions) (*ForwardResp, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	opts := ForwardOptions{}
	if len(options) > 0 {
		if err := options[0].validate(); err != nil {
			return nil, err
		}
		opts = options[0]
	}

	operation, err := constructReader(newForwardOperation(c, &cfg, &opts))
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct a reader from newForwardOperation")
	}

	result := mappedForwardResp{}
	if err := run(operation, &result, c.url); err != nil {
		return nil, err
	}

	out := ForwardResp(result)
	return &out, nil
}
