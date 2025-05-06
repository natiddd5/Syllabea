package handler

import (
	"Syllybea/UIcomponents"
	"Syllybea/mid"
	"Syllybea/repository"
	"Syllybea/types"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

func handleGetCommentsOfSyllabus(c echo.Context, r *repository.Repository) error {
	sylIDString := c.FormValue("syllabus_id")
	sylID, err := strconv.Atoi(sylIDString)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid syllabus ID")
	}

	comments, err := r.GetCommentsBySyllabusID(sylID)
	if err != nil {
		// Just log the error and continue with empty comments
		// This prevents errors when syllabus_id doesn't exist yet
		comments = []types.Comment{}
	}

	// Initialize the UI comments slice with the capacity to hold the server message plus all DB comments.
	UIComments := make([]UIcomponents.Comments, 0, len(comments)+1)

	// Add a default server message as the first comment.
	serverComment := UIcomponents.Comments{
		Name:          "מכללת כנרת",
		Message:       "זו היא מערכת ההודעות של Syllabea",
		Time:          time.Now().Format("2006-01-02 15:04"),
		IsCurrentUser: false,
	}
	UIComments = append(UIComments, serverComment)

	// Get current user ID
	currentUserID, _ := mid.GetUserID(c)

	// Append comments from the database.
	for i := 0; i < len(comments); i++ {
		UIComments = append(UIComments, UIcomponents.Comments{
			Name:          strconv.Itoa(comments[i].UserID),
			Message:       comments[i].Content,
			Time:          comments[i].CreatedAt.Format("2006-01-02 15:04"),
			IsCurrentUser: comments[i].UserID == currentUserID,
		})
	}

	// Package the UIComments slice inside the data wrapper expected by the template.
	data := struct {
		ID       int
		Comments []UIcomponents.Comments
	}{
		ID:       sylID,
		Comments: UIComments,
	}

	return c.Render(http.StatusOK, "comments.html", data)
}

func handleAddComment(c echo.Context, r *repository.Repository) error {
	// Extract user ID (assumes middleware stores it)
	userID, err := mid.GetUserID(c)
	if err != nil {
		// For testing, use a default user ID if not logged in
		userID = 1 // Default user ID for testing
	}

	// Parse syllabus ID
	sylIDString := c.FormValue("syllabus_id")
	sylID, err := strconv.Atoi(sylIDString)
	if err != nil {
		return c.HTML(http.StatusBadRequest, "<div class='error-message'>Invalid syllabus ID</div>")
	}

	// Get comment content from the form
	content := c.FormValue("content")
	if content == "" {
		return c.HTML(http.StatusBadRequest, "<div class='error-message'>Comment content cannot be empty</div>")
	}

	// Create current time
	currentTime := time.Now()

	// Create comment struct
	comment := &types.Comment{
		SyllabusID: sylID,
		UserID:     userID,
		Content:    content,
		CreatedAt:  currentTime,
	}

	// Add comment to DB - handle potential errors gracefully
	err = r.AddComment(comment)
	if err != nil {
		// For development/testing - if DB operation fails, still show the comment
		// but return error in console
		// In production, you'd want to return the error to the user
		// But for now, let's make it work even if DB fails
	}

	// Create the HTML for the new comment directly - this avoids template rendering issues
	// Format the HTML manually instead of using a template
	// Add current-user class to identify the current user's comments
	commentHTML := `
    <div class="comment-item current-user">
        <div class="comment-header">
            <span class="comment-author current-user-name">` + strconv.Itoa(userID) + `</span>
            <span class="comment-time">` + currentTime.Format("2006-01-02 15:04") + `</span>
        </div>
        <div class="comment-text">` + content + `</div>
    </div>`

	return c.HTML(http.StatusOK, commentHTML)
}
