package pkg

import (
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"golang.org/x/mod/semver"

	"github.com/rs/zerolog/log"

	"github.com/go-git/go-git/v5/config"

	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-git/v5"
)

var VersionRegex = regexp.MustCompile(`v\d+(\.\d+)+`)

//GitInterface all of the functions we use from the third party client
// to be able to mock them in the tests.
type GitInterface interface {
	List(*git.ListOptions) ([]*plumbing.Reference, error)
	Config() *config.RemoteConfig
	Push(*git.PushOptions) error
}

//RemoteBranch to implement the interface
type RemoteBranch struct {
	gitClient GitInterface
	repo      *git.Repository
}

//New constructor
func New(client GitInterface) RemoteBranch {
	return RemoteBranch{client, nil}
}

//AddRepo to add a remote repo to RemoteBranch struct
func (m *RemoteBranch) AddRepo(repo *git.Repository) {
	m.repo = repo
}

//GetRemoteBranches get remote branches from GitHub using the repoURL and the branchFilter
func (m *RemoteBranch) GetAllRemoteBranchesAndTags(repoURL string) []*plumbing.Reference {
	if m.gitClient == nil {
		m.gitClient = git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
			Name: "origin",
			URLs: []string{repoURL},
		})
	}

	// We can then use every Remote functions to retrieve wanted information
	refs, err := m.gitClient.List(&git.ListOptions{})
	if err != nil {
		log.Err(err).Msg("")
	}

	// Filters the references list and only branches which apply to the filter
	var branches []*plumbing.Reference
	for _, ref := range refs {
		if ref.Name().IsBranch() || ref.Name().IsTag() {
			branches = append(branches, ref)
		}
	}
	sortBySemVer(branches)
	log.Info().Msgf("Remote branches and tags found: %v for repo %s", branches, repoURL)
	return branches
}

func (m *RemoteBranch) CreateBranch(sourceBranch, targetBranch *plumbing.Reference, createTag bool) error {
	repoURL := m.gitClient.Config().URLs[0]
	// Clone repo temp
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(sourceBranch.Name().Short()),
	})
	if err != nil {
		log.Err(err).Msg("")
		os.Exit(1)
	}
	if m.repo == nil {
		m.AddRepo(r)
	}

	// Create new branch
	ref := plumbing.NewHashReference(targetBranch.Name(), sourceBranch.Hash())
	// The created reference is saved in the storage.
	err = r.Storer.SetReference(ref)
	if err != nil {
		return err
	}
	err = m.repo.Push(&git.PushOptions{})
	if err != nil {
		return err
	}
	log.Info().Msgf("Successfully created branch %s", targetBranch.Name().Short())

	if createTag {
		tag := VersionRegex.FindString(targetBranch.Name().Short())
		_, err = r.CreateTag(tag, sourceBranch.Hash(), &git.CreateTagOptions{
			Message: tag,
			Tagger:  &object.Signature{When: time.Now()},
		})
		if err != nil {
			return err
		}
		err = m.repo.Push(&git.PushOptions{RefSpecs: []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")}})
		if err != nil {
			return err
		}
		log.Info().Msgf("Successfully created tag: %s", tag)
	}

	return nil
}

func sortBySemVer(s []*plumbing.Reference) {
	sort.SliceStable(s, func(i, j int) bool {
		branchA := semver.Canonical(VersionRegex.FindString(s[i].Name().Short()))
		branchB := semver.Canonical(VersionRegex.FindString(s[j].Name().Short()))

		switch semver.Compare(branchA, branchB) {
		case -1:
			return true
		case 0:
			return false
		default:
			return false
		}
	})
}
