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
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/spf13/viper"

	"github.com/fhopfensperger/git-releaser/pkg/repo"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var pat string

// createCmd represents the branch command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a tag or version",
	Long:  `Creates a tag or version`,
	Run: func(cmd *cobra.Command, args []string) {
		pat = viper.GetString("pat")

		for _, r := range repos {
			branchName, err := createNewReleaseVersion(r)
			if err != nil {
				log.Err(err).Msgf("For %s", r)
			}
			if branchName != "" {
				log.Info().Msgf("Successfully completed %s", r)
			}

		}
	},
}

func init() {
	flags := createCmd.Flags()
	flags.StringP("pat", "p", "", `Use a Git Personal Access Token instead of the default private certificate! You could also set a environment variable. "export PAT=123456789" `)
	_ = viper.BindPFlag("pat", flags.Lookup("pat"))
	rootCmd.AddCommand(createCmd)
}

func createNewReleaseVersion(repoUrl string) (string, error) {
	if strings.Contains(repoUrl, "https://") {
		log.Info().Msgf(`Using PAT "-p" instead of ssh private certificate for repo %s`, repoUrl)
	}

	r := repo.New(repoUrl, &http.BasicAuth{
		Username: "123", // Using a PAT this can be anything except an empty string
		Password: pat,
	})
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
