# Prospector

[![pipeline status](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/badges/master/pipeline.svg)](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/-/commits/master)
[![Latest Release](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/-/badges/release.svg)](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/-/releases)

Prospector is a user management and infrastructure-as-a-service tool that enables easy, on demand deployment of jobs in the form of containers and virtual machines.

You can view the site here: https://prospector.ie

## API

[![pipeline status](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/badges/master/pipeline.svg)](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/-/commits/master)
[![Latest Release](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/-/badges/release.svg)](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/-/releases)
[![coverage report](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/badges/master/coverage.svg)](https://gitlab.computing.dcu.ie/hacketj5/2024-ca400-proj/-/commits/master) 

## CI/CD

This project uses Gitlab CI/CD to build and test the project. The CI/CD pipeline is defined in `.gitlab-ci.yml`. The pipeline is run on every commit to the repository. The pipeline consists of the following stages:

- `build`: Builds the project using respective build tools.
- `test`: Runs tests for the project.
- `review`: Creates a deployment at a subdomain of prospector.ie for the commit. This is used for code review.
- `deploy-canary`: Deploys the project to the canary environment. This happens on merge to master
- `deploy-production`: Deploys the project to the production environment. This is a manual step.
