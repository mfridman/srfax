package srfax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// ResultError represents an error when Result returns Failed.
// Caller can access the Status and Raw (Result error message) fields.
type ResultError struct {
	Status string // Status value from Failed Response
	Raw    string // Unformatted Result error message from Failed Response
}

func (r *ResultError) Error() string { return fmt.Sprintf("%v: %v", r.Status, r.Raw) }

// sendPost is a wrapper around http.Post method.
// Sends a JSON encoded request to SRFax and decodes the response body.
func sendPost(req interface{}) (map[string]interface{}, error) {

	client := http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	by, err := json.Marshal(&req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request")
	}

	resp, err := client.Post(apiURL, "application/json", bytes.NewReader(by))
	if err != nil {
		return nil, errors.Wrap(err, "error with POST request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body from POST")
	}

	// DEBUG only, show the raw body coming across the wire.
	// fmt.Println("DEBUG RAW: ", string(b))

	var ms map[string]interface{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&ms)
	if err != nil {
		return nil, errors.Wrap(err, "failed decoding response from POST")
	}

	return ms, nil
}

/*
	The API returns a "Result" that, depending on action:
		on success .. []interface{} or map[string]interface{} or []map[string]interface{}
		on failure .. string, where string is an error message. Not a great design, but here we are.
*/

// checkStatus checks existence for "Status" and "Result" keys in map and returns
// the status and an error message from Result.
//
// If Status equals Success the function will return Success and nil error
// This is going on the assumption (and doc examples) that successful operations do not return
// error messages.
func checkStatus(ms map[string]interface{}) (string, error) {
	if _, ok := ms["Status"]; !ok {
		return "", errors.New(`missing "Status" key in response`)
	}
	status, ok := ms["Status"].(string)
	if !ok {
		return "", errors.Errorf(`failed "Status" type assertion; expecting String but got %T`, ms["Status"])
	}
	if !(strings.ToLower(status) == "success") {
		if _, ok := ms["Result"]; !ok {
			return "", errors.New(`missing "Result" key in response`)
		}
		result, ok := ms["Result"].(string)
		if !ok {
			return "", errors.Errorf(`failed "Result" type assertion; expecting String but got %T`, ms["Result"])
		}
		return status, errors.New(result)
	}
	return status, nil
}

// PP is a convenience function to pretty print JSON.
func PP(i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}

// IDFromName parses a filename string and returns the ID.
// File name format is usually "20180101230101-8812-34_0|31524120", where the
// ID follows the pipe symbol.
func IDFromName(s string) (int, error) {
	ss := strings.SplitAfter(s, "|")
	last := ss[len(ss)-1]
	n, err := strconv.Atoi(last)
	if err != nil {
		return 0, errors.Wrap(err, "could not get ID from filename")
	}
	return n, nil
}

// decodeResp decodes response map into the underlying response type (rt).
// It is a wrapper around Mitchell's mapstructure pkg.
func decodeResp(resp map[string]interface{}, rt interface{}) error {
	if st, err := checkStatus(resp); err != nil {
		return &ResultError{Status: st, Raw: fmt.Sprint(err)}
	}

	var md mapstructure.Metadata
	cfg := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Metadata:         &md,
		Result:           rt,
	}

	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return errors.Wrapf(err, "mapstructure new decoder config error: [%+v]", cfg)
	}
	if err := decoder.Decode(resp); err != nil {
		return errors.Wrapf(err, "mapstructure Decode error: [%+v]", resp)
	}

	// DEBUG only
	// fmt.Println("DEBUG unused keys: ", cfg.Metadata.Unused)

	return nil
}

// used to validate dates and times, caller must supply format layout.
func validDateOrTime(layout string, values ...string) bool {
	for _, val := range values {
		if _, err := time.Parse(layout, val); err != nil {
			return false
		}
	}

	return true
}

// Check a struct to make sure fields are not set to their zero value.
// Supports string and int checking on non-embedded struct. Will skip fields with omitempty tag.
func hasEmpty(i interface{}) error {

	val := reflect.ValueOf(i)

	n := val.NumField()

	var empty []string

	for i := 0; i < n; i++ {
		switch val.Field(i).Kind() {
		case reflect.String:
			tags := val.Type().Field(i).Tag.Get("json")
			if !strings.Contains(tags, "omitempty") && val.Field(i).String() == "" {
				empty = append(empty, val.Type().Field(i).Name)
			}
		case reflect.Int:
			tags := val.Type().Field(i).Tag.Get("json")
			if !strings.Contains(tags, "omitempty") && val.Field(i).Int() == 0 {
				empty = append(empty, val.Type().Field(i).Name)
			}
		}
	}
	if len(empty) > 0 {
		s := fmt.Sprintf("the following fields cannot be empty: %v", strings.Join(empty, ", "))
		return errors.New(s)
	}

	return nil
}

func isNChars(s string, length int) bool {
	if len(s) != length {
		return false
	}
	return true
}

func run(req, resp interface{}) error {
	msg, err := sendPost(req)
	if err != nil {
		return errors.Wrap(err, "SendPost error")
	}

	if err := decodeResp(msg, resp); err != nil {
		return err
	}

	return nil
}
