package cmd

import (
	"testing"

	"github.com/fhopfensperger/git-releaser/pkg/repo"
)

func Test_createNewReleaseVersion(t *testing.T) {

	type args struct {
		repoUrl string
		force   bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "auth required",
			args:    args{"https://github.com/fhopfensperger/amqp-sb-client.git", false},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test repo doesnt exists",
			args:    args{"https://github.com/fhopfensperger/i-do-not-exist.git", false},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test master branch doesnt exists",
			args:    args{"https://github.com/fhopfensperger/git-releaser.git", false},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test master branch doesnt exists",
			args:    args{"https://github.com/fhopfensperger/git-releaser.git", true},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceBranch = "master"
			got, err := createNewReleaseVersion(tt.args.repoUrl, tt.args.force)
			if (err != nil) != tt.wantErr {
				t.Errorf("createNewReleaseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createNewReleaseVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setNextVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "MAJOR",
			args: args{version: "MAJOR"},
			want: repo.MAJOR,
		},
		{
			name: "MINOR",
			args: args{version: "MINOR"},
			want: repo.MINOR,
		},
		{
			name: "PATCH",
			args: args{version: "PATCH"},
			want: repo.PATCH,
		},
		{
			name: "DEFAULT",
			args: args{version: "IdontKnow"},
			want: repo.MINOR,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setNextVersion(tt.args.version); got != tt.want {
				t.Errorf("setNextVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
