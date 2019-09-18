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

// listCmd represents the run command
var listPipelinesCmd = &cobra.Command{
	Use:   "list",
	Short: "List a project pipelines",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// runCmd represents the run command
var runPipelinesCmd = &cobra.Command{
	Use:   "run",
	Short: "Run pipelines",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		branches, _ := cmd.Flags().GetStringArray("branches")
		project, _ := cmd.Flags().GetString("project")

		for _, branch := range branches {
			opt := &gitlab.CreatePipelineOptions{Ref: gitlab.String(branch)}
			_, _, err := gitlabClient.Pipelines.CreatePipeline(project, opt)
			if err != nil {
				fmt.Println(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(pipelinesCmd)

	pipelinesCmd.AddCommand(runPipelinesCmd)

	runPipelinesCmd.Flags().StringArrayP("branches", "b", []string{}, "Set the branches list")
	runPipelinesCmd.Flags().StringP("project", "p", "", "Set the project name or project ID")
	cobra.MarkFlagRequired(runPipelinesCmd.Flags(), "branches")
	cobra.MarkFlagRequired(runPipelinesCmd.Flags(), "project")

}
