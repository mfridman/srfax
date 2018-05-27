package srfax

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// ClientCfg specifies parameters required for establishing an SRFax client.
// Both ID and Pwd are unique to an SRFax account.
type ClientCfg struct {
	// access_id
	ID int

	// access_pwd
	Pwd string
}

func (cfg ClientCfg) validate() error {
	if cfg.ID <= 0 {
		return errors.New("access id (ID) must be a positive number")
	}
	if cfg.Pwd == "" {
		return errors.New("password (Pwd) cannot be blank")
	}
	return nil
}

// NewClient returns an SRFax client based on supplied configuration.
func NewClient(cfg ClientCfg) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	// apiURL is the SRFax API url.
	apiURL := "https://www.srfax.com/SRF_SecWebSvc.php"

	return &Client{account{AccessID: cfg.ID, AccessPwd: cfg.Pwd}, apiURL}, nil
}

// Client is an SRFax client.
type Client struct {
	account
	url string
}

type account struct {
	AccessID  int    `json:"access_id"`
	AccessPwd string `json:"access_pwd"`
}

// CheckAuth verifies client's ability to authenticate with SRFax. It is a wrapper around the
// GetFaxUsage method. Convenience method to quickly check if access ID & Pwd are valid.
func (c *Client) CheckAuth() (bool, error) {
	if _, err := c.GetFaxUsage(); err != nil {
		return false, err
	}
	return true, nil
}

// Usage reports the total pages used and the date range.
type Usage struct {
	TotalPages int
	AccessID   int
	StartDate  time.Time
	EndDate    time.Time
}

func (u *Usage) String() string {
	return fmt.Sprintf("Account %d used %d pages from %v to %v", u.AccessID, u.TotalPages, u.StartDate.Format("Jan 02 2006"), u.EndDate.Format("Jan 02 2006"))
}

// UsageCounter reports the number of pages used by ALL users of the account in
// the current period based on the account's reset day. Each account will have a unique reset day.
// To find the reset day navigate to My Account > Summary > look for "fax usage counter will reset on March 19, 2018",
// for this example the reset day would be 19.
func (c *Client) UsageCounter(day int) (*Usage, error) {

	var start, end time.Time
	layout := "20060102"
	now := time.Now()

	// current date before the reset day for the current month
	reset := time.Date(now.Year(), time.Month(int(now.Month())), day, 0, 0, 0, 0, time.Local)
	switch now.Before(reset) {
	case true:
		// in the current month, BEFORE the reset date
		end = reset
		start = end.AddDate(0, -1, 0) // go back one month from end date
	case false:
		// in the current month, AFTER the reset date
		end = time.Date(now.Year(), time.Month(int(now.Month()+1)), day, 0, 0, 0, 0, time.Local)
		start = end.AddDate(0, -1, 0) // go back one month from end date
	}

	opts := FaxUsageOptions{
		IncludeSubUsers: yes,
		Period:          "RANGE",
		StartDate:       start.Format(layout),
		EndDate:         end.Format(layout),
	}

	resp, err := c.GetFaxUsage(opts)
	if err != nil {
		return nil, err
	}

	u := Usage{
		AccessID:  c.AccessID,
		StartDate: start,
		EndDate:   end,
	}

	for i := range resp.Result {
		u.TotalPages += resp.Result[i].NumberOfPages
	}

	return &u, nil
}
