package toggl

import (
	"encoding/json"
	"fmt"
	"time"
)

const getProjectPath = "api/v9/workspaces/%d/projects/%d"
const getProjectsPath = "api/v9/workspaces/%d/projects"

type Project struct {
	ID                  int           `json:"id,omitempty"`
	WorkspaceID         int           `json:"workspace_id,omitempty"`
	ClientID            *int          `json:"client_id,omitempty"`
	Name                string        `json:"name,omitempty"`
	IsPrivate           bool          `json:"is_private,omitempty"`
	Active              bool          `json:"active,omitempty"`
	At                  time.Time     `json:"at,omitempty"`
	CreatedAt           time.Time     `json:"created_at,omitempty"`
	ServerDeletedAt     *time.Time    `json:"server_deleted_at,omitempty"`
	Color               string        `json:"color,omitempty"`
	Billable            bool          `json:"billable,omitempty"`
	Template            bool          `json:"template,omitempty"`
	AutoEstimates       bool          `json:"auto_estimates,omitempty"`
	EstimatedHours      *float64      `json:"estimated_hours,omitempty"`
	Rate                *float64      `json:"rate,omitempty"`
	RateLastUpdated     *time.Time    `json:"rate_last_updated,omitempty"`
	Currency            string        `json:"currency,omitempty"`
	Recurring           bool          `json:"recurring,omitempty"`
	RecurringParameters []interface{} `json:"recurring_parameters,omitempty"`
	CurrentPeriod       *time.Time    `json:"current_period,omitempty"`
	FixedFee            *float64      `json:"fixed_fee,omitempty"`
	ActualHours         float64       `json:"actual_hours,omitempty"`
	WID                 int           `json:"wid,omitempty"`
	CID                 *int          `json:"cid,omitempty"`
}

func (c TogglClient) GetProject(workspaceId int, projectId int) (Project, error) {
	var project Project
	projPath := fmt.Sprintf(getProjectPath, workspaceId, projectId)
	resp, err := c.httpGet(projPath)
	defer resp.Body.Close()
	if err != nil {
		return project, err
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&project); err != nil {
		return project, err
	}
	return project, nil
}

func (c TogglClient) GetProjects(workspaceId int) ([]Project, error) {
	var projects []Project
	projPath := fmt.Sprintf(getProjectsPath, workspaceId)
	resp, err := c.httpGet(projPath)
	defer resp.Body.Close()
	if err != nil {
		return projects, err
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&projects); err != nil {
		return projects, err
	}

	return projects, nil
}
