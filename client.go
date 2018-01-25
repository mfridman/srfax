package srfax

import "github.com/pkg/errors"

// ClientCfg specifies parameters required for establishing an SRFax client.
type ClientCfg struct {
	ID  int
	PWD string
	URL string // Optional
}

// NewClient returns an SRFax client based on configuration.
// If URL is unspecified it will default to "https://www.srfax.com/SRF_SecWebSvc.php"
func NewClient(c ClientCfg) (*Client, error) {
	if c.ID == 0 {
		return nil, errors.New("must specify access id. User Number")
	}
	if c.PWD == "" {
		return nil, errors.New("must specify access pwd. Password on the users account")
	}
	u := c.URL
	if c.URL == "" {
		// enable end-user to override default URL in case endpoint changes
		u = "https://www.srfax.com/SRF_SecWebSvc.php"
	}
	return &Client{AccessID: c.ID, AccessPwd: c.PWD, url: u}, nil
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
	if !(s.Status == "Success" && s.ResultError == "") {
		return false, errors.Errorf("Authentication [%s]: [%s]", s.Status, s.ResultError)
	}
	return true, nil
}
