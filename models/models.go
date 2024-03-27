package models

import (
	"time"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Element struct {
	gorm.Model
	ID          uint         `gorm:"primaryKey;autoIncrement:true" json:"id"`
	ElementUUID string       `gorm:"unique;not null" json:"element_uuid"`
	UserID      uint         `gorm:"not null" json:"user_id"`
	User        User         `gorm:"not null"`
	PageID      uint         `gorm:"foreignKey:ID" json:"page_id"`
	Page        Page         `gorm:"foreignKey:ID"`
	Type        string       `gorm:"not null;type:varchar(20)" json:"type"`
	Content     pgtype.JSONB `gorm:"type:jsonb;default: '{}'" json:"content"`
	Etc         pgtype.JSONB `gorm:"type:jsonb;default: '{}'" json:"etc"`
	Size        string       `gorm:"default:null;check:Size IN (null, 'small', 'medium', 'large')" json:"size"`
}

type Page struct {
	gorm.Model
	ID               uint         `gorm:"primaryKey;autoIncrement:true" json:"id"`
	UserID           uint         `gorm:"not null" json:"user_id"`
	User             User         `gorm:"foreignKey:ID"`
	PageUUID         string       `gorm:"unique;not null;type:text;" json:"page_uuid"`
	PageName         string       `gorm:"not null;default:''" json:"page_name"`
	IsRoot           bool         `gorm:"not null;default:true" json:"is_root"`
	ElementPositions pgtype.JSONB `gorm:"type:jsonb" json:"element_positions"`
	ParentPageUUID   string       `gorm:"default:null;type:text;" json:"parent_page_uuid"`
	PublicPage       bool         `gorm:"not null;default:false" json:"public_page"`
	PageUUIDURL      string       `gorm:"default:null" json:"page_uuid_url"`
	IsFavourite      bool         `gorm:"not null;default:false" json:"is_favourite"`
	LastUpdatedAt    time.Time    `gorm:"default:null" json:"last_updated_at"`
	Etc              pgtype.JSONB `gorm:"type:jsonb;default: '{}'" json:"etc"`
	ViewCount        uint         `gorm:"not null;default:0" json:"view_count"`
	DateViewCount    pgtype.JSONB `gorm:"type:jsonb;default: '{}'" json:"date_view_count"`
}

type User struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement:true"`
	GoogleID    string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	Name        string `gorm:"not null"`
	Picture     string `gorm:"not null"`
	Credentials []byte `gorm:"type:jsonb;default: '{}'"`
	Status      string `gorm:"not null;default:'freemium';check:Status IN ('freemium', 'premium', 'enterprise')" json:"status"`
}
