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
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fhopfensperger/git-releaser/pkg"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
	"golang.org/x/mod/semver"

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

func createNewReleaseVersion(repo string) (string, error) {
	gitService := pkg.RemoteBranch{}
	branches := gitService.GetAllRemoteBranchesAndTags(repo)
	latestBranch := getLatestVersion(branches, targetBranch)
	branchToCopy := getBranchByName(branches, sourceBranch)
	var newReleaseBranchName string
	var err error
	if latestBranch == nil {
		newReleaseBranchName, err = getNextReleaseVersion("", false)
		if err != nil {
			return "", err
		}
	} else if latestBranch.Hash() == branchToCopy.Hash() {
		log.Info().Msgf("Nothing to do, %s and latestBranch %s are equals, commit hash: %s", branchToCopy.Name().Short(), latestBranch.Name().Short(), latestBranch.Hash())
		return "", nil
	} else {
		newReleaseBranchName, err = getNextReleaseVersion(latestBranch.Name().Short(), nextMinor)
		if err != nil {
			return "", err
		}
	}
	fullNewReleaseBranchName := fmt.Sprintf("refs/heads/%s/%s", targetBranch, newReleaseBranchName)
	newReleaseBranchRef := plumbing.NewReferenceFromStrings(fullNewReleaseBranchName, fullNewReleaseBranchName)
	err = gitService.CreateBranch(branchToCopy, newReleaseBranchRef, createTag)
	if err != nil {
		return "", err
	}
	return newReleaseBranchName, nil
}

func getLatestVersion(refs []*plumbing.Reference, branchFilter string) *plumbing.Reference {
	if branchFilter == "" {
		log.Warn().Msg("No branchfilter defined")
		os.Exit(1)
	}

	for i := len(refs) - 1; i >= 0; i-- {
		if pkg.VersionRegex.FindString(refs[i].Name().Short()) != "" && strings.Contains(refs[1].Name().Short(), branchFilter) {
			return refs[i]
		}
	}
	return nil
}

func getBranchByName(refs []*plumbing.Reference, branchName string) *plumbing.Reference {
	for _, ref := range refs {
		if ref.Name().Short() == branchName {
			return ref
		}
	}
	return nil
}

func getNextReleaseVersion(latestVersion string, nextMinor bool) (string, error) {
	if latestVersion == "" {
		return "v1.0.0", nil
	}
	semLatestVersion := semver.Canonical(pkg.VersionRegex.FindString(latestVersion))
	latestVersionSlice := strings.Split(semLatestVersion, ".")

	var versionNumber = regexp.MustCompile(`\d`)

	major, err := strconv.Atoi(versionNumber.FindString(latestVersionSlice[len(latestVersionSlice)-3]))
	if err != nil {
		return "", err
	}
	minor, err := strconv.Atoi(latestVersionSlice[len(latestVersionSlice)-2])
	if err != nil {
		return "", err
	}
	patch, err := strconv.Atoi(latestVersionSlice[len(latestVersionSlice)-1])
	if err != nil {
		return "", err
	}
	if nextMinor {
		return semver.Canonical(fmt.Sprintf("v%v.%v.%v", major, minor+1, 0)), nil
	} else {
		return semver.Canonical(fmt.Sprintf("v%v.%v.%v", major, minor, patch+1)), nil
	}
}
