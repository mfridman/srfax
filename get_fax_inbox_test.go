package srfax

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewInboxOperation(t *testing.T) {

	t.Parallel()

	t.Run("valid with empty options", func(t *testing.T) {
		test := struct {
			c    *Client
			o    *InboxOptions
			want map[string]interface{}
		}{
			&Client{account{925, "abc"}, ""}, &InboxOptions{}, map[string]interface{}{
				"action": "Get_Fax_Inbox", "access_id": 925, "access_pwd": "abc"},
		}

		got, err := constructFromStruct(newInboxOperation(test.c, test.o))
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
		if _, ok := ms["sViewedStatus"]; ok {
			t.Errorf("want nil; got %v", ms["sViewedStatus"])
		}

	})

	t.Run("valid with options", func(t *testing.T) {
		test := struct {
			c    *Client
			o    *InboxOptions
			want map[string]interface{}
		}{
			&Client{account{925, "abc"}, ""}, &InboxOptions{ViewedStatus: "Y"}, map[string]interface{}{
				"action": "Get_Fax_Inbox", "access_id": 925, "access_pwd": "abc", "sViewedStatus": "Y"},
		}

		got, err := constructFromStruct(newInboxOperation(test.c, test.o))
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

		if test.want["sViewedStatus"] != ms["sViewedStatus"].(string) {
			t.Errorf("want sViewedStatus=%q; got sViewedStatus=%q", test.want["sViewedStatus"], ms["sViewedStatus"])
		}
	})

}

func TestNewInboxOptions(t *testing.T) {

	t.Parallel()

	t.Run("test nil", func(t *testing.T) {
		got, err := newInboxOptions()
		if err != nil {
			t.Error("should not get an error when no options supplied")
		}

		var want InboxOptions
		if want != *got {
			t.Error("want an empty struct")
		}
	})

	t.Run("test multiple empty options", func(t *testing.T) {
		got, err := newInboxOptions([]InboxOptions{InboxOptions{}, InboxOptions{}}...)
		if err != nil {
			t.Fatal("should not get an error when multiple options supplied")
		}

		var want InboxOptions
		if want != *got {
			t.Error("want an empty struct")
		}
	})

	t.Run("test valid ViewedStatus options", func(t *testing.T) {
		option := "sViewedStatus"
		tests := []struct {
			in   InboxOptions
			want map[string]interface{}
		}{
			{InboxOptions{ViewedStatus: "ALL"}, map[string]interface{}{option: "ALL"}},
			{InboxOptions{ViewedStatus: "UNREAD"}, map[string]interface{}{option: "UNREAD"}},
			{InboxOptions{ViewedStatus: "READ"}, map[string]interface{}{option: "READ"}},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
				got, err := newInboxOptions(test.in)
				if err != nil {
					t.Fatal(err)
				}
				if test.want[option] != got.ViewedStatus {
					t.Errorf("want %q; got %q", test.want[option], got.ViewedStatus)
				}
			})
		}
	})

	t.Run("test invalid options", func(t *testing.T) {
		tests := []struct {
			in InboxOptions
		}{
			{InboxOptions{ViewedStatus: "OOPS"}},
			{InboxOptions{Period: "OOPS"}},
			{InboxOptions{IncludeSubUsers: "OOPS"}},
			{InboxOptions{Period: "RANGE"}},
			{InboxOptions{Period: "RANGE", EndDate: "20180202"}},
			{InboxOptions{Period: "RANGE", StartDate: "20180101"}},
			{InboxOptions{Period: "ALL", StartDate: "20180101", EndDate: "20180202"}},
			{InboxOptions{Period: "ALL", EndDate: "20180202"}},
			{InboxOptions{Period: "ALL", StartDate: "20180101"}},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
				_, err := newInboxOptions(test.in)
				if err == nil {
					t.Error("check option validity")
					t.Fatalf("expecting an error when passing invalid option: %+v; got %v", test.in, err)
				}
			})
		}
	})

	t.Run("test valid Period options", func(t *testing.T) {
		option := "sPeriod"
		tests := []struct {
			in   InboxOptions
			want map[string]interface{}
		}{
			{InboxOptions{Period: "RANGE", StartDate: "20180101", EndDate: "20180201"}, map[string]interface{}{option: "RANGE"}},
			{InboxOptions{Period: "ALL"}, map[string]interface{}{option: "ALL"}},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
				got, err := newInboxOptions(test.in)
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
