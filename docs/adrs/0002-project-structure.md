# Project Structure

## Context and Problem Statement

Given that we are starting a new project, we need to establish a folder structure that give clarity about the
application goal, makes easier for developers regarding where should they place new components and let the application
evolve in a sustainable way.

## Decision Outcome

**internal/**: keep core packages scoped to this project. Since we have two different applications,
we have a package for each of then on the top level: **loader** and **api**.
**cmd/**: keep project entry points and main executables. Modules here will be responsible to init configurations and
glue all required modules to run a specific application.
**deployments/:**: keep deployment related configurations. e.g: docker-compose files.
**docs/**: project documentations, like ADRs.
**test/**: keeps __non-unit tests__. e.g: end-to-end tests should be placed here.

## References

- [https://blog.boot.dev/golang/golang-project-structure/]()
