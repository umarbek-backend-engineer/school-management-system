# Swagger/OpenAPI Documentation Setup

This project now includes comprehensive Swagger/OpenAPI documentation for all REST API endpoints.

## Features

- ✅ Full API documentation with Swagger UI
- ✅ Automated API specification generation
- ✅ All endpoints documented with descriptions, parameters, and responses
- ✅ JWT Bearer token authentication documentation
- ✅ Interactive API testing through Swagger UI

## Accessing Swagger UI

Once the server is running, access the Swagger documentation at:

```
https://localhost:8080/swagger
```

The Swagger JSON specification is available at:
```
https://localhost:8080/swagger.json
```

## API Endpoints Documentation

All endpoints are organized by tags:

### Teachers
- `GET /teachers/` - Get all teachers
- `POST /teachers/` - Add new teachers
- `PATCH /teachers/` - Patch multiple teachers
- `DELETE /teachers/` - Delete multiple teachers
- `GET /teachers/{id}` - Get teacher by ID
- `PUT /teachers/{id}` - Update teacher
- `PATCH /teachers/{id}` - Patch specific teacher
- `DELETE /teachers/{id}` - Delete teacher
- `GET /teachers/{id}/students/` - Get students by teacher ID
- `GET /teachers/{id}/studentcount/` - Get student count for teacher
- `DELETE /allteachers/` - Delete all teachers

### Students
- `GET /students/` - Get all students (with pagination)
- `POST /students/` - Add new students
- `PATCH /students/` - Patch multiple students
- `DELETE /students/` - Delete multiple students
- `GET /students/{id}` - Get student by ID
- `PUT /students/{id}` - Update student
- `PATCH /students/{id}` - Patch specific student
- `DELETE /students/{id}` - Delete student
- `DELETE /allstudents/` - Delete all students

### Executives (Admin)
- `GET /execs/` - Get all executives
- `POST /execs/` - Add new executives
- `PATCH /execs/` - Patch multiple executives
- `GET /execs/{id}` - Get executive by ID
- `PATCH /execs/{id}` - Patch specific executive
- `DELETE /execs/{id}` - Delete executive
- `POST /execs/login` - Login and get JWT token
- `POST /execs/logout` - Logout
- `POST /execs/{id}/updatepassword` - Update executive password
- `POST /execs/forgotpassword` - Request password reset link
- `POST /execs/resetpassword/reset/{resetcode}` - Reset password with token

### Authentication
- `POST /execs/login` - Executive login
- `POST /execs/logout` - Executive logout
- `POST /execs/{id}/updatepassword` - Update password
- `POST /execs/forgotpassword` - Forgot password
- `POST /execs/resetpassword/reset/{resetcode}` - Reset password

## Documentation Comments

All handler functions include complete Swagger documentation comments in the following format:

```go
// GetTeacherHandler retrieves a teacher by ID
// @Summary Get teacher by ID
// @Description Get a specific teacher by their ID
// @Tags Teachers
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} models.Teacher
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /teachers/{id} [get]
// @Security Bearer
func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {
```

## Swagger Decorators Used

- `@title` - API title
- `@version` - API version
- `@description` - API description
- `@host` - Host address
- `@BasePath` - Base path for API
- `@schemes` - Supported schemes (https)
- `@Summary` - Brief handler summary
- `@Description` - Detailed handler description
- `@Tags` - Endpoint categories
- `@Accept` - Request content types
- `@Produce` - Response content types
- `@Param` - Request parameters
- `@Success` - Successful responses
- `@Failure` - Error responses
- `@Router` - Endpoint route and method
- `@Security` - Security requirements (Bearer JWT)

## Models Documentation

All Go models include JSON tags for proper serialization:

### Student
```go
type Student struct {
    ID        int    `json:"id,omitempty"`
    FirstName string `json:"first_name,omitempty"`
    LastName  string `json:"last_name,omitempty"`
    Email     string `json:"email,omitempty"`
    Class     string `json:"class,omitempty"`
}
```

### Teacher
```go
type Teacher struct {
    ID        int    `json:"id,omitempty"`
    FirstName string `json:"first_name,omitempty"`
    LastName  string `json:"last_name,omitempty"`
    Email     string `json:"email,omitempty"`
    Class     string `json:"class,omitempty"`
    Subject   string `json:"subject,omitempty"`
}
```

### Executive (Admin)
```go
type Exec struct {
    ID                        int            `json:"id,omitempty"`
    FirstName                 string         `json:"first_name,omitempty"`
    LastName                  string         `json:"last_name,omitempty"`
    Email                     string         `json:"email,omitempty"`
    Username                  string         `json:"username,omitempty"`
    // ... other fields
}
```

## Project Structure

```
RestAPI/
├── cmd/api/
│   └── server.go (with Swagger annotations)
├── internal/api/
│   ├── docs.go (Swagger doc main file)
│   ├── handler/
│   │   ├── swagger.go (Swagger UI handlers)
│   │   ├── teachers.go (with Swagger docs)
│   │   ├── students.go (with Swagger docs)
│   │   └── execs.go (with Swagger docs)
│   ├── middlerwares/
│   ├── router/
│   │   └── routerr.go (with Swagger routes)
│   └── models/
├── go.mod (with Swagger dependencies)
└── SWAGGER_SETUP.md (this file)
```

## Dependencies Added

The following dependencies have been added to `go.mod`:

```go
require (
    github.com/swaggo/swag v1.16.3
    github.com/swaggo/files v1.0.1
    github.com/swaggo/gin-swagger v1.6.0
)
```

## Next Steps

To use the complete Swagger documentation generation pipeline:

1. **Install swag CLI** (optional for auto-generation):
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Generate API docs** (optional):
   ```bash
   swag init
   ```

3. **Run the server**:
   ```bash
   go run cmd/api/server.go
   ```

4. **Access Swagger UI**:
   Open your browser and navigate to `https://localhost:8080/swagger`

## Security

- All endpoints that require authentication are marked with `@Security Bearer`
- The API uses JWT Bearer tokens for authentication
- Endpoints are protected with middleware that validates user roles

## Example Usage

### Getting All Teachers
```bash
curl -X GET "https://localhost:8080/teachers/" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "accept: application/json"
```

### Creating a New Teacher
```bash
curl -X POST "https://localhost:8080/teachers/" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '[{"first_name":"John","last_name":"Doe","email":"john@example.com","class":"10A","subject":"Math"}]'
```

### Logging In
```bash
curl -X POST "https://localhost:8080/execs/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'
```

## Notes

- The Swagger UI is served at `/swagger` endpoint
- The JSON specification is available at `/swagger.json`
- All request/response examples use JSON format
- The API supports pagination on list endpoints with `page` and `limit` query parameters
- All timestamps are in ISO 8601 format
