package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opalescencelabs/backend/api"
	"github.com/opalescencelabs/backend/controllers/auth"
	"github.com/opalescencelabs/backend/controllers/templates"
	"github.com/opalescencelabs/backend/database"
	"github.com/opalescencelabs/backend/models"
)

// UserLogin is the handler for POST /user-login.
// New/Returning users are authenticated via Google OAuth2.
// User information is fetched from Google's userinfo endpoint,
// new users are added to the database, returning users have changed
// information updated in the database.
// Returns 200 on success, 400 on bad request, 401 on unauthorized, 500 on error.
func UserLogin(c *gin.Context) {
	var request api.UserLoginRequest

	if err := c.BindJSON(&request); err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// Validate access_token
	if err := auth.ValidateAccessToken(request.AccessToken); err != nil {
		fmt.Println("Validation failed:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Use access_token to get userInfo
	info, err := auth.GetUserInfo(request.AccessToken)
	if err != nil {
		fmt.Println("Failed to get user info: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	email := info["email"].(string)
	// is_verified := info["email_verified"].(bool)
	name := info["given_name"].(string) + " " + info["family_name"].(string)
	photo_url := info["picture"].(string)
	google_id := info["sub"].(string)
	credentials, err := auth.MakeCredentialsJson(request.AccessToken, request.RefreshToken, request.ExpiresIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal oauth credentials"})
		return
	}

	var user models.User
	database.DB.First(&user, "google_id = ?", google_id)

	// Add new user to DB if they don't exist -> Using a trnsaction to ensure that the user and welcome page are created together
	is_new_user := user.ID == 0
	if is_new_user {
		user = models.User{
			GoogleID:    google_id,
			Email:       email,
			Name:        name,
			Picture:     photo_url,
			Credentials: credentials,
		}
		// Using a transaction to ensure that the user and welcome page are created together
		tx := database.DB.Begin()
		if err := tx.Omit("id", "status").Create(&user).Error; err != nil {
			fmt.Println("Failed to create new user: ", err.Error())
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new user"})
			return
		}
		if err := templates.CreateWelcomePage(user.ID, tx); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create welcome page"})
			fmt.Println("Failed to create welcome page: ", err.Error())
			return
		}
		// Commit the transaction
		tx.Commit()
	} else {
		// Update user info if it has changed
		user.Credentials = credentials
		if user.Name != name {
			user.Name = name
		}
		if user.Picture != photo_url {
			user.Picture = photo_url
		}
		database.DB.Save(&user)
	}

	tokenString, err := auth.MakeTokenString(request.AccessToken, request.RefreshToken, request.ExpiresIn)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Attach token to browser
	c.SetSameSite(http.SameSiteLaxMode)

	// Trim the protocol and port from the frontend URL
	url := strings.TrimPrefix(auth.GetFrontendURL(), "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.Split(url, ":")[0]

	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", url, os.Getenv("APP_ENV") != "development", true)
	// c.JSON(http.StatusOK, api.UserLoginResp{})
	c.JSON(http.StatusOK, gin.H{"token": tokenString}) // TODO: Fix issue where cookie is not being set
}

// UserLogout is the handler for POST /user-logout.
// Revokes the access_token using Google's token revocation endpoint and
// deletes the tokens from the client's browser.
// Returns 200 on success, 401 on unauthorized, 500 on error.
func UserLogout(c *gin.Context) {
	userID, err := auth.AuthenticateUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Find the user in DB
	var user models.User
	if err := database.DB.First(&user, userID); err.Error != nil {
		fmt.Println("ERROR: ", err.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user does not exist"})
		return
	}

	var credentials map[string]interface{}
	if err := json.Unmarshal(user.Credentials, &credentials); err != nil {
		fmt.Println("ERROR: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unmarshal credentials"})
		return
	}

	// Revoke tokens
	if err := auth.RevokeAccessToken(credentials["access_token"].(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "access token revocation failure"})
		return
	}

	// Delete token from browser

	// Trim the protocol and port from the frontend URL
	url := strings.TrimPrefix(auth.GetFrontendURL(), "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.Split(url, ":")[0]

	c.SetCookie("Authorization", "", -1, "/", url, os.Getenv("APP_ENV") != "development", true)

	fmt.Printf("user-logout: %d\n", userID)
	c.JSON(http.StatusOK, api.UserLogoutResp{})
}

// UserGet is the handler for GET /user-get.
// Returns the user's information.
// Returns 200 on success, 401 on unauthorized, 500 on error.
func UserGet(c *gin.Context) {
	userID, err := auth.AuthenticateUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return user data in the response
	c.JSON(http.StatusOK, api.UserGetResp{
		ID:       user.ID,
		GoogleID: user.GoogleID,
		Name:     user.Name,
		Email:    user.Email,
		Picture:  user.Picture,
		Status:   user.Status,
		// Add other user data as needed
	})
}
