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
	"bufio"
	"os"

	"github.com/fhopfensperger/git-releaser/pkg/repo"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var repos []string
var targetBranch string
var sourceBranch string
var fileName string
var createTag bool
var createBranch bool
var nextVersion int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-releaser",
	Short: "Simple command line utility to create a new release branch or tag based on semver",
	Long: `Simple command line utility to create a new release branch or tag based on semver. 
More information can be found here: https://github.com/fhopfensperger/git-releaser
Author: Florian Hopfensperger <f.hopfensperger@gmail.com>`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		log.Err(err).Msg("")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git-releaser.yaml)")
	pf := rootCmd.PersistentFlags()
	pf.StringSliceP("repos", "r", []string{}, "Git Repo urls e.g. git@github.com:fhopfensperger/my-repo.git")
	_ = viper.BindPFlag("repos", pf.Lookup("repos"))

	pf.StringP("source", "s", "main", "Source reference branch")
	_ = viper.BindPFlag("source", pf.Lookup("source"))

	pf.StringP("target", "b", "release", "Which target branches to check for version")
	_ = viper.BindPFlag("target", pf.Lookup("target"))

	pf.BoolP("tag", "t", false, "Create a release version tag")
	_ = viper.BindPFlag("tag", pf.Lookup("tag"))

	pf.BoolP("branch", "c", false, "Create a release version branch")
	_ = viper.BindPFlag("branch", pf.Lookup("branch"))

	pf.StringP("nextversion", "n", "PATCH", "Which number should be incremented by 1. Possible values: PATCH, MINOR, MAJOR")
	_ = viper.BindPFlag("nextversion", pf.Lookup("nextversion"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	pf.StringP("file", "f", "", "Use repos from file (one repo per line, line with a leading # will be ignored)")
	_ = viper.BindPFlag("file", pf.Lookup("file"))
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.SetVersionTemplate(`{{printf "v%s\n" .Version}}`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	repos = viper.GetStringSlice("repos")
	sourceBranch = viper.GetString("source")
	fileName = viper.GetString("file")
	targetBranch = viper.GetString("target")
	createBranch = viper.GetBool("branch")
	createTag = viper.GetBool("tag")

	nv := viper.GetString("nextversion")

	switch nv {
	case "PATCH":
		log.Info().Msg("New PATCH version will be created")
		nextVersion = repo.PATCH
	case "MINOR":
		log.Info().Msg("New MINOR version will be created")
		nextVersion = repo.MINOR
	case "MAJOR":
		log.Info().Msg("New MAJOR version will be created")
		nextVersion = repo.MAJOR
	default:
		log.Info().Msgf("New MINOR version will be created, as %s is unknown", nv)
		nextVersion = repo.MINOR
	}

	if fileName != "" {
		repos = getReposFromFile(fileName)
	}
	if len(repos) == 0 && fileName == "" {
		log.Err(nil).Msg("Either -f (file) or -r (repos) must be set")
		os.Exit(1)
	}
}

func getReposFromFile(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Err(err).Msgf("Could not open file %s", fileName)
		return nil
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 0 && string(line[0]) != "#" {
			lines = append(lines, line)
		}
	}
	return lines
}
