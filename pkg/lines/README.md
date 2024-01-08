# LINES
[![Build Status](https://jenkins-gf3.xxx.cn/buildStatus/icon?job=GF3MOS-lines%2Fmain)](https://jenkins-gf3.xxx.cn/view/MOS/job/GF3MOS-lines/job/main/)


GO common framework

* [Design](./design.md)
* [FAQ](./FAQ.md)
  
## Pen
* [Use `Pen`](./tools/pen/README.md)

## Setup

```
go env -w GOPROXY=https://arf.xxx.cn/artifactory/api/go/gfsh-sdk-go-virtual,https://goproxy.cn,direct
go env -w GONOSUMDB="*.xxxmotors.com,*.xxx.cn"


go get -u -v github.startlite.cn/itapp/startlite/pkg/lines
```


## Code Read From

read from [here](example/petstorebyhand/cmd/petstore.pen.go)

## Commit Principles
* Branch: Generate branch on JIRA
* Commits:
  * Make sure every commit contains card number
  * Squash to one commit per PR unless you make every commit meaningful (`git rebase -i`)

#### put this file to `.git/hooks/prepare-commit-msg`, it will help add card number from branch name automatically

```
#!/bin/bash

# Get the current branch name
BRANCH_NAME=$(git symbolic-ref --short HEAD)

# Get the JIRA number
JIRA=$(echo $BRANCH_NAME | grep -o -E '[0-9A-Z]+-[0-9]+')

#Check if this a normal commit, or an amend
IS_NORMAL_COMMIT=$(grep -c "$JIRA" $1)

# Prepend the JIRA number to the commit message
if [ -n "$JIRA" ] && [ "$IS_NORMAL_COMMIT" -eq 0 ]; then
  sed -i.bak -e "1s/^/$JIRA /" $1
fi
```

## Go Guide
* use `errorx` as more as possible, it will add stack info in log
* follow format [here](./tools/importformat/README.md)
* [Uber Go Guide](https://github.com/uber-go/guide)

