# API Documentation

## Base URL
```
http://localhost:8000/api
```

## Base URL for Production
```
https://zocket-task-manager-backend.onrender.com
```

## Authentication

### Sign Up
```http
POST `/auth/signup`
Content-Type: `application/json`

```
{
    "email": "user@example.com",
    "password": "secure_password",
    "first_name": "John",
    "last_name": "Doe"
}
```

Response (201 Created):
{
    "user_id": 1,
    "token": "jwt_token",
    "message": "User created successfully"
}
```

### Sign In
```http
POST /auth/signin
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "secure_password"
}

Response (200 OK):
{
    "user_id": 1,
    "token": "jwt_token",
    "message": "User signed in successfully"
}
```

## Protected Routes

All protected routes require the JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

### Tasks

#### Create Task
```http
POST /v1/tasks
Content-Type: application/json

{
    "title": "Implement WebSocket",
    "description": "Add real-time updates using WebSocket",
    "priority": "High",
    "status": "ToDo",
    "assigned_to": 2
}

Response (201 Created):
{
    "id": 1,
    "title": "Implement WebSocket",
    "description": "Add real-time updates using WebSocket",
    "priority": "High",
    "status": "ToDo",
    "assigned_to": 2,
    "created_by": 1,
    "created_at": "2024-03-14T12:00:00Z",
    "updated_at": "2024-03-14T12:00:00Z"
}
```

#### Get Task
```http
GET /v1/tasks/:id

Response (200 OK):
{
    "id": 1,
    "title": "Implement WebSocket",
    "description": "Add real-time updates using WebSocket",
    "priority": "High",
    "status": "ToDo",
    "assigned_to": 2,
    "created_by": 1,
    "created_at": "2024-03-14T12:00:00Z",
    "updated_at": "2024-03-14T12:00:00Z"
}
```

#### Update Task
```http
PUT /v1/tasks/:id
Content-Type: application/json

{
    "title": "Implement WebSocket",
    "description": "Updated description",
    "priority": "Medium",
    "status": "InProgress",
    "assigned_to": 3
}

Response (200 OK):
{
    "id": 1,
    "title": "Implement WebSocket",
    "description": "Updated description",
    "priority": "Medium",
    "status": "InProgress",
    "assigned_to": 3,
    "created_by": 1,
    "created_at": "2024-03-14T12:00:00Z",
    "updated_at": "2024-03-14T12:30:00Z"
}
```

#### Delete Task
```http
DELETE /v1/tasks/:id

Response (200 OK):
{
    "message": "Task deleted successfully"
}
```

#### List Tasks
```http
GET /v1/tasks

Response (200 OK):
[
    {
    "id": 6,
    "title": "play game",
    "priority": "Low",
    "status": "ToDo",
    "assigned_to": 1,
    "assigned_to_name": "Adarsh Jaiswal",
    "description": "want to play game",
    "created_by": 1,
    "created_at": "2025-03-20T22:02:44.380179Z",
    "updated_at": "2025-03-20T22:02:44.380179Z"
  },
   {
    "id": 3,
    "title": "Implement api",
    "priority": "Low",
    "status": "Done",
    "assigned_to": 2,
    "assigned_to_name": "ritu sharma",
    "description": "Add real-time api",
    "created_by": 2,
    "created_at": "2025-03-20T20:03:55.49163Z",
    "updated_at": "2025-03-20T20:41:53.953241Z"
  },
    // ... more tasks
]
```

### Task Analysis

#### Analyze Task with AI
```http
POST /v1/tasks/:id/analyze
Content-Type: application/json

{
    "description": "Optional additional context for the task",
    "context": "Optional background information"
}

Response (200 OK):
{
    "task_id": 1,
    "analysis": "Detailed analysis of the task complexity and requirements",
    "suggestions": [
        {
            "id": 1,
            "task_id": 1,
            "user_id": 1,
            "suggestion_text": "Detailed breakdown and recommendation",
            "sub_tasks": [
                {
                    "title": "Subtask 1",
                    "description": "Implementation details",
                    "priority": "High"
                },
                {
                    "title": "Subtask 2",
                    "description": "Implementation details",
                    "priority": "Medium"
                }
            ],
            "accepted": false,
            "created_at": "2024-03-14T12:00:00Z"
        }
    ]
}
```

### WebSocket Events

In addition to the existing WebSocket events, the following event is added for task suggestions:

```json
{
    "type": "suggestion_created",
    "data": {
        "id": 1,
        "task_id": 1,
        "user_id": 1,
        "suggestion_text": "AI-generated suggestion",
        "sub_tasks": [
            {
                "title": "Subtask title",
                "description": "Subtask description",
                "priority": "High"
            }
        ],
        "accepted": false,
        "created_at": "2024-03-14T12:00:00Z"
    }
}
```

### Users

#### Get User
```http
GET /v1/user/:id

Response (200 OK):
{
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-03-14T12:00:00Z",
    "logged_in_at": "2024-03-14T12:00:00Z"
}
```

#### Get all User
```http
GET /v1/user

Response (200 OK):
[
  {
    "id": 1,
    "email": "adarsh@gmail.com",
    "password": "",
    "first_name": "Adarsh",
    "last_name": "Jaiswal",
    "logged_in_at": "",
    "created_at": ""
  },
  {
    "id": 2,
    "email": "ritu@gmail.com",
    "password": "",
    "first_name": "ritu",
    "last_name": "sharma",
    "logged_in_at": "",
    "created_at": ""
  }
]
```

### WebSocket

#### Connect to WebSocket
```http
WebSocket: ws://localhost:8000/api/v1/ws
Authorization: Bearer <jwt_token>
```

#### WebSocket Message Types

1. Task Created:
```json
{
    "type": "task_created",
    "data": {
        // Full task object
    }
}
```

2. Task Updated:
```json
{
    "type": "task_updated",
    "data": {
        // Full task object
    }
}
```

3. Task Deleted:
```json
{
    "type": "task_deleted",
    "data": 1  // task ID
}
```

## Error Responses

### 400 Bad Request
```json
{
    "error": "Invalid request body"
}
```

### 401 Unauthorized
```json
{
    "error": "Unauthorized",
    "message": "Invalid or expired JWT"
}
```

### 403 Forbidden
```json
{
    "error": "Not authorized to update this task"
}
```

### 404 Not Found
```json
{
    "error": "Task not found"
}
```

### 500 Internal Server Error
```json
{
    "error": "Failed to create task"
}
```

## Task Properties

| Field       | Type     | Description                                |
|-------------|----------|--------------------------------------------|
| id          | int      | Unique task identifier                     |
| title       | string   | Task title                                |
| description | string   | Detailed task description                 |
| priority    | string   | "High", "Medium", or "Low"               |
| status      | string   | "ToDo", "InProgress", or "Done"          |
| assigned_to | int      | User ID of assignee                      |
| created_by  | int      | User ID of creator                       |
| created_at  | string   | Creation timestamp (ISO 8601)            |
| updated_at  | string   | Last update timestamp (ISO 8601)         | 