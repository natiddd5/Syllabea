package types

import (
	"encoding/json"
	"time"
)

// User represents a row in the 'users' table.
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // Possible values: "Instructor", "Manager"
	CreatedAt time.Time `json:"created_at"`
}

// Department represents a row in the 'departments' table.
type Department struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Course represents a row in the 'courses' table.
type Course struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DepartmentID int    `json:"department_id"`
}

// Syllabus represents a row in the 'syllabi' table.
type Syllabus struct {
	ID             int             `json:"id"`
	CourseID       int             `json:"course_id"`
	LecturerID     int             `json:"lecturer_id"`
	Status         string          `json:"status"`          // Possible values: "Draft", "Pending", "In Review", "Approved"
	SubmissionDate time.Time       `json:"submission_date"` // Maps to the DATE column in MySQL
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Data           json.RawMessage `json:"data"` // Stored as JSON in the DB
}

// SyllabusStatus is a custom type for the status of a syllabus.
type SyllabusStatus string

type Comment struct {
	ID         int       `json:"id" db:"id"`                   // Unique identifier for the comment.
	SyllabusID int       `json:"syllabus_id" db:"syllabus_id"` // The ID of the associated syllabus.
	UserID     int       `json:"user_id" db:"user_id"`         // The ID of the commenting user.
	Content    string    `json:"content" db:"content"`         // The text content of the comment.
	CreatedAt  time.Time `json:"created_at" db:"created_at"`   // Timestamp of when the comment was created.
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`   // Timestamp of the last update (auto-updated on change).
}
