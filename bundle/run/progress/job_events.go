package progress

import "fmt"

type JobRunUrlEvent struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

func NewJobRunUrlEvent(url string) *JobRunUrlEvent {
	return &JobRunUrlEvent{
		Type: "job_run_url",
		Url:  url,
	}
}

func (event *JobRunUrlEvent) String() string {
	return fmt.Sprintf("Run URL: %s\n", event.Url)
}

func (event *JobRunUrlEvent) IsInplaceSupported() bool {
	return false
}
