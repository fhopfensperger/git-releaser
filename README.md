# git-releaser
![Go](https://github.com/fhopfensperger/git-releaser/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fhopfensperger/git-releaser)](https://goreportcard.com/report/github.com/fhopfensperger/git-releaser)
[![Coverage Status](https://coveralls.io/repos/github/fhopfensperger/git-releaser/badge.svg?branch=main)](https://coveralls.io/github/fhopfensperger/git-releaser?branch=main)
[![Release](https://img.shields.io/github/release/fhopfensperger/git-releaser.svg?style=flat-square)](https://github.com//fhopfensperger/git-releaser/releases/latest)

This simple command line tool can be used to create a `release branch` and/ or a `tag` with [semantic versioning](https://semver.org) e.g. branch: `release/v1.7.5` and/or tag: `v1.7.5`.
# Usage

`-t` flag is used to create a new version tag

`-c` flag is used to create a new release branch

Based on the version of the latest release `branch` or `tag`, the version number of the patch is incremented by one, if the `-s (--source-brach)` `default: main` branch is newer (based on the commit hash) than for the latest release.

Set the `-n` `--next-version` flag to release a new `PATCH`, `MINOR` or `MAJOR` version, for example, `-n MINOR` will create a `release/v1.8.0` for `release/v1.7.4`

Note: If no version `tag` or `branch` could be found, a new version based on `-n` will be created.

## All flags

```
-c, --create-branch          Create a release version branch
-t, --create-tag             Create a release version tag
-f, --file string            Uses repos from file (one repo per line)
-n, --next-version string    Which number should be incremented by 1. Possible values: PATCH, MINOR, MAJOR (default "PATCH")
-r, --repos strings          Git Repo urls e.g. git@github.com:fhopfensperger/my-repo.git
-s, --source-branch string   Source reference branch (default: "main")
-b, --target-branch string   Which target branches to check for version (default: "release")
```
---
## Demonstration
(Updated 17.03.2021)

Situation: 

Git log of git@github.com:fhopfensperger/test-repo.git
```bash
$ git log --graph --abbrev-commit --decorate --format=format:'%C(bold blue)%h%C(reset) - %C(bold green)(%ar)%C(reset) %C(white)%s%C(reset) %C(dim white)- %an%C(reset)%C(bold yellow)%d%C(reset)' --all

* 60e84e7 - (18 hours ago) empty-commit - Florian Hopfensperger (HEAD -> main, tag: v1.0.4, origin/release/v1.0.4, origin/main)
* e67d3ea - (18 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.3, origin/release/v1.0.3)
* 97728f4 - (18 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.2)
* d9ff672 - (18 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.1)
* 6b73814 - (19 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.0)
* b8b72a0 - (2 days ago) empty-commit - Florian Hopfensperger
* 40448c7 - (2 days ago) empty-commit - Florian Hopfensperger
...
```

We want to create a new branch and tag `patch` version release

```log
$ git-releaser create -r git@github.com:fhopfensperger/test-repo.git -s main -n PATCH -t -c

2021-03-17T17:08:42+01:00 INF New PATCH version will be created
2021-03-17T17:08:43+01:00 INF Remote branches and tags found: [60e84e77d8b06276d06349579e6532e6bbb8b200 refs/heads/main 6b73814074af56706ff5a40bee48a0e9e6a8f770 refs/tags/v1.0.0 d9ff672d0eaf2f5e4bd51b181168d6970c4cbd7e refs/tags/v1.0.1 97728f458c714ebfc38fdacc890f869cfef172ca refs/tags/v1.0.2 e67d3ead8022f77d11ffae34d669f8302c8cb4da refs/tags/v1.0.3 e67d3ead8022f77d11ffae34d669f8302c8cb4da refs/heads/release/v1.0.3 60e84e77d8b06276d06349579e6532e6bbb8b200 refs/tags/v1.0.4 60e84e77d8b06276d06349579e6532e6bbb8b200 refs/heads/release/v1.0.4] for repo git@github.com:fhopfensperger/test-repo.git
2021-03-17T17:08:43+01:00 INF Nothing to do, main and latest branch version release/v1.0.4 are equals, commit hash: 60e84e77d8b06276d06349579e6532e6bbb8b200
2021-03-17T17:08:43+01:00 INF Successfully completed git@github.com:fhopfensperger/test-repo.git
```


Lets create a new dummy commit on the main branch, to make the main branch the latest branch.
```bash
$ git commit -m "empty-commit" --allow-empty && git push origin main
```

Git log of git@github.com:fhopfensperger/test-repo.git
```
* f18d8fe - (7 seconds ago) empty-commit - Florian Hopfensperger (HEAD -> main, origin/main)
* 60e84e7 - (18 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.4, origin/release/v1.0.4)
* e67d3ea - (18 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.3, origin/release/v1.0.3)
* 97728f4 - (18 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.2)
* d9ff672 - (18 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.1)
* 6b73814 - (20 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.0)
* b8b72a0 - (2 days ago) empty-commit - Florian Hopfensperger
```

Lets run `git-releaser` again
```log
$ git-releaser create -r git@github.com:fhopfensperger/test-repo.git -s main -n PATCH -t -c

2021-03-17T17:15:35+01:00 INF New PATCH version will be created
2021-03-17T17:15:37+01:00 INF Remote branches and tags found: [f18d8fe595c42cfa2fbf3d416444a4ff816ae9a0 refs/heads/main 6b73814074af56706ff5a40bee48a0e9e6a8f770 refs/tags/v1.0.0 d9ff672d0eaf2f5e4bd51b181168d6970c4cbd7e refs/tags/v1.0.1 97728f458c714ebfc38fdacc890f869cfef172ca refs/tags/v1.0.2 e67d3ead8022f77d11ffae34d669f8302c8cb4da refs/tags/v1.0.3 e67d3ead8022f77d11ffae34d669f8302c8cb4da refs/heads/release/v1.0.3 60e84e77d8b06276d06349579e6532e6bbb8b200 refs/tags/v1.0.4 60e84e77d8b06276d06349579e6532e6bbb8b200 refs/heads/release/v1.0.4] for repo git@github.com:fhopfensperger/test-repo.git
2021-03-17T17:15:44+01:00 INF Successfully created branch release/v1.0.5
2021-03-17T17:15:47+01:00 INF Successfully created tag: v1.0.5
2021-03-17T17:15:47+01:00 INF Successfully completed git@github.com:fhopfensperger/test-repo.git
```

Git log after pulling: 
```
* f18d8fe - (7 minutes ago) empty-commit - Florian Hopfensperger (HEAD -> main, tag: v1.0.5, origin/release/v1.0.5, origin/main)
* 60e84e7 - (18 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.4, origin/release/v1.0.4)
* e67d3ea - (18 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.3, origin/release/v1.0.3)
* 97728f4 - (18 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.2)
* d9ff672 - (18 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.1)
* 6b73814 - (20 hours ago) empty-commit - Florian Hopfensperger (tag: v1.0.0)
* b8b72a0 - (2 days ago) empty-commit - Florian Hopfensperger
```

As you can see we have now a new `release` branch with version `v1.0.5` and a tag with `v1.0.5`



---

### Use a file name to create new versions for multiple repos

`Content of repos.txt`
```txt
git@github.com:fhopfensperger/test-repo.git
git@github.com:fhopfensperger/git-releaser.git
# git@github.com:fhopfensperger/amqp-sb-client.git # Lines with a leading `#` wont be used
```
Command to create a new release branch (increment patch version)
```bash
$ git-releaser create -f repos1.txt -s main -n PATCH -c
```
---

## Installation

### Option 1 (script)

```bash
curl https://raw.githubusercontent.com/fhopfensperger/git-releaser/main/get.sh | bash
```

### Option 2 (manually)

Go to [Releases](https://github.com/fhopfensperger/git-releaser/releases) download the latest release according to your processor architecture and operating system, and unarchive it.

```bash
tar xvfz git-releaser_x.x.x_darwin_amd64.tar.gz
cd git-releaser_x.x.x_darwin_amd64
chmod +x git-releaser
```
