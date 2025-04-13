package handler

import (
	"Syllybea/repository"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

var r repository.Repository

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
		return handleDashboard(c, repo)
	})

	// Filter endpoint.
	e.POST("/filter", func(c echo.Context) error {
		return filterCards(c, repo)
	})

	e.GET("/create-syllabus", func(c echo.Context) error {
		return HandleCreateSyllabus(c, repo)
	})

	e.POST("/submit-syllabus", func(c echo.Context) error {
		return handleSubmitSyllabus(c, repo)
	})

	e.POST("/update-syllabus", updateSyllabusHandler)

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
}
