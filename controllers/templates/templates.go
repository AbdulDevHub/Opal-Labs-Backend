package templates

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/opalescencelabs/backend/models"
	"gorm.io/gorm"
)

// CreateWelcomePage creates a welcome page with a title and a welcome message.
// The template is found in the templates directory.
// -> /controllers/templates/welcome
// Returns an error if the page cannot be created for the current user.
func CreateWelcomePage(userID uint, tx *gorm.DB) error {
	// Read the template data
	page, elements, err := ReadTemplateData("./controllers/templates/welcome")
	if err != nil {
		return err
	}

	// Assign uuid and id to page and elements
	page.PageUUID = uuid.New().String()
	page.UserID = userID
	var elementPositions []string
	for _, element := range elements {
		element.ElementUUID = uuid.New().String()
		elementPositions = append(elementPositions, element.ElementUUID)
	}

	// Update element positions
	// Convert element positions to json
	jsonElementPositions, err := json.Marshal(elementPositions)
	if err != nil {
		return err
	}
	// Convert bytes to pgtype.JSONB
	var elementPositionsJSONB pgtype.JSONB
	if err := elementPositionsJSONB.Set(string(jsonElementPositions)); err != nil {
		return err
	}

	// Commit page and elements to the database
	page.ElementPositions = elementPositionsJSONB
	if err := tx.Omit("id").Create(&page).Error; err != nil {
		return err
	}
	for i, element := range elements {
		element.PageID = page.ID
		element.UserID = userID
		element.ElementUUID = elementPositions[i]
		if err := tx.Omit("id").Create(&element).Error; err != nil {
			return err
		}
	}

	return nil
}

// ReadTemplateData reads the page and elements from a template directory.
// The files in the template directory are the page and elements for the page.
// The page is the singular file (page.json) in the directory templateDirectory.
// The files in templateDirectory/elements are the elements for the page.
// The order of the elements is alphabetical by element file name (for now).
func ReadTemplateData(templateDirectory string) (models.Page, []models.Element, error) {
	// Check if the template directory exists
	if _, err := os.ReadDir(templateDirectory); err != nil {
		fmt.Println("Error: ", err.Error())
		fmt.Println("Template directory does not exist")
		fmt.Println("Template directory: ", templateDirectory)
		return models.Page{}, nil, err
	}

	// If the template directory exists, read the page and elements
	// The page is the singular file in the directory
	// The elements are the files in the elements directory
	var page models.Page
	var elements []models.Element

	// Read page
	pageFile := templateDirectory + "/page.json"
	page, err := ReadPage(pageFile)
	if err != nil {
		return models.Page{}, nil, err
	}

	// Read elements
	elementsDirectory := templateDirectory + "/elements"
	elementsDir, err := os.ReadDir(elementsDirectory)
	if err != nil {
		return models.Page{}, nil, err
	}

	for _, elementFile := range elementsDir {
		element, err := ReadElement(elementsDirectory + "/" + elementFile.Name())
		if err != nil {
			fmt.Printf("Failed to read element: %s\n", elementFile.Name())
			return models.Page{}, nil, err
		}
		elements = append(elements, element)
	}

	return page, elements, nil
}

// ReadPageTemplate reads a page template and returns the page as a models.Page object,
func ReadPage(pageFile string) (models.Page, error) {
	// Read file
	data, err := os.ReadFile(pageFile)
	if err != nil {
		return models.Page{}, err
	}

	var page models.Page
	err = json.Unmarshal(data, &page)
	if err != nil {
		return models.Page{}, err
	}

	return page, nil
}

// ReadElement reads an element and returns the element as a models.Element object.
func ReadElement(elementFile string) (models.Element, error) {
	// Read file
	data, err := os.ReadFile(elementFile)
	if err != nil {
		return models.Element{}, err
	}

	var element models.Element
	err = json.Unmarshal(data, &element)
	if err != nil {
		return models.Element{}, err
	}

	return element, nil
}
