package cmd

import (
	"testing"
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
