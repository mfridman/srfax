package srfax

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// QueueOptions specify optional arguments when sending faxes.
//
// If the default cover page on the account is set to "Attachments ONLY" the cover page will
// not be created irrespective of this variable
// If a cover page is not provided then all cover page variables will be ignored
//
// "Basic", "Standard", "Company", or "Personal"
//
// json tag used for reflection.
type QueueOptions struct {
	Retries       int    `json:"sRetries"`
	AccountCode   string `json:"sAccountCode"`
	FaxFromHeader string `json:"sFaxFromHeader"`

	// cover page
	CoverPage      string `json:"sCoverPage"`
	CPFromName     string `json:"sCPFromName"`
	CPToName       string `json:"sCPToName"`
	CPOrganization string `json:"sCPOrganization"`
	CPSubject      string `json:"sCPSubject"`
	CPComments     string `json:"sCPComments"`

	NotifyURL string `json:"sNotifyURL"`

	// YYYY-MM-DD
	QueueFaxDate string `json:"sQueueFaxDate"`
	// HH:MM, using 24 hour time
	QueueFaxTime string `json:"sQueueFaxTime"`
}

// QueueCfg specify mandatory arguments when sending faxes.
//
// If sending to a single number use SINGLE and pass in a slice of len 1.
// Otherwise use BROADCAST and pass in a slice of numbers (as string)
type QueueCfg struct {
	CallerID    int      // sender's fax number (must be 10 digits)
	SenderEmail string   // sender's email address
	FaxType     string   // "SINGLE" or "BROADCAST"
	ToFaxNumber []string // each number must be 11 digits represented as a String
}

// File represents a queueable fax item.
// It is the callers responsibility to ensure that Content is base64-encoded.
type File struct {
	// filename
	Name string

	// base64-encoded string
	Content string
}

// QueueFaxResp represents information about faxes added to the queue.
type QueueFaxResp struct {
	Status string
	Result string
}

type mappedQueueFaxResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// QueueFax adds fax item(s) to a queue for delivery.
//
// If Files is nil, the CoverPage option must be enabled. Otherwise will receive error: No Files to Fax
func (c *Client) QueueFax(files []File, cfg QueueCfg, options ...QueueOptions) (*QueueFaxResp, error) {
	opr := map[string]interface{}{
		"action":       actionQueueFax,
		"access_id":    c.AccessID,
		"access_pwd":   c.AccessPwd,
		"sCallerID":    cfg.CallerID,
		"sSenderEmail": cfg.SenderEmail,
		"sFaxType":     cfg.FaxType,
		"sToFaxNumber": strings.Join(cfg.ToFaxNumber, "|"),
	}

	// fail early if any of the above mandatory values are empty
	// potential error if the above are empty:
	/*
		"ResultError": "Invalid Fax Type / "
		"ResultError": "Invalid Senders Email Address /"
		"ResultError": "Invalid CallerID provided / "
		"ResultError": "Forbidden: Access is denied / Invalid Authentication."
	*/
	if err := failIfEmpty(opr); err != nil {
		return nil, err
	}

	// build up optional, non-empty, options based on srfax tags through reflection.
	// TODO this may not be the best approach. Hard to test.
	// Think about writing a function to parse optional args, build a map and merge with existing operation map.
	if len(options) >= 1 {
		v := reflect.ValueOf(options[0])

		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Interface().(type) {
			case string:
				if v.Field(i).String() == "" {
					continue
				}
				s, ok := reflect.TypeOf(options[0]).Field(i).Tag.Lookup("json")
				if !ok {
					return nil, errors.Errorf("QueueFax: failed string reflection on optional arguments")
				}
				opr[s] = v.Field(i).String()
			case int:
				if v.Field(i).Int() <= 0 || v.Field(i).Int() > 6 {
					continue
				}
				s, ok := reflect.TypeOf(options[0]).Field(i).Tag.Lookup("json")
				if !ok {
					return nil, errors.Errorf("QueueFax: failed int reflection on optional arguments")
				}
				opr[s] = v.Field(i).Int()
			}
		}
	}

	const (
		prefixName    = "sFileName_"
		prefixContent = "sFileContent_"
	)

	// Don't fail if len == 0, because SRFax can queue a cover page only,
	// this is why this method accepts nil as an argument to Files. If file(s) are missing
	// name or content they get stored in emptyFiles and if slice is not zero return error.
	emptyFiles := make([]File, 0)
	if len(files) > 0 {
		for i, f := range files {
			if f.Name == "" || f.Content == "" {
				emptyFiles = append(emptyFiles, f)
			}
			opr[prefixName+strconv.Itoa(i)] = f.Name
			opr[prefixContent+strconv.Itoa(i)] = f.Content
		}
		if len(emptyFiles) > 0 {
			return nil, errors.Errorf("skipping empty file(s), check name or content: %+v", emptyFiles)
		}
	}

	operation, err := constructFromMap(opr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct a reader for queue fax")
	}

	result := mappedQueueFaxResp{}
	if err := run(operation, &result, c.url); err != nil {
		return nil, err
	}

	out := QueueFaxResp(result)
	return &out, nil
}

// failIfEmpty takes a map to check if a value is zero-value (string and int). If a value is empty
// the corresponding key is stored in slice and an error is returned.
//
// Convenience func to fail early if a mandatory field(s) empty.
func failIfEmpty(m map[string]interface{}) error {
	errs := make([]string, 0)
	for k, v := range m {
		if v == "" || v == 0 {
			errs = append(errs, k)
		}
	}
	if len(errs) != 0 {
		return errors.Errorf("check QueueCfg, the following fields cannot be empty: %s", strings.Join(errs, ", "))
	}
	return nil
}
