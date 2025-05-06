package handler

import (
	"Syllybea/UIcomponents"
	"Syllybea/mid"
	"Syllybea/repository"
	"Syllybea/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func handleHome(c echo.Context) error {
	_, err := mid.GetUserID(c)
	if err != nil {
		// Redirect to login
		c.Logger().Info("Redirecting to /login")
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	// Redirect to dashboard
	c.Logger().Info("Redirecting to /dashboard")
	return c.Redirect(http.StatusSeeOther, "/dashboard")
}

func handleDashboard(c echo.Context, repo *repository.Repository) error {

	// Use our common getUserID function.
	userID, err := mid.GetUserID(c)
	if err != nil {
		c.Logger().Error("getUserID error:", err)
		// Set HX-Redirect header for HTMX so that it knows to navigate back to login.
		c.Response().Header().Set("HX-Redirect", "/login")
		return c.String(http.StatusOK, "Redirecting to login page...")
	}

	user, err := repo.GetUserByID(userID)
	if err != nil {
		c.Logger().Error("GetUserByID error:", err)
		return c.String(http.StatusInternalServerError, "Error retrieving user data")
	}

	// Fetch cards from repository.
	c.Logger().Info("Fetching cards for lecturer...")
	// Now rawCards is of type []UIcomponents.Card
	rawCards, err := repo.GetCardsByLecturer(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching cards")
	}

	// Group cards by month/year.
	groupStart := time.Now()
	cardsByMonth := make(map[string][]UIcomponents.Card)
	monthOrder := make(map[string]time.Time)
	total, attempts, inReview, approved := 0, 0, 0, 0

	for _, card := range rawCards {

		// Parse the date (card.Date is expected in "02/01/2006" format).
		parsedDate, err := time.Parse("02/01/2006", card.Date)
		if err != nil {
			c.Logger().Warn("Skipping card due to date parse error:", err)
			continue
		}
		monthLabel := utils.HebrewMonths[parsedDate.Month()]
		monthYearKey := monthLabel + " " + strconv.Itoa(parsedDate.Year())

		// Append the card (which already contains the ID, Title, Lecturer, Field, etc.)
		cardsByMonth[monthYearKey] = append(cardsByMonth[monthYearKey], card)

		// Record the first day of the month for sorting.
		if _, exists := monthOrder[monthYearKey]; !exists {
			monthOrder[monthYearKey] = time.Date(parsedDate.Year(), parsedDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		}

		total++
		switch card.Status {
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

	return c.Render(http.StatusOK, "base", pageData)
}

func filterCards(c echo.Context, repo *repository.Repository) error {
	userID, err := mid.GetUserID(c)
	if err != nil {
		c.Logger().Warn("Authentication error in filter endpoint:", err)
		c.Response().Header().Set("HX-Redirect", "/login")
		return c.String(http.StatusOK, "Redirecting to login page...")
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

		monthLabel := utils.HebrewMonths[parsedDate.Month()]
		monthYearKey := monthLabel + " " + strconv.Itoa(parsedDate.Year())

		id, _ := cardMap["id"].(int) // You need this line!
		card := UIcomponents.Card{
			ID:          id,
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
}

// handleDeleteSyllabus updates a syllabus status to "Deleted"
func handleDeleteSyllabus(c echo.Context, repo *repository.Repository) error {
	// Get the syllabus ID from the URL parameter
	syllabusID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error("Invalid syllabus ID:", err)
		return c.String(http.StatusBadRequest, "Invalid syllabus ID")
	}

	// Get the syllabus from the database
	syl, err := repo.GetSyllabusByID(syllabusID)
	if err != nil {
		c.Logger().Error("Error retrieving syllabus:", err)
		return c.String(http.StatusInternalServerError, "Error retrieving syllabus")
	}

	// Update the syllabus status to "Deleted"
	syl.Status = "Deleted"
	err = repo.UpdateSyllabus(syl)
	if err != nil {
		c.Logger().Error("Error updating syllabus status:", err)
		return c.String(http.StatusInternalServerError, "Error updating syllabus status")
	}

	// Return an empty response to indicate success (the card will be removed from the UI)
	return c.NoContent(http.StatusOK)
}

// handleTrashPage displays the trash page with deleted syllabi
func handleTrashPage(c echo.Context, repo *repository.Repository) error {
	// Use our common getUserID function.
	userID, err := mid.GetUserID(c)
	if err != nil {
		c.Logger().Error("getUserID error:", err)
		// Set HX-Redirect header for HTMX so that it knows to navigate back to login.
		c.Response().Header().Set("HX-Redirect", "/login")
		return c.String(http.StatusOK, "Redirecting to login page...")
	}

	user, err := repo.GetUserByID(userID)
	if err != nil {
		c.Logger().Error("GetUserByID error:", err)
		return c.String(http.StatusInternalServerError, "Error retrieving user data")
	}

	// Fetch deleted cards from repository.
	c.Logger().Info("Fetching deleted cards for lecturer...")
	rawCards, err := repo.GetDeletedCardsByLecturer(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching deleted cards")
	}

	// Group cards by month/year.
	groupStart := time.Now()
	cardsByMonth := make(map[string][]UIcomponents.Card)
	monthOrder := make(map[string]time.Time)
	total := len(rawCards)

	for _, card := range rawCards {
		// Parse the date (card.Date is expected in "02/01/2006" format).
		parsedDate, err := time.Parse("02/01/2006", card.Date)
		if err != nil {
			c.Logger().Warn("Skipping card due to date parse error:", err)
			continue
		}
		monthLabel := utils.HebrewMonths[parsedDate.Month()]
		monthYearKey := monthLabel + " " + strconv.Itoa(parsedDate.Year())

		// Append the card (which already contains the ID, Title, Lecturer, Field, etc.)
		cardsByMonth[monthYearKey] = append(cardsByMonth[monthYearKey], card)

		// Record the first day of the month for sorting.
		if _, exists := monthOrder[monthYearKey]; !exists {
			monthOrder[monthYearKey] = time.Date(parsedDate.Year(), parsedDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		}
	}
	c.Logger().Infof("Grouped deleted cards by month in %s", time.Since(groupStart))
	c.Logger().Infof("Total deleted cards processed: %d", total)

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
		Title: "Trash",
		Name:  user.Name,
	}
	content := UIcomponents.CoursesData{
		Total:        total,
		Attempts:     0,
		InReview:     0,
		Approved:     0,
		DateSections: dateSections,
	}
	pageData := UIcomponents.PageData{
		Header:  header,
		Content: content,
	}

	return c.Render(http.StatusOK, "trash-page", pageData)
}

// handlePermanentDeleteSyllabus permanently deletes a syllabus from the database
func handlePermanentDeleteSyllabus(c echo.Context, repo *repository.Repository) error {
	// Get the syllabus ID from the URL parameter
	syllabusID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error("Invalid syllabus ID:", err)
		return c.String(http.StatusBadRequest, "Invalid syllabus ID")
	}

	// Delete the syllabus from the database
	err = repo.DeleteSyllabus(syllabusID)
	if err != nil {
		c.Logger().Error("Error deleting syllabus:", err)
		return c.String(http.StatusInternalServerError, "Error deleting syllabus")
	}

	// Return an empty response to indicate success (the card will be removed from the UI)
	return c.NoContent(http.StatusOK)
}
