# Contributing to cfd

Patches and contributions to this project are welcome!!

## Getting started

First you will need to setup your GitHub account and create a fork:

1. Create [a GitHub account](https://github.com/join)
1. Setup [GitHub access via
   SSH](https://help.github.com/articles/connecting-to-github-with-ssh/)
1. Create your own [fork of this
  repo](https://help.github.com/articles/fork-a-repo/)
1. [Clone it to your machine](https://docs.github.com/en/github/creating-cloning-and-archiving-repositories/cloning-a-repository)

## Testing

Ensure test pass correctly before create the pull request.

```shell
> go test ./... -cover -count=1
```

## Building

Final artifacts are building by using Github Actions. If you want to compile the binary on your machine:

```shell
> cd cmd/cfd
> go build
```

## Creating a PR

When you have changes you would like to propose to `cdf`, you will need to:

1. Ensure the commit message(s) describe what issue you are fixing and how you are fixing it
   (include references to [issue numbers](https://help.github.com/articles/closing-issues-using-keywords/)
   if appropriate)
1. [Create a pull request](https://help.github.com/articles/creating-a-pull-request-from-a-fork/)
1. Please follow the guidelines on the pull request template.

## Code reviews

All submissions, including submissions by project members, require review. We
use GitHub pull requests for this purpose. Consult
[GitHub Help](https://help.github.com/articles/about-pull-requests/) for more
information on using pull requests.