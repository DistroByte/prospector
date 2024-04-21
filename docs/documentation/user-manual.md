
# User Manual - Prospector

- **Project Title:** Prospector
- **Student 1:** James Hackett - 20308896
- **Student 2:** Alexandru Dorofte - 20414772
- **Supervisor:** Dr Stephen Blott
- **Date Completed:** 2024-04-21

## Table of Contents

- [User Interface](#User%20Interface)
	- [Setup](#Setup)
	- [Getting Started](#Getting%20Started)
	- [Creating a Project](#Creating%20a%20Project)
		- [Container Project Type](#Container%20Project%20Type)
		- [Virtual Machine Project Type](#Virtual%20Machine%20Project%20Type)
- [Project Options](#Project%20Options)
	- [Interacting with a Project](#Interacting%20with%20a%20Project)
	- [Stopping a Project/Component](#Stopping%20a%20Project/Component)
	- [Starting a Project/Component](#Starting%20a%20Project/Component)
	- [Restarting a Project/Component](#Restarting%20a%20Project/Component)
	- [Deleting a Project/Component](#Deleting%20a%20Project/Component)
- [Project Configuration](#Project%20Configuration)
	- [Exposing a Container Type Component](#Exposing%20a%20Container%20Type%20Component)

## User Interface

### Setup

To get started with the web application, simply visit [https://prospector.ie](https://prospector.ie).

The web application is hosted entirely inside the Prospector ecosystem, so it is always available, provided you have access to the internet. There are no additional software or hardware needs to make full use of the application.

### Getting Started

Once you open the web application, you'll be greeted with a landing page. The login button is in the top right corner. Clicking that will bring you to the login page.

If you've never used Prospector before, you can register for an account by clicking the "Register" button will guide you through setting up a new account. If you already have an account, then you can enter your credentials and you'll be taken to your dashboard.

From your dashboard you can manage your projects using Prospector.

### Creating a Project

To get started, open the hamburger menu in the top left and click the "Create Project" button. This will bring you to a form where you can specify the configuration of your project.

First, give your project a name. This can be anything you like and will be used to identify the project in your list of projects.

Next, you'll need to select a project type. This can either be "container" or "virtual machine".

#### Container Project Type

First, give your container component a name, this will be used to identify the component and will also be used if you want to expose the container to the internet.

Next, you will need to choose your container image. You can select from a variety of popular container images, or specify a custom one if you prefer.

If you want to make your container publically accessible, you can toggle the "expose" switch. Note, this requires the port number the container uses to work. See [exposing a container](#Exposing-a-Container-Type-Component) for more information.

Finally, you'll need to specify the hardware requirements of the container. The default values will work well for most containers, but if you're running a lot of applications, you may hit your quota, so make sure to only request what you'll use!

Once you've entered all the information required for that component, you can follow the same steps again to add another component to the project, or submit the project as is. Once you submit the project, you'll be redirected back to the dashboard where you'll be able to see your project listed with the name you gave it earlier!

Congratulations, you've successfully created your first container project with Prospector.

#### Virtual Machine Project Type

First, give your virtual machine component a name. This will be the hostname assigned to the virtual machine and will also be used to identify the component.

Next, you will need to choose from a list of supported virtual machine operating system images. The default supported image is Debian.

If you want to be able to connect to your virtual machine over ssh, you should provide an ssh key.

Finally, you’ll need to specify the hardware requirements of the virtual machine. The default values will work, but if you need a bigger virtual machine (or lots of smaller ones), then tune it as you see fit.

Once you've entered all the information required for your component, submit the project and you'll be redirected back to the dashboard where you'll be able to see your project listed with the name you gave it earlier!

## Project Options

### Interacting with a Project

You can interact with your project by clicking on it's name in the dashboard. This will bring you to a detailed view where you can see a project's components and resource usage, and stop/start/restart/delete a project or any of it's components. You can then click on each component and view more detailed resource usage and logs for that component.

### Stopping a Project/Component

You can stop a project or component with ease by clicking on the "stop" button beside the project's name on the dashboard or in the detailed view. This will leave the project/component registered with Prospector, allowing you to restart it easily later.

### Starting a Project/Component

To start a project or component that is stopped, simply click on the "start" button beside it's name in the dashboard or in the detailed view.

### Restarting a Project/Component

To restart a project or component, simply click on the "restart" button beside it's name in the dashboard or in the detailed view.

### Deleting a Project/Component

To delete a project or component, simply click on the "delete" button beside it's name in the dashboard or in the detailed view. This is a non-recoverable delete. There is no way to retrieve any data that is inside your project once this action has been performed!

## Project Configuration

### Exposing a Container Type Component

To expose your container, specify the port you'd like to have exposed and enable the expose option. The container will then be available at `https://<username>-<component_name>-<project_name>.prospector.ie`, a link to which can be found on your dashboard beside the project name.
