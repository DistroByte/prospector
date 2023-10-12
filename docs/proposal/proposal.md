# School of Computing &mdash; Year 4 Project Proposal Form

## SECTION A

|                     |                   |
|---------------------|-------------------|
|Project Title:       | Prospector        |
|Student 1 Name:      | James Hackett     |
|Student 1 ID:        | 20308896          |
|Student 2 Name:      | Alexandru Dorofte |
|Student 2 ID:        | 20414772          |
|Project Supervisor:  | Stephen Blott     |

## SECTION B

### Introduction

Prospector will serve as a user management and infrastructure-as-a-service tool, enabling easy on-demand deployment of containers and virtual machines. This will be achieved by automating the administration aspects of deploying on scalable clusters. The tool will incorporate user management, service deployment and automated provisioning capabilities.

### Outline

Prospector will consist of a web based dashboard where users can create either virtual machines or containers quickly and easily. There will be an admin dashboard where administrators can modify user accounts and add limits to accounts.

- Users can choose to have their files stored on a network storage device and have those files mounted into containers or VMs when they are created. 
- Jobs can be exposed to the internet if the user desires.

The frontend shall talk exclusively to a backend REST API which will orchestrate the various jobs, store user information and configuration, and handle authentication.

### Background

The idea for this project came from the desire to run multiple job types (projects, websites, etc) quickly and easily on Redbrick hardware. It was augmented slightly to solve the problem of School of Computing student home directories filling up with garbage files.

The aim for this project is to make it easy for both administrators and users to create and destroy environments on demand, allowing for persistence of important files via a networked file system if desired. Provisioning both Docker containers and virtual machines will allow almost any project to be deployed. In this way, we are attempting to create a self-hosted version of a cloud hosting provider.

Our project will allow users to run any project they like, provided it can be run in a container (simple tasks) or a virtual machine (complex tasks).

### Achievements

Prospector will provide the following functions:

- User management
- Workload orchestration
- Mounting of networked storage 
- Exposing jobs to the internet
- CLI tool for interacting with the API

It will provide quick and easy access to user files for students, a platform for developers to run applications, system administrators to run services and manage users, and an interface for researchers to deploy applications that aid in their work.

### Justification

Prospector will be useful for any administrator in need of a complete solution for minimizing the overhead of managing all aspects of running a cluster for number of users.

Some use cases include:

- A tool for the School of Computing that can provide students with a platform to run their own environments, without the need for the admin to manually provision resources for each student. As well as having the option to mount networked storage.
- Researchers that need to train machine learning models could use this tool to deploy a model to be trained, "pre-baked" recipies can be provided.

### Programming language(s)

The backend and CLI tool of this project will be written in Go.

The frontend of this project will be written in Typescript using Angular. 

### Programming tools / Tech stack

This project will make use of the following tools:
- Angular - Web framework to aid development and design
- Docker - tool for running containers
- Qemu - tool for running virtual machines
- Traefik - reverse proxy to handle routing of traffic to user jobs
- Nomad - job orchestration tool
- Consul - key value store and service mesh
- OpenLDAP - database designed for user management
- Ansible - automation tool for configuration management

### Learning Challenges

The main learning challenges of this project will be:

- Learning how to orchestrate using Nomad 
- Interacting with network storage devices
- Backend will be written mainly in Go which is a new language for both of us

Overall the project contains familiar aspects for us as we have both worked with self-hosting and containers, however we'd both like to expand our horizons to orchestration and a language.

### Breakdown of work

#### Student 1 (James)

I will be focused on setting up Docker, Nomad, Consul, QUMU and Traefik. This is mainly down to my experience with some of the tools mentioned. I'll also ensure that our repo and associated components are kept in a good state (branch protection rules, CI/CD, etc).

I'll also focus on the backend of Prospector, in equal part with Alex. This will ensure that requirements and limitations of the infrastructure are bubbled up to the backend and frontend.

We'll both be writing unit tests as we go, so the CI/CD will be pretty important to add to checks on merge requests.

#### Student 2

I will be focused on the frontend of Prospector, while also working on the backend with James. Mainly in Docker since I have worked with it before and am familiar in working with containers.

For the front-end I will be using Angular to create the frontend, and will be using Typescript to write the code. As well as testing the frontend code using unit tests as I go and manual testing.
