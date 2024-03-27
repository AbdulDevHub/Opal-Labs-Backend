package tests

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"log"

	"github.com/joho/godotenv"
	"github.com/opalescencelabs/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

var DB *gorm.DB

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env", err.Error())
	}
	log.Println("Connecting to DB...")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		return
	}

	// Check if user with user_id 0 exists in the users table
	var user models.User
	result := DB.First(&user, "id = ?", 0)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		// User with user_id 1 not found, create a new user
		err := DB.Exec("INSERT INTO users (id, google_id, email, name, picture, credentials) VALUES (?, ?, ?, ?, ?, ?)", 0, "0000", "opalesencetest@gmail.com", "Opalescence Test", "https://lh3.googleusercontent.com/a-/AOh14Gh", "{}").Error
		if err != nil {
			return
		}
	}
}

func TestCreatePage(t *testing.T) {

	token := "Only_for_testing1200332"
	client := http.Client{}

	body := strings.NewReader(`{"page_uuid":"1234PageCreateTest", "page_name":"PageCreateTest", "is_root":true, "element_positions":["1", "2"]}`)

	req, err := http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-create", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}

	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	client = http.Client{}

	body = strings.NewReader(`{ "page_uuid":"1234PageCreateTest"}`)

	req, err = http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-delete", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie = http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

}

func TestCreatePageFail(t *testing.T) {

	token := "Only_for_testing1200332"
	client := http.Client{}

	body := strings.NewReader(`{}`)

	req, err := http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-create", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}

	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func TestGetPage(t *testing.T) {

	token := "Only_for_testing1200332"
	client := http.Client{}

	body := strings.NewReader(`{"page_uuid":"12234PageGettest", "page_name":"PageGettest", "is_root":true, "element_positions":["1", "2", "3"]}`)

	req, err := http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-create", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	client = http.Client{}

	req, err = http.NewRequest("GET", os.Getenv("DOMAIN")+"/page-get/12234PageGettest", nil)

	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie = http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	client = http.Client{}

	body = strings.NewReader(`{"page_uuid":"12234PageGettest"}`)

	req, err = http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-delete", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie = http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

}

func TestGetPageFail(t *testing.T) {

	token := "Only_for_testing1200332"
	client := http.Client{}

	body := strings.NewReader(`{"page_uuid":"999999"}`)

	req, err := http.NewRequest("GET", os.Getenv("DOMAIN")+"/page-get", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}

	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

}

func TestListPageFail(t *testing.T) {

	token := "Only_for_testing1200332"
	client := http.Client{}

	req, err := http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-list", nil)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

}

func TestListPage(t *testing.T) {

	token := "Only_for_testing1200332"
	client := http.Client{}

	body := strings.NewReader(`{"page_uuid":"12234PageListtest", "page_name":"PageListtest", "is_root":true, "element_positions":["1", "2", "3"]}`)

	req, err := http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-create", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	client = http.Client{}

	req, err = http.NewRequest("GET", os.Getenv("DOMAIN")+"/page-list", nil)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie = http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	client = http.Client{}

	body = strings.NewReader(`{ "page_uuid":"12234PageListtest"}`)

	req, err = http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-delete", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie = http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

}

func TestUpdatePage(t *testing.T) {

	token := "Only_for_testing1200332"
	client := http.Client{}

	body := strings.NewReader(`{
		"page_uuid":"12234PageUpdtest",
		"page_name":"PageUpdtest",
		"is_root":true,
		"element_positions":[]
		}`)
	req, err := http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-create", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	client = http.Client{}

	body = strings.NewReader(`{
		"page": {
			"page_uuid":"12234PageUpdtest",
			"page_name":"PageUpdtestNew",
			"is_root":true,
			"element_positions":[]
			},
		"elements": [
			{
				"element_uuid":"1234ElementUpdtest",
				"user_id":0,
				"page_id":0,
				"type":"text",
				"content":{},
				"etc":{},
				"size":"small"
			}
		]
	}`)

	req, err = http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-update", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie = http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	client = http.Client{}

	body = strings.NewReader(`{ "page_uuid":"12234PageUpdtest"}`)

	req, err = http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-delete", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie = http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

}

func TestUpdatePageFail(t *testing.T) {

	token := "Only_for_testing1200332"
	client := http.Client{}

	body := strings.NewReader(`{
		"page": {
			"page_uuid":"12234PageUpdtestFail", // This page does not exist
			"page_name":"PageUpdtestNew",
			"is_root":true,
			"element_positions":[]
			},
		"elements": [
			{
				"element_uuid":"1234ElementUpdtest",
				"user_id":0,
				"page_id":0,
				"type":"text",
				"content":{},
				"etc":{},
				"size":"small"
			}
		]
	}`)

	req, err := http.NewRequest("GET", os.Getenv("DOMAIN")+"/page-update", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}

	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

}

func TestDeletePage(t *testing.T) {

	token := "Only_for_testing1200332"
	client := http.Client{}

	body := strings.NewReader(`{"page_uuid":"12234PageDeltest", "page_name":"PageDeltest", "is_root":true, "element_positions":["1", "2", "3"]}`)

	req, err := http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-create", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	client = http.Client{}

	body = strings.NewReader(`{ "page_uuid":"12234PageDeltest"}`)

	req, err = http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-delete", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie = http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeletePageFail(t *testing.T) {

	token := "Only_for_testing1200332"
	client := http.Client{}

	body := strings.NewReader(`{"page_uuid":"12234ShouldNotExistDeltest"}`)

	req, err := http.NewRequest("POST", os.Getenv("DOMAIN")+"/page-delete", body)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUserGet(t *testing.T) {
	token := "Only_for_testing1200332"
	client := http.Client{}

	req, err := http.NewRequest("GET", os.Getenv("DOMAIN")+"/user-get", nil)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "Authorization", Value: token, HttpOnly: true, Secure: false, Domain: "localhost", Path: "/"}
	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
