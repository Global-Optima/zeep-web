# ZeepApp

ZeepApp is a real-time drink ordering application designed for seamless interaction between customers and our coffee shops. This repository contains the complete web flow, including a Golang backend and a Vite Vue 3 frontend, facilitating efficient order management and tracking.

## Table of Contents
- [Description](#description)
- [Getting Started](#getting-started)
  - [Backend](#backend)
  - [Frontend](#frontend)
- [Documentation](#documentation)
- [Branching Strategy](#branching-strategy)
- [Pull Request Guidelines](#pull-request-guidelines)
- [Testing](#testing)

## Description

ZeepApp allows users to conveniently order beverages from our coffee shops while providing real-time tracking of their orders. The application is structured to enhance user experience, ensuring that customers receive immediate updates on order statuses. This project aims to streamline the ordering process for both in-store and remote users.

## Getting Started

### Backend
The backend is built using Golang with the Gin framework and is located in the `/backend` folder.

#### Commands to Run the Backend
1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the application:
   ```bash
   go run main.go
   ```

### Frontend
The frontend is developed with Vite and Vue 3, using Yarn as the package manager, and is located in the `/frontend` folder.

#### Commands to Run the Backend
1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   yarn install
   ```
3. Run the application:
   ```bash
   yarn dev
   ```
Access the application at http://localhost:3000.

## Documentation

Each folder contains a `docs` directory with a README file that provides additional guidelines and documentation for specific parts of the project. Developers are encouraged to refer to these documents for more detailed information about the architecture, setup, and usage.

## Branching Strategy

This project follows a clear branching strategy:

- **`dev`**: This branch is used for merging functionality that has been tested. It serves as the main development branch.
- **`main`**: This branch is reserved for production use only. Merges to this branch should be performed following thorough testing and approval.

## Pull Request Guidelines

When creating a pull request (PR), adhere to the following naming conventions and templates:

### Branch Naming

Use the following prefixes for branch names according to Git Flow conventions:

- **`feature/`**: For new features (e.g., `feature/new-order-form`)
- **`hotfix/`**: For quick fixes (e.g., `hotfix/fix-bug-in-order`)
- **`bugfix/`**: For bug fixes (e.g., `bugfix/fix-ui-issue`)

### Pull Request Template

When creating a PR, include the following information:

- **Sprint Number**: Specify the current sprint number.
- **Functionality Description**: Briefly describe the changes and functionality introduced.
- **Testing Guidelines**: Include any relevant testing instructions.
- **Related Issue**: Provide a link to the related issue in Asana.

Ensure all sections are filled out to streamline the review process.

## Testing

It is crucial to cover the written code with appropriate tests:

- **Unit Tests**: Implement unit tests for individual components and functions.
- **Integration Tests**: Create integration tests to verify the interaction between different parts of the application.

Refer to the documentation in the `docs` folders of each project to see recommended libraries and practices for testing.

