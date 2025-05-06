package handler

import (
	"Syllybea/UIcomponents"
	"Syllybea/mid"
	"Syllybea/repository"
	"Syllybea/types"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

func HandleEditSyllabus(c echo.Context, repo *repository.Repository) error {
	// Parse the syllabus ID from the URL parameter.
	syllabusID := c.Param("id")
	id, err := strconv.Atoi(syllabusID)
	if err != nil {
		c.Logger().Error("Invalid syllabus ID: ", err)
		return c.String(http.StatusBadRequest, "Invalid syllabus ID")
	}

	// Get the edited draft directly from the repository
	draft, err := repo.GetEditedSyllabus(id)
	if err != nil {
		c.Logger().Error("Error retrieving syllabus: ", err)
		return c.String(http.StatusInternalServerError, "Error retrieving syllabus")
	}

	if draft == nil {
		c.Logger().Warn("Syllabus not found with id: ", id)
		return c.String(http.StatusNotFound, "Syllabus not found")
	}

	return c.Render(http.StatusOK, "create-syllabus", draft)
}

func HandleCreateSyllabus(c echo.Context, repo *repository.Repository) error {
	userID, _ := mid.GetUserID(c)
	draft, err := repo.CreateNewUserDraft(userID)
	if err != nil {
		c.Logger().Error("Error creating new user draft: ", err)
		return c.String(http.StatusInternalServerError, "Error creating new user draft")
	}

	departments, _ := r.GetAllDepartments()
	deptNames := make([]string, 0, len(departments))
	for _, dept := range departments {
		deptNames = append(deptNames, dept.Name)
	}
	draft.Departments = deptNames
	if draft.SyllabusDepartment == "" && len(deptNames) > 0 {
		draft.SyllabusDepartment = deptNames[0]
	}

	courses, _ := r.GetAllCourses()
	courseNames := make([]string, 0, len(courses))
	for _, course := range courses {
		courseNames = append(courseNames, course.Name)
	}
	draft.Courses = courseNames
	if draft.SelectedCourse == "" && len(courseNames) > 0 {
		draft.SelectedCourse = courseNames[0]
	}

	// Save the updated draft
	err = repo.SaveUserDraft(userID, draft)
	if err != nil {
		c.Logger().Error("Error saving user draft: ", err)
		return c.String(http.StatusInternalServerError, "Error saving user draft")
	}

	return c.Render(http.StatusOK, "create-syllabus", draft)
}

func handleSaveSyllabus(c echo.Context, repo *repository.Repository) error {
	userID, _ := mid.GetUserID(c)

	// Get the user's draft from the database
	draft, err := repo.GetUserDraft(userID)
	if err != nil {
		c.Logger().Error("Error getting user draft: ", err)
		return c.String(http.StatusInternalServerError, "Error getting user draft")
	}

	// Update the draft with form data
	if err := c.Request().ParseForm(); err == nil {
		// Update lecturer details
		if lecturerName := c.FormValue("LecturerName"); lecturerName != "" {
			draft.LecturerName = lecturerName
		}
		if lecturerEmail := c.FormValue("LecturerEmail"); lecturerEmail != "" {
			draft.LecturerEmail = lecturerEmail
		}
		if officeDay := c.FormValue("office-day"); officeDay != "" {
			draft.OfficeDay = officeDay
		}
		if officeStart := c.FormValue("office-start"); officeStart != "" {
			draft.OfficeStart = officeStart
		}
		if officeEnd := c.FormValue("office-end"); officeEnd != "" {
			draft.OfficeEnd = officeEnd
		}

		// Update course details
		if selectedCourse := c.FormValue("course-dropdown"); selectedCourse != "" {
			draft.SelectedCourse = selectedCourse
		}
		if syllabusDepartment := c.FormValue("syllabus-department"); syllabusDepartment != "" {
			draft.SyllabusDepartment = syllabusDepartment
		}
		if credits := c.FormValue("credits"); credits != "" {
			draft.Credits = credits
		}
		if weeklyHours := c.FormValue("weekly-hours"); weeklyHours != "" {
			draft.WeeklyHours = weeklyHours
		}
		if year := c.FormValue("year"); year != "" {
			draft.Year = year
		}
		if semester := c.FormValue("semester"); semester != "" {
			draft.Semester = semester
		}
		if prerequisites := c.FormValue("prerequisites"); prerequisites != "" {
			draft.Prerequisites = prerequisites
		}

		// Update arrays
		draft.CourseRequirements = c.Request().Form["course-requirements[]"]
		draft.LearningOutcomes = c.Request().Form["learning-outcomes[]"]
		draft.CourseObjectives = c.Request().Form["course-objectives[]"]
		draft.BibliographyRequired = c.Request().Form["bibliography-required[]"]
		draft.BibliographyRecommended = c.Request().Form["bibliography-recommended[]"]
	}

	// Save the draft as a syllabus with "Draft" status
	// First, marshal the draft to JSON
	jsonData, err := json.Marshal(draft)
	if err != nil {
		c.Logger().Error("Error marshaling draft: ", err)
		return c.String(http.StatusInternalServerError, "Error processing syllabus data")
	}

	// Check if this syllabus already exists in the database
	var existingSyllabus *types.Syllabus
	if draft.ID > 0 {
		existingSyllabus, err = repo.GetSyllabusByID(draft.ID)
		if err != nil {
			c.Logger().Error("Error checking for existing syllabus: ", err)
			// Continue with creating a new syllabus
			existingSyllabus = nil
		}
	}

	// Look up the CourseID from the available courses
	courses, err := repo.GetAllCourses()
	if err != nil {
		c.Logger().Error("Error getting courses: ", err)
		return c.String(http.StatusInternalServerError, "Error processing syllabus data")
	}

	var courseID int
	for _, course := range courses {
		if course.Name == draft.SelectedCourse {
			courseID = course.ID
			break
		}
	}

	now := time.Now()

	if existingSyllabus != nil {
		// Update the existing syllabus
		existingSyllabus.CourseID = courseID
		existingSyllabus.LecturerID = userID // Set the lecturer ID to the current user's ID
		existingSyllabus.Status = "Draft"
		existingSyllabus.UpdatedAt = now
		existingSyllabus.Data = jsonData

		err = repo.UpdateSyllabus(existingSyllabus)
		if err != nil {
			c.Logger().Error("Error updating syllabus: ", err)
			return c.String(http.StatusInternalServerError, "Error saving syllabus")
		}
	} else {
		// Create a new syllabus with the "Draft" status
		syl := types.Syllabus{
			CourseID:       courseID,
			LecturerID:     userID,
			Status:         "Draft",
			SubmissionDate: now,
			CreatedAt:      now,
			UpdatedAt:      now,
			Data:           jsonData,
		}

		// Create the syllabus in the database
		err = repo.CreateSyllabus(&syl)
		if err != nil {
			c.Logger().Error("Error saving syllabus: ", err)
			return c.String(http.StatusInternalServerError, "Error saving syllabus")
		}
	}

	c.Response().Header().Set("HX-Redirect", "/dashboard")
	return c.String(http.StatusOK, "Redirecting...")
}

// hasEmptyFields checks if any required fields in the draft are empty
func hasEmptyFields(draft *UIcomponents.Draft) bool {
	// Check lecturer details
	if draft.LecturerName == "" || draft.LecturerEmail == "" {
		return true
	}

	// Check course details
	if draft.SelectedCourse == "" || draft.SyllabusDepartment == "" {
		return true
	}

	// Check if there are any course requirements
	if len(draft.CourseRequirements) == 0 || (len(draft.CourseRequirements) == 1 && draft.CourseRequirements[0] == "") {
		return true
	}

	// Check if there are any learning outcomes
	if len(draft.LearningOutcomes) == 0 || (len(draft.LearningOutcomes) == 1 && draft.LearningOutcomes[0] == "") {
		return true
	}

	// Check if there are any course objectives
	if len(draft.CourseObjectives) == 0 || (len(draft.CourseObjectives) == 1 && draft.CourseObjectives[0] == "") {
		return true
	}

	// Check active learning fields
	if draft.ActiveLearning1 == "" || draft.ActiveLearning2 == "" ||
		draft.ActiveLearning3 == "" || draft.ActiveLearning4 == "" {
		return true
	}

	// Check if there are any syllabus rows with content
	hasContent := false
	for _, row := range draft.SyllabusRows {
		if row.MainTopic != "" || row.LessonTopics != "" || row.Subtopics != "" {
			hasContent = true
			break
		}
	}
	if !hasContent {
		return true
	}

	// Check if there are any grade components
	if len(draft.GradeComponents) == 0 {
		return true
	}

	return false
}

func handleSubmitSyllabus(c echo.Context, repo *repository.Repository) error {
	userID, _ := mid.GetUserID(c)

	// Get the user's draft from the database
	draft, err := repo.GetUserDraft(userID)
	if err != nil {
		c.Logger().Error("Error getting user draft: ", err)
		return c.String(http.StatusInternalServerError, "Error getting user draft")
	}

	//// Validate that all required fields are filled
	//if hasEmptyFields(draft) {
	//	return c.String(http.StatusBadRequest, "some fields are empty, can't send")
	//}

	// Marshal the draft to JSON
	jsonData, err := json.Marshal(draft)
	if err != nil {
		c.Logger().Error("Error marshaling draft: ", err)
		return c.String(http.StatusInternalServerError, "Error processing syllabus data")
	}

	// Look up the CourseID from the available courses
	courseID := 0
	courses, err := repo.GetAllCourses()
	if err != nil {
		c.Logger().Error("Error getting courses: ", err)
		return c.String(http.StatusInternalServerError, "Error processing syllabus data")
	}

	for _, course := range courses {
		if course.Name == draft.SelectedCourse {
			courseID = course.ID
			break
		}
	}

	// Update the existing draft syllabus to "In Review" status
	now := time.Now()
	syl := types.Syllabus{
		ID:             draft.ID,
		CourseID:       courseID,
		LecturerID:     userID,
		Status:         "In Review", // Change status from "Draft" to "In Review"
		SubmissionDate: now,
		UpdatedAt:      now,
		Data:           jsonData,
	}

	// Update the syllabus in the database
	err = repo.UpdateSyllabus(&syl)
	if err != nil {
		c.Logger().Error("Error updating syllabus: ", err)
		return c.String(http.StatusInternalServerError, "Error updating syllabus")
	}

	c.Response().Header().Set("HX-Redirect", "/dashboard")
	return c.String(http.StatusOK, "Redirecting...")
}

func removeBibliographyRecommended(c echo.Context, draft *UIcomponents.Draft) error {
	indexStr := c.FormValue("index")
	if index, err := strconv.Atoi(indexStr); err == nil && index >= 0 && index < len(draft.BibliographyRecommended) {
		draft.BibliographyRecommended = append(draft.BibliographyRecommended[:index], draft.BibliographyRecommended[index+1:]...)
	}
	return c.Render(http.StatusOK, "bibliographyRecommended", draft)
}

func addBibliographyRecommended(c echo.Context, draft *UIcomponents.Draft) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	draft.BibliographyRecommended = c.Request().Form["bibliography-recommended[]"]
	draft.BibliographyRecommended = append(draft.BibliographyRecommended, "")
	return c.Render(http.StatusOK, "bibliographyRecommended", draft)
}

func removeBibliographyRequired(c echo.Context, draft *UIcomponents.Draft) error {
	indexStr := c.FormValue("index")
	if index, err := strconv.Atoi(indexStr); err == nil && index >= 0 && index < len(draft.BibliographyRequired) {
		draft.BibliographyRequired = append(draft.BibliographyRequired[:index], draft.BibliographyRequired[index+1:]...)
	}
	return c.Render(http.StatusOK, "bibliographyRequired", draft)
}

func addBibliographyRequired(c echo.Context, draft *UIcomponents.Draft) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	// Update the slice based on current form values
	draft.BibliographyRequired = c.Request().Form["bibliography-required[]"]
	// Add a new empty entry
	draft.BibliographyRequired = append(draft.BibliographyRequired, "")
	// Render only the partial template for required bibliography
	return c.Render(http.StatusOK, "bibliographyRequired", draft)
}

func removeAssignmentStructure(c echo.Context, draft *UIcomponents.Draft) error {
	indexStr := c.FormValue("index")
	if index, err := strconv.Atoi(indexStr); err == nil && index >= 0 && index < len(draft.AssignmentsStructure) {
		draft.AssignmentsStructure = append(draft.AssignmentsStructure[:index], draft.AssignmentsStructure[index+1:]...)
	}
	return c.Render(http.StatusOK, "assignmentsStructure", draft)
}

func addAssignmentStructure(c echo.Context, draft *UIcomponents.Draft) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	draft.AssignmentsStructure = c.Request().Form["assignments-structure[]"]
	// Append an empty entry for a new assignment.
	draft.AssignmentsStructure = append(draft.AssignmentsStructure, "")
	return c.Render(http.StatusOK, "assignmentsStructure", draft)
}

func removeGradeComponent(c echo.Context, draft *UIcomponents.Draft) error {
	indexStr := c.FormValue("index")
	if index, err := strconv.Atoi(indexStr); err == nil && index >= 0 && index < len(draft.GradeComponents) {
		draft.GradeComponents = append(draft.GradeComponents[:index], draft.GradeComponents[index+1:]...)
	}
	return c.Render(http.StatusOK, "gradeComponents", draft)
}

func addGradeComponent(c echo.Context, draft *UIcomponents.Draft) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	names := c.Request().Form["grade-component-name[]"]
	percents := c.Request().Form["grade-component-percentage[]"]
	var comps []UIcomponents.GradeComponent
	count := len(names)
	for i := 0; i < count; i++ {
		comps = append(comps, UIcomponents.GradeComponent{
			PartName:   names[i],
			Percentage: percents[i],
		})
	}
	draft.GradeComponents = comps
	// Append an extra empty row.
	draft.GradeComponents = append(draft.GradeComponents, UIcomponents.GradeComponent{})
	return c.Render(http.StatusOK, "gradeComponents", draft)
}

func removeCourseObjective(c echo.Context, draft *UIcomponents.Draft) error {
	indexStr := c.FormValue("index")
	if index, err := strconv.Atoi(indexStr); err == nil && index >= 0 && index < len(draft.CourseObjectives) {
		draft.CourseObjectives = append(draft.CourseObjectives[:index], draft.CourseObjectives[index+1:]...)
	}
	return c.Render(http.StatusOK, "courseObjectives", draft)

}

func addCourseObjective(c echo.Context, draft *UIcomponents.Draft) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	draft.CourseObjectives = c.Request().Form["course-objectives[]"]
	draft.CourseObjectives = append(draft.CourseObjectives, "")
	return c.Render(http.StatusOK, "courseObjectives", draft)
}

func removeLearningOutcome(c echo.Context, draft *UIcomponents.Draft) error {
	indexStr := c.FormValue("index")
	if index, err := strconv.Atoi(indexStr); err == nil && index >= 0 && index < len(draft.LearningOutcomes) {
		draft.LearningOutcomes = append(draft.LearningOutcomes[:index], draft.LearningOutcomes[index+1:]...)
	}
	return c.Render(http.StatusOK, "learningOutcomes", draft)
}

func addLearningOutcome(c echo.Context, draft *UIcomponents.Draft) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	draft.LearningOutcomes = c.Request().Form["learning-outcomes[]"]
	draft.LearningOutcomes = append(draft.LearningOutcomes, "")
	return c.Render(http.StatusOK, "learningOutcomes", draft)
}

func insertSyllabusRow(c echo.Context, draft *UIcomponents.Draft) error {
	// Get the insertion index from the form values.
	indexStr := c.FormValue("index")
	insertIndex, err := strconv.Atoi(indexStr)
	if err != nil {
		return err
	}
	// Update the current rows from the form data.
	lessonNumbers := c.Request().Form["lesson-number[]"]
	mainTopics := c.Request().Form["main-topic[]"]
	lessonTopics := c.Request().Form["lesson-topics[]"]
	subtopics := c.Request().Form["subtopics[]"]
	readingMaterials := c.Request().Form["reading-material[]"]

	var rows []UIcomponents.SyllabusRow
	count := len(lessonNumbers)
	for i := 0; i < count; i++ {
		rows = append(rows, UIcomponents.SyllabusRow{
			LessonNumber:    lessonNumbers[i],
			MainTopic:       mainTopics[i],
			LessonTopics:    lessonTopics[i],
			Subtopics:       subtopics[i],
			ReadingMaterial: readingMaterials[i],
		})
	}
	// Create a new empty row.
	newRow := UIcomponents.SyllabusRow{}
	// Insert the new row after the specified index.
	if insertIndex+1 >= len(rows) {
		rows = append(rows, newRow)
	} else {
		rows = append(rows[:insertIndex+1], append([]UIcomponents.SyllabusRow{newRow}, rows[insertIndex+1:]...)...)
	}
	draft.SyllabusRows = rows
	return c.Render(http.StatusOK, "syllabusRows", draft)
}

func removeSyllabusRow(c echo.Context, draft *UIcomponents.Draft) error {
	// Prevent deletion if there's only one row left
	if len(draft.SyllabusRows) <= 1 {
		return c.Render(http.StatusOK, "syllabusRows", draft) // No changes, just re-render
	}

	indexStr := c.FormValue("index")
	if index, err := strconv.Atoi(indexStr); err == nil && index >= 0 && index < len(draft.SyllabusRows) {
		draft.SyllabusRows = append(draft.SyllabusRows[:index], draft.SyllabusRows[index+1:]...)
	}
	return c.Render(http.StatusOK, "syllabusRows", draft)
}

func updateSyllabusRow(c echo.Context, draft *UIcomponents.Draft) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	lessonNumbers := c.Request().Form["lesson-number[]"]
	mainTopics := c.Request().Form["main-topic[]"]
	lessonTopics := c.Request().Form["lesson-topics[]"]
	subtopics := c.Request().Form["subtopics[]"]
	readingMaterials := c.Request().Form["reading-material[]"]

	var rows []UIcomponents.SyllabusRow
	count := len(lessonNumbers)
	for i := 0; i < count; i++ {
		var rows []UIcomponents.SyllabusRow
		rows = append(rows, UIcomponents.SyllabusRow{
			LessonNumber:    lessonNumbers[i],
			MainTopic:       mainTopics[i],
			LessonTopics:    lessonTopics[i],
			Subtopics:       subtopics[i],
			ReadingMaterial: readingMaterials[i],
		})
	}
	draft.SyllabusRows = rows
	return c.Render(http.StatusOK, "syllabusRows", draft)
}

func removeCourseRequirement(c echo.Context, draft *UIcomponents.Draft) error {
	indexStr := c.FormValue("index")
	if index, err := strconv.Atoi(indexStr); err == nil && index >= 0 && index < len(draft.CourseRequirements) {
		draft.CourseRequirements = append(draft.CourseRequirements[:index], draft.CourseRequirements[index+1:]...)
	}
	return c.Render(http.StatusOK, "courseRequirements", draft)
}

func addCourseRequirement(c echo.Context, draft *UIcomponents.Draft) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	draft.CourseRequirements = c.Request().Form["course-requirements[]"]
	draft.CourseRequirements = append(draft.CourseRequirements, "")
	return c.Render(http.StatusOK, "courseRequirements", draft)
}

func updateSyllabusHandler(c echo.Context) error {
	userID, _ := mid.GetUserID(c)

	// Get the user's draft from the database
	draft, err := r.GetUserDraft(userID)
	if err != nil {
		c.Logger().Error("Error getting user draft: ", err)
		return c.String(http.StatusInternalServerError, "Error getting user draft")
	}

	var result error

	switch c.FormValue("action") {
	case "addCourseRequirement":
		result = addCourseRequirement(c, draft)
	case "removeCourseRequirement":
		result = removeCourseRequirement(c, draft)
	case "updateSyllabusRow":
		result = updateSyllabusRow(c, draft)
	case "removeSyllabusRow":
		result = removeSyllabusRow(c, draft)
	case "insertSyllabusRow":
		result = insertSyllabusRow(c, draft)
	case "addLearningOutcome":
		result = addLearningOutcome(c, draft)
	case "removeLearningOutcome":
		result = removeLearningOutcome(c, draft)
	case "addCourseObjective":
		result = addCourseObjective(c, draft)
	case "removeCourseObjective":
		result = removeCourseObjective(c, draft)
	case "addGradeComponent":
		result = addGradeComponent(c, draft)
	case "removeGradeComponent":
		result = removeGradeComponent(c, draft)
	case "addAssignmentStructure":
		result = addAssignmentStructure(c, draft)
	case "removeAssignmentStructure":
		result = removeAssignmentStructure(c, draft)
	case "addBibliographyRequired":
		result = addBibliographyRequired(c, draft)
	case "removeBibliographyRequired":
		result = removeBibliographyRequired(c, draft)
	case "addBibliographyRecommended":
		result = addBibliographyRecommended(c, draft)
	case "removeBibliographyRecommended":
		result = removeBibliographyRecommended(c, draft)
	default:
		result = handleGeneralUpdate(c, draft)
	}

	// Save the updated draft to the database
	if err := r.SaveUserDraft(userID, draft); err != nil {
		c.Logger().Error("Error saving user draft: ", err)
		return c.String(http.StatusInternalServerError, "Error saving user draft")
	}

	return result
}

func handleGeneralUpdate(c echo.Context, draft *UIcomponents.Draft) error {
	c.Logger().Print(draft)
	updateField := c.FormValue("updateField")
	switch updateField {
	case "syllabus-department":
		c.Logger().Print(c.FormValue("syllabus-department"))
		draft.SyllabusDepartment = c.FormValue("syllabus-department")
		return c.Render(http.StatusOK, "syllabusDepartment", draft)

	case "bibliographyRequired":
		if err := c.Request().ParseForm(); err == nil {
			draft.BibliographyRequired = c.Request().Form["bibliography-required[]"]
		}
		return c.Render(http.StatusOK, "bibliographyRequired", draft)
	case "bibliographyRecommended":
		if err := c.Request().ParseForm(); err == nil {
			draft.BibliographyRecommended = c.Request().Form["bibliography-recommended[]"]
		}
		return c.Render(http.StatusOK, "bibliographyRecommended", draft)
	case "course-dropdown":
		draft.SelectedCourse = c.FormValue("course-dropdown")
		return c.Render(http.StatusOK, "coursesDropdown", draft)
	case "LecturerName":
		draft.LecturerName = c.FormValue("LecturerName")
	case "LecturerEmail":
		draft.LecturerEmail = c.FormValue("LecturerEmail")
	case "officeDay":
		draft.OfficeDay = c.FormValue("office-day")
	case "officeStart":
		draft.OfficeStart = c.FormValue("office-start")
	case "officeEnd":
		draft.OfficeEnd = c.FormValue("office-end")
	case "credits":
		draft.Credits = c.FormValue("credits")
	case "weeklyHours":
		draft.WeeklyHours = c.FormValue("weekly-hours")
	case "year":
		draft.Year = c.FormValue("year")
	case "semester":
		draft.Semester = c.FormValue("semester")
	case "prerequisites":
	case "learningOutcomes":
		if err := c.Request().ParseForm(); err == nil {
			draft.LearningOutcomes = c.Request().Form["learning-outcomes[]"]
		}
		return c.Render(http.StatusOK, "learningOutcomes", draft)
	case "courseObjectives":
		draft.Prerequisites = c.FormValue("prerequisites")
		if err := c.Request().ParseForm(); err == nil {
			draft.CourseObjectives = c.Request().Form["course-objectives[]"]
		}
		return c.Render(http.StatusOK, "courseObjectives", draft)
	case "courseRequirements":
		if err := c.Request().ParseForm(); err == nil {
			draft.CourseRequirements = c.Request().Form["course-requirements[]"]
		}
		return c.Render(http.StatusOK, "courseRequirements", draft)
	case "activeLearning1":
		draft.ActiveLearning1 = c.FormValue("active-learning-1")
	case "activeLearning2":
		draft.ActiveLearning2 = c.FormValue("active-learning-2")
	case "activeLearning3":
		draft.ActiveLearning3 = c.FormValue("active-learning-3")
	case "activeLearning4":
		draft.ActiveLearning4 = c.FormValue("active-learning-4")
	case "courseStructure":
		if err := c.Request().ParseForm(); err == nil {
			draft.CourseStructure = c.Request().Form["lecture-type"]
		}
		return c.Render(http.StatusOK, "course-structure-container", draft)
	case "otherCourseStructure":
		draft.OtherCourseStructure = c.FormValue("other-course-structure")
		return c.Render(http.StatusOK, "otherCourseInput", draft)

	case "gradeComponents":
		// New case for grade composition update.
		if err := c.Request().ParseForm(); err == nil {
			names := c.Request().Form["grade-component-name[]"]
			percents := c.Request().Form["grade-component-percentage[]"]
			var comps []UIcomponents.GradeComponent
			count := len(names)
			for i := 0; i < count; i++ {
				comps = append(comps, UIcomponents.GradeComponent{
					PartName:   names[i],
					Percentage: percents[i],
				})
			}
			draft.GradeComponents = comps
		}
		return c.Render(http.StatusOK, "gradeComponents", draft)

	}

	return c.NoContent(http.StatusNoContent)
}
