package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/fhopfensperger/git-releaser/pkg/repo"

	"github.com/stretchr/testify/assert"
)

func Test_getReposFromFile(t *testing.T) {
	repo1 := "https://github.com/fhopfensperger/amqp-sb-client.git"
	repo2 := "git@github.com:fhopfensperger/json-log-to-human-readable.git"
	fileName := "test.txt"
	f, _ := os.Create(fileName)
	f.WriteString(fmt.Sprintln(repo1))
	f.WriteString(fmt.Sprintln(repo2))
	repos := getReposFromFile(fileName)
	assert.Equal(t, []string{repo1, repo2}, repos)

	os.Remove(fileName)
}

func Test_getReposFromFile_ignore_empty_and_hashtag_lines(t *testing.T) {
	repo1 := "https://github.com/fhopfensperger/amqp-sb-client.git"
	repo2 := ""
	repo3 := "https://github.com/fhopfensperger/json-log-to-human-readable.git"
	repo4 := "#git@github.com:fhopfensperger/json-log-to-human-readable.git"
	fileName := "test.txt"
	f, _ := os.Create(fileName)
	f.WriteString(fmt.Sprintln(repo1))
	f.WriteString(fmt.Sprintln(repo2))
	f.WriteString(fmt.Sprintln(repo3))
	f.WriteString(fmt.Sprintln(repo4))
	repos := getReposFromFile(fileName)
	assert.Equal(t, []string{repo1, repo3}, repos)

	os.Remove(fileName)
}

func Test_getReposFromFile_failed(t *testing.T) {
	fileName := "test.txt"
	repos := getReposFromFile(fileName)
	assert.Equal(t, []string([]string(nil)), repos)

}

func Test_initConfig(t *testing.T) {
	os.Setenv("REPOS", "repos123")

	initConfig()

	assert.Equal(t, []string{"repos123"}, repos)
}
func Test_initConfig_major(t *testing.T) {
	os.Setenv("REPOS", "repos123")
	os.Setenv("NEXTVERSION", "MAJOR")

	initConfig()

	assert.Equal(t, nextVersion, repo.MAJOR)
}
func Test_initConfig_minor(t *testing.T) {
	os.Setenv("REPOS", "repos123")
	os.Setenv("NEXTVERSION", "MINOR")

	initConfig()

	assert.Equal(t, nextVersion, repo.MINOR)
}
func Test_initConfig_patch(t *testing.T) {
	os.Setenv("REPOS", "repos123")
	os.Setenv("NEXTVERSION", "PATCH")

	initConfig()

	assert.Equal(t, nextVersion, repo.PATCH)
}

func Test_initConfig_default(t *testing.T) {
	os.Setenv("REPOS", "repos123")
	os.Setenv("NEXTVERSION", "DEFAULT")

	initConfig()

	assert.Equal(t, nextVersion, repo.MINOR)
}

func TestExecute(t *testing.T) {
	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--version"})
	Execute("0.0.0")
	out, _ := ioutil.ReadAll(b)
	assert.Equal(t, "v0.0.0\n", string(out))
}

func TestExecute_repos_from_args(t *testing.T) {
	cmd := rootCmd
	testRepos := []string{"git@github.com:fhopfensperger/git-releaser.git"}
	cmd.SetArgs([]string{"create", "-s", "main", "release", "-r", testRepos[0]})
	Execute("0.0.0")

	assert.Equal(t, repos, testRepos)
	assert.Equal(t, targetBranch, "release")
}

func TestExecute_repos_from_args_not_existing(t *testing.T) {
	cmd := rootCmd
	testRepos := []string{"git@github.com:fhopfensperger/i-dont-exist.git"}
	cmd.SetArgs([]string{"create", "-s", "main", "release", "-r", testRepos[0]})
	Execute("0.0.0")
}
