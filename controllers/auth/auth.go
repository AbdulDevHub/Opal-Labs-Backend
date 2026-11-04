package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/opalescencelabs/backend/database"
	"github.com/opalescencelabs/backend/models"
)

// TokenLifetime is how long a session token issued by /user-login stays valid.
const TokenLifetime = 30 * 24 * time.Hour

// testToken lets the integration tests in /tests skip Google OAuth.
// It is only accepted when APP_ENV=development.
const testToken = "Only_for_testing1200332"

// GetDomain returns the domain of the app.
// Return "localhost" if the app is in development mode otherwise return the DOMAIN env variable.
func GetDomain() string {
	if os.Getenv("APP_ENV") == "development" {
		return "http://localhost:8000"
	}
	return os.Getenv("DOMAIN")
}

// GetFrontendURLs returns the frontend origins allowed by CORS.
// FRONTEND_URL may hold several comma-separated URLs; spaces and trailing
// slashes are trimmed because browsers send the origin without them.
func GetFrontendURLs() []string {
	if os.Getenv("APP_ENV") == "development" {
		return []string{"http://localhost:3000"}
	}
	var origins []string
	for _, url := range strings.Split(os.Getenv("FRONTEND_URL"), ",") {
		url = strings.TrimRight(strings.TrimSpace(url), "/")
		if url != "" {
			origins = append(origins, url)
		}
	}
	return origins
}

// GetSecretKey returns the secret key from the environment variables.
func GetSecretKey() string {
	return os.Getenv("SECRET")
}

// ValidateAccessToken sends a request to Google's tokeninfo endpoint to validate the access token.
// Returns an error if the access token is invalid or there was an error with the api. Otherwise, returns nil.
func ValidateAccessToken(accessToken string) error {
	resp, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?access_token=%s", accessToken))
	if err != nil {
		return errors.New("Unable to verify access token: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error verifying Google access token, status code: %d", resp.StatusCode)
	}

	return nil
}

// RevokeAccessToken sends a request to Google's revocation endpoint to revoke the access token.
// Returns an error if the access token is invalid or there was an error with the api. Otherwise, returns nil.
func RevokeAccessToken(accessToken string) error {
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/revoke?token="+accessToken, nil)
	if err != nil {
		return errors.New("Unable to create request: " + err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("Unable to send request: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error revoking Google access token, status code: %d", resp.StatusCode)
	}

	return nil
}

// GetUserInfo sends a request to Google's userinfo endpoint to get the user's information.
// Returns a map containing the user's information or an error if there was an error with the api.
func GetUserInfo(accessToken string) (map[string]interface{}, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, errors.New("Failed to create userinfo request: " + err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to fetch user information: " + err.Error())
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, errors.New("Failed to parse user information: " + err.Error())
	}

	return userInfo, nil
}

// MakeCredentialsJson returns a json byte array containing the access token, refresh token, and expiry time.
// Returns an error if there was an error marshalling the credentials. Otherwise, returns the credentials.
func MakeCredentialsJson(accessToken string, refreshToken string, expiry int) ([]byte, error) {
	credentials, err := json.Marshal(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiry,
	})
	if err != nil {
		return nil, errors.New("failed to marshal oauth credentials")
	}
	return credentials, nil
}

// MakeTokenString returns a signed session JWT for the given user.
// The Google OAuth tokens are not embedded; they stay in the users.credentials column.
func MakeTokenString(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(TokenLifetime).Unix(),
	})
	tokenString, err := token.SignedString([]byte(GetSecretKey()))
	if err != nil {
		return "", errors.New("Failed to create token string: " + err.Error())
	}
	return tokenString, nil
}

// getRequestToken reads the session token from the "Authorization: Bearer <token>" header,
// falling back to an "Authorization" cookie (used by local development and the tests).
func getRequestToken(c *gin.Context) (string, error) {
	if header := c.GetHeader("Authorization"); header != "" {
		if token, ok := strings.CutPrefix(header, "Bearer "); ok && token != "" {
			return token, nil
		}
		return "", errors.New("malformed Authorization header")
	}
	return c.Cookie("Authorization")
}

// AuthenticateUser validates the session token on the request.
// Returns the user's ID if the token is valid. Otherwise, returns an error.
func AuthenticateUser(c *gin.Context) (uint, error) {
	token, err := getRequestToken(c)
	if err != nil {
		return 0, err
	}
	if token == testToken && os.Getenv("APP_ENV") == "development" {
		return 0, nil
	}

	claims := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(GetSecretKey()), nil
	})
	if err != nil {
		return 0, err
	}
	if !parsed.Valid {
		return 0, errors.New("invalid token")
	}
	// Tokens without an expiry (e.g. from older versions of the app) are rejected.
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return 0, errors.New("token expired")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok || userID <= 0 {
		return 0, errors.New("unable to parse user_id from JWT claims")
	}

	// Make sure the account still exists
	var user models.User
	if err := database.DB.Select("id").First(&user, uint(userID)).Error; err != nil {
		return 0, errors.New("user doesn't exist in DB")
	}

	return user.ID, nil
}

// AuthenticateMiddleware is middleware that checks if the user is authenticated.
func AuthenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := AuthenticateUser(c)
		if err != nil {
			fmt.Println("unauthenticated: " + err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
