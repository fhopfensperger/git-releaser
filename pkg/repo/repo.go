package repo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/fhopfensperger/git-releaser/pkg/remote"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
	"golang.org/x/mod/semver"
)

const (
	MAJOR = iota // MAJOR == 0
	MINOR = iota // MINOR == 1
	PATCH = iota // PATCH == 2
)

type Repo struct {
	remoteUrl     string
	allReferences []*plumbing.Reference
	// typically the `main` branch
	sourceBranch *plumbing.Reference
	// typically the `release` branches
	versionBranches []*plumbing.Reference
	// Latest released version, either branch or tag
	latestVersionReference *plumbing.Reference
	versionTags            []*plumbing.Reference
	nextReleaseVersion     string
	remoteBranch           remote.CreateBranchAndTager
	branchFilter           string
}

func New(remoteUrl string) *Repo {
	r := Repo{}
	r.remoteUrl = remoteUrl
	r.remoteBranch = &remote.GitRepo{}
	r.allReferences = r.remoteBranch.GetAllRemoteBranchesAndTags(remoteUrl)
	return &r
}

func (r *Repo) GetVersionBranches(branchFilter string) []*plumbing.Reference {
	r.branchFilter = branchFilter
	for _, b := range r.allReferences {
		if b.Name().IsBranch() && strings.Contains(b.Name().Short(), branchFilter) && remote.VersionRegex.FindString(b.Name().Short()) != "" {
			r.versionBranches = append(r.versionBranches, b)
		}
	}
	return r.versionBranches
}

func (r *Repo) GetVersionTags() []*plumbing.Reference {
	for _, b := range r.allReferences {
		if b.Name().IsTag() && remote.VersionRegex.FindString(b.Name().Short()) != "" {
			r.versionTags = append(r.versionTags, b)
		}
	}
	return r.versionTags
}

func (r *Repo) GetLatestVersionReference() *plumbing.Reference {
	var latestBranchVersion, latestTagVersion string
	var latestBranch, latestTag *plumbing.Reference

	if len(r.versionBranches) > 0 {
		latestBranch = r.versionBranches[len(r.versionBranches)-1]
		latestBranchVersion = remote.VersionRegex.FindString(latestBranch.Name().Short())
	}
	if len(r.versionTags) > 0 {
		latestTag = r.versionTags[len(r.versionTags)-1]
		latestTagVersion = remote.VersionRegex.FindString(latestTag.Name().Short())
	}

	if latestBranch == nil && latestTag == nil {
		return nil
	}

	switch semver.Compare(latestBranchVersion, latestTagVersion) {
	case -1:
		r.latestVersionReference = latestTag
		return latestTag
	case 0:
		r.latestVersionReference = latestBranch
		return latestBranch
	default:
		r.latestVersionReference = latestBranch
		return latestBranch
	}
}

func (r *Repo) GetSourceBranch(name string) *plumbing.Reference {
	for _, ref := range r.allReferences {
		if ref.Name().Short() == name {
			r.sourceBranch = ref
			return ref
		}
	}
	return nil
}

func (r *Repo) NextReleaseVersion(nextVersion int) (string, error) {
	if r.latestVersionReference == nil {
		r.nextReleaseVersion = fallBackVersion(nextVersion)
		return r.nextReleaseVersion, nil
	}

	semLatestVersion := semver.Canonical(remote.VersionRegex.FindString(r.latestVersionReference.Name().Short()))
	latestVersionSlice := strings.Split(semLatestVersion, ".")

	if len(latestVersionSlice) == 1 {
		r.nextReleaseVersion = fallBackVersion(nextVersion)
		return r.nextReleaseVersion, nil
	}

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
	switch nextVersion {
	case MAJOR:
		r.nextReleaseVersion = semver.Canonical(fmt.Sprintf("v%v.%v.%v", major+1, 0, 0))
	case MINOR:
		r.nextReleaseVersion = semver.Canonical(fmt.Sprintf("v%v.%v.%v", major, minor+1, 0))
	case PATCH:
		r.nextReleaseVersion = semver.Canonical(fmt.Sprintf("v%v.%v.%v", major, minor, patch+1))
	default:
		r.nextReleaseVersion = fallBackVersion(nextVersion)
	}
	return r.nextReleaseVersion, nil
}

func fallBackVersion(nextVersion int) string {
	switch nextVersion {
	case MAJOR:
		return semver.Canonical(fmt.Sprintf("v%v.%v.%v", 1, 0, 0))
	case MINOR:
		return semver.Canonical(fmt.Sprintf("v%v.%v.%v", 0, 1, 0))
	case PATCH:
		return semver.Canonical(fmt.Sprintf("v%v.%v.%v", 0, 0, 1))
	default:
		return semver.Canonical(fmt.Sprintf("v%v.%v.%v", 0, 1, 0))
	}
}

func (r *Repo) CreateNewRelease(branch, tag bool) error {
	if r.latestVersionReference == nil {
		log.Info().Msg("No current version branches / tags found")
		return r.remoteBranch.CreateBranchAndTag(r.sourceBranch, r.branchFilter, r.nextReleaseVersion, branch, tag)
	}

	if r.latestVersionReference.Name().IsTag() {
		if r.latestVersionReference.Hash().String() == r.sourceBranch.Hash().String() {
			log.Info().Msgf("Nothing to do, %s branch and latest tag version %s are equals, commit hash: %s", r.sourceBranch.Name().Short(), r.latestVersionReference.Name().Short(), r.sourceBranch.Hash())
			return nil
		}
	}

	if r.latestVersionReference.Hash().String() == r.sourceBranch.Hash().String() {
		log.Info().Msgf("Nothing to do, %s and latest branch version %s are equals, commit hash: %s", r.sourceBranch.Name().Short(), r.latestVersionReference.Name().Short(), r.sourceBranch.Hash())
		return nil
	}

	return r.remoteBranch.CreateBranchAndTag(r.sourceBranch, r.branchFilter, r.nextReleaseVersion, branch, tag)
}
