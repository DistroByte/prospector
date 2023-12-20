# Frontend

## CI/CD

Every push to a branch that is _not_ `master` will trigger a build, test and review pipeline.

The build stage will build a docker image and push it to the registry at `git.dbyte.xyz`.

It will then trigger a review application to be deployed and viewed at the link from the PR, using gitlab's "environment" feature.
