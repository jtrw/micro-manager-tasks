# Micro manager tasks

This application serves as a REST API server designed to manage tasks and subtasks using a MongoDB database. It provides a set of endpoints for creating tasks, adding subtasks to existing tasks, checking task status, and retrieving task information. The server leverages GoChi for routing and handles HTTP requests, integrating seamlessly with MongoDB for task management operations.

---

## Table of Contents

1. [Application Overview](#application-overview)
2. [Server Configuration](#server-configuration)
3. [Routes and Handlers](#routes-and-handlers)
    - [Create Task](#create-task)
    - [Add Subtask](#add-subtask)
    - [Check Task Status](#check-task-status)
    - [Show Task Information](#show-task-information)
4. [Handler Functions](#handler-functions)
5. [Constants](#constants)

---

## Server Configuration

The server is configured with various parameters, including the server address, pin size, maximum pin attempts, expiration time for tasks, web root directory, secret key, version, MongoDB client, database name, and collection name.

---

## Routes and Handlers

### Create Task

**Endpoint:** `/api/v1/tasks` (POST)

- **Description:** Creates a new task.
- **Request Body:** JSON object containing task details (UUID, Count, Type, Callback, Subtasks).
- **Response:** JSON object containing the UUID of the created task.

### Add Subtask

**Endpoint:** `/api/v1/tasks/{uuid}/subtask` (POST)

- **Description:** Adds a subtask to an existing task identified by UUID.
- **Request Body:** JSON object containing subtask details (UUID, Type, Status).
- **Response:** JSON object containing the UUID of the added subtask.

### Check Task Status

**Endpoint:** `/api/v1/tasks/{uuid}/status` (GET)

- **Description:** Retrieves the status of a task identified by UUID.
- **Response:** JSON object containing the UUID and status of the task.

### Show Task Information

**Endpoint:** `/api/v1/tasks/{uuid}` (GET)

- **Description:** Retrieves detailed information about a task identified by UUID.
- **Response:** JSON object containing all details of the task.

---

## Handler Functions

### `CreateTask`

- **Description:** Creates a new task in the MongoDB collection.
- **Request Body:** Task details (UUID, Count, Type, Callback, Subtasks).
- **Response:** JSON response with the UUID of the created task.

### `AddSubTask`

- **Description:** Adds a subtask to an existing task in the MongoDB collection.
- **Request Body:** Subtask details (UUID, Type, Status).
- **Response:** JSON response with the UUID of the added subtask.

### `CheckStatus`

- **Description:** Retrieves the status of a task from the MongoDB collection.
- **Response:** JSON response with the UUID and status of the task.

### `ShowTaskInfo`

- **Description:** Retrieves detailed information about a task from the MongoDB collection.
- **Response:** JSON response containing all details of the task.

---

## Constants

- **COLLECTION_TASKS**: Name of the MongoDB collection used to store tasks.

---

This application functions as a robust REST API server, employing GoChi for routing and handling HTTP requests. It seamlessly integrates with MongoDB to manage tasks and subtasks efficiently, providing a comprehensive set of endpoints for various task management operations.
