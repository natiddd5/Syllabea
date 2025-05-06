package mid

import (
	"errors"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// GetUserID retrieves the user ID from the session.
func GetUserID(c echo.Context) (int, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		c.Logger().Error("Failed to get session:", err)
		return 0, err
	}

	userIDVal, ok := sess.Values["user_id"]
	if !ok {
		c.Logger().Warn("No user_id found in session")
		return 0, errors.New("user not logged in")
	}

	var userID int
	switch v := userIDVal.(type) {
	case int:
		userID = v
	case int64:
		userID = int(v)
	case string:
		userID, err = strconv.Atoi(v)
		if err != nil {
			c.Logger().Error("Error converting user_id to int:", err)
			return 0, errors.New("invalid user id")
		}
	default:
		c.Logger().Error("Invalid type for user_id:", v)
		return 0, errors.New("invalid user id type")
	}

	return userID, nil
}

// AuthMiddleware checks for a valid session and redirects to login if not found.
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := GetUserID(c)
		if err != nil {
			c.Logger().Warn("Unauthorized access attempt, redirecting to /login")
			c.Response().Header().Set("HX-Redirect", "/login")
			return c.String(http.StatusOK, "Redirecting to login page...")
		}
		return next(c)
	}
}
