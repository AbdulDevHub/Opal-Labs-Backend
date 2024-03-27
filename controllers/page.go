package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	"github.com/opalescencelabs/backend/api"
	"github.com/opalescencelabs/backend/api/caching"
	"github.com/opalescencelabs/backend/controllers/auth"
	"github.com/opalescencelabs/backend/database"
	"github.com/opalescencelabs/backend/models"
	"gorm.io/gorm"
)

// PageCreate is the handler for POST /page/create
// Creates a new page in the database given the request and authentication.
// Returns 200 on success, 400 on bad request, 401 on unauthorized, 500 on error.
func PageCreate(c *gin.Context) {
	userID, err := auth.AuthenticateUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request api.PageCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Check if page_uuid already exists
	var existingPage models.Page

	if request.PageUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page UUID cannot be empty"})
		return
	}

	if err := database.DB.Where("page_uuid = ?", request.PageUUID).First(&existingPage).Error; err == nil {
		// If the record is found, return an error response
		c.JSON(http.StatusConflict, gin.H{"error": "Page UUID already in use"})
		return
	}

	if len(request.ElementPositions.Bytes) == 0 {
		request.ElementPositions = pgtype.JSONB{Status: pgtype.Null}
	} else if string(request.ElementPositions.Bytes) == "[]" {
		request.ElementPositions.Status = pgtype.Null
	} else {
		request.ElementPositions.Status = pgtype.Present
	}

	if len(request.Etc.Bytes) == 0 {
		request.Etc = pgtype.JSONB{Status: pgtype.Null}
	} else if string(request.Etc.Bytes) == "{}" {
		request.Etc.Status = pgtype.Null
	} else {
		request.Etc.Status = pgtype.Present
	}

	var newPage models.Page
	saveresult := database.DB.Model(&newPage).Create(map[string]interface{}{
		"created_at":        time.Now(),
		"page_uuid":         request.PageUUID,
		"page_name":         request.PageName,
		"is_root":           request.IsRoot,
		"UserID":            userID,
		"parent_page_uuid":  request.ParentPageUUID,
		"element_positions": request.ElementPositions,
		"public_page":       request.PublicPage,
		"is_favourite":      request.IsFavourite,
		"etc":               request.Etc,
		"view_count":        0,
		"last_updated_at":   time.Now(),
	})
	if saveresult.Error != nil {
		fmt.Println("Failed to create page: ", saveresult.Error)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, api.PageCreateResp{})
}

// PageGet is the handler for POST /page-get/:page_uuid.
// Returns the page with the given UUID from the database.
// Uses the cache if available, otherwise fetches from the database and stores in the cache.
// If the page is public, increments the page view count.
// Returns 200 on success, 400 on bad request, 401 on unauthorized, 500 on error.
func PageGet(c *gin.Context) {
	pageUUID := c.Param("page_uuid")
	fmt.Println("Page UUID: ", pageUUID)
	if pageUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page UUID cannot be empty"})
		return
	}

	userID, user_err := auth.AuthenticateUser(c)

	// Check if the response is already cached
	cachedResponse, err := caching.Check(c, fmt.Sprintf("/page-get/%s", pageUUID))
	if err != nil {
		fmt.Printf("Response for /page-get/%s not found in cache\n", pageUUID)
	} else {
		// Return the cached response as JSON
		fmt.Printf("Response for /page-get/%s found in cache\n", pageUUID)

		// If the page is public, increment the view count and update the cache
		// Unmarshmal the cached response
		var response api.PageGetRespOwner
		var pub_response api.PageGetResp
		var flag = 0
		err := json.Unmarshal([]byte(cachedResponse.(string)), &response)
		if err == nil {
			if response.Page.PublicPage {
				// If page is public, increment the view count and update the cache, database
				fmt.Printf("Response /page-get/%s -> page:%s is public, incrementing view count\n", pageUUID, pageUUID)
				response.Page.ViewCount++
				response.Page.DateViewCount[time.Now().Format("2006-01-02")]++

				responseMarshaled, err := json.Marshal(response)
				if err != nil {
					return
				}
				if err := caching.Store(c, fmt.Sprintf("/page-get/%s", pageUUID), responseMarshaled); err != nil {
					// (failed to update cache)
					fmt.Printf("Failed to store /page-get/%s in cache: %s\n", pageUUID, err)
					// Not a critical error, so we can continue with the rest of the handler
				}
				// Update the database
				var page models.Page
				if err := database.DB.Where("page_uuid = ?", pageUUID).First(&page).Error; err != nil {
				}
				page.ViewCount++
				type ViewCountData map[string]int
				var viewCounts ViewCountData

				currentDate := time.Now().Format("2006-01-02")
				// Increment view count
				if page.DateViewCount.Status != pgtype.Present {
					viewCounts = ViewCountData{currentDate: 1}
				} else {
					err := json.Unmarshal(page.DateViewCount.Bytes, &viewCounts)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process view count data", "details": err.Error()})
						return
					}

					viewCounts[currentDate]++
				}

				updatedBytes, err := json.Marshal(viewCounts)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process view count data", "details": err.Error()})
					return
				}

				page.DateViewCount = pgtype.JSONB{Bytes: updatedBytes, Status: pgtype.Present}

				if err := database.DB.Save(&page).Error; err != nil {
					fmt.Printf("Failed to update view count for /page-get/%s: %s\n", pageUUID, err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update view count"})
					return
				}

				if page.UserID != userID {
					flag = 1
					pub_response = api.PageGetResp{
						Page: api.PageResp{
							ID:             response.Page.ID,
							CreatedAt:      response.Page.CreatedAt,
							UpdatedAt:      response.Page.UpdatedAt,
							PageUUID:       response.Page.PageUUID,
							PageName:       response.Page.PageName,
							IsRoot:         response.Page.IsRoot,
							ParentPageUUID: response.Page.ParentPageUUID,
							PublicPage:     response.Page.PublicPage,
							PageUUIDURL:    response.Page.PageUUIDURL,
							IsFavourite:    response.Page.IsFavourite,
							ViewCount:      response.Page.ViewCount,
							LastUpdatedAt:  response.Page.LastUpdatedAt,
							Etc:            response.Page.Etc,
						},
						Elements: response.Elements,
						SubPages: response.SubPages,
					}

				}

			} else {
				if user_err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
					return
				}
			}
			// Return the cached response, updated or not
			if flag == 1 {
				c.JSON(http.StatusOK, pub_response)
				return
			} else {
				c.JSON(http.StatusOK, response)
				return
			}
		} else { // (unmarshalling failed)
			// If unmarshalling fails, invalidate the cache and return to the rest of the handler
			fmt.Printf("Failed to unmarshal cached response for /page-get/%s: %s\n", pageUUID, err)
			if err := caching.Invalidate(c, fmt.Sprintf("/page-get/%s", pageUUID)); err != nil {
				fmt.Printf("Failed to invalidate /page-get/%s in cache: %s\n", pageUUID, err)
			}
			// Not a critical error, so we can continue
		}

	}

	var page models.Page
	result := database.DB.Where("page_uuid = ?", pageUUID).First(&page)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	if user_err != nil && !page.PublicPage {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userID != page.UserID && !page.PublicPage {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if page.PublicPage {
		type ViewCountData map[string]int
		var viewCounts ViewCountData

		currentDate := time.Now().Format("2006-01-02")
		// Increment view count
		if page.DateViewCount.Status != pgtype.Present {
			viewCounts = make(map[string]int)
			viewCounts = ViewCountData{currentDate: 1}
		} else {
			err := json.Unmarshal(page.DateViewCount.Bytes, &viewCounts)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process view count data", "details": err.Error()})
				return
			}
			if viewCounts == nil {
				viewCounts = make(map[string]int) // Initialize the map if it was nil
			}

			viewCounts[currentDate]++
		}

		updatedBytes, err := json.Marshal(viewCounts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process view count data", "details": err.Error()})
			return
		}

		page.DateViewCount = pgtype.JSONB{Bytes: updatedBytes, Status: pgtype.Present}
		page.ViewCount++
		updateResult := database.DB.Save(&page)
		if updateResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update view count"})
			return
		}
	}

	// Unmarshal DateViewCount if present and user ID matches
	var dateViewCountData map[string]int
	if page.DateViewCount.Status == pgtype.Present && userID == page.UserID {
		err := json.Unmarshal(page.DateViewCount.Bytes, &dateViewCountData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process date view count data", "details": err.Error()})
			return
		}
	}

	var elementPositions []string

	// Unmarshal ElementPositions if present
	if page.ElementPositions.Status == pgtype.Present {
		err := json.Unmarshal(page.ElementPositions.Bytes, &elementPositions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process page data", "details": err.Error()})
			return
		}
	}
	var PageEtc map[string]interface{}
	if page.Etc.Status == pgtype.Present {
		err := json.Unmarshal(page.Etc.Bytes, &PageEtc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process page data", "details": err.Error()})
			return
		}
	}

	var elements []models.Element
	// Replace 'ForeignKeyColumn' with the actual foreign key column name in your Element model that references Page
	result = database.DB.Where("page_id = ?", page.ID).Find(&elements)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch elements for the page"})
		return
	}

	responseElements := make([]api.ElementsResponseObject, 0, len(elements))
	// Add the elements to the response in the order of the element positions
	positionMap := make(map[string]int)
	for i, uuid := range elementPositions {
		positionMap[uuid] = i
	}
	// Sort the elements by their position in the elementPositions slice
	sort.Slice(elements, func(i, j int) bool {
		positionI := positionMap[elements[i].ElementUUID]
		positionJ := positionMap[elements[j].ElementUUID]
		return positionI < positionJ
	})

	for _, element := range elements {
		var content map[string]interface{}
		var etc map[string]interface{}

		// Unmarshal Content if present
		if element.Content.Status == pgtype.Present {
			if err := json.Unmarshal(element.Content.Bytes, &content); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal element content", "details": err.Error()})
				return
			}
		}

		// Unmarshal Etc if present
		if element.Etc.Status == pgtype.Present {
			if err := json.Unmarshal(element.Etc.Bytes, &etc); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal element etc", "details": err.Error()})
				return
			}
		}

		// Append the processed element to the response slice
		responseElements = append(responseElements, api.ElementsResponseObject{
			ID:          element.ID,
			ElementUUID: element.ElementUUID,
			Type:        element.Type,
			Content:     content,
			Etc:         etc,
			Size:        element.Size,
		})
	}
	fmt.Println("Page UUID: ", pageUUID)

	// Query all Sub-Pages
	var subPages []models.Page
	subPagesResult := database.DB.Where("parent_page_uuid = ?", pageUUID).Find(&subPages)
	if subPagesResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sub-pages"})
		return
	}

	// Create a map of sub-pages for easy access
	subPagesMap := make(map[string]string)
	for _, subPage := range subPages {
		fmt.Println("SubPage: ", subPage.PageUUID)
		subPagesMap[subPage.PageUUID] = subPage.PageName
	}

	var response any
	// Construct the response based on the user's ownership of the page
	if userID == page.UserID {
		response = api.PageGetRespOwner{
			Page: api.PageRespOwner{
				ID:             page.ID,
				CreatedAt:      page.CreatedAt,
				UpdatedAt:      page.UpdatedAt,
				PageUUID:       page.PageUUID,
				PageName:       page.PageName,
				IsRoot:         page.IsRoot,
				ParentPageUUID: page.ParentPageUUID,
				PublicPage:     page.PublicPage,
				PageUUIDURL:    page.PageUUIDURL,
				IsFavourite:    page.IsFavourite,
				ViewCount:      page.ViewCount,
				DateViewCount:  dateViewCountData,
				LastUpdatedAt:  page.LastUpdatedAt,
				Etc:            PageEtc,
			},
			Elements: responseElements,
			SubPages: subPagesMap,
		}
	} else {
		response = api.PageGetResp{
			Page: api.PageResp{
				ID:             page.ID,
				CreatedAt:      page.CreatedAt,
				UpdatedAt:      page.UpdatedAt,
				PageUUID:       page.PageUUID,
				PageName:       page.PageName,
				IsRoot:         page.IsRoot,
				ParentPageUUID: page.ParentPageUUID,
				PublicPage:     page.PublicPage,
				PageUUIDURL:    page.PageUUIDURL,
				IsFavourite:    page.IsFavourite,
				ViewCount:      page.ViewCount,
				LastUpdatedAt:  page.LastUpdatedAt,
				Etc:            PageEtc,
			},
			Elements: responseElements,
			SubPages: subPagesMap,
		}
	}

	// Only store the owner's response in the cache
	if userID == page.UserID {
		if err := caching.Store(c, fmt.Sprintf("/page-get/%s", pageUUID), response); err != nil {
			fmt.Printf("Failed to store /page-get/%s in cache: %s\n", pageUUID, err)
			// Not a critical error, so we can continue
		}
	}

	c.JSON(http.StatusOK, response)
}

// PageList is the handler for POST /page-list.
// Returns a list of the user's pages from the database.
// Returns 200 on success, 401 on unauthorized, 404 on no pages found, 500 on error.
func PageList(c *gin.Context) {
	var pages []models.Page
	userID, err := auth.AuthenticateUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	result := database.DB.Where("user_id = ? AND deleted_at IS NULL", userID).Order("is_favourite DESC, last_updated_at DESC").Find(&pages)
	if result.Error != nil || len(pages) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No pages found for this user"})
		return
	}

	// Convert pages to PageResp, including unmarshalling ElementPositions
	pageResps := make([]api.PageResp, len(pages))
	for i, page := range pages {
		// Initialize the slice to store the unmarshalled element positions
		var elementPositions []string

		// Check if the JSONB data is present and not null
		if page.ElementPositions.Status == pgtype.Present && len(page.ElementPositions.Bytes) > 0 {
			// Directly unmarshal the JSONB Bytes into the []string slice
			err := json.Unmarshal(page.ElementPositions.Bytes, &elementPositions)
			if err != nil {
				// Handle JSON unmarshal error if necessary
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process page data (elementPositions)", "details": err.Error()})
				return
			}
		}

		var etc map[string]interface{}
		if page.Etc.Status == pgtype.Present && len(page.Etc.Bytes) > 0 {
			// Directly unmarshal the JSONB Bytes into the []string slice
			err := json.Unmarshal(page.Etc.Bytes, &etc)
			if err != nil {
				// Handle JSON unmarshal error if necessary
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process page data (etc)", "details": err.Error()})
				return
			}
		}
		// Manually map each field from Page to PageResp, excluding User
		pageResps[i] = api.PageResp{
			ID:        page.ID,
			CreatedAt: page.CreatedAt,
			UpdatedAt: page.UpdatedAt,
			PageUUID:  page.PageUUID,
			PageName:  page.PageName,
			IsRoot:    page.IsRoot,
			// ElementPositions: elementPositions,
			ParentPageUUID: page.ParentPageUUID,
			PublicPage:     page.PublicPage,
			PageUUIDURL:    page.PageUUIDURL,
			IsFavourite:    page.IsFavourite,
			ViewCount:      page.ViewCount,
			Etc:            etc,
			LastUpdatedAt:  page.LastUpdatedAt,
		}
	}

	c.JSON(http.StatusOK, api.PageListResp{Pages: pageResps})
}

// PageUpdate is the handler for POST /page-update.
// Updates a page in the database given the request and authentication.
// Invalidates the Page cache for the updated page.
// Returns 200 on success, 400 on bad request, 401 on unauthorized, 500 on error.
func PageUpdate(c *gin.Context) {
	var request api.PageUpdateRequest
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	userID, err := auth.AuthenticateUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var page models.Page
	PageUUID := request.Page.PageUUID
	if err := database.DB.Where("page_uuid = ? AND user_id = ?", PageUUID, userID).First(&page).Error; err != nil {
		// Page does not exist / user doesn't own
		fmt.Printf("Page not found for user %d (UUID %s)", userID, PageUUID)
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	parent_page_uuid := page.ParentPageUUID

	page.LastUpdatedAt = time.Now()
	saveResult := database.DB.Save(&page)
	if saveResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update last visited time"})
		return
	}

	// Wrap changes in a transaction in event of error
	tx := database.DB.Begin()
	// Rollback in case of an unexpected error...
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update page
	pageUpdate := request.Page

	pageUpdateQuery := tx.Model(&models.Page{})

	// Set the page ID, UUID to the one from the request
	pageUpdateQuery = pageUpdateQuery.Where("page_uuid = ? AND user_id = ?", page.PageUUID, userID)

	// Tack on Omit statements for fields that should not be updated (i.e. they are missing from the request)
	if pageUpdate.PageName == "" {
		pageUpdateQuery = pageUpdateQuery.Omit("page_name")
	}
	if pageUpdate.IsRoot == nil {
		pageUpdateQuery = pageUpdateQuery.Omit("is_root")
	}
	if pageUpdate.PublicPage == nil {
		pageUpdateQuery = pageUpdateQuery.Omit("public_page")
	} else {
		// Reset view counts if toggling public_page status
		if *pageUpdate.PublicPage != page.PublicPage {
			pageUpdateQuery = pageUpdateQuery.Update("view_count", 0)
			pageUpdateQuery = pageUpdateQuery.Update("date_view_count", pgtype.JSONB{Bytes: []byte("{}"), Status: pgtype.Null})
		} else {
			pageUpdateQuery = pageUpdateQuery.Omit("view_count")
			pageUpdateQuery = pageUpdateQuery.Omit("date_view_count")
		}

	}
	if pageUpdate.IsFavourite == nil {
		pageUpdateQuery = pageUpdateQuery.Omit("is_favourite")
	}
	// Rest of these fields are never updated
	pageUpdateQuery.Omit("etc", "page_uuid", "element_positions", "user_id")

	if err := pageUpdateQuery.Updates(pageUpdate).Error; err != nil {
		fmt.Println("Failed to update page", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update page" + err.Error()})
		tx.Rollback()
		return
	}

	PageID := page.ID

	elementQuery := "element_uuid = ? AND page_id = ? AND user_id = ?"
	var newElementPositions []string

	// update, create new elements
	for _, update := range request.Elements {
		var element models.Element
		if err := database.DB.Where(elementQuery, update.ElementUUID, PageID, userID).First(&element).Error; err != nil {
			// element not found -> new element
			element.ElementUUID = update.ElementUUID
			element.PageID = PageID
			element.UserID = userID
		}
		element.Type = update.Type
		element.Content = update.Content
		element.Etc = update.Etc
		element.Size = update.Size
		// create or update
		if err := tx.Save(&element).Error; err != nil {
			fmt.Println("Failed to update elements", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update elements"})
			tx.Rollback()
			return
		}
		// Add the element UUID to the newElementPositions slice
		newElementPositions = append(newElementPositions, update.ElementUUID)
	}

	// Delete elements that are not present in the newElementPositions slice

	// Unmarshal the existing element positions
	var existingElementPositions []string
	if page.ElementPositions.Status == pgtype.Present {
		if err := json.Unmarshal(page.ElementPositions.Bytes, &existingElementPositions); err != nil {
			fmt.Println("Failed to unmarshal existing element positions", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process page data (elementPositions)", "details": err.Error()})
			tx.Rollback()
			return
		}
	}

	// Iterate over the existing element positions and delete the ones that are not present in the newElementPositions slice
	for _, existingElementUUID := range existingElementPositions {
		if !slices.Contains(newElementPositions, existingElementUUID) {
			// Delete the element
			var element models.Element
			if err := tx.Where("element_uuid = ? AND page_id = ? AND user_id = ?", existingElementUUID, PageID, userID).First(&element).Error; err != nil {
				fmt.Println("Element not found:", existingElementUUID, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "delete element that doesn't exist"})
				tx.Rollback()
				return
			}
			if err := tx.Delete(&element).Error; err != nil {
				fmt.Println("Failed to delete elements", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "remove elements"})
				tx.Rollback()
				return
			}
		}
	}

	// Marshal the newElementPositions slice and update the page's ElementPositions
	newElementPositionsJSON, err := json.Marshal(newElementPositions)
	if err != nil {
		fmt.Println("Failed to marshal new element positions", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process page data (elementPositions)", "details": err.Error()})
		tx.Rollback()
		return
	}
	if err := tx.Model(&page).Update("element_positions", newElementPositionsJSON).Error; err != nil {
		fmt.Println("Failed to update element positions", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update element positions"})
		tx.Rollback()
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		fmt.Println("Unexpected error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "enexpected"})
		return
	}

	// Invalidate cache for this page post update
	caching.Invalidate(c, fmt.Sprintf("/page-get/%s", PageUUID))

	if parent_page_uuid != "" {
		caching.Invalidate(c, fmt.Sprintf("/page-get/%s", parent_page_uuid))
	}

	// Success
	c.JSON(http.StatusOK, api.PageUpdateResp{})
}

// Helper function to recursively delete a page and its children.
// This function is called by PageDelete.
// Invalidates the Page cache for the deleted page and its children.
// Returns an error if one occurs, nil otherwise.
func deletePageAndChildren(c *gin.Context, tx *gorm.DB, pageUUID string, userID uint) error {
	var childPages []models.Page
	// Find all pages that have the parent page UUID of the page we are deleting
	if err := tx.Unscoped().Where("parent_page_uuid = ? AND user_id = ?", pageUUID, userID).Find(&childPages).Error; err != nil {
		return err
	}

	// Recursively delete child pages and their elements
	for _, childPage := range childPages {
		if err := deletePageAndChildren(c, tx, childPage.PageUUID, userID); err != nil {
			return err
		}
	}

	// Delete elements associated with the current pageUUID before deleting the page itself
	if err := tx.Unscoped().Where("page_id = (SELECT id FROM pages WHERE page_uuid = ? AND user_id = ?)", pageUUID, userID).Delete(&models.Element{}).Error; err != nil {
		return err
	}

	// Finally, delete the page itself
	if err := tx.Unscoped().Where("page_uuid = ? AND user_id = ?", pageUUID, userID).Delete(&models.Page{}).Error; err != nil {
		return err
	}

	// Invalidate the cache for the deleted page
	if err := caching.Invalidate(c, fmt.Sprintf("/page-get/%s", pageUUID)); err != nil {
		return err
	}

	return nil // Successfully deleted the page and its children
}

// PageDelete is the handler for POST /page-delete.
// Deletes a page and its associated elements and sub-pages from the database.
// Invalidates the Page cache for the deleted page.
// Returns 200 on success, 400 on bad request, 401 on unauthorized, 404 on not found, 500 on error.
func PageDelete(c *gin.Context) {
	var req api.PageDeleteReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := auth.AuthenticateUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err = database.DB.Where("page_uuid = ? AND user_id = ?", req.PageUUID, userID).First(&models.Page{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found or does not belong to the current user"})
		return
	}

	// Start a transaction
	tx := database.DB.Begin()

	// Pass the transaction to the recursive deletion process
	if err := deletePageAndChildren(c, tx, req.PageUUID, userID); err != nil {
		tx.Rollback() // Rollback the transaction if an error occurs
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found or does not belong to the current user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page and its associated elements/children"})
		}
		return
	}

	// Commit the transaction if no errors occurred
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error committing transaction"})
		return
	}

	c.JSON(http.StatusOK, api.PageDeleteResp{})
}
