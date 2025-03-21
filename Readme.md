# Task Management System with Real-Time Updates

A modern task management system built with Go Fiber, featuring real-time updates through WebSockets and secure authentication.

## Features

- 🔐 JWT-based Authentication
- 🚀 Real-time Task Updates
- 📋 Complete Task CRUD Operations
- 👥 User Management
- 🔄 WebSocket Integration
- 🛡️ Role-based Access Control

## Project Structure

```
.
├── db/                 # Database connection and schema
├── internal/
│   ├── middleware/     # JWT authentication middleware
│   ├── tasks/         # Task-related handlers and logic
│   ├── user/          # User-related handlers and logic
│   └── websocket/     # WebSocket manager for real-time updates
├── types/             # Shared types and constants
└── main.go           # Application entry point
```

## Prerequisites

- Go 1.23 or higher
- PostgreSQL
- Make (optional, for using Makefile commands)

## Setup

1. Clone the repository:
```bash
git clone https://github.com/your-username/task-management.git
```

2. Create a `.env` file in the root directory:
```env
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=your_db_name
JWT_SECRET=your_jwt_secret
```

3. Initialize the database:
```bash
make table
```

4. Run the application:
```bash
make run
```

## WebSocket Implementation

### Why WebSockets?

Traditional REST APIs require clients to poll the server for updates. WebSockets provide:
- Real-time updates without polling
- Reduced server load
- Better user experience
- Instant notifications for task changes

### How it Works

1. **Server-side**:
   - WebSocket manager maintains active connections
   - Task changes trigger broadcasts to all connected clients
   - JWT authentication ensures secure connections

2. **Client-side Integration**:
```javascript
// Connect to WebSocket with JWT
const token = 'your_jwt_token';
const ws = new WebSocket(`ws://localhost:8000/api/v1/ws`);

// Set up connection
ws.onopen = () => {
    console.log('Connected to WebSocket');
};

// Handle incoming messages
ws.onmessage = (event) => {
    const update = JSON.parse(event.data);
    switch (update.type) {
        case 'task_created':
            handleNewTask(update.data);
            break;
        case 'task_updated':
            handleTaskUpdate(update.data);
            break;
        case 'task_deleted':
            handleTaskDeletion(update.data);
            break;
    }
};

// Handle errors
ws.onerror = (error) => {
    console.error('WebSocket error:', error);
};

// Handle disconnection
ws.onclose = () => {
    console.log('Disconnected from WebSocket');
    // Implement reconnection logic if needed
};
```

### Message Types

1. **Task Created**:
```json
{
    "type": "task_created",
    "data": {
        "id": 1,
        "title": "New Task",
        "priority": "High",
        "status": "ToDo",
        "assigned_to": 2,
        "description": "Task description",
        "created_by": 1,
        "created_at": "2024-03-14T12:00:00Z",
        "updated_at": "2024-03-14T12:00:00Z"
    }
}
```

2. **Task Updated**:
```json
{
    "type": "task_updated",
    "data": {
        "id": 1,
        "title": "Updated Task",
        "priority": "Medium",
        "status": "InProgress"
        // ... other fields
    }
}
```

3. **Task Deleted**:
```json
{
    "type": "task_deleted",
    "data": 1  // task ID
}
```

## Security

- JWT authentication for all API endpoints
- WebSocket connections require valid JWT
- Role-based access control for task operations
- Input validation and sanitization

## Error Handling

The application uses standard HTTP status codes and consistent error responses:
```json
{
    "error": "Error message here"
}
```

## Development

- Run tests: `go test ./...`
- Format code: `go fmt ./...`
- Lint code: `golangci-lint run`

## API Endpoints

- Check the [API documentation](API.md) for detailed information on all endpoints.

## DockerFile

- Build the Docker image:
```bash
docker build -t task-management .
```

- Run the Docker container:
```bash
docker run -p 8000:8000 task-management
```

