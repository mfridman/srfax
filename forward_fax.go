package srfax

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ForwardOptions specify optional arguments when forwarding a fax.
type ForwardOptions struct {
	SubUserID     int    `json:"sSubUserID,omitempty"`
	AccountCode   string `json:"sAccountCode,omitempty"`
	Retries       int    `json:"sRetries,omitempty"`
	FaxFromHeader string `json:"sFaxFromHeader,omitempty"`
	NotifyURL     string `json:"sNotifyURL,omitempty"`

	// YYYY-MM-DD
	QueueFaxDate string `json:"sQueueFaxDate,omitempty"`

	// HH:MM, using 24 hour time
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
		if !strings.HasPrefix(o.NotifyURL, "http://") && !strings.HasPrefix(o.NotifyURL, "https://") {
			return errors.New(`NotifyURL must have prefix "http://" or "https://"`)
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
	FaxDetailsID string `json:"sFaxDetailsID,omitempty"`
	FaxFileName  string `json:"sFaxFileName,omitempty"`

	// "IN" or "OUT" for inbound or outbound
	Direction string `json:"sDirection"`

	// sender's fax number (must be 10 digits)
	CallerID int `json:"sCallerID"`

	// sender's email address
	SenderEmail string `json:"sSenderEmail"`

	// "SINGLE" or "BROADCAST"
	FaxType string `json:"sFaxType"`

	// slice of string where each string must be an 11 digit number
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
	if len(options) >= 1 {
		if err := options[0].validate(); err != nil {
			return nil, err
		}
		opts = options[0]
	}

	operation, err := constructFromStruct(newForwardOperation(c, &cfg, &opts))
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
