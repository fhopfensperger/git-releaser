/*
Copyright © 2021 Florian Hopfensperger <f.hopfensperger@gmail.com>

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
	"os"

	"github.com/fhopfensperger/git-releaser/pkg/repo"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

// createCmd represents the branch command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Get remote create",
	Long:  `Get remote create`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, r := range repos {
			branchName, err := createNewReleaseVersion(r)
			if err != nil {
				log.Err(err).Msgf("For %s", r)
				os.Exit(1)
			}
			if branchName != "" {
				log.Info().Msgf("Successfully completed %s", r)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func createNewReleaseVersion(repoUrl string) (string, error) {
	r := repo.New(repoUrl)
	if r == nil {
		return "", errors.New("could not get repo")
	}

	if ref := r.GetSourceBranch(sourceBranch); ref == nil {
		return "", errors.New("could not get source branch")
	}
	if createBranch {
		r.GetVersionBranches(targetBranch)
	}

	if createTag {
		r.GetVersionTags()
	}

	r.GetLatestVersionReference()

	_, err := r.NextReleaseVersion(nextVersion)
	if err != nil {
		return "", err
	}
	if err := r.CreateNewRelease(createBranch, createTag); err != nil {
		return "", err
	}
	return repoUrl, nil
}
