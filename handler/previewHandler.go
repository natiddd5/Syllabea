package handler

import (
	"Syllybea/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// HandleSyllabusPreview handles the request to preview a syllabus
func HandleSyllabusPreview(c echo.Context, repo *repository.Repository) error {
	// Get the syllabus ID from the request
	syllabusID := c.Param("id")
	if syllabusID == "" {
		return c.String(http.StatusBadRequest, "Syllabus ID is required")
	}

	// Convert the ID to an integer
	id, err := strconv.Atoi(syllabusID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid syllabus ID")
	}

	// Get the draft from the repository
	draft, err := repo.GetEditedSyllabus(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get syllabus")
	}

	// Render the preview template with the draft data
	return c.Render(http.StatusOK, "syllabus-preview.html", draft)
}

// HandleSyllabusPreviewFromForm handles the request to preview a syllabus from the form
func HandleSyllabusPreviewFromForm(c echo.Context, repo *repository.Repository) error {
	// Get the syllabus ID from the form
	syllabusIDStr := c.FormValue("syllabus_id")
	if syllabusIDStr == "" {
		return c.String(http.StatusBadRequest, "Syllabus ID is required")
	}

	// Convert the ID to an integer
	syllabusID, err := strconv.Atoi(syllabusIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid syllabus ID")
	}

	// Get the draft from the repository
	draft, err := repo.GetEditedSyllabus(syllabusID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get syllabus")
	}

	// Render the preview template with the draft data
	return c.Render(http.StatusOK, "syllabus-preview.html", draft)
}
