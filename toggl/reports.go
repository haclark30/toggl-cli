package toggl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ProjectReportRequest struct {
	EndDate   string `json:"end_date"`
	StartDate string `json:"start_date"`
}

type ProjectSummary struct {
	BillableSeconds int `json:"billable_seconds,omitempty"`
	ProjectId       int `json:"project_id"`
	TrackedSeconds  int `json:"tracked_seconds"`
	UserId          int `json:"user_id"`
}

func (c TogglClient) GetProjectSummary(workspaceId int, startDate time.Time, endDate time.Time) ([]ProjectSummary, error) {
	var report []ProjectSummary
	url := fmt.Sprintf("reports/api/v3/workspace/%d/projects/summary", workspaceId)
	req := ProjectReportRequest{
		StartDate: startDate.Format(time.DateOnly),
		EndDate:   endDate.Format(time.DateOnly),
	}
	resp, err := c.httpPost(url, req)
	if err != nil {
		return report, err
	}

	if resp.StatusCode != http.StatusOK {
		return report, errors.New(fmt.Sprintf("got bad status code %d", resp.StatusCode))
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&report); err != nil {
		return report, err
	}

	return report, nil
}
