# Energy Management System
To resolve the “exec format error,” you need to build a multi-arch Docker image that is suitable for both Mac and Linux.

This will ensure that the binary format of the executable file inside the container is compatible with both Mac and Linux architectures.

Building a Multi-Arch Docker Image on Mac

docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64 -t <image-name> .

docker build -t your_image_name .
docker build --progress=plain --platform="linux/amd64" -t 19860314/starlite-first:latest --build-arg="MODULE_NAME=first" .

docker build --progress=plain --platform="linux/amd64" -t 19860314/starlite-front:latest 

docker build --progress=plain -t 19860314/starlite-first-mac:latest --build-arg="MODULE_NAME=first" .

<!-- docker tag your-image-name your-dockerhub-username/your-repo-name:your-tag
docker tag starlite-first 19860314/starlite-first:latest -->

docker push your-dockerhub-username/your-repo-name:your-tag
docker push 19860314/starlite-first:latest

docker pull 19860314/starlite-first:latest


host
cloud.sealos.io
port
30536
connection
postgresql://postgres:zcs9jsmp@cloud.sealos.io:30536?directConnection=true
postgresql://postgres:zcs9jsmp@starlite-postgresql.ns-xx7juor7.svc:5432
## Overall Arch
[xxx Charging Hub 充电互联互通开放平台](https://confluence.xxxmotors.com/pages/viewpage.action?pageId=1898501740)

## Development Setup

### Follow Lines Setup

https://github.startlite.cn/itapp/startlite/ITApp/lines#pen

Copy config.yaml.example for the service you want to start

## OpenAPI management
### Adding a submodule

You can add a submodule to a repository like this:

    git submodule add git@github.com:url_to/awesome_submodule.git path_to_awesome_submodule

With default configuration, this will check out the **code** of the
`awesome_submodule.git` repository to the `path_to_awesome_submodule`
directory, and will **add information to the main repository** about
this submodule, which contains the **commit the submodule points to**,
which will be the **current** commit of the default branch (usually the
`master` branch) at the time this command is executed.

After this operation, if you do a `git status` you'll see two files in
the `Changes to be committed` list: the `.gitmodules` file and the path
to the submodule. When you commit and push these files you commit/push
the submodule to the origin.


### Getting the submodule's code

If a new submodule is created by one person, the other people in the
team need to initiate this submodule. First you have to get the
**information** about the submodule, this is retrieved by a normal
`git pull`. If there are new submodules you'll see it in the output of
`git pull`. Then you'll have to initiate them with:

    git submodule init

This will pull all the **code** from the submodule and place it in the
directory that it's configured to.

If you've cloned a repository that makes use of submodules, you should
also run this command to get the submodule's code. This is not
automatically done by `git clone`. However, if you add the
`--recurse-submodules` flag, it will.


### Pushing updates in the submodule

The submodule is just a separate repository. If you want to make changes
to it, you should make the changes in its repository and push them like
in a regular Git repository (just execute the git commands in the
submodule's directory). However, you should also let the **main**
repository know that you've updated the submodule's repository, and make
it use the new commit of the repository of the submodule. Because if
you make new commits inside a submodule, the **main** repository will
still **point to the old commit**.

If there are changes in the submodule's repository, then `git status` in
the **main** repository, will show `Changes not staged for commit` and it
has the text `(modified content)` behind it. This means that the **code**
of the submodule is checked out on a different commit than the **main**
repository is **pointing to**. To make the **main** repository **point
to** this new commit, you should create another commit in the **main**
repository.

The next sections describe different scenarios on doing this.


#### Make changes inside a submodule

- `cd` inside the submodule directory.
- Make the desired changes.
- `git commit` the new changes.
- `git push` the new commit.
- `cd` back to the main repository.
- In `git status` you'll see that the submodule directory is modified.
- In `git diff` you'll see the old and new commit pointers.
- When you `git commit` in the main repository, it will update the
  pointer.


#### Update the submodule pointer to a different commit

- `cd` inside the submodule directory.
- `git checkout` the branch/commit you want to point to.
- `cd` back to the main repository.
- In `git status` you'll see that the submodule directory is modified.
- In `git diff` you'll see the old and new commit pointers.
- When you `git commit` in the main repository, it will update the
  pointer.


#### If someone else updated the submodule pointer

If someone updated a submodule, the other team-members should update
the code of their submodules. This is not automatically done by
`git pull`, because with `git pull` it only retrieves the
**information** that the submodule is **pointing** to another
**commit**, but doesn't update the submodule's **code**. To update the
**code** of your submodules, you should run:

    // retrieves the commit of the submodule recorded in the main project repo
    git submodule update
    
    //fetches the latest commit of the submodule repo
    git submodule update --remote
    

If a submodule is not initiated yet, add the `--init` flag. If any
submodule has submodules itself, you can add the `--recursive` flag to
recursively init and update submodules.

##### What happens if you don't run this command?

If you don't run this command, the **code** of your submodule is checked
out to an **old** commit. When you do `git status` you will see the
submodule in the `Changes not staged for commit` list with the text
`(modified content)` behind it. If you would do a `git status` **inside**
the submodule, it would say `HEAD detached at <commit-hash>`. This is not
because you changed the submodule's code, but because its **code** is
checked out to a different **commit**. So Git sees this as a change, but
actually you just didn't update the submodule's code. So if you're working
with submodules, don't forget to keep your submodules up-to-date.


### Making it easier for everyone

It is sometimes annoying if you forget to initiate and update your
submodules. Fortunately, there are some tricks to make it easier:

    git clone --recurse-submodules

This will clone a repository and also check out and init any possible
submodules the repository has.

    git pull --recurse-submodules

This will pull the main repository and also it's submodules.

And you can make it easier with aliases:

    git config --global alias.clone-all 'clone --recurse-submodules'
    git config --global alias.pull-all 'pull --recurse-submodules'


## Git Branch Name & Commit Message Pattern Rules

### Branch Name Rules
There are several branchName patterns here:
1. General Dev branch:
  * ```bugfix/GFSHBJM-<card num>-<description>```
  * ```feature/GFSHBJM-<card num>-<description>```
  * **Example:** feature/GFSHBJM-123-testForCommit

2. Need upload to ```eng``` server for testing
* ```try/GFSHBJM-<card num>-<description>```
* **Example:** try/GFSHBJM-123-testInEng

3. Cooperative develop branch
* ```eng/<year>.<week>-<description (optional)>```
* **Example:**
  * eng/2023.26
  * eng/2023.26-coWork

4. Release branch
* ```release/IEM.<year>.<week>```
* **Example:** release/IEM.2023.25


### Commit Message Rules
We follow [Conventioanl Commit](https://www.conventionalcommits.org/en/v1.0.0/)

A typical commit is like: ```<type>(<scope>): GFSHBJM-123 <description>```. The `<scope>` is optional and usually be used to mark which main scope you are changing(eg. cpo)
* feat: Represents a new feature.
* fix: Represents a bug fix.
* docs: Represents documentation updates.
* style: Represents changes related to code style, formatting, or missing semi-colons.
* refactor: Represents code refactoring, does not include bug fixes or new features.
* perf: Represents performance improvements.
* test: Represents changes related to tests, adding new tests or fixing existing ones.
* chore: Represents mundane tasks or changes not affecting functionality or production code.
* build: Represents changes impacting the build system or external dependencies.
* ci: Represents changes to Continuous Integration pipeline or configuration.

## Database Migrations

### Schema migration with Goose
We adopt a tool called Goose to manage database migrations.

Install `goose` using:
```bash
$ go get -u bitbucket.org/liamstask/goose/cmd/goose
```

### Migration Naming Convention
The naming convention that tends to be most effective for us should be:
`<Incrementing ID>_<Jira board>_<Jira number>_<brief description of what it does>.sql`

#### Running schema migration locally
Start a Postgres database in local development environment:
```
$ docker compose -f docker-compose.local.yaml up -d
```
Run the goose migration:
```
$ goose -env local -path ./database/db up
```
Generate ORM structs with the help of gen tool:
```bash
gen --sqltype=postgres \
    --connstr "host=127.0.0.1 port=5432 user=postgres dbname=iem password=dummy-password sslmode=disable" \
    --database iem --model=dbmodel  --gorm --guregu --overwrite \
    --table operation_log,charging_session,partner_data,trt_partner \
    --out ./
```

#### Running schema migration in ENG
Set up a ssh tunnel through the jump box to the ENG RDS instance.
```
$ ssh -p 22 -i ~/.ssh/id_rsa.pub -L 5432:db-iem-eng.czktmrdhhxjh.rds.cn-north-1.amazonaws.com.cn:5432 $(whoami)@jmpbx1.it.cloud.xxx.cn -fNg
```

Run the goose migration. To get the password look in AWS Secrets Manager. If you don't have access, ask someone on
gfsh-application-charging@xxx.com.
```
$ export PGUSER={user}
$ export PGPASS={password}
$ goose -env eng -path ./database/db up
```

#### Running schema migration in PRD
Set up an ssh tunnel through the jump box to the PRD RDS instance.
```
$ ssh -fNL 5433:{to_be_created}:5432 $(whoami)@jmpbx1.it.cloud.xxx.cn
(in case you still can't connect to the db, please try to go with `sudo`, sometimes this happen in OSX)

```

Run the goose migration. To get the password look in AWS Secrets Manager. If you don't have access, ask someone on
eps-charging-devs@xxx.com.
```
$ export PGUSER={user}
$ export PGPASS={password}
$ goose -env eng -path ./database/db up
```
### mockgen
example : 
mockgen -source=internal/pkg/features/clients/charging/command/command.go -destination=internal/pkg/test/mocks/mockcommand.go -package=mocks