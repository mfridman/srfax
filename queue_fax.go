package srfax

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// QueueFaxOpts contains optional arguments when sending faxes.
//
// If the default cover page on the account is set to "Attachments ONLY" the cover page will
// not be created irrespective of this variable
// If a cover page is not provided then all cover page variables will be ignored
//
// "Basic", "Standard", "Company", or "Personal"
//
// srfax is a custom tag used for reflection (specific to this struct).
type QueueFaxOpts struct {
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

	NotifyURL    string `srfax:"sNotifyURL"`
	QueueFaxDate string `srfax:"sQueueFaxDate"` // YYYY-MM-DD
	QueueFaxTime string `srfax:"sQueueFaxTime"` // HH:MM, using 24 hour time
}

// QueueFaxCfg contains mandatory arguments when sending faxes.
//
// If sending to a single number use SINGLE and pass in a slice of len 1.
// Otherwise use BROADCAST and pass in a slice of numbers (as string)
type QueueFaxCfg struct {
	CallerID    int      // sender's fax number (must be 10 digits)
	SenderEmail string   // sender's email address
	FaxType     string   // "SINGLE" or "BROADCAST"
	ToFaxNumber []string // each number must be 11 digits represented as a String
}

// File represents a queueable fax item.
// It is the callers responsibility to ensure that Content is base64-encoded.
// TODO think about adding a convenience function that specifies an "outbox", i.e., a directory,
// and generates all files in that directory as a []File, which can be passed directly to QueueFax
type File struct {
	Name    string // filename
	Content string // base64-encoded string
}

// QueueFaxResp represents information about faxes added to the queue.
type QueueFaxResp struct {
	Status string `mapstructure:"Status"`
	Result string `mapstructure:"Result"`
}

// QueueFax adds faxes to the queue of items to send.
//
// if files is an empty slice, the CoverPage opts must be enabled. Otherwise will receive
// error: No Files to Fax /
func (c *Client) QueueFax(files []File, q QueueFaxCfg, optArgs ...QueueFaxOpts) (map[string]interface{}, error) {
	msg := map[string]interface{}{
		"action":       actionQueueFax,
		"access_id":    c.AccessID,
		"access_pwd":   c.AccessPwd,
		"sCallerID":    q.CallerID,
		"sSenderEmail": q.SenderEmail,
		"sFaxType":     q.FaxType,
		"sToFaxNumber": strings.Join(q.ToFaxNumber, "|"),
	}

	// fail early if any of the above mandatory values are empty
	// potential error if the above are empty:
	/*
		"ResultError": "Invalid Fax Type / "
		"ResultError": "Invalid Senders Email Address /"
		"ResultError": "Invalid CallerID provided / "
		"ResultError": "Forbidden: Access is denied / Invalid Authentication."
	*/
	if err := failIfEmpty(msg); err != nil {
		return nil, err
	}

	// build up optional, non-empty, options based on srfax tags through reflection.
	// TODO this may not be the best approach. Hard to test and may be prone to error.
	// Think about writing a function to parse optional args, build a map and merging with existing msg map.
	if len(optArgs) > 0 {
		v := reflect.ValueOf(optArgs[0])

		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Interface().(type) {
			case string:
				if v.Field(i).String() == "" {
					continue
				}
				s, ok := reflect.TypeOf(optArgs[0]).Field(i).Tag.Lookup("srfax")
				if !ok {
					return nil, errors.Errorf("QueueFax: failed string reflection on optional arguments")
				}
				msg[s] = v.Field(i).String()
			case int:
				if v.Field(i).Int() <= 0 || v.Field(i).Int() > 6 {
					continue
				}
				s, ok := reflect.TypeOf(optArgs[0]).Field(i).Tag.Lookup("srfax")
				if !ok {
					return nil, errors.Errorf("QueueFax: failed int reflection on optional arguments")
				}
				msg[s] = v.Field(i).Int()
			}
		}
	}

	const (
		prefixName    = "sFileName_"
		prefixContent = "sFileContent_"
	)
	// No need to fail if len == 0, because SRFax can queue cover page only.
	if len(files) > 0 {
		for i, f := range files {
			if f.Name == "" || f.Content == "" {
				log.Printf("skipping empty file, check name or content: %+v\n", f)
				continue
			}
			msg[prefixName+strconv.Itoa(i)] = f.Name
			msg[prefixContent+strconv.Itoa(i)] = f.Content
		}
	}

	return msg, nil
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
		return errors.Errorf("check QueueFaxCfg, the following cannot be empty: %s", strings.Join(em, ","))
	}

	return nil
}
