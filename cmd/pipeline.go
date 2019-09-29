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
	"errors"
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

//jobTraceCmd represents the trace job command
var jobTraceCmd = &cobra.Command{
	Use:   "trace",
	Short: "Show job trace",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		project, _ := cmd.Flags().GetString("project")
		job, _ := cmd.Flags().GetInt("job")

		traceFile, _, _ := gitlabClient.Jobs.GetTraceFile(project, job)

		fmt.Print(traceFile)
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

func displayPipelineStatus(gitlabClient *client.Client, pid string, pipeline *gitlab.Pipeline, watch bool) error {

	if pipeline == nil {
		return errors.New("pipeline is required")
	}

	tw := util.NewTableWriter()
	tw.AppendHeader(table.Row{"ID", "REF", "STATUS"})
	tw.AppendRow(table.Row{pipeline.ID, pipeline.Ref, pipeline.Status})
	fmt.Println(tw.Render())

	if watch {
		opt := &gitlab.ListJobsOptions{}
		jobs, _, _ := gitlabClient.Jobs.ListPipelineJobs(pid, pipeline.ID, opt)
		jobsStats, _ := gitlabClient.GetProjectJobsStats(pid)
		jobsStatsMap := make(map[string]*client.JobStats)
		for _, stat := range jobsStats {
			jobsStatsMap[stat.Name] = stat
		}

		pw := util.NewProgressWriter()
		sw := util.NewStatusWriter()

		pw.SetNumTrackersExpected(len(jobs))
		pw.SetUpdateFrequency(WatchUpdateSleep)
		fmt.Print("\n\nPipeline progress:\n")
		go pw.Render()
		trackersMap := make(map[int]*progress.Tracker)
		done := false
		for !done {
			jobs, _, _ := gitlabClient.Jobs.ListPipelineJobs(pid, pipeline.ID, opt)
			for _, job := range jobs {
				if _, ok := trackersMap[job.ID]; !ok {
					total := int64(100)
					if stat, ok := jobsStatsMap[job.Name]; ok {
						total = int64(stat.AvgDuration)
					}
					trackersMap[job.ID] = &progress.Tracker{Message: fmt.Sprintf("%d) %s", job.ID, job.Name), Total: total, Units: util.UnitTime}
					pw.AppendTracker(trackersMap[job.ID])
				} else {
					if job.Status == "success" {
						duration := int64(job.FinishedAt.Sub(*job.StartedAt).Seconds())
						trackersMap[job.ID].SetValue(duration)
						trackersMap[job.ID].MarkAsDone()
					} else if job.Status == "running" {
						duration := int64(time.Since(*job.StartedAt) / time.Second)
						trackersMap[job.ID].SetValue(duration)
					}
				}
			}

			pipeline, _, err := gitlabClient.Pipelines.GetPipeline(pid, pipeline.ID)
			if pipeline.Status != "running" && pipeline.Status != "pending" || err != nil {
				done = true
				fmt.Printf("The pipeline %d exit with status: %s \n", pipeline.ID, sw.Sprintf(pipeline.Status))
			} else {
				time.Sleep(WatchUpdateSleep)
			}
		}
	}

	return nil
}

// runCmd represents the run command
var pipelineStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Diplay a pipeline status",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		pid, _ := cmd.Flags().GetString("project")
		pipelineID, _ := cmd.Flags().GetInt("pipeline")

		pipeline, _, _ := gitlabClient.Pipelines.GetPipeline(pid, pipelineID)
		watch := pipeline.Status == "running" || pipeline.Status == "pending"
		displayPipelineStatus(gitlabClient, pid, pipeline, watch)
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
		pid, _ := cmd.Flags().GetString("project")
		watch, _ := cmd.Flags().GetBool("watch")

		opt := &gitlab.CreatePipelineOptions{Ref: gitlab.String(ref)}
		pipeline, _, err := gitlabClient.Pipelines.CreatePipeline(pid, opt)
		if err != nil {
			log.Fatal(err)
		}

		displayPipelineStatus(gitlabClient, pid, pipeline, watch)
	},
}

func init() {
	rootCmd.AddCommand(pipelineCmd)
	pipelineCmd.AddCommand(runCmd)
	pipelineCmd.AddCommand(listPipelinesCmd)
	pipelineCmd.AddCommand(jobsCmd)
	pipelineCmd.AddCommand(pipelineStatusCmd)
	jobsCmd.AddCommand(jobStatsCmd)
	jobsCmd.AddCommand(jobTraceCmd)

	pipelineCmd.PersistentFlags().StringP("project", "p", "", "Set the project name or project ID")
	cobra.MarkFlagRequired(pipelineCmd.PersistentFlags(), "project")

	runCmd.Flags().StringP("ref", "r", "", "Set the ref")
	cobra.MarkFlagRequired(runCmd.Flags(), "ref")

	pipelineStatusCmd.Flags().IntP("pipeline", "l", -1, "Set the pipeline id")
	runCmd.Flags().BoolP("watch", "w", false, "Watch the pipeline execution")

	jobsCmd.Flags().IntP("pipeline", "l", -1, "Set the pipeline id")
	jobTraceCmd.Flags().IntP("job", "j", -1, "Set the job id")
	cobra.MarkFlagRequired(jobTraceCmd.Flags(), "job")

}
