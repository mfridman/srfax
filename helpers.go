package srfax

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
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
func sendPost(r io.Reader, url string) (map[string]interface{}, error) {

	client := http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	resp, err := client.Post(url, "application/json", r)
	if err != nil {
		return nil, errors.Wrap(err, "failed POST request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading response body from POST")
	}

	// DEBUG only, show the raw body coming across the wire.
	// fmt.Println("DEBUG RAW: ", string(b))

	var ms map[string]interface{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&ms); err != nil {
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
// If Status equals Success the function will return nil.
// The SRFax docs state that successful operations do not return error messages.
func checkStatus(ms map[string]interface{}) error {
	if ok := hasKeys(ms, []string{"Status", "Result"}); !ok {
		return &ResultError{Status: "", Raw: `missing "Status" or "Result" key in response`}
	}
	status, ok := ms["Status"].(string)
	if !ok {
		return &ResultError{Status: "", Raw: fmt.Sprintf(`failed "Status" type assertion; expecting String but got %T`, ms["Status"])}
	}
	if strings.ToLower(status) != "success" {
		result, ok := ms["Result"].(string)
		if !ok {
			return &ResultError{Status: "", Raw: fmt.Sprintf(`failed "Result" type assertion; expecting String but got %T`, ms["Status"])}
		}
		return &ResultError{Status: status, Raw: result}
	}
	return nil
}

// hasKeys iterates over a given slice and checks string existence as a key in given map.
func hasKeys(ms map[string]interface{}, ss []string) bool {
	if len(ss) == 0 {
		return false
	}
	for _, s := range ss {
		if _, ok := ms[s]; !ok {
			return false
		}
	}
	return true
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
// File name format expected "20180101230101-8812-34_0|31524120", where the
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

// decodeMap decodes a map into the underlying result type.
// It is a wrapper around Mitchell's mapstructure pkg.
func decodeMap(msi map[string]interface{}, resultType interface{}) error {
	if err := checkStatus(msi); err != nil {
		return err
	}
	var md mapstructure.Metadata
	cfg := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Metadata:         &md,
		Result:           resultType, // this MUST be a pointer to a struct
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return errors.Wrapf(err, "mapstructure new decoder config error: [%+v]", cfg)
	}
	if err := decoder.Decode(msi); err != nil {
		return errors.Wrapf(err, "mapstructure Decode error: [%+v]", msi)
	}
	// DEBUG only
	// fmt.Println("DEBUG unused keys: ", cfg.Metadata.Unused)
	return nil
}

// used to validate dates and times with time.Parse, caller must supply format layout.
func validDateOrTime(layout string, values ...string) bool {
	if len(values) == 0 {
		return false
	}
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
	if !val.IsValid() {
		return errors.Errorf("%v is not a valid value", i)
	}
	n := val.NumField()
	if n == 0 {
		return errors.Errorf("struct cannot have %d fields", n)
	}
	empty := make([]string, 0)
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

func run(r io.Reader, resultType interface{}, url string) error {
	msi, err := sendPost(r, url)
	if err != nil {
		return errors.Wrap(err, "failed sendPost")
	}
	if err := decodeMap(msi, resultType); err != nil {
		return errors.Wrap(err, "failed decodeResp")
	}
	return nil
}

func constructFromMap(i interface{}) (io.Reader, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(i); err != nil {
		return nil, errors.Wrap(err, "failed encode")
	}
	return bytes.NewReader(buf.Bytes()), nil
}

func constructFromStruct(i interface{}) (io.Reader, error) {
	by, err := json.Marshal(i)
	if err != nil {
		return nil, errors.Wrap(err, "failed json marshal")
	}
	return bytes.NewReader(by), nil
}
