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
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

// pipelinesCmd represents the pipelines command
var pipelinesCmd = &cobra.Command{
	Use:   "pipelines",
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

		tw := table.NewWriter()
		tw.Style().Options.DrawBorder = false
		tw.Style().Options.SeparateColumns = false
		tw.Style().Options.SeparateHeader = false
		tw.Style().Options.SeparateRows = false
		tw.AppendHeader(table.Row{"ID", "Branch", "Status", "SHA", "URL"})
		for _, pipeline := range pipelines {
			tw.AppendRow(table.Row{pipeline.ID, pipeline.Ref, pipeline.Status, pipeline.SHA, pipeline.WebURL})
		}
		fmt.Println(tw.Render())
	},
}

// listJobsCmd represents the list pipeline jobs command
var listJobsCmd = &cobra.Command{
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

		tw := table.NewWriter()
		tw.Style().Options.DrawBorder = false
		tw.Style().Options.SeparateColumns = false
		tw.Style().Options.SeparateHeader = false
		tw.Style().Options.SeparateRows = false

		tw.AppendHeader(table.Row{"ID", "NAME", "STAGE", "STATUS", "STARTED AT"})
		for _, job := range jobs {
			tw.AppendRow(table.Row{job.ID, job.Name, job.Stage, job.Status, job.StartedAt})
		}
		fmt.Println(tw.Render())
	},
}

// listJobsStatsCmd represents the list jobs stats command
var listJobsStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "List the stats of jobs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		project, _ := cmd.Flags().GetString("project")

		stats := gitlabClient.GetProjectJobsStats(project)

		tw := table.NewWriter()
		tw.Style().Options.DrawBorder = false
		tw.Style().Options.SeparateColumns = false
		tw.Style().Options.SeparateHeader = false
		tw.Style().Options.SeparateRows = false

		tw.AppendHeader(table.Row{"NAME", "MIN DURATION", "MAX DURATION", "AVG DURATION"})
		for _, stat := range stats {
			tw.AppendRow(table.Row{stat.Name, stat.MinDuration, stat.MaxDuration, stat.AvgDuration})
		}
		fmt.Println(tw.Render())
	},
}

// runPipelinesCmd represents the run command
var runPipelinesCmd = &cobra.Command{
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
			pw := progress.NewWriter()
			pw.SetAutoStop(false)
			pw.SetTrackerLength(25)
			pw.ShowOverallTracker(true)
			pw.ShowTime(true)
			pw.ShowTracker(true)
			pw.ShowValue(true)
			pw.SetMessageWidth(24)
			pw.SetNumTrackersExpected(len(jobs))
			pw.SetSortBy(progress.SortByPercentDsc)
			pw.SetStyle(progress.StyleDefault)
			pw.SetTrackerPosition(progress.PositionRight)
			pw.SetUpdateFrequency(time.Millisecond * 1000)
			pw.Style().Colors = progress.StyleColorsExample
			pw.Style().Options.PercentFormat = "%4.1f%%"
			go pw.Render()
			trackersMap := make(map[int]*progress.Tracker)
			done := false
			for !done {
				jobs, _, _ := gitlabClient.Jobs.ListPipelineJobs(project, pipeline.ID, opt)
				for _, job := range jobs {
					if _, ok := trackersMap[job.ID]; !ok {
						if job.Status == "running" || job.Status == "success" {
							trackersMap[job.ID] = &progress.Tracker{Message: job.Name, Total: 100, Units: progress.UnitsDefault}
							pw.AppendTracker(trackersMap[job.ID])
						}
					} else {
						if job.Status == "success" {
							trackersMap[job.ID].SetValue(100)
						} else {
							trackersMap[job.ID].SetValue(int64(job.Coverage))
						}
					}
				}
				time.Sleep(time.Millisecond * 1000)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(pipelinesCmd)
	pipelinesCmd.AddCommand(runPipelinesCmd)
	pipelinesCmd.AddCommand(listPipelinesCmd)
	pipelinesCmd.AddCommand(listJobsCmd)
	listJobsCmd.AddCommand(listJobsStatsCmd)

	pipelinesCmd.PersistentFlags().StringP("project", "p", "", "Set the project name or project ID")
	cobra.MarkFlagRequired(pipelinesCmd.PersistentFlags(), "project")

	runPipelinesCmd.Flags().StringP("ref", "r", "", "Set the ref")
	cobra.MarkFlagRequired(runPipelinesCmd.Flags(), "ref")

	runPipelinesCmd.Flags().BoolP("watch", "w", false, "Watch the pipeline execution")

	listJobsCmd.Flags().Int("pipeline", -1, "List pipeline jobs")
}
