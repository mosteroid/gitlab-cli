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
	"os"

	"github.com/mosteroid/gitlab-cli/client"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage the projects",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) {},
}

// listProjectsCmd represents the list command
var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "Display the list of projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient := client.GetClient()

		optSearchString, _ := cmd.Flags().GetString("search")
		optLimit, _ := cmd.Flags().GetInt("limit")
		opt := &gitlab.ListProjectsOptions{
			Membership:  gitlab.Bool(true),
			Search:      gitlab.String(optSearchString),
			ListOptions: gitlab.ListOptions{1, optLimit}}

		projects, _, err := gitlabClient.Projects.ListProjects(opt)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for index, project := range projects {
			fmt.Printf("%d) %s\n", index, project.PathWithNamespace)
		}
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(listProjectsCmd)

	listProjectsCmd.Flags().Int("limit", 10, "Sets the maximun number of results. The default value is 10")

	listProjectsCmd.Flags().String("search", "", "Search a project")

}
