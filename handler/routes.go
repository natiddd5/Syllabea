package handler

import (
	"Syllybea/mid"
	"Syllybea/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

var r repository.Repository

// handleLogout removes the JWT cookie and redirects to the login page
func handleLogout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(c.Response(), cookie)

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.String(http.StatusOK, "Redirecting to login page...")
}

// RegisterRoutes registers all endpoints.
func RegisterRoutes(e *echo.Echo, repo *repository.Repository) {
	r = *repo
	// Login page.
	e.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "login.html", nil)
	})

	e.GET("/", handleHome)

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

		token, err := mid.GenerateToken(user.ID)
		if err != nil {
			c.Logger().Error("Failed to generate token:", err)
			return c.String(http.StatusInternalServerError, "error")
		}

		cookie := &http.Cookie{
			Name:     "jwt",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(c.Response(), cookie)
		c.Logger().Info("JWT cookie set for user:", user.Email)

		// Set the HX-Redirect header for HTMX and perform redirect.
		c.Response().Header().Set("HX-Redirect", "/dashboard")
		c.Logger().Info("Redirecting to /dashboard")
		return c.String(http.StatusOK, "Redirecting...")
	})

	// Dashboard.
	e.GET("/dashboard", func(c echo.Context) error {
		return handleDashboard(c, repo)
	})

	// Filter endpoint.
	e.POST("/filter", func(c echo.Context) error {
		return filterCards(c, repo)
	})

	e.GET("/syllabus/create", func(c echo.Context) error {
		return HandleCreateSyllabus(c, repo)
	})

	e.POST("/syllabus/submit", func(c echo.Context) error {
		return handleSubmitSyllabus(c, repo)
	})

	e.POST("/syllabus/save", func(c echo.Context) error {
		return handleSaveSyllabus(c, repo)
	})

	e.POST("/syllabus/update", updateSyllabusHandler)
	e.POST("/update-syllabus", updateSyllabusHandler)

	//does not match HTMX request
	e.GET("/edit-syllabus/:id", func(c echo.Context) error {
		return HandleEditSyllabus(c, repo)
	})

	// New POST route for fetching comments.
	e.GET("/syllabus/comments", func(c echo.Context) error {
		return handleGetCommentsOfSyllabus(c, repo)
	})

	e.POST("/add-comment", func(c echo.Context) error {
		return handleAddComment(c, repo)
	})

	// Preview syllabus routes
	e.GET("/syllabus/preview/:id", func(c echo.Context) error {
		return HandleSyllabusPreview(c, repo)
	})

	e.POST("/syllabus/preview", func(c echo.Context) error {
		return HandleSyllabusPreviewFromForm(c, repo)
	})

	// Logout endpoint
	e.POST("/logout", func(c echo.Context) error {
		return handleLogout(c)
	})

	// Delete syllabus endpoint
	e.DELETE("/delete-syllabus/:id", func(c echo.Context) error {
		return handleDeleteSyllabus(c, repo)
	})

	// Trash page endpoint
	e.GET("/trash", func(c echo.Context) error {
		return handleTrashPage(c, repo)
	})

	// Permanent delete syllabus endpoint
	e.DELETE("/permanent-delete-syllabus/:id", func(c echo.Context) error {
		return handlePermanentDeleteSyllabus(c, repo)
	})
}
