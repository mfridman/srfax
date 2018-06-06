package srfax

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewOutboxOperation(t *testing.T) {

	t.Parallel()

	t.Run("valid with empty options", func(t *testing.T) {
		test := struct {
			c    *Client
			o    *OutboxOptions
			want map[string]interface{}
		}{
			&Client{account{925, "abc"}, ""}, &OutboxOptions{}, map[string]interface{}{
				"action": "Get_Fax_Outbox", "access_id": 925, "access_pwd": "abc"},
		}

		got, err := constructReader(newOutboxOperation(test.c, test.o))
		if err != nil {
			t.Fatal(err)
		}

		var ms map[string]interface{}
		if err := json.NewDecoder(got).Decode(&ms); err != nil {
			t.Fatal(err)
		}

		if test.want["action"] != ms["action"].(string) {
			t.Errorf("want %q; got %q", test.want["action"], ms["action"])
		}

		if test.want["access_pwd"] != ms["access_pwd"].(string) {
			t.Errorf("want %q; got %q", test.want["access_pwd"], ms["access_pwd"])
		}

		if test.want["access_id"] != int(ms["access_id"].(float64)) {
			t.Errorf("want %[1]v with type %[1]T; got %[2]v with type %[2]T", test.want["access_id"], ms["access_id"])
		}

		// test absence of empty options, i.e., when no options supplied the JSON-encoded message should not contain options
		// sViewedStatus is a valid key but not present in the decoded JSON, good!
		if _, ok := ms["sViewedStatus"]; ok {
			t.Errorf("want nil; got %v", ms["sViewedStatus"])
		}

	})

	t.Run("valid with options", func(t *testing.T) {
		test := struct {
			c    *Client
			o    *OutboxOptions
			want map[string]interface{}
		}{
			&Client{account{925, "abc"}, ""}, &OutboxOptions{Period: "ALL"}, map[string]interface{}{
				"action": "Get_Fax_Outbox", "access_id": 925, "access_pwd": "abc", "sPeriod": "ALL"},
		}

		got, err := constructReader(newOutboxOperation(test.c, test.o))
		if err != nil {
			t.Fatal(err)
		}

		var ms map[string]interface{}
		if err := json.NewDecoder(got).Decode(&ms); err != nil {
			t.Fatal(err)
		}

		if test.want["action"] != ms["action"].(string) {
			t.Errorf("want %q; got %q", test.want["action"], ms["action"])
		}

		if test.want["sPeriod"] != ms["sPeriod"].(string) {
			t.Errorf("want sPeriod=%q; got sPeriod=%q", test.want["sPeriod"], ms["sPeriod"])
		}
	})

}

func TestNewOutboxOptions(t *testing.T) {

	t.Parallel()

	t.Run("test nil", func(t *testing.T) {

		got, err := newOutboxOptions()
		if err != nil {
			t.Error("should not get an error when no options supplied")
		}

		var want OutboxOptions
		if want != *got {
			t.Error("want an empty struct")
		}
	})

	t.Run("test multiple empty options", func(t *testing.T) {

		got, err := newOutboxOptions([]OutboxOptions{OutboxOptions{}, OutboxOptions{}}...)
		if err != nil {
			t.Fatal("should not get an error when multiple options supplied")
		}

		var want OutboxOptions
		if want != *got {
			t.Error("want an empty struct")
		}
	})

	t.Run("test valid IncludeSubUsers options", func(t *testing.T) {

		option := "sIncludeSubUsers"
		tests := []struct {
			in   OutboxOptions
			want map[string]interface{}
		}{
			{OutboxOptions{IncludeSubUsers: "Y"}, map[string]interface{}{option: "Y"}},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
				got, err := newOutboxOptions(test.in)
				if err != nil {
					t.Fatal(err)
				}
				if test.want[option] != got.IncludeSubUsers {
					t.Errorf("want %q; got %q", test.want[option], got.IncludeSubUsers)
				}
			})
		}
	})

	t.Run("test invalid options", func(t *testing.T) {
		tests := []struct {
			in OutboxOptions
		}{
			{OutboxOptions{Period: "OOPS"}},
			{OutboxOptions{IncludeSubUsers: "OOPS"}},
			{OutboxOptions{Period: "RANGE"}},
			{OutboxOptions{Period: "RANGE", EndDate: "20180202"}},
			{OutboxOptions{Period: "RANGE", StartDate: "20180101"}},
			{OutboxOptions{Period: "ALL", StartDate: "20180101", EndDate: "20180202"}},
			{OutboxOptions{Period: "ALL", EndDate: "20180202"}},
			{OutboxOptions{Period: "ALL", StartDate: "20180101"}},
			{OutboxOptions{Period: "ALL", StartDate: "2018-01-01"}},
			{OutboxOptions{Period: "RANGE", StartDate: "2018-01-01"}},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
				_, err := newOutboxOptions(test.in)
				if err == nil {
					t.Error("check option validity")
					t.Fatalf("expecting to receieve an error when passing invalid option(s): %+v; got %v", test.in, err)
				}
			})
		}
	})

	t.Run("test valid Period options", func(t *testing.T) {
		option := "sPeriod"
		tests := []struct {
			in   OutboxOptions
			want map[string]interface{}
		}{
			{OutboxOptions{Period: "RANGE", StartDate: "20180101", EndDate: "20180201"}, map[string]interface{}{option: "RANGE"}},
			{OutboxOptions{Period: "ALL"}, map[string]interface{}{option: "ALL"}},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
				got, err := newOutboxOptions(test.in)
				if err != nil {
					t.Fatal(err)
				}
				if test.want[option] != got.Period {
					t.Errorf("want %q; got %q", test.want[option], got.Period)
				}
			})
		}
	})

}

func TestGetFaxOutbox(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		type msi map[string]interface{}
		data := struct {
			Status string `json:"Status"`
			Result []msi  `json:"Result"`
		}{
			Status: "Success",
			Result: []msi{msi{"SentStatus": "Sent"}},
		}
		if err := json.NewEncoder(w).Encode(data); err != nil {
			t.Fatal(err)
		}
	}))
	defer srv.Close()

	t.Run("valid response, no options", func(t *testing.T) {

		// no need to call NewClient, because we want to pass a mock URL
		client := Client{account{9090, "abc"}, srv.URL}

		outbox, err := client.GetFaxOutbox()
		if err != nil {
			t.Fatal(err)
		}

		if outbox.Status != "Success" {
			t.Fatal("expecting Status to be Success")
		}

		if outbox.Result[0].SentStatus != "Sent" {
			t.Fatal("expecting SentStatus to be Sent")
		}
	})

}

func TestGetFaxOutboxInvalidOptions(t *testing.T) {

	t.Parallel()

	client := &Client{}

	tests := []struct {
		in OutboxOptions
	}{
		{OutboxOptions{Period: "OOPS"}},
		{OutboxOptions{IncludeSubUsers: "No thanks!"}},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			_, err := client.GetFaxOutbox(test.in)
			if err == nil {
				t.Fatalf("expecting to receieve an error when passing invalid option(s): %+v; got %v", test.in, err)
			}
		})
	}

}
