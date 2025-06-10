package mid

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

var jwtSecret = []byte("secret")

type customClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT containing the user's ID.
func GenerateToken(userID int) (string, error) {
	claims := customClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func parseToken(tokenStr string) (*customClaims, error) {
	claims := &customClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// GetUserID retrieves the user ID from a JWT placed in the Authorization header or cookie.
func GetUserID(c echo.Context) (int, error) {
	tokenStr := ""

	authHeader := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	} else if cookie, err := c.Cookie("jwt"); err == nil {
		tokenStr = cookie.Value
	}

	if tokenStr == "" {
		c.Logger().Warn("missing JWT token")
		return 0, errors.New("token missing")
	}

	claims, err := parseToken(tokenStr)
	if err != nil {
		c.Logger().Warn("invalid token: ", err)
		return 0, err
	}

	return claims.UserID, nil
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
