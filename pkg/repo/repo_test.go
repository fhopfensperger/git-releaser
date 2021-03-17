package repo

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-git/v5/storage"

	"github.com/stretchr/testify/mock"

	"github.com/fhopfensperger/git-releaser/pkg/remote"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
)

var (
	a    = plumbing.NewHashReference(plumbing.NewBranchReferenceName("release/v0.100.0"), plumbing.NewHash("a48656c6c6f20476f7068657221"))
	b    = plumbing.NewHashReference(plumbing.NewBranchReferenceName("release/v1.0.9"), plumbing.NewHash("b48656c6c6f20476f7068657221"))
	c    = plumbing.NewHashReference(plumbing.NewBranchReferenceName("release/v1.10.9"), plumbing.NewHash("c48656c6c6f20476f7068657221"))
	d    = plumbing.NewHashReference(plumbing.NewBranchReferenceName("release/v2.10.79"), plumbing.NewHash("d48656c6c6f20476f7068657221"))
	main = plumbing.NewHashReference(plumbing.NewBranchReferenceName("main"), plumbing.NewHash("aa48656c6c6f20476f7068657221"))
	e    = plumbing.NewHashReference(plumbing.NewTagReferenceName("v1.10.9"), plumbing.NewHash("e48656c6c6f20476f7068657221"))
	f    = plumbing.NewHashReference(plumbing.NewTagReferenceName("v2.1.79"), plumbing.NewHash("f48656c6c6f20476f7068657221"))
	g    = plumbing.NewHashReference(plumbing.NewTagReferenceName("v2.10.79"), plumbing.NewHash("ff48656c6c6f20476f7068657221"))

	dev  = plumbing.NewHashReference(plumbing.NewBranchReferenceName("dev"), plumbing.NewHash("afd48656c6c6f20476f7068657221"))
	test = plumbing.NewHashReference(plumbing.NewBranchReferenceName("testing"), plumbing.NewHash("afad48656c6c6f20476f7068657221"))
)

func generateBranchPlumbReferences() []*plumbing.Reference {
	return []*plumbing.Reference{a, b, c, d, main, dev, test}
}

func generateTagsPlumbReferences() []*plumbing.Reference {
	return []*plumbing.Reference{e, f, g}
}

func TestRepo_GetVersionBranches(t *testing.T) {
	type fields struct {
		remoteUrl              string
		allReferences          []*plumbing.Reference
		sourceBranch           *plumbing.Reference
		versionBranches        []*plumbing.Reference
		latestVersionReference *plumbing.Reference
		versionTags            []*plumbing.Reference
		nextReleaseVersion     string
		remoteBranch           *remote.GitRepo
		branchFilter           string
	}
	type args struct {
		branchFilter string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*plumbing.Reference
	}{
		{
			name: "test1",
			fields: fields{
				remoteUrl:     "",
				allReferences: append(generateBranchPlumbReferences(), generateTagsPlumbReferences()...),
				sourceBranch:  main,
			},
			args: args{branchFilter: ""},
			want: []*plumbing.Reference{a, b, c, d},
		},
		{
			name: "test2",
			fields: fields{
				remoteUrl:     "",
				allReferences: generateTagsPlumbReferences(),
				sourceBranch:  main,
			},
			args: args{branchFilter: ""},
			want: nil,
		},
		{
			name: "test3",
			fields: fields{
				remoteUrl:     "",
				allReferences: []*plumbing.Reference{main, dev, test},
				sourceBranch:  main,
			},
			args: args{branchFilter: ""},
			want: nil,
		},
		{
			name: "test4",
			fields: fields{
				remoteUrl:     "",
				allReferences: append(generateBranchPlumbReferences(), generateTagsPlumbReferences()...),
				sourceBranch:  main,
			},
			args: args{branchFilter: "NOT_RELEASE"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				remoteUrl:              tt.fields.remoteUrl,
				allReferences:          tt.fields.allReferences,
				sourceBranch:           tt.fields.sourceBranch,
				versionBranches:        tt.fields.versionBranches,
				latestVersionReference: tt.fields.latestVersionReference,
				versionTags:            tt.fields.versionTags,
				nextReleaseVersion:     tt.fields.nextReleaseVersion,
				remoteBranch:           tt.fields.remoteBranch,
				branchFilter:           tt.fields.branchFilter,
			}
			assert.Equal(t, tt.want, r.GetVersionBranches(tt.args.branchFilter))
		})
	}
}

func TestRepo_GetVersionTags(t *testing.T) {
	type fields struct {
		remoteUrl              string
		allReferences          []*plumbing.Reference
		sourceBranch           *plumbing.Reference
		versionBranches        []*plumbing.Reference
		latestVersionReference *plumbing.Reference
		versionTags            []*plumbing.Reference
		nextReleaseVersion     string
		remoteBranch           *remote.GitRepo
		branchFilter           string
	}
	tests := []struct {
		name   string
		fields fields
		want   []*plumbing.Reference
	}{
		{
			name: "test1",
			fields: fields{
				remoteUrl:     "",
				allReferences: append(generateBranchPlumbReferences(), generateTagsPlumbReferences()...),
				sourceBranch:  main,
			},
			want: generateTagsPlumbReferences(),
		},
		{
			name: "test1",
			fields: fields{
				remoteUrl:     "",
				allReferences: generateBranchPlumbReferences(),
				sourceBranch:  main,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				remoteUrl:              tt.fields.remoteUrl,
				allReferences:          tt.fields.allReferences,
				sourceBranch:           tt.fields.sourceBranch,
				versionBranches:        tt.fields.versionBranches,
				latestVersionReference: tt.fields.latestVersionReference,
				versionTags:            tt.fields.versionTags,
				nextReleaseVersion:     tt.fields.nextReleaseVersion,
				remoteBranch:           tt.fields.remoteBranch,
				branchFilter:           tt.fields.branchFilter,
			}
			assert.Equal(t, tt.want, r.GetVersionTags())
		})
	}
}

func TestRepo_GetLatestVersionReference(t *testing.T) {
	type fields struct {
		remoteUrl          string
		allReferences      []*plumbing.Reference
		sourceBranch       *plumbing.Reference
		nextReleaseVersion string
		remoteBranch       *remote.GitRepo
		branchFilter       string
	}
	tests := []struct {
		name   string
		fields fields
		want   *plumbing.Reference
	}{
		{
			name: "Only Branches",
			fields: fields{
				remoteUrl:     "",
				allReferences: generateBranchPlumbReferences(),
				sourceBranch:  main,
			},
			want: d,
		},
		{
			name: "Branches and Tags (Branch will be used if branch==tag)",
			fields: fields{
				remoteUrl:     "",
				allReferences: append(generateBranchPlumbReferences(), generateTagsPlumbReferences()...),
				sourceBranch:  main,
			},
			want: d,
		},
		{
			name: "Tag the latest",
			fields: fields{
				remoteUrl:     "",
				allReferences: append(generateTagsPlumbReferences(), a, b, c),
				sourceBranch:  main,
			},
			want: g,
		},
		{
			name: "Only Tags",
			fields: fields{
				remoteUrl:     "",
				allReferences: generateTagsPlumbReferences(),
				sourceBranch:  main,
			},
			want: g,
		},
		{
			name: "No Branch found",
			fields: fields{
				remoteUrl:    "",
				sourceBranch: main,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				remoteUrl:          tt.fields.remoteUrl,
				allReferences:      tt.fields.allReferences,
				sourceBranch:       tt.fields.sourceBranch,
				nextReleaseVersion: tt.fields.nextReleaseVersion,
				remoteBranch:       tt.fields.remoteBranch,
				branchFilter:       tt.fields.branchFilter,
			}
			r.GetVersionBranches("")
			r.GetVersionTags()
			assert.Equal(t, tt.want, r.GetLatestVersionReference())
		})
	}
}

func TestRepo_GetSourceBranch(t *testing.T) {
	type fields struct {
		remoteUrl              string
		allReferences          []*plumbing.Reference
		sourceBranch           *plumbing.Reference
		versionBranches        []*plumbing.Reference
		latestVersionReference *plumbing.Reference
		versionTags            []*plumbing.Reference
		nextReleaseVersion     string
		remoteBranch           *remote.GitRepo
		branchFilter           string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *plumbing.Reference
	}{
		{
			name: "Main branch",
			fields: fields{
				remoteUrl:     "",
				allReferences: append(generateBranchPlumbReferences(), generateTagsPlumbReferences()...),
			},
			args: args{name: "main"},
			want: main,
		},
		{
			name: "Dev branch",
			fields: fields{
				remoteUrl:     "",
				allReferences: append(generateBranchPlumbReferences(), generateTagsPlumbReferences()...),
			},
			args: args{name: "dev"},
			want: dev,
		},
		{
			name: "No branch found for tags only (unrealistic)",
			fields: fields{
				remoteUrl:     "",
				allReferences: generateTagsPlumbReferences(),
			},
			args: args{name: "main"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				remoteUrl:              tt.fields.remoteUrl,
				allReferences:          tt.fields.allReferences,
				sourceBranch:           tt.fields.sourceBranch,
				versionBranches:        tt.fields.versionBranches,
				latestVersionReference: tt.fields.latestVersionReference,
				versionTags:            tt.fields.versionTags,
				nextReleaseVersion:     tt.fields.nextReleaseVersion,
				remoteBranch:           tt.fields.remoteBranch,
				branchFilter:           tt.fields.branchFilter,
			}
			if got := r.GetSourceBranch(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSourceBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_NextReleaseVersion(t *testing.T) {
	type fields struct {
		latestVersionReference *plumbing.Reference
	}
	type args struct {
		nextVersion int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{a},
			args:    args{MAJOR},
			want:    "v1.0.0",
			wantErr: false,
		},
		{
			name:    "test2",
			fields:  fields{a},
			args:    args{MINOR},
			want:    "v0.101.0",
			wantErr: false,
		},
		{
			name:    "test3",
			fields:  fields{a},
			args:    args{PATCH},
			want:    "v0.100.1",
			wantErr: false,
		},
		{
			name:    "default, as no latest version found",
			args:    args{PATCH},
			want:    "v0.0.1",
			wantErr: false,
		},
		{
			name:    "invalid version branch reference",
			fields:  fields{plumbing.NewHashReference("vmain.main.main", plumbing.Hash{})},
			args:    args{MINOR},
			want:    "v0.1.0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				latestVersionReference: tt.fields.latestVersionReference,
			}
			got, err := r.NextReleaseVersion(tt.args.nextVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextReleaseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NextReleaseVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type repoMock struct {
	mock.Mock
}

func (m *repoMock) CreateBranchAndTag(sourceBranch *plumbing.Reference, targetBranch, version string, createBranch, createTag bool) error {
	fmt.Println("Mocked CreateBranchAndTag() function")
	args := m.Called(sourceBranch, targetBranch, version, createBranch, createTag)
	return args.Error(0)
}

func (m *repoMock) GetAllRemoteBranchesAndTags(repoURL string) []*plumbing.Reference {
	fmt.Println("Mocked GetAllRemoteBranchesAndTags() function")
	args := m.Called(repoURL)
	return args.Get(0).([]*plumbing.Reference)
}

func (m *repoMock) GetStorer() storage.Storer {
	fmt.Println("Mocked GetStorer() function")
	args := m.Called()
	return args.Get(0).(storage.Storer)
}

func TestRepo_CreateNewRelease(t *testing.T) {
	tag_v1_0_0 := plumbing.NewHashReference(plumbing.NewTagReferenceName("v1.0.0"), plumbing.NewHash("40448c70cf1ac313d22aa2b2454ca68baa122542"))
	tag_v2_0_0 := plumbing.NewHashReference(plumbing.NewTagReferenceName("v2.0.0"), plumbing.NewHash("12448"))

	remoteBranchMock := new(repoMock)
	remoteBranchMock.On("CreateBranchAndTag", main, "release", "v1.0.0", false, false).Return(nil)
	remoteBranchMock.On("CreateBranchAndTag", main, "release", "v1.0.0", true, false).Return(nil)
	remoteBranchMock.On("CreateBranchAndTag", main, "release", "v1.0.0", true, true).Return(nil)
	remoteBranchMock.On("CreateBranchAndTag", main, "release", "v1.0.0", false, true).Return(nil)
	remoteBranchMock.On("CreateBranchAndTag", main, "release", "v1.0.1", true, true).Return(nil)

	// Create mock storage commit and associated tag...
	stor := memory.NewStorage()

	commit := object.Commit{Message: "commit123", Hash: plumbing.NewHash("123abc")}
	tag := object.Tag{Name: tag_v1_0_0.Name().String(), Hash: tag_v1_0_0.Hash(), TargetType: plumbing.CommitObject, Target: commit.Hash}

	co := stor.NewEncodedObject()
	commit.EncodeWithoutSignature(co)

	eo := stor.NewEncodedObject()
	tag.EncodeWithoutSignature(eo)

	stor.Objects[commit.Hash] = co
	stor.Objects[tag_v1_0_0.Hash()] = eo

	// main branch & commit and tag same hash
	hash := plumbing.NewHash("a0dacb3d48b64358760871c73a02b6c4962a9d28")
	mainCommit1 := plumbing.NewHashReference(plumbing.NewBranchReferenceName("main"), plumbing.NewHash(hash.String()))
	commit1 := object.Commit{Hash: hash}
	tag1 := object.Tag{Name: tag_v2_0_0.Name().String(), Hash: tag_v2_0_0.Hash(), TargetType: plumbing.CommitObject, Target: hash}
	remoteBranchMock.On("CreateBranchAndTag", mainCommit1, "release", "v2.0.1", true, true).Return(nil)

	co1 := stor.NewEncodedObject()
	commit1.EncodeWithoutSignature(co1)

	eo1 := stor.NewEncodedObject()
	tag1.EncodeWithoutSignature(eo1)

	stor.Objects[commit1.Hash] = co1
	stor.Objects[tag_v2_0_0.Hash()] = eo1

	remoteBranchMock.On("GetStorer").Return(stor)

	type fields struct {
		remoteUrl              string
		allReferences          []*plumbing.Reference
		sourceBranch           *plumbing.Reference
		versionBranches        []*plumbing.Reference
		latestVersionReference *plumbing.Reference
		versionTags            []*plumbing.Reference
		nextReleaseVersion     string
		remoteBranch           remote.CreateBranchAndTager
		branchFilter           string
	}
	type args struct {
		branch bool
		tag    bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "invalid version branch reference",
			fields: fields{
				sourceBranch:           main,
				latestVersionReference: nil,
				branchFilter:           "release",
				nextReleaseVersion:     "v1.0.0",
				remoteBranch:           remoteBranchMock,
			},
			args: args{
				branch: false,
				tag:    false,
			},
			wantErr: false,
		},
		{
			name: "test1",
			fields: fields{
				sourceBranch:       main,
				branchFilter:       "release",
				nextReleaseVersion: "v1.0.0",
				remoteBranch:       remoteBranchMock,
			},
			args: args{
				branch: true,
				tag:    false,
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				sourceBranch:       main,
				branchFilter:       "release",
				nextReleaseVersion: "v1.0.0",
				remoteBranch:       remoteBranchMock,
			},
			args: args{
				branch: true,
				tag:    false,
			},
			wantErr: false,
		},
		{
			name: "test3",
			fields: fields{
				sourceBranch:       main,
				branchFilter:       "release",
				nextReleaseVersion: "v1.0.0",
				remoteBranch:       remoteBranchMock,
			},
			args: args{
				branch: false,
				tag:    false,
			},
			wantErr: false,
		},
		{
			name: "test4",
			fields: fields{
				sourceBranch:       main,
				branchFilter:       "release",
				nextReleaseVersion: "v1.0.0",
				remoteBranch:       remoteBranchMock,
			},
			args: args{
				branch: false,
				tag:    true,
			},
			wantErr: false,
		},
		{
			name: "test5",
			fields: fields{
				sourceBranch:           main,
				latestVersionReference: a,
				branchFilter:           "release",
				nextReleaseVersion:     "v1.0.0",
				remoteBranch:           remoteBranchMock,
			},
			args: args{
				branch: false,
				tag:    true,
			},
			wantErr: false,
		},
		{
			name: "test6",
			fields: fields{
				sourceBranch:           main,
				latestVersionReference: tag_v1_0_0,
				branchFilter:           "release",
				nextReleaseVersion:     "v1.0.1",
				remoteBranch:           remoteBranchMock,
			},
			args: args{
				branch: true,
				tag:    true,
			},
			wantErr: false,
		},
		{
			name: "test7",
			fields: fields{
				sourceBranch:           mainCommit1,
				latestVersionReference: tag_v2_0_0,
				branchFilter:           "release",
				nextReleaseVersion:     "v2.0.1",
				remoteBranch:           remoteBranchMock,
			},
			args: args{
				branch: true,
				tag:    true,
			},
			wantErr: false,
		},
		{
			name: "test8",
			fields: fields{
				sourceBranch:           main,
				latestVersionReference: main,
				branchFilter:           "release",
				nextReleaseVersion:     "v2.0.1",
				remoteBranch:           remoteBranchMock,
			},
			args: args{
				branch: true,
				tag:    true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				remoteUrl:              tt.fields.remoteUrl,
				allReferences:          tt.fields.allReferences,
				sourceBranch:           tt.fields.sourceBranch,
				versionBranches:        tt.fields.versionBranches,
				latestVersionReference: tt.fields.latestVersionReference,
				versionTags:            tt.fields.versionTags,
				nextReleaseVersion:     tt.fields.nextReleaseVersion,
				remoteBranch:           tt.fields.remoteBranch,
				branchFilter:           tt.fields.branchFilter,
			}
			if err := r.CreateNewRelease(tt.args.branch, tt.args.tag); (err != nil) != tt.wantErr {
				t.Errorf("CreateNewRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fallBackVersion(t *testing.T) {
	type args struct {
		nextVersion int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Initial Major version",
			args: args{MAJOR},
			want: "v1.0.0",
		},
		{
			name: "Initial Minor version",
			args: args{MINOR},
			want: "v0.1.0",
		},
		{
			name: "Initial Patch version",
			args: args{PATCH},
			want: "v0.0.1",
		},
		{
			name: "Initial default version",
			args: args{123},
			want: "v0.1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fallBackVersion(tt.args.nextVersion); got != tt.want {
				t.Errorf("fallBackVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
