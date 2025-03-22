package handler

import (
	"Syllybea/UIcomponents"
	"Syllybea/mid"
	"Syllybea/repository"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var r repository.Repository

// RegisterRoutes registers all endpoints.
func RegisterRoutes(e *echo.Echo, repo *repository.Repository) {
	r = *repo
	// Login page.
	e.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "login.html", nil)
	})

	e.GET("/", func(c echo.Context) error {
		_, err := mid.GetUserID(c)
		if err != nil {
			// Redirect to login
			c.Logger().Info("Redirecting to /login")
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Redirect to dashboard
		c.Logger().Info("Redirecting to /dashboard")
		return c.Redirect(http.StatusSeeOther, "/dashboard")
	})

	// Login submission.
	e.POST("/login", func(c echo.Context) error {
		email := c.FormValue("email")
		if email == "" {
			c.Logger().Warn("Empty email provided")
			return c.String(http.StatusBadRequest, "אנא ספק אימייל")
		}

		// Check if the user exists.
		user, err := repo.GetUserByEmail(email)
		if err != nil {
			c.Logger().Warn("User not found for email:", email)
			return c.String(http.StatusUnauthorized, "אימייל לא קיים")
		}

		// Create a session and store the user's ID and email.
		sess, err := session.Get("session", c)
		if err != nil {
			c.Logger().Error("Failed to get session during login:", err)
			return err
		}
		sess.Values["user_id"] = user.ID
		sess.Values["email"] = user.Email

		if err := sess.Save(c.Request(), c.Response()); err != nil {
			c.Logger().Error("Failed to save session:", err)
			return err
		}
		c.Logger().Info("Session saved successfully for user:", user.Email)

		// Set the HX-Redirect header for HTMX and perform redirect.
		c.Response().Header().Set("HX-Redirect", "/dashboard")
		c.Logger().Info("Redirecting to /dashboard")
		return c.String(http.StatusOK, "Redirecting...")
	})

	// Dashboard.
	e.GET("/dashboard", func(c echo.Context) error {
		startOverall := time.Now()
		c.Logger().Info("Dashboard route hit")

		// Use our common getUserID function.
		userID, err := mid.GetUserID(c)
		if err != nil {
			c.Logger().Error("getUserID error:", err)
			// Set HX-Redirect header for HTMX so that it knows to navigate back to login.
			c.Response().Header().Set("HX-Redirect", "/login")
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Get user data.
		userStart := time.Now()
		user, err := repo.GetUserByID(userID)
		if err != nil {
			c.Logger().Error("GetUserByID error:", err)
			return c.String(http.StatusInternalServerError, "Error retrieving user data")
		}
		c.Logger().Infof("Retrieved user data in %s", time.Since(userStart))

		// Fetch cards from repository.
		c.Logger().Info("Fetching cards for lecturer...")
		cardsStart := time.Now()
		rawCards, err := repo.GetCardsByLecturer(userID)
		if err != nil {
			c.Logger().Error("GetCardsByLecturer error:", err)
			return c.String(http.StatusInternalServerError, "Error fetching cards")
		}
		c.Logger().Infof("Fetched %d cards in %s", len(rawCards), time.Since(cardsStart))

		// Hebrew month names mapping.
		hebrewMonths := map[time.Month]string{
			time.January:   "ינואר",
			time.February:  "פברואר",
			time.March:     "מרץ",
			time.April:     "אפריל",
			time.May:       "מאי",
			time.June:      "יוני",
			time.July:      "יולי",
			time.August:    "אוגוסט",
			time.September: "ספטמבר",
			time.October:   "אוקטובר",
			time.November:  "נובמבר",
			time.December:  "דצמבר",
		}

		// Group cards by month/year.
		groupStart := time.Now()
		cardsByMonth := make(map[string][]UIcomponents.Card)
		monthOrder := make(map[string]time.Time)
		total, attempts, inReview, approved := 0, 0, 0, 0

		for i, cardMap := range rawCards {
			// Log every 50 cards processed for insight into progress.
			if i%50 == 0 {
				c.Logger().Infof("Processing card %d/%d", i+1, len(rawCards))
			}

			// Extract fields.
			dateStr, _ := cardMap["date"].(string)
			title, _ := cardMap["title"].(string)
			lecturer, _ := cardMap["lecturer"].(string)
			field, _ := cardMap["field"].(string)
			status, _ := cardMap["status"].(string)

			// Parse date.
			parsedDate, err := time.Parse("02/01/2006", dateStr)
			if err != nil {
				c.Logger().Warn("Skipping card due to date parse error:", err)
				continue
			}
			monthLabel := hebrewMonths[parsedDate.Month()]
			monthYearKey := monthLabel + " " + strconv.Itoa(parsedDate.Year())

			card := UIcomponents.Card{
				Title:       title,
				Date:        dateStr,
				Lecturer:    lecturer,
				Field:       field,
				Status:      status,
				StatusLabel: status,
			}
			cardsByMonth[monthYearKey] = append(cardsByMonth[monthYearKey], card)

			// Record first day of month for sorting.
			if _, exists := monthOrder[monthYearKey]; !exists {
				monthOrder[monthYearKey] = time.Date(parsedDate.Year(), parsedDate.Month(), 1, 0, 0, 0, 0, time.UTC)
			}

			total++
			switch status {
			case "Draft":
				attempts++
			case "In Review":
				inReview++
			case "Approved":
				approved++
			}
		}
		c.Logger().Infof("Grouped cards by month in %s", time.Since(groupStart))
		c.Logger().Infof("Total cards processed: %d; Draft: %d, In Review: %d, Approved: %d", total, attempts, inReview, approved)

		// Build sorted date sections.
		sortStart := time.Now()
		var dateSections []UIcomponents.DateSection
		for key, cards := range cardsByMonth {
			dateSections = append(dateSections, UIcomponents.DateSection{
				DateLabel: key,
				Cards:     cards,
			})
		}
		sort.Slice(dateSections, func(i, j int) bool {
			return monthOrder[dateSections[i].DateLabel].After(monthOrder[dateSections[j].DateLabel])
		})
		c.Logger().Infof("Sorted date sections in %s", time.Since(sortStart))

		header := UIcomponents.HeaderData{
			Title: "Dashboard",
			Name:  user.Name,
		}
		content := UIcomponents.CoursesData{
			Total:        total,
			Attempts:     attempts,
			InReview:     inReview,
			Approved:     approved,
			DateSections: dateSections,
		}
		pageData := UIcomponents.PageData{
			Header:  header,
			Content: content,
		}

		c.Logger().Infof("Dashboard rendered in %s", time.Since(startOverall))
		return c.Render(http.StatusOK, "base", pageData)
	})

	// Filter endpoint.
	e.POST("/filter", func(c echo.Context) error {
		userID, err := mid.GetUserID(c)
		if err != nil {
			c.Logger().Warn("Authentication error in filter endpoint:", err)
			c.Response().Header().Set("HX-Redirect", "/login")
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Retrieve filter parameters.
		search := c.FormValue("search")
		fromDate := c.FormValue("from-date")
		toDate := c.FormValue("to-date")
		statuses := c.Request().Form["status"]

		// Map Hebrew statuses to English.
		englishStatuses := []string{}
		for _, s := range statuses {
			switch s {
			case "סילבוס":
				englishStatuses = append(englishStatuses, "Draft")
			case "בבחינה":
				englishStatuses = append(englishStatuses, "In Review")
			case "מאושר":
				englishStatuses = append(englishStatuses, "Approved")
			}
		}

		// Get filtered cards.
		rawCards, err := repo.FilterCardsByLecturer(userID, search, fromDate, toDate, englishStatuses)
		if err != nil {
			c.Logger().Error("FilterCardsByLecturer error:", err)
			return c.String(http.StatusInternalServerError, "Error fetching filtered cards")
		}

		// Map for Hebrew month names.
		hebrewMonths := map[time.Month]string{
			time.January:   "ינואר",
			time.February:  "פברואר",
			time.March:     "מרץ",
			time.April:     "אפריל",
			time.May:       "מאי",
			time.June:      "יוני",
			time.July:      "יולי",
			time.August:    "אוגוסט",
			time.September: "ספטמבר",
			time.October:   "אוקטובר",
			time.November:  "נובמבר",
			time.December:  "דצמבר",
		}

		// Group cards by month/year.
		cardsByMonth := make(map[string][]UIcomponents.Card)
		monthOrder := make(map[string]time.Time)
		total, attempts, inReview, approved := 0, 0, 0, 0

		for _, cardMap := range rawCards {
			dateStr, _ := cardMap["date"].(string)
			title, _ := cardMap["title"].(string)
			lecturer, _ := cardMap["lecturer"].(string)
			field, _ := cardMap["field"].(string)
			status, _ := cardMap["status"].(string)

			parsedDate, err := time.Parse("02/01/2006", dateStr)
			if err != nil {
				continue
			}

			monthLabel := hebrewMonths[parsedDate.Month()]
			monthYearKey := monthLabel + " " + strconv.Itoa(parsedDate.Year())

			card := UIcomponents.Card{
				Title:       title,
				Date:        dateStr,
				Lecturer:    lecturer,
				Field:       field,
				Status:      status,
				StatusLabel: status,
			}
			cardsByMonth[monthYearKey] = append(cardsByMonth[monthYearKey], card)

			if _, exists := monthOrder[monthYearKey]; !exists {
				monthOrder[monthYearKey] = time.Date(parsedDate.Year(), parsedDate.Month(), 1, 0, 0, 0, 0, time.UTC)
			}

			total++
			switch status {
			case "Draft":
				attempts++
			case "In Review":
				inReview++
			case "Approved":
				approved++
			}
		}

		// Build sorted date sections.
		var dateSections []UIcomponents.DateSection
		for key, cards := range cardsByMonth {
			dateSections = append(dateSections, UIcomponents.DateSection{
				DateLabel: key,
				Cards:     cards,
			})
		}
		sort.Slice(dateSections, func(i, j int) bool {
			return monthOrder[dateSections[i].DateLabel].After(monthOrder[dateSections[j].DateLabel])
		})

		return c.Render(http.StatusOK, "outer-container.html", dateSections)
	})

	e.GET("create-syllabus", func(c echo.Context) error {
		userID, _ := mid.GetUserID(c)
		user, _ := repo.GetUserByID(userID)
		draft := GetUserDraft(user)
		return c.Render(http.StatusOK, "create-syllabus", draft)
	})

	e.POST("update-syllabus", updateSyllabusHandler)

}

func updateSyllabusHandler(c echo.Context) error {
	userID, _ := mid.GetUserID(c)
	user, _ := r.GetUserByID(userID)
	draft := GetUserDraft(user)
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
		draft.Prerequisites = c.FormValue("prerequisites")
	case "learningOutcomes":
		if err := c.Request().ParseForm(); err == nil {
			draft.LearningOutcomes = c.Request().Form["learning-outcomes[]"]
		}
		return c.Render(http.StatusOK, "learningOutcomes", draft)
	case "courseObjectives":
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
