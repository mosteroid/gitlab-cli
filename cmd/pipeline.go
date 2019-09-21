/*
Copyright © 2019 The Mosteroid Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/jedib0t/go-pretty/progress"
	"github.com/jedib0t/go-pretty/table"
	"github.com/mosteroid/gitlabctl/client"
	"github.com/mosteroid/gitlabctl/util"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

const (
	// WatchUpdateSleep the watch sleep
	WatchUpdateSleep = 1000 * time.Millisecond
)

type jobTracker struct {
	Tracker   *progress.Tracker
	StartTime time.Time
}

// pipelineCmd represents the pipelines command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Manage pipelines",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) {},
}

// listPipelinesCmd represents the run command
var listPipelinesCmd = &cobra.Command{
	Use:   "list",
	Short: "List a project pipelines",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		project, _ := cmd.Flags().GetString("project")
		opt := &gitlab.ListProjectPipelinesOptions{}

		pipelines, _, _ := gitlabClient.Pipelines.ListProjectPipelines(project, opt)

		tw := util.NewTableWriter()

		tw.AppendHeader(table.Row{"ID", "Branch", "Status", "SHA", "URL"})
		for _, pipeline := range pipelines {
			tw.AppendRow(table.Row{pipeline.ID, pipeline.Ref, pipeline.Status, pipeline.SHA, pipeline.WebURL})
		}
		fmt.Println(tw.Render())
	},
}

// jobsCmd represents the list pipeline jobs command
var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "List the jobs of a pipelines",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		project, _ := cmd.Flags().GetString("project")
		pipeline, _ := cmd.Flags().GetInt("pipeline")

		opt := &gitlab.ListJobsOptions{}

		var jobs []gitlab.Job

		if pipeline != -1 {
			jobsPointers, _, _ := gitlabClient.Jobs.ListPipelineJobs(project, pipeline, opt)
			for _, jobPointer := range jobsPointers {
				jobs = append(jobs, *jobPointer)
			}
		} else {
			jobs, _, _ = gitlabClient.Jobs.ListProjectJobs(project, opt)
		}

		tw := util.NewTableWriter()

		tw.AppendHeader(table.Row{"ID", "NAME", "STAGE", "STATUS", "STARTED AT"})
		for _, job := range jobs {
			tw.AppendRow(table.Row{job.ID, job.Name, job.Stage, job.Status, job.StartedAt})
		}
		fmt.Println(tw.Render())
	},
}

// jobStatsCmd represents the list jobs stats command
var jobStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "List the stats of jobs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		project, _ := cmd.Flags().GetString("project")

		stats, _ := gitlabClient.GetProjectJobsStats(project)

		tw := util.NewTableWriter()

		tw.AppendHeader(table.Row{"NAME", "MIN DURATION", "MAX DURATION", "AVG DURATION"})
		for _, stat := range stats {
			tw.AppendRow(table.Row{stat.Name, stat.MinDuration, stat.MaxDuration, stat.AvgDuration})
		}
		fmt.Println(tw.Render())
	},
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run pipelines",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		ref, _ := cmd.Flags().GetString("ref")
		project, _ := cmd.Flags().GetString("project")
		watch, _ := cmd.Flags().GetBool("watch")

		opt := &gitlab.CreatePipelineOptions{Ref: gitlab.String(ref)}
		pipeline, _, err := gitlabClient.Pipelines.CreatePipeline(project, opt)
		if err != nil {
			log.Fatal(err)
		}

		if watch {
			opt := &gitlab.ListJobsOptions{}
			jobs, _, _ := gitlabClient.Jobs.ListPipelineJobs(project, pipeline.ID, opt)
			jobsStats, _ := gitlabClient.GetProjectJobsStats(project)
			jobsStatsMap := make(map[string]*client.JobStats)
			for _, stat := range jobsStats {
				jobsStatsMap[stat.Name] = stat
			}

			pw := util.NewProgressWriter()
			sw := util.NewStatusWriter()

			pw.SetNumTrackersExpected(len(jobs))
			pw.SetUpdateFrequency(500 * time.Millisecond)
			go pw.Render()
			trackersMap := make(map[int]*jobTracker)
			done := false
			for !done {
				jobs, _, _ := gitlabClient.Jobs.ListPipelineJobs(project, pipeline.ID, opt)
				for _, job := range jobs {
					if _, ok := trackersMap[job.ID]; !ok {
						if job.Status == "running" || job.Status == "success" {
							total := int64(100)
							if stat, ok := jobsStatsMap[job.Name]; ok {
								total = int64(stat.AvgDuration)
							}
							trackersMap[job.ID] = &jobTracker{Tracker: &progress.Tracker{Message: fmt.Sprintf("%d) %s", job.ID, job.Name), Total: total, Units: progress.UnitsDefault}, StartTime: time.Now().UTC()}
							pw.AppendTracker(trackersMap[job.ID].Tracker)
						}
					} else {
						if job.Status == "success" {
							trackersMap[job.ID].Tracker.MarkAsDone()
						} else {
							duration := int64(time.Since(trackersMap[job.ID].StartTime) / time.Second)
							trackersMap[job.ID].Tracker.SetValue(duration)
						}
					}
				}

				pipeline, _, err := gitlabClient.Pipelines.GetPipeline(project, pipeline.ID)
				if pipeline.Status != "running" && pipeline.Status != "pending" || err != nil {
					done = true
					fmt.Printf("The pipeline %d exit with status: %s \n", pipeline.ID, sw.Sprintf(pipeline.Status))
				} else {
					time.Sleep(WatchUpdateSleep)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(pipelineCmd)
	pipelineCmd.AddCommand(runCmd)
	pipelineCmd.AddCommand(listPipelinesCmd)
	pipelineCmd.AddCommand(jobsCmd)
	jobsCmd.AddCommand(jobStatsCmd)

	pipelineCmd.PersistentFlags().StringP("project", "p", "", "Set the project name or project ID")
	cobra.MarkFlagRequired(pipelineCmd.PersistentFlags(), "project")

	runCmd.Flags().StringP("ref", "r", "", "Set the ref")
	cobra.MarkFlagRequired(runCmd.Flags(), "ref")

	runCmd.Flags().BoolP("watch", "w", false, "Watch the pipeline execution")

	jobsCmd.Flags().Int("pipeline", -1, "List pipeline jobs")
}
