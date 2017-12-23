# dvcs-link

`dvcs-link` is a utility to generate a GitHub or GitLab link from a file under version control.

It uses the URL under `git config --get remote.origin.url` to select the host repository, and `origin` can be overriden with the `-remote` flag.

Currently, only GitHub and GitLab are supported.

## Usage

```bash
$ dvcs-link foo.go
https://gitlab.com/user/repo/blob/commit/foo.go

$ dvcs-link foo.go 5
https://gitlab.com/user/repo/blob/commit/foo.go#L5

$ dvcs-link foo.go 5 10
https://gitlab.com/user/repo/blob/commit/foo.go#L5-10

$ dvcs-link -remote github foo.go 5 10
https://github.com/user/repo/blob/commit/foo.go#L5-L10
```
