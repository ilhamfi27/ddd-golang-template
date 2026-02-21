package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	UserID int      `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// JWTConfig returns the JWT middleware configuration
func JWTConfig(secret string) echojwt.Config {
	return echojwt.Config{
		SigningKey:     []byte(secret),
		SigningMethod:  "HS256",
		ContextKey:     "user",
		ErrorHandler:   jwtErrorHandler,
		SuccessHandler: jwtSuccessHandler,
		Skipper:        defaultSkipper,
	}
}

// jwtErrorHandler handles JWT errors
func jwtErrorHandler(c echo.Context, err error) error {
	return c.JSON(401, map[string]string{
		"error": "invalid or missing token",
	})
}

// jwtSuccessHandler handles successful JWT validation
func jwtSuccessHandler(c echo.Context) {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*CustomClaims)

	// Store user information in context for later use
	c.Set("user_id", claims.UserID)
	c.Set("email", claims.Email)
	c.Set("roles", claims.Roles)
}

// defaultSkipper skips authentication for public endpoints
func defaultSkipper(c echo.Context) bool {
	// List of public endpoints that don't need authentication
	publicPaths := map[string]bool{
		"/":              true,
		"/healthcheck":   true,
		"/swagger/index": true,
		"/auth/login":    true,
		"/auth/register": true,
	}
	return publicPaths[c.Path()]
}

// GetUserID returns the user ID from context
func GetUserID(c echo.Context) int {
	if userID, ok := c.Get("user_id").(int); ok {
		return userID
	}
	return 0
}

// GetEmail returns the email from context
func GetEmail(c echo.Context) string {
	if email, ok := c.Get("email").(string); ok {
		return email
	}
	return ""
}

// GetRoles returns the roles from context
func GetRoles(c echo.Context) []string {
	if roles, ok := c.Get("roles").([]string); ok {
		return roles
	}
	return []string{}
}
