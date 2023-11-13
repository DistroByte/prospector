# Prospector

## API

[![pipeline status](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/badges/master/pipeline.svg)](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/-/commits/master)
[![Latest Release](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/-/badges/release.svg)](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/-/releases)

This is a template for CA400 projects.

## CI/CD

This project uses Gitlab CI/CD to build and test the project. The CI/CD pipeline is defined in `.gitlab-ci.yml`. The pipeline is run on every commit to the repository. The pipeline consists of the following stages:

- `build`: Builds the project using respective build tools.
- `test`: Runs tests for the project.
- `review`: Creates a deployment at a subdomain of prospector.ie for the commit. This is used for code review.
- `deploy-canary`: Deploys the project to the canary environment. This happens on merge to master
- `deploy-production`: Deploys the project to the production environment. This is a manual step.

### Using the CI/CD pipeline

The pipeline will run on every commit pushed to the repo. Once the build has passed, it will create a review URL for you to view. This happens for merge requests.
![Screenshot of merge request](https://i.dbyte.xyz/firefox_40gRN6MdK.png)

Clicking on the "view app" button will bring you to the review environment for the commit. This is useful for code review.

Once you merge your branch into master, the pipeline will run again. This time, it will deploy to the canary environment. This is a full deployment, and will be available at [https://canary.prospector.ie](https://canary.prospector.ie).

Once you are happy with the canary deployment, you can manually deploy to production. This is done by clicking on the play button beside the `deploy-production` job in the pipeline view. This will deploy to [https://prospector.ie](https://prospector.ie).
