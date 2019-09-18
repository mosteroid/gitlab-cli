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
	"github.com/mosteroid/gitlab-cli/client"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

// pipelinesCmd represents the pipelines command
var pipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "Manage the pipelines",
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

// runPipelinesCmd represents the run command
var runPipelinesCmd = &cobra.Command{
	Use:   "run",
	Short: "Run pipelines",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		refs, _ := cmd.Flags().GetStringArray("refs")
		project, _ := cmd.Flags().GetString("project")

		for _, ref := range refs {
			opt := &gitlab.CreatePipelineOptions{Ref: gitlab.String(ref)}
			_, _, err := gitlabClient.Pipelines.CreatePipeline(project, opt)
			if err != nil {
				log.Fatal(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(pipelinesCmd)
	pipelinesCmd.AddCommand(runPipelinesCmd)
	pipelinesCmd.AddCommand(listPipelinesCmd)

	pipelinesCmd.PersistentFlags().StringP("project", "p", "", "Set the project name or project ID")
	cobra.MarkFlagRequired(pipelinesCmd.PersistentFlags(), "project")

	runPipelinesCmd.Flags().StringArrayP("refs", "r", []string{}, "Set the refs list")
	cobra.MarkFlagRequired(runPipelinesCmd.Flags(), "refs")
}
