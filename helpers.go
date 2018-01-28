package srfax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const (
	// SRFax specific action verbs. Every POST request will use one of the following:
	actionQueueFax           = "Queue_Fax"
	actionGetFaxStatus       = "Get_FaxStatus"
	actionGetMulFaxStatus    = "Get_MultiFaxStatus"
	actionGetFaxInbox        = "Get_Fax_Inbox"
	actionGetFaxOutbox       = "Get_Fax_Outbox"
	actionForwardFax         = "Forward_Fax"
	actionRetrieveFax        = "Retrieve_Fax"
	actionUpdateViewedStatus = "Update_Viewed_Status"
	actionDeleteFax          = "Delete_Fax"
	actionStopFax            = "Stop_Fax"
	actionGetFaxUsage        = "Get_Fax_Usage"
)

// decodeResp is a wrapper around mitchell's mapstructure pkg. Mainly used for debugging
// parts of the API as the docs don't always line up with what comes across the wire.
func decodeResp(resp map[string]interface{}, cfg *mapstructure.DecoderConfig) error {
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

// sendPost is a wrapper around Post. Sends JSON encoded string to SRFax and attempts
// to decode response.
func sendPost(msg interface{}, url string) (map[string]interface{}, error) {

	// write out msg to bytes buffer
	r := new(bytes.Buffer)
	err := json.NewEncoder(r).Encode(msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed writing out to bytes buffer")
	}

	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Post(url, "application/json", r)
	if err != nil {
		return nil, errors.Wrap(err, "posting message")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %v", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body from POST")
	}

	// DEBUG only
	// fmt.Println("DEBUG RAW: ", string(b))

	var v map[string]interface{}
	err = json.NewDecoder(bytes.NewReader(b)).Decode(&v)
	if err != nil {
		return nil, errors.Wrap(err, "failed decoding response from POST")
	}

	return v, nil
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

// DEBUG ONLY, not exported and should not be used.
func whereami() {
	var outF string
	pc, _, _, ok := runtime.Caller(0)
	if !ok {
		outF = "unnamed"
	}
	me := runtime.FuncForPC(pc)
	if me == nil {
		outF = "unnamed"
	} else {
		outF = me.Name()
	}
	log.Println(outF)
}

// PP is a convenience function to pretty print JSON.
func PP(i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", string(b))
}

// DEBUG ONLY, not exported and should not be used.
func nameOf(f interface{}) string {
	v := reflect.ValueOf(f)
	if v.Kind() == reflect.Func {
		if rf := runtime.FuncForPC(v.Pointer()); rf != nil {
			return rf.Name()
		}
	}
	return v.String()
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
