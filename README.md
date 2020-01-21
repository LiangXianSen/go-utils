# Go Utils

go-utils provides several tools for agile development


# Contribution

## lint

```
make lint
```

## test

```
make test
```

## build

```
make build
```

more from [Makefile](Makefile)



## develop

Branchs:

- master
- dev
- k8s



Each feature checkout from dev, through [lint](#lint), [test](#test), [build](#build) all action then merge. 

branch naming rules:

```
feat/xxx
bugfix/xxx
hotfix/xxx
new/xxx
var/xxx
```

each featuer from dev need to merge to master and k8s both.

> you only can run lint & test locally with dev branch.



## commit

make a `.gitignore` in your project to prevent non-meaning file check in. like:

```
.idea/
.vscode/
.DS_Store
```

more from https://github.com/github/gitignore



commit naming rules:

```
{package_name}:{one_blank}{Upper-case commit string} 
```

if you wanna add detailed comments, keep the plain and clear title, then add a new line below to writes more infomation.

samples:

```
*: Add .gitlab-ci.yaml
```

> without specific package, for all project use *



```
elastic: Add elasticsearch user & password
```

elastic package modifies



```
elastic: Add elasticsearch user & password

Due to dev env has no elastic authorization, the default elastic client have not do it yet. Add user & password for deployment
```

add commit string with detailed info if needs