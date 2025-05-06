package repository

import (
	"Syllybea/UIcomponents"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	// Importing types from the types directory.
	"Syllybea/types"
)

// Repository wraps the DB connection.
type Repository struct {
	DB *sql.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// =======================
//       USERS CRUD
// =======================

// CreateUser inserts a new user into the DB.
func (r *Repository) CreateUser(u *types.User) error {
	query := `INSERT INTO users (name, email, role) VALUES (?, ?, ?)`
	result, err := r.DB.Exec(query, u.Name, u.Email, u.Role)
	if err != nil {
		return fmt.Errorf("CreateUser: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("CreateUser (retrieve id): %w", err)
	}
	u.ID = int(id)
	return nil
}

func (r *Repository) GetUserByID(id int) (*types.User, error) {
	query := `SELECT id, name, email, role, created_at FROM users WHERE id = ?`
	u := &types.User{}
	var createdAtStr string
	if err := r.DB.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &createdAtStr); err != nil {
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}
	// Parse the created_at string into a time.Time.
	createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: parsing created_at: %w", err)
	}
	u.CreatedAt = createdAt
	return u, nil
}

// GetUserByEmail retrieves a user by their email address.
func (r *Repository) GetUserByEmail(email string) (*types.User, error) {
	query := `SELECT id, name, email, role, created_at FROM users WHERE email = ?`
	user := &types.User{}
	var createdAtStr string
	if err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &createdAtStr); err != nil {
		return nil, fmt.Errorf("GetUserByEmail: %w", err)
	}

	// Parse the created_at string into a time.Time value.
	createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("GetUserByEmail (parse created_at): %w", err)
	}
	user.CreatedAt = createdAt

	return user, nil
}

// GetAllUsers retrieves all users.
func (r *Repository) GetAllUsers() ([]types.User, error) {
	query := `SELECT id, name, email, role, created_at FROM users`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: %w", err)
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var u types.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetAllUsers scan: %w", err)
		}
		users = append(users, u)
	}
	return users, nil
}

// UpdateUser updates an existing user.
func (r *Repository) UpdateUser(u *types.User) error {
	query := `UPDATE users SET name = ?, email = ?, role = ? WHERE id = ?`
	_, err := r.DB.Exec(query, u.Name, u.Email, u.Role, u.ID)
	if err != nil {
		return fmt.Errorf("UpdateUser: %w", err)
	}
	return nil
}

// DeleteUser removes a user by ID.
func (r *Repository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("DeleteUser: %w", err)
	}
	return nil
}

// =============================
//    DEPARTMENTS CRUD
// =============================

// CreateDepartment inserts a new department.
func (r *Repository) CreateDepartment(d *types.Department) error {
	query := `INSERT INTO departments (name) VALUES (?)`
	result, err := r.DB.Exec(query, d.Name)
	if err != nil {
		return fmt.Errorf("CreateDepartment: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("CreateDepartment (retrieve id): %w", err)
	}
	d.ID = int(id)
	return nil
}

// GetDepartmentByID retrieves a department by ID.
func (r *Repository) GetDepartmentByID(id int) (*types.Department, error) {
	query := `SELECT id, name FROM departments WHERE id = ?`
	d := &types.Department{}
	if err := r.DB.QueryRow(query, id).Scan(&d.ID, &d.Name); err != nil {
		return nil, fmt.Errorf("GetDepartmentByID: %w", err)
	}
	return d, nil
}

// GetAllDepartments retrieves all departments.
func (r *Repository) GetAllDepartments() ([]types.Department, error) {
	query := `SELECT id, name FROM departments`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllDepartments: %w", err)
	}
	defer rows.Close()

	var departments []types.Department
	for rows.Next() {
		var d types.Department
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, fmt.Errorf("GetAllDepartments scan: %w", err)
		}
		departments = append(departments, d)
	}
	return departments, nil
}

// UpdateDepartment updates an existing department.
func (r *Repository) UpdateDepartment(d *types.Department) error {
	query := `UPDATE departments SET name = ? WHERE id = ?`
	_, err := r.DB.Exec(query, d.Name, d.ID)
	if err != nil {
		return fmt.Errorf("UpdateDepartment: %w", err)
	}
	return nil
}

// DeleteDepartment removes a department by ID.
func (r *Repository) DeleteDepartment(id int) error {
	query := `DELETE FROM departments WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("DeleteDepartment: %w", err)
	}
	return nil
}

// =============================
//        COURSES CRUD
// =============================

// CreateCourse inserts a new course.
func (r *Repository) CreateCourse(c *types.Course) error {
	query := `INSERT INTO courses (name, department_id) VALUES (?, ?)`
	result, err := r.DB.Exec(query, c.Name, c.DepartmentID)
	if err != nil {
		return fmt.Errorf("CreateCourse: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("CreateCourse (retrieve id): %w", err)
	}
	c.ID = int(id)
	return nil
}

// GetCourseByID retrieves a course by ID.
func (r *Repository) GetCourseByID(id int) (*types.Course, error) {
	query := `SELECT id, name, department_id FROM courses WHERE id = ?`
	c := &types.Course{}
	if err := r.DB.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.DepartmentID); err != nil {
		return nil, fmt.Errorf("GetCourseByID: %w", err)
	}
	return c, nil
}

// GetAllCourses retrieves all courses.
func (r *Repository) GetAllCourses() ([]types.Course, error) {
	query := `SELECT id, name, department_id FROM courses`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllCourses: %w", err)
	}
	defer rows.Close()

	var courses []types.Course
	for rows.Next() {
		var c types.Course
		if err := rows.Scan(&c.ID, &c.Name, &c.DepartmentID); err != nil {
			return nil, fmt.Errorf("GetAllCourses scan: %w", err)
		}
		courses = append(courses, c)
	}
	return courses, nil
}

// UpdateCourse updates an existing course.
func (r *Repository) UpdateCourse(c *types.Course) error {
	query := `UPDATE courses SET name = ?, department_id = ? WHERE id = ?`
	_, err := r.DB.Exec(query, c.Name, c.DepartmentID, c.ID)
	if err != nil {
		return fmt.Errorf("UpdateCourse: %w", err)
	}
	return nil
}

// DeleteCourse removes a course by ID.
func (r *Repository) DeleteCourse(id int) error {
	query := `DELETE FROM courses WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("DeleteCourse: %w", err)
	}
	return nil
}

// =============================
//      SYLLABI CRUD
// =============================

// CreateSyllabus inserts a new syllabus.
// Note: submission_date is stored as DATE; we format the time accordingly.
func (r *Repository) CreateSyllabus(s *types.Syllabus) error {
	query := `INSERT INTO syllabi (course_id, lecturer_id, status, submission_date, data) VALUES (?, ?, ?, ?, ?)`
	result, err := r.DB.Exec(query, s.CourseID, s.LecturerID, s.Status, s.SubmissionDate.Format("2006-01-02"), s.Data)
	if err != nil {
		return fmt.Errorf("CreateSyllabus: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("CreateSyllabus (retrieve id): %w", err)
	}
	s.ID = int(id)
	return nil
}
func (r *Repository) GetSyllabusByID(id int) (*types.Syllabus, error) {
	query := "SELECT id, course_id, lecturer_id, status, submission_date, created_at, updated_at, data FROM syllabi WHERE id = ?"
	row := r.DB.QueryRow(query, id)

	var syl types.Syllabus
	var submissionDateStr, createdAtStr, updatedAtStr string

	// Scan into all fields of the syllabus
	if err := row.Scan(&syl.ID, &syl.CourseID, &syl.LecturerID, &syl.Status, &submissionDateStr, &createdAtStr, &updatedAtStr, &syl.Data); err != nil {
		return nil, fmt.Errorf("GetSyllabusByID: %w", err)
	}

	// Parse the date strings into time.Time
	submissionDate, err := time.Parse("2006-01-02", submissionDateStr)
	if err != nil {
		return nil, fmt.Errorf("GetSyllabusByID: submission_date parse error: %w", err)
	}
	syl.SubmissionDate = submissionDate

	createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("GetSyllabusByID: created_at parse error: %w", err)
	}
	syl.CreatedAt = createdAt

	updatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("GetSyllabusByID: updated_at parse error: %w", err)
	}
	syl.UpdatedAt = updatedAt

	return &syl, nil
}

// GetSyllabiByLecturer fetches all syllabi for the given lecturer (user) ID.
func (r *Repository) GetSyllabiByLecturer(lecturerID int) ([]types.Syllabus, error) {
	query := `
        SELECT id, course_id, lecturer_id, status, submission_date, created_at, updated_at, data
        FROM syllabi
        WHERE lecturer_id = ? AND status != 'Deleted'
        ORDER BY submission_date DESC
    `
	rows, err := r.DB.Query(query, lecturerID)
	if err != nil {
		return nil, fmt.Errorf("GetSyllabiByLecturer: %w", err)
	}
	defer rows.Close()

	var syllabi []types.Syllabus
	for rows.Next() {
		var s types.Syllabus
		// If your driver does not directly scan into time.Time,
		// scan date values into temporary strings and parse them.
		var submissionDateStr, createdAtStr, updatedAtStr string
		if err := rows.Scan(&s.ID, &s.CourseID, &s.LecturerID, &s.Status, &submissionDateStr, &createdAtStr, &updatedAtStr, &s.Data); err != nil {
			return nil, fmt.Errorf("GetSyllabiByLecturer scan: %w", err)
		}
		s.SubmissionDate, err = time.Parse("2006-01-02", submissionDateStr)
		if err != nil {
			return nil, fmt.Errorf("parsing submission_date: %w", err)
		}
		s.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("parsing created_at: %w", err)
		}
		s.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("parsing updated_at: %w", err)
		}
		syllabi = append(syllabi, s)
	}
	return syllabi, nil
}

// GetAllSyllabi retrieves all syllabi.
func (r *Repository) GetAllSyllabi() ([]types.Syllabus, error) {
	query := `SELECT id, course_id, lecturer_id, status, submission_date, created_at, updated_at, data FROM syllabi WHERE status != 'Deleted'`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllSyllabi: %w", err)
	}
	defer rows.Close()

	var syllabi []types.Syllabus
	for rows.Next() {
		var s types.Syllabus
		var submissionDate string
		if err := rows.Scan(&s.ID, &s.CourseID, &s.LecturerID, &s.Status, &submissionDate, &s.CreatedAt, &s.UpdatedAt, &s.Data); err != nil {
			return nil, fmt.Errorf("GetAllSyllabi scan: %w", err)
		}
		parsed, err := time.Parse("2006-01-02", submissionDate)
		if err != nil {
			return nil, fmt.Errorf("GetAllSyllabi (parse submission_date): %w", err)
		}
		s.SubmissionDate = parsed
		syllabi = append(syllabi, s)
	}
	return syllabi, nil
}

// UpdateSyllabus updates an existing syllabus.
func (r *Repository) UpdateSyllabus(s *types.Syllabus) error {
	query := `UPDATE syllabi SET course_id = ?, lecturer_id = ?, status = ?, submission_date = ?, data = ? WHERE id = ?`
	_, err := r.DB.Exec(query, s.CourseID, s.LecturerID, s.Status, s.SubmissionDate.Format("2006-01-02"), s.Data, s.ID)
	if err != nil {
		return fmt.Errorf("UpdateSyllabus: %w", err)
	}
	return nil
}

// DeleteSyllabus removes a syllabus by ID.
func (r *Repository) DeleteSyllabus(id int) error {
	query := `DELETE FROM syllabi WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("DeleteSyllabus: %w", err)
	}
	return nil
}

// RETURN UI COMPONENTS
func (r *Repository) GetCardsByLecturer(lecturerID int) ([]UIcomponents.Card, error) {
	query := `
		SELECT s.id, s.status, s.submission_date, c.name AS courseName, d.name AS departmentName, u.name AS lecturerName
		FROM syllabi s
		JOIN courses c ON s.course_id = c.id
		JOIN departments d ON c.department_id = d.id
		JOIN users u ON s.lecturer_id = u.id
		WHERE s.lecturer_id = ? and status != 'Deleted'
		ORDER BY s.submission_date DESC
	`
	rows, err := r.DB.Query(query, lecturerID)
	if err != nil {
		return nil, fmt.Errorf("GetCardsByLecturer: %w", err)
	}
	defer rows.Close()

	var cards []UIcomponents.Card
	for rows.Next() {
		var (
			id                int
			status            string
			submissionDateStr string
			courseName        string
			departmentName    string
			lecturerName      string
		)

		if err := rows.Scan(&id, &status, &submissionDateStr, &courseName, &departmentName, &lecturerName); err != nil {
			return nil, fmt.Errorf("GetCardsByLecturer scan: %w", err)
		}

		submissionDate, err := time.Parse("2006-01-02", submissionDateStr)
		if err != nil {
			return nil, fmt.Errorf("parsing date: %w", err)
		}

		card := UIcomponents.Card{
			ID:          id,
			Date:        submissionDate.Format("02/01/2006"),
			Title:       courseName,
			Lecturer:    lecturerName,
			Field:       departmentName,
			Status:      status,
			StatusLabel: status,
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// GetDeletedCardsByLecturer fetches all deleted syllabi for the given lecturer (user) ID.
func (r *Repository) GetDeletedCardsByLecturer(lecturerID int) ([]UIcomponents.Card, error) {
	query := `
		SELECT s.id, s.status, s.submission_date, c.name AS courseName, d.name AS departmentName, u.name AS lecturerName
		FROM syllabi s
		JOIN courses c ON s.course_id = c.id
		JOIN departments d ON c.department_id = d.id
		JOIN users u ON s.lecturer_id = u.id
		WHERE s.lecturer_id = ? and status = 'Deleted'
		ORDER BY s.submission_date DESC
	`
	rows, err := r.DB.Query(query, lecturerID)
	if err != nil {
		return nil, fmt.Errorf("GetDeletedCardsByLecturer: %w", err)
	}
	defer rows.Close()

	var cards []UIcomponents.Card
	for rows.Next() {
		var (
			id                int
			status            string
			submissionDateStr string
			courseName        string
			departmentName    string
			lecturerName      string
		)

		if err := rows.Scan(&id, &status, &submissionDateStr, &courseName, &departmentName, &lecturerName); err != nil {
			return nil, fmt.Errorf("GetDeletedCardsByLecturer scan: %w", err)
		}

		submissionDate, err := time.Parse("2006-01-02", submissionDateStr)
		if err != nil {
			return nil, fmt.Errorf("parsing date: %w", err)
		}

		card := UIcomponents.Card{
			ID:          id,
			Date:        submissionDate.Format("02/01/2006"),
			Title:       courseName,
			Lecturer:    lecturerName,
			Field:       departmentName,
			Status:      status,
			StatusLabel: status,
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (r *Repository) FilterCardsByLecturer(lecturerID int, search, fromDate, toDate string, statuses []string) ([]map[string]interface{}, error) {
	baseQuery := `
		SELECT s.status, s.submission_date, c.name AS courseName, d.name AS departmentName, u.name AS lecturerName
		FROM syllabi s
		JOIN courses c ON s.course_id = c.id
		JOIN departments d ON c.department_id = d.id
		JOIN users u ON s.lecturer_id = u.id
		WHERE s.lecturer_id = ? and s.status != 'Deleted'
	`
	params := []interface{}{lecturerID}

	// Add search filter if provided.
	if search != "" {
		baseQuery += " AND (c.name LIKE ? OR d.name LIKE ? OR u.name LIKE ?)"
		searchParam := "%" + search + "%"
		params = append(params, searchParam, searchParam, searchParam)
	}

	// Filter by from-date if provided.
	if fromDate != "" {
		parsedFrom, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			return nil, fmt.Errorf("invalid from-date format: %w", err)
		}
		baseQuery += " AND s.submission_date >= ?"
		params = append(params, parsedFrom.Format("2006-01-02"))
	}

	// Filter by to-date if provided.
	if toDate != "" {
		parsedTo, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			return nil, fmt.Errorf("invalid to-date format: %w", err)
		}
		baseQuery += " AND s.submission_date <= ?"
		params = append(params, parsedTo.Format("2006-01-02"))
	}

	// Filter by statuses if provided.
	if len(statuses) > 0 {
		placeholders := make([]string, len(statuses))
		for i := range statuses {
			placeholders[i] = "?"
		}
		baseQuery += " AND s.status IN (" + strings.Join(placeholders, ",") + ")"
		for _, s := range statuses {
			params = append(params, s)
		}
	}

	// Order results by submission_date in descending order.
	baseQuery += " ORDER BY s.submission_date DESC"

	// Execute the query.
	rows, err := r.DB.Query(baseQuery, params...)
	if err != nil {
		return nil, fmt.Errorf("FilterCardsByLecturer: %w", err)
	}
	defer rows.Close()

	var cards []map[string]interface{}
	for rows.Next() {
		var status, courseName, departmentName, lecturerName string
		var submissionDateStr string

		if err := rows.Scan(&status, &submissionDateStr, &courseName, &departmentName, &lecturerName); err != nil {
			return nil, fmt.Errorf("FilterCardsByLecturer scan: %w", err)
		}

		// Parse the submission date.
		dt, err := time.Parse("2006-01-02", submissionDateStr)
		if err != nil {
			return nil, fmt.Errorf("parsing submission_date: %w", err)
		}

		card := map[string]interface{}{
			"date":     dt.Format("02/01/2006"),
			"title":    courseName,
			"lecturer": lecturerName,
			"field":    departmentName,
			"status":   status,
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// InsertNewSyllabus converts the draft into a syllabus record and inserts it into the database.
func (r *Repository) InsertSyllabusFromDraft(userID int, draft *UIcomponents.Draft) error {
	// Create the base syllabus record.
	now := time.Now()
	syl := types.Syllabus{
		ID:             0, // dummy value; DB will auto-increment
		CourseID:       0,
		LecturerID:     userID,
		Status:         "Draft",
		SubmissionDate: now, // Save the current time
		CreatedAt:      now,
		UpdatedAt:      now,
		Data:           nil, // to be set below
	}

	// Look up the CourseID from the available courses.
	courses, err := r.GetAllCourses()
	if err != nil {
		return err
	}

	for _, course := range courses {
		if course.Name == draft.SelectedCourse {
			syl.CourseID = course.ID
			break
		}
	}

	// Encode the draft into JSON.
	jsonData, err := json.Marshal(draft)
	if err != nil {
		return err
	}
	syl.Data = jsonData

	// When inserting, format the submission_date to "YYYY-MM-DD"
	err = r.CreateSyllabus(&syl)
	if err != nil {
		return err
	}

	return nil
}

// GetCommentsBySyllabusID fetches all comments associated with the given syllabus ID.
func (r *Repository) GetCommentsBySyllabusID(syllabusID int) ([]types.Comment, error) {
	query := `
        SELECT id, syllabus_id, user_id, content, created_at, updated_at 
        FROM comments 
        WHERE syllabus_id = ? 
        ORDER BY created_at ASC
    `
	rows, err := r.DB.Query(query, syllabusID)
	if err != nil {
		return nil, fmt.Errorf("GetCommentsBySyllabusID: %w", err)
	}
	defer rows.Close()

	var comments []types.Comment

	for rows.Next() {
		var c types.Comment
		var createdAtStr, updatedAtStr string

		if err := rows.Scan(&c.ID, &c.SyllabusID, &c.UserID, &c.Content, &createdAtStr, &updatedAtStr); err != nil {
			return nil, fmt.Errorf("GetCommentsBySyllabusID scan: %w", err)
		}

		// Parse the timestamps from string to time.Time.
		c.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("GetCommentsBySyllabusID: parsing created_at: %w", err)
		}
		c.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("GetCommentsBySyllabusID: parsing updated_at: %w", err)
		}

		comments = append(comments, c)
	}

	return comments, nil
}

// AddComment inserts a new comment associated with a syllabus into the DB.
func (r *Repository) AddComment(c *types.Comment) error {
	query := `INSERT INTO comments (syllabus_id, user_id, content) VALUES (?, ?, ?)`
	result, err := r.DB.Exec(query, c.SyllabusID, c.UserID, c.Content)
	if err != nil {
		return fmt.Errorf("AddComment: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("AddComment (retrieve id): %w", err)
	}
	c.ID = int(id)
	return nil
}

// GetUserDraft retrieves a user's draft from the database.
func (r *Repository) GetUserDraft(userID int) (*UIcomponents.Draft, error) {
	// Check if there's a draft syllabus for this user
	query := `
		SELECT id, data 
		FROM syllabi 
		WHERE lecturer_id = ? AND status = 'Draft' 
		ORDER BY updated_at DESC 
		LIMIT 1
	`
	row := r.DB.QueryRow(query, userID)

	var id int
	var data []byte
	err := row.Scan(&id, &data)
	if err != nil {
		if err == sql.ErrNoRows {
			// No draft found, create a new one
			user, err := r.GetUserByID(userID)
			if err != nil {
				return nil, fmt.Errorf("GetUserDraft (get user): %w", err)
			}

			draft := &UIcomponents.Draft{
				ID:                      -1,
				LecturerName:            user.Name,
				LecturerEmail:           user.Email,
				CourseRequirements:      []string{},
				LearningOutcomes:        []string{},
				CourseObjectives:        []string{},
				CourseStructure:         []string{},
				AssignmentsStructure:    []string{},
				SyllabusRows:            []UIcomponents.SyllabusRow{{}},
				GradeComponents:         []UIcomponents.GradeComponent{},
				BibliographyRequired:    []string{},
				BibliographyRecommended: []string{},
			}

			// Store the new draft in the database
			jsonData, err := json.Marshal(draft)
			if err != nil {
				return nil, fmt.Errorf("GetUserDraft (marshal): %w", err)
			}

			now := time.Now()
			syl := types.Syllabus{
				ID:             0,
				CourseID:       1,
				LecturerID:     userID,
				Status:         "Draft",
				SubmissionDate: now,
				CreatedAt:      now,
				UpdatedAt:      now,
				Data:           jsonData,
			}

			err = r.CreateSyllabus(&syl)
			if err != nil {
				return nil, fmt.Errorf("GetUserDraft (create): %w", err)
			}

			draft.ID = syl.ID
			return draft, nil
		}
		return nil, fmt.Errorf("GetUserDraft: %w", err)
	}

	// Unmarshal the data into a Draft
	var draft UIcomponents.Draft
	err = json.Unmarshal(data, &draft)
	if err != nil {
		return nil, fmt.Errorf("GetUserDraft (unmarshal): %w", err)
	}

	draft.ID = id
	return &draft, nil
}

// CreateNewUserDraft creates a new draft for a user regardless of whether they already have one.
func (r *Repository) CreateNewUserDraft(userID int) (*UIcomponents.Draft, error) {
	// Get user information
	user, err := r.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("CreateNewUserDraft (get user): %w", err)
	}

	// Create a new draft
	draft := &UIcomponents.Draft{
		ID:                      -1,
		LecturerName:            user.Name,
		LecturerEmail:           user.Email,
		CourseRequirements:      []string{},
		LearningOutcomes:        []string{},
		CourseObjectives:        []string{},
		CourseStructure:         []string{},
		AssignmentsStructure:    []string{},
		SyllabusRows:            []UIcomponents.SyllabusRow{{}},
		GradeComponents:         []UIcomponents.GradeComponent{},
		BibliographyRequired:    []string{},
		BibliographyRecommended: []string{},
	}

	// Store the new draft in the database
	jsonData, err := json.Marshal(draft)
	if err != nil {
		return nil, fmt.Errorf("CreateNewUserDraft (marshal): %w", err)
	}

	now := time.Now()
	syl := types.Syllabus{
		ID:             0,
		CourseID:       1,
		LecturerID:     userID,
		Status:         "Draft",
		SubmissionDate: now,
		CreatedAt:      now,
		UpdatedAt:      now,
		Data:           jsonData,
	}

	err = r.CreateSyllabus(&syl)
	if err != nil {
		return nil, fmt.Errorf("CreateNewUserDraft (create): %w", err)
	}

	draft.ID = syl.ID
	return draft, nil
}

// SaveUserDraft saves a user's draft to the database.
func (r *Repository) SaveUserDraft(userID int, draft *UIcomponents.Draft) error {
	// Marshal the draft to JSON
	jsonData, err := json.Marshal(draft)
	if err != nil {
		return fmt.Errorf("SaveUserDraft (marshal): %w", err)
	}

	// Check if this is a new draft or an existing one
	if draft.ID == -1 {
		// Create a new draft
		now := time.Now()
		syl := types.Syllabus{
			ID:             0,
			CourseID:       0,
			LecturerID:     userID,
			Status:         "Draft",
			SubmissionDate: now,
			CreatedAt:      now,
			UpdatedAt:      now,
			Data:           jsonData,
		}

		err = r.CreateSyllabus(&syl)
		if err != nil {
			return fmt.Errorf("SaveUserDraft (create): %w", err)
		}

		draft.ID = syl.ID
	} else {
		// Update existing draft
		query := `UPDATE syllabi SET data = ?, updated_at = ? WHERE id = ?`
		_, err := r.DB.Exec(query, jsonData, time.Now().Format("2006-01-02 15:04:05"), draft.ID)
		if err != nil {
			return fmt.Errorf("SaveUserDraft (update): %w", err)
		}
	}

	return nil
}

// GetEditedSyllabus retrieves a syllabus by ID and unmarshals its data into a Draft.
func (r *Repository) GetEditedSyllabus(syllabusID int) (*UIcomponents.Draft, error) {
	syl, err := r.GetSyllabusByID(syllabusID)
	if err != nil {
		return nil, fmt.Errorf("GetEditedSyllabus: %w", err)
	}

	var draft UIcomponents.Draft
	err = json.Unmarshal(syl.Data, &draft)
	if err != nil {
		return nil, fmt.Errorf("GetEditedSyllabus (unmarshal): %w", err)
	}

	draft.ID = syl.ID
	return &draft, nil
}
