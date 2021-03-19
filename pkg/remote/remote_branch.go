package remote

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

type GitRemoter interface {
	List(o *git.ListOptions) (rfs []*plumbing.Reference, err error)
	Push(o *git.PushOptions) error
}
type GitRepo struct {
	remote GitRemoter
	storer storage.Storer
	Auth   transport.AuthMethod
}

func (m *GitRepo) GetStorer() storage.Storer {
	return m.storer
}

//GetRemoteBranches get remote branches from GitHub using the repoURL
func (m *GitRepo) GetAllRemoteBranchesAndTags(repoURL string) []*plumbing.Reference {
	var stor *memory.Storage
	if m.storer == nil {
		stor = memory.NewStorage()
		m.storer = stor
	}
	if m.remote == nil {
		rem := git.NewRemote(m.storer, &config.RemoteConfig{
			Name: "origin",
			URLs: []string{repoURL},
		})
		m.remote = rem
	}

	// We can then use every Remote functions to retrieve wanted information
	refs, err := m.remote.List(&git.ListOptions{Auth: m.Auth})
	if err != nil {
		log.Err(err).Msg("")
	}

	// Filters the references list and only keeps tags
	var branches []*plumbing.Reference
	var tags []*plumbing.Reference

	for _, ref := range refs {
		eo := m.storer.NewEncodedObject()
		if ref.Name().IsTag() {
			err := m.storer.SetReference(ref)
			if err != nil {
				log.Err(err).Msg("")
				continue
			}
			tags = append(tags, ref)
		} else if ref.Name().IsBranch() {
			err := m.storer.SetReference(ref)
			if err != nil {
				log.Err(err).Msg("")
				continue
			}
			commit := object.Commit{Hash: ref.Hash()}
			err = commit.EncodeWithoutSignature(eo)
			if err != nil {
				log.Err(err).Msg("")
				continue
			}
			stor.Objects[ref.Hash()] = eo
			branches = append(branches, ref)
		}
	}
	branchesAndTags := append(tags, branches...)
	sortBySemVer(branchesAndTags)
	log.Info().Msgf("Remote branches and tags found: %v for repo %s", branchesAndTags, repoURL)

	return branchesAndTags
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
			log.Err(err).Msg("")
			return err
		}
		refspec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)
		err = m.remote.Push(&git.PushOptions{RefSpecs: []config.RefSpec{config.RefSpec(refspec)}, Auth: m.Auth})
		if err != nil {
			log.Err(err).Msg("")
			return err
		}
		log.Info().Msgf("Successfully created branch %s/%s", targetBranch, version)
	}

	if createTag {
		ref := plumbing.NewHashReference(plumbing.NewTagReferenceName(version), sourceBranch.Hash())
		err := m.storer.SetReference(ref)
		if err != nil {
			log.Err(err).Msg("")
			return err
		}
		refspecs := fmt.Sprintf("refs/tags/%s:refs/tags/%s", version, version)
		err = m.remote.Push(&git.PushOptions{RefSpecs: []config.RefSpec{config.RefSpec(refspecs)}})
		if err != nil {
			log.Err(err).Msg("")
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
