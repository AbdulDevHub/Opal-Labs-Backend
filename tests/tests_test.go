package tests

import (
	"encoding/json"
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
	// CI sets variables directly; locally they come from ../.env
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No ../.env file found, using environment variables")
	}
	if os.Getenv("DOMAIN") == "" {
		os.Setenv("DOMAIN", "http://localhost:8000")
	}
	log.Println("Connecting to DB...")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_SSL_MODE"))
	}
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		log.Fatal("Failed to connect to the test database: ", err)
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

	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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

	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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

	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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

	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Page-level etc (icon, cover) is saved by page-update and returned by page-get
func TestUpdatePageEtc(t *testing.T) {
	send := func(method, path, payload string) *http.Response {
		req, err := http.NewRequest(method, os.Getenv("DOMAIN")+path, strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer Only_for_testing1200332")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	created := send("POST", "/page-create", `{"page_uuid":"12234PageEtctest","page_name":"PageEtctest","is_root":true,"element_positions":[]}`)
	created.Body.Close()
	assert.Equal(t, http.StatusOK, created.StatusCode)
	// Clean up at the end (a plain defer would send the delete straight away)
	defer func() { send("POST", "/page-delete", `{"page_uuid":"12234PageEtctest"}`).Body.Close() }()

	resp := send("POST", "/page-update", `{"page":{"page_uuid":"12234PageEtctest","etc":{"icon":"🚀","cover":"aurora"}}}`)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// An update without etc keeps the saved value
	send("POST", "/page-update", `{"page":{"page_uuid":"12234PageEtctest","page_name":"Renamed"}}`).Body.Close()

	resp = send("GET", "/page-get/12234PageEtctest", "")
	defer resp.Body.Close()
	var got struct {
		Page struct {
			PageName string                 `json:"page_name"`
			Etc      map[string]interface{} `json:"etc"`
		} `json:"page"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Renamed", got.Page.PageName)
	assert.Equal(t, "🚀", got.Page.Etc["icon"])
	assert.Equal(t, "aurora", got.Page.Etc["cover"])
}
