package api

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgtype"
)

// Holds all page related api request and response structs

// Page Create

type PageCreateRequest struct {
	PageUUID         string       `json:"page_uuid"`
	PageName         string       `json:"page_name,omitempty"`
	IsRoot           bool         `json:"is_root,omitempty"`
	ParentPageUUID   *string      `json:"parent_page_uuid,omitempty"`
	ElementPositions pgtype.JSONB `json:"element_positions,omitempty"`
	PublicPage       bool         `json:"public_page,omitempty"`
	IsFavourite      bool         `json:"is_favourite,omitempty"`
	Etc              pgtype.JSONB `json:"etc,omitempty"`
}

type PageCreateResp struct{}

type ElementsResponseObject struct {
	ID          uint                   `json:"id"`
	ElementUUID string                 `json:"element_uuid"`
	Type        string                 `json:"type"`
	Content     map[string]interface{} `json:"content"`
	Etc         map[string]interface{} `json:"etc"`
	Size        string                 `json:"size"`
}

// Page Get
type PageGetResp struct {
	Page     PageResp                 `json:"page"`
	Elements []ElementsResponseObject `json:"elements"` // Should be in order
	SubPages map[string]string        `json:"sub_pages,omitempty"`
}

type PageGetRespOwner struct {
	Page     PageRespOwner            `json:"page"`
	Elements []ElementsResponseObject `json:"elements"` // Should be in order
	SubPages map[string]string        `json:"sub_pages,omitempty"`
}

type PageRespOwner struct {
	ID               uint                   `json:"id"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	DeletedAt        time.Time              `json:"deleted_at,omitempty"`
	PageUUID         string                 `json:"page_uuid"`
	PageName         string                 `json:"page_name"`
	IsRoot           bool                   `json:"is_root"`
	ElementPositions []string               `json:"element_positions,omitempty"`
	ParentPageUUID   string                 `json:"parent_page_uuid,omitempty"`
	PublicPage       bool                   `json:"public_page"`
	PageUUIDURL      string                 `json:"page_uuid_url,omitempty"`
	IsFavourite      bool                   `json:"is_favourite"`
	LastUpdatedAt    time.Time              `json:"last_updated_at"`
	ViewCount        uint                   `json:"view_count"`
	DateViewCount    map[string]int         `json:"date_view_count"`
	Etc              map[string]interface{} `json:"etc"`
}

// Implement encoding.BinaryMarshaler to store api.PageGetResp in the cache
// MarshalBinary encodes the response to store in the cache
func (p PageGetResp) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

// Implement encoding.BinaryMarshaler to store api.PageGetRespOwner in the cache
// MarshalBinary encodes the response to store in the cache
func (p PageGetRespOwner) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

// Page List
type PageListResp struct {
	Pages []PageResp `json:"pages"`
}

type PageResp struct {
	ID               uint                   `json:"id"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	PageUUID         string                 `json:"page_uuid"`
	PageName         string                 `json:"page_name"`
	IsRoot           bool                   `json:"is_root"`
	ElementPositions []string               `json:"element_positions,omitempty"`
	ParentPageUUID   string                 `json:"parent_page_uuid,omitempty"`
	PublicPage       bool                   `json:"public_page"`
	PageUUIDURL      string                 `json:"page_uuid_url,omitempty"`
	IsFavourite      bool                   `json:"is_favourite"`
	LastUpdatedAt    time.Time              `json:"last_updated_at"`
	ViewCount        uint                   `json:"view_count"`
	Etc              map[string]interface{} `json:"etc"`
}

// Page Update
type PageUpdateRequest struct {
	Page     PageUpdateObject       `json:"page"`
	Elements []ElementsUpdateObject `json:"elements,omitempty"`
	// RemoveElements []string               `json:"remove_elements,omitempty"`
}

type PageUpdateObject struct {
	PageUUID       string `json:"page_uuid"`
	PageName       string `json:"page_name"`
	IsRoot         *bool  `json:"is_root"`
	ParentPageUUID string `json:"parent_page_uuid,omitempty"`
	PublicPage     *bool  `json:"public_page"`
	IsFavourite    *bool  `json:"is_favourite,omitempty"`
}

type ElementsUpdateObject struct {
	ElementUUID string       `json:"element_uuid"`
	Type        string       `json:"type"`
	Content     pgtype.JSONB `json:"content"`
	Etc         pgtype.JSONB `json:"etc"`
	Size        string       `json:"size,omitempty"`
}

type PageUpdateResp struct{}

// Page Delete

type PageDeleteReq struct {
	PageUUID string `json:"page_uuid"`
}

type PageDeleteResp struct{}
