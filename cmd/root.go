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
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var repos []string
var targetBranch string
var sourceBranch string
var fileName string
var nextMinor bool
var createTag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-releaser",
	Short: "Simple command line utility to create a new release branch",
	Long:  `Simple command line utility to create a new release branch`,
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

	pf.StringP("source-branch", "s", "", "Which branch should be used to create release branch & tag")
	_ = viper.BindPFlag("source-branch", pf.Lookup("source-branch"))
	_ = cobra.MarkFlagRequired(pf, "source-branch")

	pf.StringP("target-branch", "b", "release", "Which target branches to check for version")
	_ = viper.BindPFlag("target-branch", pf.Lookup("target-branch"))

	pf.BoolP("next-minor-release", "n", false, "Next Version should be a minor release")
	_ = viper.BindPFlag("next-minor-release", pf.Lookup("next-minor-release"))

	pf.BoolP("create-tag", "t", false, "Create a tag also")
	_ = viper.BindPFlag("create-tag", pf.Lookup("create-tag"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	pf.StringP("file", "f", "", "Uses repos from file (one repo per line)")
	_ = viper.BindPFlag("file", pf.Lookup("file"))
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.SetVersionTemplate(`{{printf "v%s\n" .Version}}`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.AutomaticEnv() // read in environment variables that match

	repos = viper.GetStringSlice("repos")
	sourceBranch = viper.GetString("source-branch")
	fileName = viper.GetString("file")
	targetBranch = viper.GetString("target-branch")
	nextMinor = viper.GetBool("next-minor-release")
	createTag = viper.GetBool("create-tag")

	if fileName != "" {
		repos = getReposFromFile(fileName)
	}
	if len(repos) == 0 && fileName == "" {
		fmt.Println("Either -f (file) or -r (repos) must be set")
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