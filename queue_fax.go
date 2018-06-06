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
// not be created irrespective of the CoverPage variable.
// If CoverPage is not provided all cover page variables are ignored
// For more info see docs: https://www.srfax.com/api-page/queue_fax/
//
// json tag used for reflection.
type QueueOptions struct {
	// Number of times the system is to retry a number if busy or
	// an error is encountered. Must be a number from 0 to 6
	Retries int `json:"sRetries"`

	// Internal Reference Number (Maximum of 20 characters)
	AccountCode string `json:"sAccountCode"`

	// From: On the Fax Header Line(Maximum of 30 characters)
	FaxFromHeader string `json:"sFaxFromHeader"`

	// COVER PAGE OPTIONS
	//
	// To use one of the cover pages on file, specify the cover page you wish to use:
	// Basic, Standard, Company, or Personal
	CoverPage string `json:"sCoverPage"`
	// Sender name on the Cover Page
	CPFromName string `json:"sCPFromName"`
	// Recipient name on the Cover Page
	CPToName string `json:"sCPToName"`
	// Organiation on the Cover Page
	CPOrganization string `json:"sCPOrganization"`
	// Subject Line on the Cover Page
	// The subject line details are saved in the fax record even is a cover page is
	// not requested – so the subject can be used for filtering / searching
	CPSubject string `json:"sCPSubject"`
	// Comments placed in the body of the Cover Page
	CPComments string `json:"sCPComments"`

	// Provide an absolute URL (prefixed with http:// or https://) and the SRFax
	// system will POST back the fax status record when the fax completes.
	NotifyURL string `json:"sNotifyURL"`

	// The date you want to schedule a future fax for.
	// Must be in the format YYYY-MM-DD. Required if using QueueFaxTime
	QueueFaxDate string `json:"sQueueFaxDate"`
	// The time you want to schedule a future fax for. Must be in the format HH:MM,
	// using 24 hour time (ie, 00:00 – 23:59). Required if using QueueFaxDate.
	// The timezone set on the account will be used when scheduling.
	QueueFaxTime string `json:"sQueueFaxTime"`
}

// QueueCfg specify mandatory arguments when sending faxes.
//
// If sending to a single number use SINGLE and pass in a slice of len 1.
// Otherwise use BROADCAST and pass in a slice of numbers (as string)
type QueueCfg struct {
	// Sender fax number (must be 10 digits)
	CallerID int

	// Sender email address
	SenderEmail string

	// SINGLE when sending to one number; BROADCAST when sending to multiple numbers
	FaxType string

	// Slice of string representing an 11 digit fax number
	ToFaxNumber []string
}

// File represents a queueable fax item.
// It is the callers responsibility to ensure that Content is base64-encoded.
// Check the FAQs to see a list of supported file types: https://www.srfax.com/api-page/queue_fax/
type File struct {
	// Valid File Name
	Name string

	// Base64 encoding of file contents
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
	// Think about writing a function to parse optional args, build a map and merge with existing opr map from above.
	if len(options) > 0 {
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

	operation, err := constructReader(opr)
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
