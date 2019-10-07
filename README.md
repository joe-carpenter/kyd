# kyd (Kustomize YAML Diff)
A super simple CLI tool to diff two Kustomize style YAML files, where you have a number of K8s manifests seperated by `---`.

The tool will match up and compare the individual manifests between the two files, and report differences on a per-manifest basis.

## Installation

```bash
go get github.com/joe-carpenter/kyd && \
go build ${GOPATH}/src/github.com/joe-carpenter/kyd && \
mv kyd /usr/local/bin/
```

## Usage

```bash
kyd /path/to/kustomize/yamlfile1.yml /path/to/kustomize/yamlfile2.yml
```