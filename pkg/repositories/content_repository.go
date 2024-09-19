package repositories

import (
	"context"
	"texApi/pkg/schemas/request"
	"texApi/pkg/schemas/response"

	"github.com/jackc/pgx/v5"
)

// ContentRepository handles the database operations.
type ContentRepository struct {
	DB *pgx.Conn
}

// NewContentRepository creates a new ContentRepository.
func NewContentRepository(db *pgx.Conn) *ContentRepository {
	return &ContentRepository{DB: db}
}

// GetAll retrieves all content.
func (repo *ContentRepository) GetAll(ctx context.Context) ([]response.ContentResponse, error) {
	rows, err := repo.DB.Query(ctx, "SELECT * FROM content WHERE deleted = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []response.ContentResponse
	for rows.Next() {
		var content response.ContentResponse
		err := rows.Scan(&content.ID, &content.UUID, &content.LangID, &content.ContentTypeID,
			&content.Title, &content.Subtitle, &content.Description,
			&content.ImageURL, &content.VideoURL, &content.Step,
			&content.CreatedAt, &content.UpdatedAt, &content.Deleted)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

// GetByID retrieves content by ID.
func (repo *ContentRepository) GetByID(ctx context.Context, id int) (response.ContentResponse, error) {
	var content response.ContentResponse
	err := repo.DB.QueryRow(ctx, "SELECT * FROM content WHERE id = $1 AND deleted = 0", id).
		Scan(&content.ID, &content.UUID, &content.LangID, &content.ContentTypeID,
			&content.Title, &content.Subtitle, &content.Description,
			&content.ImageURL, &content.VideoURL, &content.Step,
			&content.CreatedAt, &content.UpdatedAt, &content.Deleted)
	return content, err
}

// GetByTitle retrieves content by title.
func (repo *ContentRepository) GetByTitle(ctx context.Context, title string) ([]response.ContentResponse, error) {
	rows, err := repo.DB.Query(ctx, "SELECT * FROM content WHERE title ILIKE $1 AND deleted = 0", "%"+title+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []response.ContentResponse
	for rows.Next() {
		var content response.ContentResponse
		err := rows.Scan(&content.ID, &content.UUID, &content.LangID, &content.ContentTypeID,
			&content.Title, &content.Subtitle, &content.Description,
			&content.ImageURL, &content.VideoURL, &content.Step,
			&content.CreatedAt, &content.UpdatedAt, &content.Deleted)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

// GetByUUID retrieves content by UUID.
func (repo *ContentRepository) GetByUUID(ctx context.Context, uuid string) (response.ContentResponse, error) {
	var content response.ContentResponse
	err := repo.DB.QueryRow(ctx, "SELECT * FROM content WHERE uuid = $1 AND deleted = 0", uuid).
		Scan(&content.ID, &content.UUID, &content.LangID, &content.ContentTypeID,
			&content.Title, &content.Subtitle, &content.Description,
			&content.ImageURL, &content.VideoURL, &content.Step,
			&content.CreatedAt, &content.UpdatedAt, &content.Deleted)
	return content, err
}

// GetByContentTypeID retrieves content by content type ID.
func (repo *ContentRepository) GetByContentTypeID(ctx context.Context, contentTypeID int) ([]response.ContentResponse, error) {
	rows, err := repo.DB.Query(ctx, "SELECT * FROM content WHERE content_type_id = $1 AND deleted = 0", contentTypeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []response.ContentResponse
	for rows.Next() {
		var content response.ContentResponse
		err := rows.Scan(&content.ID, &content.UUID, &content.LangID, &content.ContentTypeID,
			&content.Title, &content.Subtitle, &content.Description,
			&content.ImageURL, &content.VideoURL, &content.Step,
			&content.CreatedAt, &content.UpdatedAt, &content.Deleted)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

// Create inserts new content into the database.
func (repo *ContentRepository) Create(ctx context.Context, content request.ContentRequest) (response.ContentResponse, error) {
	var newContent response.ContentResponse
	sql := `INSERT INTO content (uuid, lang_id, content_type_id, title, subtitle, description, image_url, video_url, step, created_at, updated_at, deleted)
            VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), 0)
            RETURNING id, uuid, lang_id, content_type_id, title, subtitle, description, image_url, video_url, step, created_at, updated_at, deleted`

	err := repo.DB.QueryRow(ctx, sql, content.LangID, content.ContentTypeID, content.Title, content.Subtitle, content.Description, content.ImageURL, content.VideoURL, content.Step).
		Scan(&newContent.ID, &newContent.UUID, &newContent.LangID, &newContent.ContentTypeID,
			&newContent.Title, &newContent.Subtitle, &newContent.Description,
			&newContent.ImageURL, &newContent.VideoURL, &newContent.Step,
			&newContent.CreatedAt, &newContent.UpdatedAt, &newContent.Deleted)

	return newContent, err
}

// Update modifies an existing content record in the database.
func (repo *ContentRepository) Update(ctx context.Context, id int, content request.ContentRequest) (response.ContentResponse, error) {
	var updatedContent response.ContentResponse
	sql := `UPDATE content 
            SET lang_id = $1, content_type_id = $2, title = $3, subtitle = $4, description = $5, image_url = $6, video_url = $7, step = $8, updated_at = NOW()
            WHERE id = $9 AND deleted = 0 
            RETURNING id, uuid, lang_id, content_type_id, title, subtitle, description, image_url, video_url, step, created_at, updated_at, deleted`

	err := repo.DB.QueryRow(ctx, sql, content.LangID, content.ContentTypeID, content.Title, content.Subtitle, content.Description, content.ImageURL, content.VideoURL, content.Step, id).
		Scan(&updatedContent.ID, &updatedContent.UUID, &updatedContent.LangID, &updatedContent.ContentTypeID,
			&updatedContent.Title, &updatedContent.Subtitle, &updatedContent.Description,
			&updatedContent.ImageURL, &updatedContent.VideoURL, &updatedContent.Step,
			&updatedContent.CreatedAt, &updatedContent.UpdatedAt, &updatedContent.Deleted)

	return updatedContent, err
}

// Delete sets the deleted flag to 1 for a content record.
func (repo *ContentRepository) Delete(ctx context.Context, id int) error {
	_, err := repo.DB.Exec(ctx, "UPDATE content SET deleted = 1 WHERE id = $1", id)
	return err
}
