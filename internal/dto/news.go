package dto

import (
	"github.com/google/uuid"
	"time"
)

type News struct {
	ID               uuid.UUID  `json:"id"`
	ExternalID       *string    `json:"external_id"`
	Slug             string     `json:"slug"`
	Title            string     `json:"title"`
	Subtitle         *string    `json:"subtitle"`
	Excerpt          *string    `json:"excerpt"`
	Content          string     `json:"content"`
	ContentPlain     *string    `json:"content_plain"`
	FeaturedImageURL *string    `json:"featured_image_url"`
	AuthorName       string     `json:"author_name"`
	CategoryPrimary  string     `json:"category_primary"`
	ContentType      string     `json:"content_type"`
	Status           string     `json:"status"`
	Priority         string     `json:"priority"`
	PublishedAt      *time.Time `json:"published_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type CreateNewsRequest struct {
	Title            string     `json:"title" binding:"max=500"`
	Subtitle         *string    `json:"subtitle" binding:"omitempty,max=1000"`
	Excerpt          *string    `json:"excerpt" binding:"omitempty"`
	Content          string     `json:"content" binding:""`
	ContentPlain     *string    `json:"content_plain" binding:"omitempty"`
	FeaturedImageURL *string    `json:"featured_image_url" binding:"omitempty,max=500"`
	AuthorName       string     `json:"author_name" binding:"max=200"`
	CategoryPrimary  string     `json:"category_primary" binding:"max=100"`
	ContentType      string     `json:"content_type" binding:"oneof=article breaking_news opinion analysis interview review"`
	Status           string     `json:"status" binding:"oneof=draft published archived deleted"`
	Priority         string     `json:"priority" binding:"oneof=low medium high urgent"`
	PublishedAt      *time.Time `json:"published_at" binding:"omitempty"`
}

type UpdateNewsRequest struct {
	Title            *string    `json:"title" binding:"omitempty,max=500"`
	Subtitle         *string    `json:"subtitle" binding:"omitempty,max=1000"`
	Excerpt          *string    `json:"excerpt" binding:"omitempty"`
	Content          *string    `json:"content" binding:"omitempty"`
	ContentPlain     *string    `json:"content_plain" binding:"omitempty"`
	FeaturedImageURL *string    `json:"featured_image_url" binding:"omitempty,max=500"`
	AuthorName       *string    `json:"author_name" binding:"omitempty,max=200"`
	CategoryPrimary  *string    `json:"category_primary" binding:"omitempty,max=100"`
	ContentType      *string    `json:"content_type" binding:"omitempty,oneof=article breaking_news opinion analysis interview review"`
	Status           *string    `json:"status" binding:"omitempty,oneof=draft published archived deleted"`
	Priority         *string    `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	PublishedAt      *time.Time `json:"published_at" binding:"omitempty"`
}

type NewsFilter struct {
	Category    *string `form:"category"`
	Status      *string `form:"status" binding:"omitempty,oneof=draft published archived deleted"`
	Priority    *string `form:"priority" binding:"omitempty,oneof=low medium high urgent"`
	ContentType *string `form:"content_type" binding:"omitempty,oneof=article breaking_news opinion analysis interview review"`
	Search      *string `form:"search"`
	Page        int     `form:"page,default=1" binding:"min=1"`
	PerPage     int     `form:"per_page,default=10" binding:"min=1,max=100"`
}
