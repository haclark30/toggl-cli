package toggl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const currentTimePath = "api/v9/me/time_entries/current"

type TimeEntry struct {
	ID              int        `json:"id,omitempty"`
	WorkspaceID     int        `json:"workspace_id,omitempty"`
	ProjectID       *int       `json:"project_id,omitempty"`
	TaskID          *int       `json:"task_id,omitempty"`
	Billable        bool       `json:"billable,omitempty"`
	Start           time.Time  `json:"start,omitempty"`
	Stop            *time.Time `json:"stop,omitempty"`
	Duration        int        `json:"duration,omitempty"`
	Description     string     `json:"description,omitempty"`
	Tags            []string   `json:"tags,omitempty"`
	TagIDs          []int      `json:"tag_ids,omitempty"`
	DurationOnly    bool       `json:"duronly,omitempty"`
	At              time.Time  `json:"at,omitempty"`
	ServerDeletedAt *time.Time `json:"server_deleted_at,omitempty"`
	UserID          int        `json:"user_id,omitempty"`
	UID             int        `json:"uid,omitempty"`
	WID             int        `json:"wid,omitempty"`
}

type CreateTimeEntry struct {
	CreatedWith string    `json:"created_with,omitempty"`
	Description string    `json:"description,omitempty"`
	Start       time.Time `json:"start,omitempty"`
	Duration    int       `json:"duration,omitempty"`
	WorkspaceID int       `json:"workspace_id,omitempty"`
}

func (t TimeEntry) IsZero() bool {
	return t.ID == 0 && t.Start.IsZero()
}

func (c TogglClient) GetCurrentTimeEntry() (TimeEntry, error) {
	var timeEntry TimeEntry
	resp, err := c.httpGet(currentTimePath)
	defer resp.Body.Close()
	if err != nil {
		return TimeEntry{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&timeEntry); err != nil {
		return timeEntry, err
	}
	if timeEntry.IsZero() {
		return timeEntry, errors.New("no timer is running")
	}
	return timeEntry, nil
}

func (c TogglClient) StartTimeEntry(timeEntry CreateTimeEntry) error {
	url := fmt.Sprintf("api/v9/workspaces/%d/time_entries", timeEntry.WorkspaceID)
	resp, err := c.httpPost(url, timeEntry)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("got bad status code %d", resp.StatusCode))
	}
	return err
}
