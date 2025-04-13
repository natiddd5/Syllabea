package handler

import (
	"Syllybea/UIcomponents"
	"Syllybea/cache"
	"Syllybea/mid"
	"Syllybea/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func HandleEditSyllabus(c echo.Context, repo *repository.Repository) error {
	// Parse the syllabus ID from the URL parameter.
	syllabusID := c.Param("id")
	id, err := strconv.Atoi(syllabusID)
	if err != nil {
		c.Logger().Error("Invalid syllabus ID: ", err)
		return c.String(http.StatusBadRequest, "Invalid syllabus ID")
	}

	// Retrieve the syllabus using the repository.
	syl, err := repo.GetSyllabusByID(id)
	if err != nil {
		c.Logger().Error("Error retrieving syllabus: ", err)
		return c.String(http.StatusInternalServerError, "Error retrieving syllabus")
	}

	if syl == nil {
		c.Logger().Warn("Syllabus not found with id: ", id)
		return c.String(http.StatusNotFound, "Syllabus not found")
	}

	// Get the edited draft from the cache.
	draft, err := cache.GetEditedSyllabus(syl)
	if err != nil {
		c.Logger().Error("Error processing syllabus data: ", err)
		return c.String(http.StatusInternalServerError, "Error processing syllabus data")
	}
	if draft == nil {
		c.Logger().Error("Draft not found for syllabus id: ", id)
		return c.String(http.StatusInternalServerError, "Draft not found")
	}

	return c.Render(http.StatusOK, "create-syllabus", draft)
}

func HandleCreateSyllabus(c echo.Context, repo *repository.Repository) error {
	userID, _ := mid.GetUserID(c)
	user, _ := repo.GetUserByID(userID)
	draft := cache.GetUserDraft(user)
	departments, _ := r.GetAllDepartments()
	deptNames := make([]string, 0, len(departments))
	for _, dept := range departments {
		deptNames = append(deptNames, dept.Name)
	}
	draft.Departments = deptNames
	draft.SyllabusDepartment = deptNames[0]
	courses, _ := r.GetAllCourses()
	courseNames := make([]string, 0, len(courses))
	for _, course := range courses {
		courseNames = append(courseNames, course.Name)
	}
	draft.Courses = courseNames
	draft.SelectedCourse = courseNames[0]

	return c.Render(http.StatusOK, "create-syllabus", draft)
}

func handleSubmitSyllabus(c echo.Context, repo *repository.Repository) error {

	userID, _ := mid.GetUserID(c)
	user, _ := r.GetUserByID(userID)
	draft := cache.GetUserDraft(user)
	err := r.InsertSyllabusFromDraft(userID, draft)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error saving syllabus")
	}
	c.Response().Header().Set("HX-Redirect", "/dashboard")
	cache.DraftCache.Delete(user.ID)
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
	user, _ := r.GetUserByID(userID)
	draft := cache.GetUserDraft(user)
	switch c.FormValue("action") {
	case "addCourseRequirement":
		return addCourseRequirement(c, draft)
	case "removeCourseRequirement":
		return removeCourseRequirement(c, draft)
	case "updateSyllabusRow":
		return updateSyllabusRow(c, draft)
	case "removeSyllabusRow":
		return removeSyllabusRow(c, draft)
	case "insertSyllabusRow":
		insertSyllabusRow(c, draft)
	case "addLearningOutcome":
		return addLearningOutcome(c, draft)
	case "removeLearningOutcome":
		return removeLearningOutcome(c, draft)
	case "addCourseObjective":
		return addCourseObjective(c, draft)
	case "removeCourseObjective":
		return removeCourseObjective(c, draft)
	case "addGradeComponent":
		return addGradeComponent(c, draft)
	case "removeGradeComponent":
		return removeGradeComponent(c, draft)
	case "addAssignmentStructure":
		return addAssignmentStructure(c, draft)
	case "removeAssignmentStructure":
		return removeAssignmentStructure(c, draft)
	case "addBibliographyRequired":
		return addBibliographyRequired(c, draft)
	case "removeBibliographyRequired":
		return removeBibliographyRequired(c, draft)
	case "addBibliographyRecommended":
		return addBibliographyRecommended(c, draft)
	case "removeBibliographyRecommended":
		return removeBibliographyRecommended(c, draft)
	default:
		return handleGeneralUpdate(c, draft)
	}
	return c.NoContent(http.StatusNoContent)

}

func handleGeneralUpdate(c echo.Context, draft *UIcomponents.Draft) error {
	updateField := c.FormValue("updateField")
	switch updateField {
	case "syllabus-department":
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
