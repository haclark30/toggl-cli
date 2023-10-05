package toggl

import (
	"encoding/json"
	"time"
)

type Me struct {
	ID                 int       `json:"id,omitempty"`
	APIToken           string    `json:"api_token,omitempty"`
	Email              string    `json:"email,omitempty"`
	Fullname           string    `json:"fullname,omitempty"`
	Timezone           string    `json:"timezone,omitempty"`
	TogglAccountsID    string    `json:"toggl_accounts_id,omitempty"`
	DefaultWorkspaceID int       `json:"default_workspace_id,omitempty"`
	BeginningOfWeek    int       `json:"beginning_of_week,omitempty"`
	ImageURL           string    `json:"image_url,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
	OpenIDEmail        string    `json:"openid_email,omitempty"`
	OpenIDEnabled      bool      `json:"openid_enabled,omitempty"`
	CountryID          int       `json:"country_id,omitempty"`
	HasPassword        bool      `json:"has_password,omitempty"`
	At                 time.Time `json:"at,omitempty"`
	IntercomHash       string    `json:"intercom_hash,omitempty"`
	OAuthProviders     []string  `json:"oauth_providers,omitempty"`
}

func (c TogglClient) Me() (Me, error) {
	var me Me
	resp, err := c.httpGet(mePath)
	defer resp.Body.Close()
	if err != nil {
		return me, err
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&me)
	return me, nil

}
