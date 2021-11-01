//go:generate mockgen -source=$GOFILE -package=mock -destination=./mocks/$GOFILE
package circleci

import (
	"context"
	"fmt"
	"time"
)

// Jobs describes all the job related methods that the CircleCI API
// supports.
//
// CircleCI API docs: https://circleci.com/docs/api/v2/#tag/Job
type Jobs interface {
	// Returns job details.
	Get(ctx context.Context, projectSlug string, jobNumber string) (*Job, error)

	// Cancel job with a given job number.
	Cancel(ctx context.Context, projectSlug string, jobNumber string) error

	// Returns a job's artifacts.
	ListArtifacts(ctx context.Context, projectSlug string, jobNumber string) (*ArtifactList, error)

	// Get test metadata for a build.
	ListTestMetadata(ctx context.Context, projectSlug string, jobNumber string) (*TestMetadataList, error)
}

// jobs implments Jobs interface
type jobs struct {
	client *Client
}

// Job represents a CircleCI job.
type Job struct {
	WebURL         string          `json:"web_url"`
	Project        *JobProject     `json:"project"`
	ParallelRuns   []*ParallelRuns `json:"parallel_runs"`
	StartedAt      time.Time       `json:"started_at"`
	LatestWorkflow *LatestWorkflow `json:"latest_workflow"`
	Name           string          `json:"name"`
	Executor       *Executor       `json:"executor"`
	Parallelism    int             `json:"parallelism"`
	Status         interface{}     `json:"status"`
	Number         int             `json:"number"`
	Pipeline       *JobPipeline    `json:"pipeline"`
	Duration       int             `json:"duration"`
	CreatedAt      time.Time       `json:"created_at"`
	Messages       []*JobMessage   `json:"messages"`
	Contexts       []*Context      `json:"contexts"`
	Organization   Organization    `json:"organization"`
	QueuedAt       time.Time       `json:"queued_at"`
	StoppedAt      time.Time       `json:"stopped_at"`
}

type JobProject struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	ExternalURL string `json:"external_url"`
}

type ParallelRuns struct {
	Index  int    `json:"index"`
	Status string `json:"status"`
}

type LatestWorkflow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Executor struct {
	Type          string `json:"type"`
	ResourceClass string `json:"resource_class"`
}

type JobPipeline struct {
	ID string `json:"id"`
}

type JobMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

type JobContext struct {
	Name string `json:"name"`
}

type Organization struct {
	Name string `json:"name"`
}

func (s *jobs) Get(ctx context.Context, projectSlug string, jobNumber string) (*Job, error) {
	if !validString(&projectSlug) {
		return nil, ErrRequiredProjectSlug
	}

	if !validString(&jobNumber) {
		return nil, ErrRequiredJobNumber
	}

	u := fmt.Sprintf("project/%s/job/%s", projectSlug, jobNumber)
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	j := &Job{}
	err = s.client.do(ctx, req, j)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func (s *jobs) Cancel(ctx context.Context, projectSlug string, jobNumber string) error {
	if !validString(&projectSlug) {
		return ErrRequiredProjectSlug
	}

	if !validString(&jobNumber) {
		return ErrRequiredJobNumber
	}

	u := fmt.Sprintf("project/%s/job/%s/cancel", projectSlug, jobNumber)
	req, err := s.client.newRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

type ArtifactList struct {
	Items         []*Artifact `json:"items"`
	NextPageToken string      `json:"next_page_token"`
}

type Artifact struct {
	Path      string `json:"path"`
	NodeIndex int64  `json:"node_index"`
	URL       string `json:"url"`
}

func (s *jobs) ListArtifacts(ctx context.Context, projectSlug string, jobNumber string) (*ArtifactList, error) {
	if !validString(&projectSlug) {
		return nil, ErrRequiredProjectSlug
	}

	if !validString(&jobNumber) {
		return nil, ErrRequiredJobNumber
	}

	u := fmt.Sprintf("project/%s/%s/artifacts", projectSlug, jobNumber)
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	al := &ArtifactList{}
	err = s.client.do(ctx, req, al)
	if err != nil {
		return nil, err
	}

	return al, nil
}

type TestMetadataList struct {
	Items         []*TestMetadata `json:"items"`
	NextPageToken string          `json:"next_page_token"`
}

type TestMetadata struct {
	Message   string `json:"message"`
	Source    string `json:"source"`
	RunTime   string `json:"run_time"`
	File      string `json:"file"`
	Result    string `json:"result"`
	Name      string `json:"name"`
	Classname string `json:"classname"`
}

func (s *jobs) ListTestMetadata(ctx context.Context, projectSlug string, jobNumber string) (*TestMetadataList, error) {
	if !validString(&projectSlug) {
		return nil, ErrRequiredProjectSlug
	}

	if !validString(&jobNumber) {
		return nil, ErrRequiredJobNumber
	}

	u := fmt.Sprintf("project/%s/%s/tests", projectSlug, jobNumber)
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	tml := &TestMetadataList{}
	err = s.client.do(ctx, req, tml)
	if err != nil {
		return nil, err
	}

	return tml, nil
}
