# buildkite-operator
the intention is for this operator to be able to act as a thin interface to the buildkite API (graciously exposed by [go-buildkite](https://github.com/buildkite/go-buildkite)

## Getting Started

Currently this project is in its infancy and barely anything works. That being said, the operator is able to create pipelines using a provided access token and return some information to populate the resource statuses.

Basically, the steps to get this running are:

* get an access token with [proper scopes](https://buildkite.com/docs/apis/managing-api-tokens)
* run `make kind` to create a new cluster, compile the code and deploy it
* edit the examples in `config/samples/`
  * the secret.data.token needs to be a token
  * the accessToken resource should point to the secret
  * the pipeline resource should reference the accessToken
* apply them and start developing!

## Example

```
$ kubectl get accesstokens
NAME                 UUID                                   SCOPES
banshee              63674a9a-665e-4cb4-ba52-5647fbea8381   ["read_agents","write_agents","read_teams","read_artifacts","read_builds","read_job_env","read_build_logs","read_organizations","read_pipelines","write_pipelines","read_user"]
```

```
$ kubectl get pipelines
NAME            ORGANIZATION   RUNNINGBUILDS   RUNNINGJOBS
omar            bonnie-doon    0               0
```

## Status
Going by the functionality provided by the [API](https://buildkite.com/docs/apis/rest-api)

done:
* nothing really...

in progress:
* pipelines
* access tokens
* organizations (might not continue as it's a bit useless - might be better as an purely controller managed resource)

planned:
* builds
* agents

maybe (extra functionality not provided by Buildkite directly):
* accsss token provisioning via user+password auth
* agent pod provisioning with scoped tokens and automatic registration
* git ops via bk??

todo: everything
