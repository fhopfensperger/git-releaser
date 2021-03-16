package remote

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/mock"
)

type gitRepoMock struct {
	mock.Mock
}

func (m *gitRepoMock) Push(options *git.PushOptions) error {
	fmt.Println("Mocked Push() function")
	return nil
}

func (m *gitRepoMock) List(o *git.ListOptions) (rfs []*plumbing.Reference, err error) {
	fmt.Println("Mocked List() function")
	args := m.Called(o)
	return args.Get(0).([]*plumbing.Reference), args.Error(1)
}

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

func TestGitRepo_GetAllRemoteBranchesAndTags(t *testing.T) {
	gitRemoteRepo := new(gitRepoMock)

	gitRepo := GitRepo{remote: gitRemoteRepo}

	gitRemoteRepo.On("List", &git.ListOptions{}).Return(append(generateTagsPlumbReferences(), generateBranchPlumbReferences()...), nil)

	refs := gitRepo.GetAllRemoteBranchesAndTags("https://github.com/just-a-repo-name")
	gitRemoteRepo.AssertExpectations(t)

	sortedRefs := []*plumbing.Reference{
		main, dev, test, a, b, e, c, f, g, d,
	}

	assert.Equal(t, refs, sortedRefs)
}

func TestGitRepo_GetAllRemoteBranchesAndTags_TagsOnly(t *testing.T) {
	gitRemoteRepo := new(gitRepoMock)

	gitRepo := GitRepo{remote: gitRemoteRepo}

	gitRemoteRepo.On("List", &git.ListOptions{}).Return(append(generateTagsPlumbReferences(), main, test, dev), nil)

	refs := gitRepo.GetAllRemoteBranchesAndTags("https://github.com/just-a-repo-name")
	gitRemoteRepo.AssertExpectations(t)

	sortedRefs := []*plumbing.Reference{
		main, test, dev, e, f, g,
	}

	assert.Equal(t, refs, sortedRefs)
}

func TestGitRepo_GetAllRemoteBranchesAndTags_BranchesOnly(t *testing.T) {
	gitRemoteRepo := new(gitRepoMock)

	gitRepo := GitRepo{remote: gitRemoteRepo}

	gitRemoteRepo.On("List", &git.ListOptions{}).Return(append(generateTagsPlumbReferences(), main, test, dev), nil)

	refs := gitRepo.GetAllRemoteBranchesAndTags("https://github.com/just-a-repo-name")
	gitRemoteRepo.AssertExpectations(t)

	sortedRefs := []*plumbing.Reference{
		main, test, dev, e, f, g,
	}

	assert.Equal(t, refs, sortedRefs)
}

func TestGitRepo_CreateBranchAndTag(t *testing.T) {
	gitRemoteRepo := new(gitRepoMock)
	m := GitRepo{remote: gitRemoteRepo, storer: memory.NewStorage()}

	gitRemoteRepo.On("Push", &git.PushOptions{}).Return(nil)

	type args struct {
		sourceBranch *plumbing.Reference
		targetBranch string
		version      string
		createBranch bool
		createTag    bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{main, "release", "v1.0.1", true, false}, false},
		{"test2", args{main, "release", "v1.0.2", false, true}, false},
		{"test3", args{main, "release", "v1.0.2", true, true}, false},
		{"test4", args{main, "release", "v1.0.2", false, false}, false},
		{"test5", args{main, "", "v1.0.2", true, false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.CreateBranchAndTag(tt.args.sourceBranch, tt.args.targetBranch, tt.args.version, tt.args.createBranch, tt.args.createTag); (err != nil) != tt.wantErr {
				t.Errorf("CreateBranchAndTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sortBySemVer(t *testing.T) {

	tests := []struct {
		name string
		s    []*plumbing.Reference
		want []*plumbing.Reference
	}{
		{"test1",
			[]*plumbing.Reference{a, main, b, c, d},
			[]*plumbing.Reference{main, a, b, c, d},
		},
		{"test2",
			[]*plumbing.Reference{d, a, b, main, c},
			[]*plumbing.Reference{main, a, b, c, d},
		},
		{"test3",
			[]*plumbing.Reference{d, c, main, b, a},
			[]*plumbing.Reference{main, a, b, c, d},
		},
		{"test4",
			[]*plumbing.Reference{main, d, d, a, a},
			[]*plumbing.Reference{main, a, a, d, d},
		},
		{"test5",
			[]*plumbing.Reference{a, b, c, d, e, f, g, main},
			[]*plumbing.Reference{main, a, b, c, e, f, d, g},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortBySemVer(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortBySemVer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitRepo_GetStorer(t *testing.T) {
	gitRemoteRepo := new(gitRepoMock)
	stor := memory.NewStorage()
	m := GitRepo{remote: gitRemoteRepo, storer: stor}
	assert.Equal(t, stor, m.GetStorer())
}
