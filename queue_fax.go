package srfax

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// QueueOptions contains optional arguments when sending faxes.
//
// If the default cover page on the account is set to "Attachments ONLY" the cover page will
// not be created irrespective of this variable
// If a cover page is not provided then all cover page variables will be ignored
//
// "Basic", "Standard", "Company", or "Personal"
//
// srfax is a custom tag used for reflection (specific to this struct).
type QueueOptions struct {
	Retries       int    `srfax:"sRetries"`
	AccountCode   string `srfax:"sAccountCode"`
	FaxFromHeader string `srfax:"sFaxFromHeader"`

	// cover page
	CoverPage      string `srfax:"sCoverPage"`
	CPFromName     string `srfax:"sCPFromName"`
	CPToName       string `srfax:"sCPToName"`
	CPOrganization string `srfax:"sCPOrganization"`
	CPSubject      string `srfax:"sCPSubject"`
	CPComments     string `srfax:"sCPComments"`

	NotifyURL string `srfax:"sNotifyURL"`

	// YYYY-MM-DD
	QueueFaxDate string `srfax:"sQueueFaxDate"`
	// HH:MM, using 24 hour time
	QueueFaxTime string `srfax:"sQueueFaxTime"`
}

// QueueCfg contains mandatory arguments when sending faxes.
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

// Files is a slice of quequeable File items.
type Files []File

// QueueFaxResp represents information about faxes added to the queue.
type QueueFaxResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// QueueFax adds fax item(s) to a queue for delivery.
//
// If Files is nil, the CoverPage option must be enabled. Otherwise will receive error: No Files to Fax
func (c *Client) QueueFax(files Files, cgf QueueCfg, options ...QueueOptions) (*QueueFaxResp, error) {
	req := map[string]interface{}{
		"action":       actionQueueFax,
		"access_id":    c.AccessID,
		"access_pwd":   c.AccessPwd,
		"sCallerID":    cgf.CallerID,
		"sSenderEmail": cgf.SenderEmail,
		"sFaxType":     cgf.FaxType,
		"sToFaxNumber": strings.Join(cgf.ToFaxNumber, "|"),
	}

	// fail early if any of the above mandatory values are empty
	// potential error if the above are empty:
	/*
		"ResultError": "Invalid Fax Type / "
		"ResultError": "Invalid Senders Email Address /"
		"ResultError": "Invalid CallerID provided / "
		"ResultError": "Forbidden: Access is denied / Invalid Authentication."
	*/
	if err := failIfEmpty(req); err != nil {
		return nil, err
	}

	// build up optional, non-empty, options based on srfax tags through reflection.
	// TODO this may not be the best approach. Hard to test and may be prone to error.
	// Think about writing a function to parse optional args, build a map and merge with existing request map.
	if len(options) > 0 {
		v := reflect.ValueOf(options[0])

		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Interface().(type) {
			case string:
				if v.Field(i).String() == "" {
					continue
				}
				s, ok := reflect.TypeOf(options[0]).Field(i).Tag.Lookup("srfax")
				if !ok {
					return nil, errors.Errorf("QueueFax: failed string reflection on optional arguments")
				}
				req[s] = v.Field(i).String()
			case int:
				if v.Field(i).Int() <= 0 || v.Field(i).Int() > 6 {
					continue
				}
				s, ok := reflect.TypeOf(options[0]).Field(i).Tag.Lookup("srfax")
				if !ok {
					return nil, errors.Errorf("QueueFax: failed int reflection on optional arguments")
				}
				req[s] = v.Field(i).Int()
			}
		}
	}

	const (
		prefixName    = "sFileName_"
		prefixContent = "sFileContent_"
	)

	emptyFiles := make([]File, 0)
	// Don't fail if len == 0, because SRFax can queue a cover page only,
	// this is why this method accepts nil as an argument to Files.
	if len(files) > 0 {
		for i, f := range files {
			if f.Name == "" || f.Content == "" {
				emptyFiles = append(emptyFiles, f)
			}
			req[prefixName+strconv.Itoa(i)] = f.Name
			req[prefixContent+strconv.Itoa(i)] = f.Content
		}

		if len(emptyFiles) > 0 {
			return nil, errors.Errorf("skipping empty file(s), check name or content: %+v", emptyFiles)
		}
	}

	var resp QueueFaxResp
	if err := run(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// failIfEmpty takes a map to check if any value is empty. If a value is empty
// the corresponding key is stored in slice and error is returned.
//
// Convenience func to fail early in the event a mandatory config field is missing.
func failIfEmpty(m map[string]interface{}) error {
	em := make([]string, 0)

	for k, v := range m {
		if v == "" {
			em = append(em, k)
		}
	}
	if len(em) != 0 {
		return errors.Errorf("check QueueCfg, the following fields cannot be empty: %s", strings.Join(em, ", "))
	}

	return nil
}
