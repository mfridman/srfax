package srfax

import (
	"strings"

	"github.com/pkg/errors"
)

// ClientCfg specifies parameters required for establishing an SRFax client.
type ClientCfg struct {
	ID  int    // access_id
	Pwd string // access_pwd
}

const (
	// SRFax API url.
	url = "https://www.srfax.com/SRF_SecWebSvc.php"
)

// NewClient returns an SRFax client based on configuration.
func NewClient(c ClientCfg) (*Client, error) {
	if c.ID == 0 {
		return nil, errors.New("must specify access id. User Number")
	}
	if c.Pwd == "" {
		return nil, errors.New("must specify access pwd. Password on the users account")
	}
	return &Client{AccessID: c.ID, AccessPwd: c.Pwd, url: url}, nil
}

// Client is an SRFax client that contains authentication information
type Client struct {
	AccessID  int    `json:"access_id"`
	AccessPwd string `json:"access_pwd"`
	url       string
}

// CheckAuth checks whether client is able to authenticate. This is a wrapper around the
// GetFaxUsage method. Used for convenience to quickly check if access ID & PWD are valid.
func (c *Client) CheckAuth() (bool, error) {
	s, err := c.GetFaxUsage()
	if err != nil {
		return false, err
	}
	if strings.ToLower(s.Status) != "success" {
		return false, errors.Errorf("No Result errors but Status was not Success: [%v]", s.Status)
	}
	return true, nil
}
