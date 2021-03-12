# git-releaser
![Go](https://github.com/fhopfensperger/git-releaser/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fhopfensperger/git-releaser)](https://goreportcard.com/report/github.com/fhopfensperger/git-releaser)
[![Coverage Status](https://coveralls.io/repos/github/fhopfensperger/git-releaser/badge.svg?branch=master)](https://coveralls.io/github/fhopfensperger/git-releaser?branch=master)
[![Release](https://img.shields.io/github/release/fhopfensperger/git-releaser.svg?style=flat-square)](https://github.com//fhopfensperger/git-releaser/releases/latest)


Creates new release branches

## Installation

### Option 1 (script)

```bash
curl https://raw.githubusercontent.com/fhopfensperger/git-releaser/master/get.sh | bash
```

### Option 2 (manually)

Go to [Releases](https://github.com/fhopfensperger/git-releaser/releases) download the latest release according to your processor architecture and operating system, then unarchive and copy it to the right location

```bash
tar xvfz git-releaser_x.x.x_darwin_amd64.tar.gz
cd git-releaser_x.x.x_darwin_amd64
chmod +x git-releaser
sudo mv git-releaser /usr/local/bin/
```
