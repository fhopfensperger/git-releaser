package remote

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"

	"golang.org/x/mod/semver"

	"github.com/rs/zerolog/log"

	"github.com/go-git/go-git/v5/config"

	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-git/v5"
)

var VersionRegex = regexp.MustCompile(`v\d+(\.\d+)+`)

type CreateBranchAndTager interface {
	CreateBranchAndTag(*plumbing.Reference, string, string, bool, bool) error
	GetAllRemoteBranchesAndTags(repoURL string) []*plumbing.Reference
	GetStorer() storage.Storer
}

// GitRepoer
type GitRepoer interface {
	Push(*git.PushOptions) error
	Tags() (storer.ReferenceIter, error)
	CreateTag(string, plumbing.Hash, *git.CreateTagOptions) (*plumbing.Reference, error)
}

type GitRemoter interface {
	List(o *git.ListOptions) (rfs []*plumbing.Reference, err error)
}
type GitRepo struct {
	repo   GitRepoer
	remote GitRemoter
	storer storage.Storer
}

//AddRepo to add a remote repo to RemoteBranch struct
func (m *GitRepo) AddRepo(repo GitRepoer) {
	m.repo = repo
}

func (m *GitRepo) GetRepo() GitRepoer {
	return m.repo
}

func (m *GitRepo) GetStorer() storage.Storer {
	return m.storer
}

//GetRemoteBranches get remote branches from GitHub using the repoURL
func (m *GitRepo) GetAllRemoteBranchesAndTags(repoURL string) []*plumbing.Reference {
	if m.storer == nil {
		m.storer = memory.NewStorage()
	}

	if m.remote == nil {
		rem := git.NewRemote(m.storer, &config.RemoteConfig{
			Name: "origin",
			URLs: []string{repoURL},
		})
		m.remote = rem
	}

	if m.repo == nil {
		r, err := git.Clone(m.storer, nil, &git.CloneOptions{
			URL: repoURL,
		})
		if err != nil {
			log.Err(err).Msg("")
			os.Exit(1)
		}
		m.repo = r
	}

	var refs []*plumbing.Reference
	references, err := m.remote.List(&git.ListOptions{})
	if err != nil {
		return nil
	}
	for _, ref := range references {
		if ref.Name().IsBranch() {
			refs = append(refs, ref)
		}
	}

	tagIter, err := m.repo.Tags()
	if err != nil {
		log.Err(err).Msg("")
	}

	err = tagIter.ForEach(func(r *plumbing.Reference) error {
		refs = append(refs, r)
		return nil
	})
	if err != nil {
		return nil
	}

	sortBySemVer(refs)
	log.Info().Msgf("Remote branches and tags found: %v for repo %s", refs, repoURL)
	return refs
}

func (m *GitRepo) CreateBranchAndTag(sourceBranch *plumbing.Reference, targetBranch, version string, createBranch, createTag bool) error {

	if createBranch {
		// Create new branch
		var branchName string
		if targetBranch == "" {
			branchName = version
		} else {
			branchName = fmt.Sprintf("%s/%s", targetBranch, version)
		}
		ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branchName), sourceBranch.Hash())
		// The created reference is saved in the storage.
		err := m.storer.SetReference(ref)
		if err != nil {
			return err
		}
		err = m.repo.Push(&git.PushOptions{})
		if err != nil {
			return err
		}
		log.Info().Msgf("Successfully created branch %s/%s", targetBranch, version)
	}

	if createTag {
		_, err := m.repo.CreateTag(version, sourceBranch.Hash(), &git.CreateTagOptions{
			Message: version,
			Tagger:  &object.Signature{When: time.Now()},
		})
		if err != nil {
			return err
		}
		err = m.repo.Push(&git.PushOptions{RefSpecs: []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")}})
		if err != nil {
			return err
		}
		log.Info().Msgf("Successfully created tag: %s", version)
	}

	return nil
}

func sortBySemVer(s []*plumbing.Reference) []*plumbing.Reference {
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
	return s
}
