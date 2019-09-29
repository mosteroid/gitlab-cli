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

	"github.com/jedib0t/go-pretty/table"
	"github.com/mosteroid/gitlabctl/client"
	"github.com/mosteroid/gitlabctl/util"
	"github.com/spf13/cobra"
)

// jobsCmd represents the pipelines command
var jobsCmd = &cobra.Command{
	Use:   "job",
	Short: "Manage jobs",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) {},
}

//jobTraceCmd represents the trace job command
var jobTraceCmd = &cobra.Command{
	Use:   "trace",
	Short: "Show a job trace",
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

// retryJobCmd represents the retry job command
var retryJobCmd = &cobra.Command{
	Use:   "retry",
	Short: "Retry a job",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		project, _ := cmd.Flags().GetString("project")
		jobID, _ := cmd.Flags().GetInt("job")

		_, _, err := gitlabClient.Jobs.RetryJob(project, jobID)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Job restarted")
	},
}

// cancelJobCmd represents the cancel job command
var cancelJobCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a job",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		project, _ := cmd.Flags().GetString("project")
		jobID, _ := cmd.Flags().GetInt("job")

		_, _, err := gitlabClient.Jobs.CancelJob(project, jobID)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Job cancelled")
	},
}

// runJobCmd represents the cancel job command
var runJobCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a job",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		project, _ := cmd.Flags().GetString("project")
		jobID, _ := cmd.Flags().GetInt("job")

		_, _, err := gitlabClient.Jobs.PlayJob(project, jobID)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Job started")
	},
}

func init() {
	rootCmd.AddCommand(jobsCmd)
	jobsCmd.AddCommand(jobStatsCmd)
	jobsCmd.AddCommand(jobTraceCmd)
	jobsCmd.AddCommand(retryJobCmd)
	jobsCmd.AddCommand(cancelJobCmd)
	jobsCmd.AddCommand(runJobCmd)

	jobsCmd.PersistentFlags().StringP("project", "p", "", "Set the project name or project ID")
	cobra.MarkFlagRequired(jobsCmd.PersistentFlags(), "project")

	jobsCmd.Flags().IntP("job", "j", -1, "Set the job id")
	cobra.MarkFlagRequired(jobsCmd.Flags(), "job")
}
