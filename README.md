# Sport Link

SportLink is an application designed to help individuals find others interested in playing a specific sport and provide
options for available venues. It simplifies the process of organizing a match by connecting players and facilitating
court or field reservations.

## Architecture Overview

This project follows the Hexagonal Architecture, also known as Ports and Adapters Architecture, which aims to create a
loosely coupled application that isolates the core logic from external concerns. Below are the main directories and
their responsibilities:

### `cmd`

Contains the main applications of the project. Each application has its own directory with its own `main.go` file, which
serves as the entry point.

- **Responsibilities**:
    - Starting up the application.
    - Setting up high-level application configurations and dependencies.

### `api`

This directory contains the internal components of the application, which are not meant to be exposed outside.

#### _app_

- **Responsibilities**:
    - Orchestrating the flow of data between the domain layer and the infrastructure.
    - Implementing application-specific logic (use cases).

#### _domain_

- **Responsibilities**:
    - Containing all the business logic and business rules.
    - Defining interfaces (ports) that describe the operations that can be performed with domain objects.

#### _infrastructure_

Contains all the external concerns and details such as database access, file handling, external APIs, and web
frameworks.

- **Subdirectories**:
    - `persistence`
        - **Responsibilities**:
            - Implementing repository interfaces defined in the domain layer.
            - Handling all database operations.
    - `rest`
        - **Responsibilities**:
            - Handling all HTTP request routing and responses.
            - Marshalling and unmarshalling of JSON data.
    - `config`
        - **Responsibilities**:
            - Managing configuration settings from files or environment variables.

### `dev`

The `dev/` directory contains all configurations and scripts specifically used for the development environment. This
includes local setups, mock services, and other development-only utilities that are not part of the production
deployment. Key components include:

- **LocalStack**: This subdirectory houses configuration files and scripts for LocalStack, which simulates AWS cloud
  services locally. It allows developers to test cloud applications by using cloud APIs offline.

    - **CloudFormation Templates**: Located under `dev/localstack/cloudformation/`, these YAML files are used to set up
      AWS resources such as DynamoDB tables in the LocalStack environment. They provide a reliable and repeatable way to
      configure necessary services for development tests.

This structure ensures that development-specific tools and configurations are neatly organized and separated from the
production codebase, reducing the risk of deploying development artifacts to production environments and enhancing
overall project maintainability.

## Getting Started

Instructions on how to build and run the project.

### Prerequisites

List of software and tools required.

### Running the application

Steps to run the application locally.

## Testing

Explanation on how to run the tests.

## Deployment

Guidelines for deploying the application in different environments.

