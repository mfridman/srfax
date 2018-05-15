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
	for _, tt := range tests {
		got, _ := IDFromName(tt.in)
		if tt.want != got {
			t.Fatalf("IDFromName(%s): want %d, got %d", tt.in, tt.want, got)
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
	for _, tt := range tests {
		err := checkStatus(tt.in)
		if err != nil {
			switch err.(type) {
			case *ResultError:
				break
			default:
				t.Fatalf("checkStatus(%v): want [%v], got [%v]", tt.in, tt.want, err)
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
	for _, tt := range tests {
		got := hasKeys(tt.m, tt.s)
		if tt.want != got {
			t.Fatalf("hasKey(%v, %v): want %t, got %t", tt.m, tt.s, tt.want, got)
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
	for _, tt := range tests {
		got := isNChars(tt.s, tt.l)
		if tt.want != got {
			t.Fatalf("isNChars(%v, %v): want %t, got %t", tt.s, tt.l, tt.want, got)
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
	for _, tt := range tests {
		got := validDateOrTime(tt.layout, tt.values...)
		if tt.want != got {
			t.Fatalf("validDateOrTime(%v, %v): want %t, got %t", tt.layout, tt.values, tt.want, got)
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
			t.Fatalf("hasEmpty(%+v): want nil, got error: %v", c, err)
		}
	})
}
