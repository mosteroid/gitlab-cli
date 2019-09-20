package client

import (
	"log"
	"sort"

	"github.com/xanzy/go-gitlab"
)

//Client rapresents gitlab client wrapper
type Client struct {
	*gitlab.Client
}

var client *Client

// InitClient initializes the gitlab client
func InitClient(baseURL, accessToken string) {
	client = &Client{gitlab.NewClient(nil, accessToken)}
	client.SetBaseURL(baseURL)
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
func (client *Client) GetProjectJobsStats(pid string) []*JobStats {

	var jobsStats = make([]*JobStats, 1)
	opt := &gitlab.ListJobsOptions{}
	jobs, _, err := client.Jobs.ListProjectJobs(pid, opt)

	if err != nil {
		log.Fatal(err)
	}

	jobsMap := make(map[string][]*gitlab.Job)
	for _, job := range jobs {
		jobsMap[job.Name] = append(jobsMap[job.Name], &job)
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
	return jobsStats
}

func calcAvgDuration(jobs []*gitlab.Job) float64 {

	if jobs == nil {
		return 0
	}

	jobsNum := len(jobs)
	if jobsNum > 0 {
		medianIndex := (jobsNum / 2) - 1
		jobsSlice := jobs[:]
		sort.Slice(jobsSlice, func(i, j int) bool {
			return jobsSlice[i].Duration < jobsSlice[j].Duration
		})

		if jobsNum%2 != 0 {
			return jobsSlice[medianIndex].Duration
		}

		return (jobsSlice[medianIndex-1].Duration + jobsSlice[medianIndex+1].Duration) / 2
	}

	return 0
}

func calcMinMaxDuration(jobs []*gitlab.Job) (float64, float64) {
	min := 0.0
	max := 0.0

	if jobs != nil {
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
