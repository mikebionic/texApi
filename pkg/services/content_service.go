package services

import (
	"context"
	"texApi/pkg/repositories"
	"texApi/pkg/schemas/request"
	"texApi/pkg/schemas/response"
)

// ContentService handles the business logic.
type ContentService struct {
	repo *repositories.ContentRepository
}

// NewContentService creates a new ContentService.
func NewContentService(repo *repositories.ContentRepository) *ContentService {
	return &ContentService{repo: repo}
}

// GetAll retrieves all content.
func (s *ContentService) GetAll(ctx context.Context) ([]response.ContentResponse, error) {
	return s.repo.GetAll(ctx)
}

// GetByID retrieves content by ID.
func (s *ContentService) GetByID(ctx context.Context, id int) (response.ContentResponse, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByTitle retrieves content by title.
func (s *ContentService) GetByTitle(ctx context.Context, title string) ([]response.ContentResponse, error) {
	return s.repo.GetByTitle(ctx, title)
}

// GetByUUID retrieves content by UUID.
func (s *ContentService) GetByUUID(ctx context.Context, uuid string) (response.ContentResponse, error) {
	return s.repo.GetByUUID(ctx, uuid)
}

// GetByContentTypeID retrieves content by content type ID.
func (s *ContentService) GetByContentTypeID(ctx context.Context, contentTypeID int) ([]response.ContentResponse, error) {
	return s.repo.GetByContentTypeID(ctx, contentTypeID)
}

// Create a new content record.
func (s *ContentService) Create(ctx context.Context, content request.ContentRequest) (response.ContentResponse, error) {
	return s.repo.Create(ctx, content)
}

// Update an existing content record.
func (s *ContentService) Update(ctx context.Context, id int, content request.ContentRequest) (response.ContentResponse, error) {
	return s.repo.Update(ctx, id, content)
}

// Delete removes a content record by setting deleted to 1.
func (s *ContentService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
