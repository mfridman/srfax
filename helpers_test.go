package srfax

import (
	"testing"
)

func TestIDFromName(t *testing.T) {
	var tests = []struct {
		in   string
		want int
	}{
		{"20180101230101-8812-34_0|31524120", 31524120},
		{"|31524120", 31524120},
		{"20180101230101", 20180101230101},
		{"20180101230101-8812-34_0|31524120|2222|", 0},
		{"20180101230101-8812-34_0|31524120|9999", 9999},
		{"31524120|", 0},
		{"|", 0},
		{"20180101230101-|", 0},
		{"", 0},
	}
	for _, test := range tests {
		got, _ := IDFromName(test.in)
		if test.want != got {
			t.Fatalf("IDFromName(%q) = %d; want %d", test.in, got, test.want)
		}
	}
}

func TestCheckStatus(t *testing.T) {
	type ms map[string]interface{}
	var tests = []struct {
		in   ms
		want *ResultError
	}{
		{ms{"Status": "Success", "Result": ""}, nil},
		{ms{"Status": "success", "Result": "123"}, nil},
		{ms{"Status": "success"}, &ResultError{}},
		{ms{"Status": 123, "Result": ""}, &ResultError{}},
		{ms{"Status": "", "Result": ""}, &ResultError{}},
		{ms{"Status": "", "Result": 123}, &ResultError{}},
		{ms{"Result": []string{}}, &ResultError{}},
	}
	for _, test := range tests {
		if err := checkStatus(test.in); err != nil {
			switch err.(type) {
			case *ResultError:
				break
			default:
				t.Fatalf("checkStatus(%+v) = %v; want %v", test.in, err, test.want)
			}
		}
	}
}

func TestHasKeys(t *testing.T) {
	type ms map[string]interface{}
	var tests = []struct {
		m    ms
		s    []string
		want bool
	}{
		{ms{"Status": "Failed", "Result": []string{}}, []string{"Status", "Result"}, true},
		{ms{"Status": "Success", "Result": ""}, []string{"Status", "Result"}, true},
		{ms{"Result": ""}, []string{"Result"}, true},
		{ms{}, []string{"Status", "Result"}, false},
		{ms{"": ""}, []string{"Status"}, false},
		{ms{"Result": "Failed"}, []string{}, false},
		{ms{}, []string{""}, false},
		{ms{}, []string{}, false},
		{ms{"abc": "", "def": ""}, []string{"Status"}, false},
	}
	for _, test := range tests {
		got := hasKeys(test.m, test.s)
		if test.want != got {
			t.Fatalf("hasKey(%+v, %v) = %t; want %t", test.m, test.s, got, test.want)
		}
	}
}

func TestSendPost(t *testing.T) {

}

func TestIsNChars(t *testing.T) {
	var tests = []struct {
		s    string
		l    int
		want bool
	}{
		{"hello", 5, true},
		{"hello", 4, false},
		{"", 0, true},
		{"", 1, false},
	}
	for _, test := range tests {
		got := isNChars(test.s, test.l)
		if test.want != got {
			t.Fatalf("isNChars(%q, %d) = %t; want %t", test.s, test.l, got, test.want)
		}
	}

}

func TestValidDateOrTime(t *testing.T) {
	var tests = []struct {
		layout string
		values []string
		want   bool
	}{
		{"2006-01-02", []string{"1987-02-20"}, true},
		{"20060102", []string{"19870220"}, true},
		{"15:04", []string{"10:20"}, true},
		{"2006-01-02", []string{"1987-02-20", "2017-01-15"}, true},
		{"2006-01-02", []string{}, false},
		{"20060102", []string{"198702-20"}, false},
	}
	for _, test := range tests {
		got := validDateOrTime(test.layout, test.values...)
		if test.want != got {
			t.Fatalf("validDateOrTime(%q, %v) = %t; want %t", test.layout, test.values, got, test.want)
		}
	}
}

func TestHasEmpty(t *testing.T) {
	t.Parallel()

	t.Run("checkNil", func(t *testing.T) {
		if err := hasEmpty(nil); err == nil {
			t.Fatal(err)
		}
	})

	t.Run("checkEmptyStruct", func(t *testing.T) {
		if err := hasEmpty(struct{}{}); err == nil {
			t.Fatal(err)
		}
	})

	t.Run("checkInitEmptyStruct", func(t *testing.T) {
		fwdCfg := &ForwardCfg{}
		if err := hasEmpty(*fwdCfg); err == nil {
			t.Fatal(err)
		}
	})

	t.Run("correctStruct", func(t *testing.T) {
		c := &ForwardCfg{
			FaxDetailsID: "30294755",
			FaxFileName:  "",
			Direction:    "OUT",
			CallerID:     4161112222,
			SenderEmail:  "email@example.com",
			FaxType:      "SINGLE",
		}

		if err := hasEmpty(*c); err != nil {
			t.Errorf("input: %+v", *c)
			t.Fatal("got err; want nil")
		}
	})
}
