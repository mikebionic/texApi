package services

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

func GetNews(ctx *gin.Context) {
	var filter dto.NewsFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid filter parameters", err.Error()))
		return
	}

	news, total, err := GetNewsList(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve news", err.Error()))
		return
	}

	response := utils.PaginatedResponse{
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
		Data:    news,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("News", response))
}

func GetNewsByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid news ID", err.Error()))
		return
	}

	news, err := GetNewsByIDInternal(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("News not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("News retrieved", news))
}

func CreateNews(ctx *gin.Context) {
	var req dto.CreateNewsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	news, err := CreateNewsInternal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to create news", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("News created successfully", news))
}

func UpdateNews(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid news ID", err.Error()))
		return
	}

	var req dto.UpdateNewsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	news, err := UpdateNewsInternal(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to update news", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("News updated successfully", news))
}

func DeleteNews(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid news ID", err.Error()))
		return
	}

	if err := DeleteNewsInternal(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to delete news", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("News deleted successfully", nil))
}

func GetNewsList(filter dto.NewsFilter) ([]dto.News, int, error) {
	var news []dto.News
	var total int

	whereClauses := []string{"deleted = 0"}
	args := []interface{}{}
	argIndex := 1

	if filter.Category != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("category_primary = $%d", argIndex))
		args = append(args, *filter.Category)
		argIndex++
	}
	if filter.Status != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}
	if filter.Priority != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, *filter.Priority)
		argIndex++
	}
	if filter.ContentType != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("content_type = $%d", argIndex))
		args = append(args, *filter.ContentType)
		argIndex++
	}
	if filter.Search != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("search_vector @@ to_tsquery('english', $%d)", argIndex))
		args = append(args, *filter.Search)
		argIndex++
	}

	whereClause := strings.Join(whereClauses, " AND ")
	if whereClause == "" {
		whereClause = "1=1"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tbl_article WHERE %s", whereClause)
	err := database.DB.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, external_id, slug, title, subtitle, excerpt, content, content_plain,
		       featured_image_url, author_name, category_primary, content_type, status,
		       priority, published_at, created_at, updated_at
		FROM tbl_article
		WHERE %s
		ORDER BY published_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1)

	args = append(args, filter.PerPage, (filter.Page-1)*filter.PerPage)
	err = pgxscan.Select(context.Background(), database.DB, &news, query, args...)
	return news, total, err
}

func GetNewsByIDInternal(id uuid.UUID) (dto.News, error) {
	var news dto.News
	query := `
		SELECT id, external_id, slug, title, subtitle, excerpt, content, content_plain,
		       featured_image_url, author_name, category_primary, content_type, status,
		       priority, published_at, created_at, updated_at
		FROM tbl_article
		WHERE id = $1 AND deleted = 0`
	err := pgxscan.Get(context.Background(), database.DB, &news, query, id)
	return news, err
}

func CreateNewsInternal(req dto.CreateNewsRequest) (dto.News, error) {
	var news dto.News
	query := `
		INSERT INTO tbl_article (
			title, subtitle, excerpt, content, content_plain, featured_image_url,
			author_name, category_primary, content_type, status, priority, published_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) RETURNING id, external_id, slug, title, subtitle, excerpt, content, content_plain,
		           featured_image_url, author_name, category_primary, content_type, status,
		           priority, published_at, created_at, updated_at`

	err := pgxscan.Get(context.Background(), database.DB, &news, query,
		req.Title, req.Subtitle, req.Excerpt, req.Content, req.ContentPlain,
		req.FeaturedImageURL, req.AuthorName, req.CategoryPrimary, req.ContentType,
		req.Status, req.Priority, req.PublishedAt)
	return news, err
}

func UpdateNewsInternal(id uuid.UUID, req dto.UpdateNewsRequest) (dto.News, error) {
	var news dto.News
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if req.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *req.Title)
		argIndex++
	}
	if req.Subtitle != nil {
		setParts = append(setParts, fmt.Sprintf("subtitle = $%d", argIndex))
		args = append(args, *req.Subtitle)
		argIndex++
	}
	if req.Excerpt != nil {
		setParts = append(setParts, fmt.Sprintf("excerpt = $%d", argIndex))
		args = append(args, *req.Excerpt)
		argIndex++
	}
	if req.Content != nil {
		setParts = append(setParts, fmt.Sprintf("content = $%d", argIndex))
		args = append(args, *req.Content)
		argIndex++
	}
	if req.ContentPlain != nil {
		setParts = append(setParts, fmt.Sprintf("content_plain = $%d", argIndex))
		args = append(args, *req.ContentPlain)
		argIndex++
	}
	if req.FeaturedImageURL != nil {
		setParts = append(setParts, fmt.Sprintf("featured_image_url = $%d", argIndex))
		args = append(args, *req.FeaturedImageURL)
		argIndex++
	}
	if req.AuthorName != nil {
		setParts = append(setParts, fmt.Sprintf("author_name = $%d", argIndex))
		args = append(args, *req.AuthorName)
		argIndex++
	}
	if req.CategoryPrimary != nil {
		setParts = append(setParts, fmt.Sprintf("category_primary = $%d", argIndex))
		args = append(args, *req.CategoryPrimary)
		argIndex++
	}
	if req.ContentType != nil {
		setParts = append(setParts, fmt.Sprintf("content_type = $%d", argIndex))
		args = append(args, *req.ContentType)
		argIndex++
	}
	if req.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}
	if req.Priority != nil {
		setParts = append(setParts, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, *req.Priority)
		argIndex++
	}
	if req.PublishedAt != nil {
		setParts = append(setParts, fmt.Sprintf("published_at = $%d", argIndex))
		args = append(args, *req.PublishedAt)
		argIndex++
	}

	args = append(args, id)
	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf(`
		UPDATE tbl_article 
		SET %s 
		WHERE id = $%d AND deleted = 0 
		RETURNING id, external_id, slug, title, subtitle, excerpt, content, content_plain,
		          featured_image_url, author_name, category_primary, content_type, status,
		          priority, published_at, created_at, updated_at`, setClause, argIndex)

	err := pgxscan.Get(context.Background(), database.DB, &news, query, args...)
	return news, err
}

func DeleteNewsInternal(id uuid.UUID) error {
	query := `UPDATE tbl_article SET deleted = 1, updated_at = NOW() WHERE id = $1 AND deleted = 0`
	_, err := database.DB.Exec(context.Background(), query, id)
	return err
}
