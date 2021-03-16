# git-releaser
![Go](https://github.com/fhopfensperger/git-releaser/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fhopfensperger/git-releaser)](https://goreportcard.com/report/github.com/fhopfensperger/git-releaser)
[![Release](https://img.shields.io/github/release/fhopfensperger/git-releaser.svg?style=flat-square)](https://github.com//fhopfensperger/git-releaser/releases/latest)


This simple command line tool can be used to create release branches with correct versioning, e.g. `release/v1.7.5`. The `-t` flag is used to create a tag in parallel with the branch. Based on the current version of the release branch, the version number of the patch is incremented by one, if the `main` branch is newer than for the latest release. Set the `-n` `--next-flag` flag to release a new PATCH, MINOR or MAJOR version, for example, `-n MINOR` will create a `release/v1.8.0` for `release/v1.7.4`

## Installation

### Option 1 (script)

```bash
curl https://raw.githubusercontent.com/fhopfensperger/git-releaser/main/get.sh | bash
```

### Option 2 (manually)

Go to [Releases](https://github.com/fhopfensperger/git-releaser/releases) download the latest release according to your processor architecture and operating system, then unarchive and copy it to the right location

```bash
tar xvfz git-releaser_x.x.x_darwin_amd64.tar.gz
cd git-releaser_x.x.x_darwin_amd64
chmod +x git-releaser
sudo mv git-releaser /usr/local/bin/
```
