package client

import (
	"sort"

	"github.com/xanzy/go-gitlab"
)

//Client rapresents gitlab client wrapper
type Client struct {
	*gitlab.Client
}

var client *Client

// InitClient initializes the gitlab client
func InitClient(baseURL, accessToken string) error {
	client = &Client{gitlab.NewClient(nil, accessToken)}
	return client.SetBaseURL(baseURL)
}

// GetClient returns the initialized gitlabe client
func GetClient() *Client {
	return client
}

//JobStats rapresents the job stats
type JobStats struct {
	Name        string
	Total       int
	MaxDuration float64
	MinDuration float64
	AvgDuration float64
}

// GetProjectJobsStats returns the project jobs stats
func (client *Client) GetProjectJobsStats(pid string) ([]*JobStats, error) {

	jobsStats := *new([]*JobStats)
	opt := &gitlab.ListJobsOptions{}
	jobs, _, err := client.Jobs.ListProjectJobs(pid, opt)

	if err != nil {
		return nil, err
	}

	jobsMap := make(map[string][]*gitlab.Job)
	for i := 0; i < len(jobs); i++ {
		jobsMap[jobs[i].Name] = append(jobsMap[jobs[i].Name], &jobs[i])
	}

	for _, jobs := range jobsMap {
		minDuration, maxDuration := calcMinMaxDuration(jobs)
		jobsStats = append(jobsStats, &JobStats{
			Name:        jobs[0].Name,
			Total:       len(jobs),
			MinDuration: minDuration,
			MaxDuration: maxDuration,
			AvgDuration: calcAvgDuration(jobs),
		})
	}

	return jobsStats, nil
}

func calcAvgDuration(jobs []*gitlab.Job) float64 {

	if len(jobs) == 0 {
		return 0
	}

	jobsNum := len(jobs)
	medianIndex := jobsNum / 2
	sort.SliceStable(jobs, func(i, j int) bool {
		return jobs[i].Duration < jobs[j].Duration
	})

	if jobsNum%2 != 0 {
		return jobs[medianIndex].Duration
	}

	return (jobs[medianIndex-1].Duration + jobs[medianIndex].Duration) / 2

}

func calcMinMaxDuration(jobs []*gitlab.Job) (float64, float64) {
	min := 0.0
	max := 0.0

	if len(jobs) > 0 {
		for _, job := range jobs {

			if job.Duration < min {
				min = job.Duration
			}
			if job.Duration > max {
				max = job.Duration
			}
		}
	}

	return min, max
}
