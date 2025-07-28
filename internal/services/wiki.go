package services

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

func GetWikis(ctx *gin.Context) {
	var filter dto.WikiFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid filter parameters", err.Error()))
		return
	}

	wikis, total, err := GetWikisInternal(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve wikis", err.Error()))
		return
	}

	response := utils.PaginatedResponse{
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
		Data:    wikis,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Wikis", response))
}

func GetWikiBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Slug is required", ""))
		return
	}

	wiki, err := GetWikiBySlugInternal(slug)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Wiki not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Wiki", wiki))
}

func GetWikiCategories(ctx *gin.Context) {
	categories, err := GetWikiCategoriesInternal()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve categories", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Categories", categories))
}

func CreateWiki(ctx *gin.Context) {
	var req dto.CreateWikiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	wiki, err := CreateWikiInternal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to create wiki", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Wiki created successfully", wiki))
}

func UpdateWiki(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid wiki ID", err.Error()))
		return
	}

	var req dto.UpdateWikiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	wiki, err := UpdateWikiInternal(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to update wiki", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Wiki updated successfully", wiki))
}

func DeleteWiki(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid wiki ID", err.Error()))
		return
	}

	if err := DeleteWikiInternal(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to delete wiki", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Wiki deleted successfully", nil))
}

func GetWikisInternal(filter dto.WikiFilter) ([]dto.Wiki, int, error) {
	var wikis []dto.Wiki
	var total int

	whereClauses := []string{"deleted = 0"}
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if filter.Category != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, *filter.Category)
		argIndex++
	}
	if filter.Subcategory != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("subcategory = $%d", argIndex))
		args = append(args, *filter.Subcategory)
		argIndex++
	}
	if filter.ContentType != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("content_type = $%d", argIndex))
		args = append(args, *filter.ContentType)
		argIndex++
	}
	if filter.IsFeatured != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_featured = $%d", argIndex))
		args = append(args, *filter.IsFeatured)
		argIndex++
	}
	if filter.IsPublic != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_public = $%d", argIndex))
		args = append(args, *filter.IsPublic)
		argIndex++
	}
	if filter.RequiresAuth != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("requires_auth = $%d", argIndex))
		args = append(args, *filter.RequiresAuth)
		argIndex++
	}
	if filter.DifficultyLevel != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("difficulty_level = $%d", argIndex))
		args = append(args, *filter.DifficultyLevel)
		argIndex++
	}
	if filter.Tags != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("tags ILIKE $%d", argIndex))
		args = append(args, "%"+*filter.Tags+"%")
		argIndex++
	}
	if filter.Active != nil {
		if *filter.Active {
			whereClauses = append(whereClauses, "active = 1")
		} else {
			whereClauses = append(whereClauses, "active = 0")
		}
	}

	if filter.Search != nil {
		searchFields := []string{}
		if filter.Language != nil {
			switch *filter.Language {
			case "en":
				searchFields = append(searchFields, fmt.Sprintf("(title_en ILIKE $%d OR COALESCE(description_en, '') ILIKE $%d OR COALESCE(text_md_en, '') ILIKE $%d)", argIndex, argIndex, argIndex))
			case "ru":
				searchFields = append(searchFields, fmt.Sprintf("(title_ru ILIKE $%d OR COALESCE(description_ru, '') ILIKE $%d OR COALESCE(text_md_ru, '') ILIKE $%d)", argIndex, argIndex, argIndex))
			case "tk":
				searchFields = append(searchFields, fmt.Sprintf("(title_tk ILIKE $%d OR COALESCE(description_tk, '') ILIKE $%d OR COALESCE(text_md_tk, '') ILIKE $%d)", argIndex, argIndex, argIndex))
			}
		} else {
			searchFields = append(searchFields, fmt.Sprintf("(title_en ILIKE $%d OR title_ru ILIKE $%d OR title_tk ILIKE $%d OR COALESCE(description_en, '') ILIKE $%d OR COALESCE(description_ru, '') ILIKE $%d OR COALESCE(description_tk, '') ILIKE $%d OR COALESCE(tags, '') ILIKE $%d)", argIndex, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex))
		}
		whereClauses = append(whereClauses, "("+strings.Join(searchFields, " OR ")+")")
		args = append(args, "%"+*filter.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tbl_wiki WHERE %s", whereClause)
	err := database.DB.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM tbl_wiki
		WHERE %s
		ORDER BY priority DESC, created_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1)

	args = append(args, filter.PerPage, (filter.Page-1)*filter.PerPage)
	err = pgxscan.Select(context.Background(), database.DB, &wikis, query, args...)
	return wikis, total, err
}

func GetWikiBySlugInternal(slug string) (dto.Wiki, error) {
	var wiki dto.Wiki
	query := `
		SELECT *
		FROM tbl_wiki
		WHERE slug = $1 AND deleted = 0 AND active = 1`

	err := pgxscan.Get(context.Background(), database.DB, &wiki, query, slug)
	if err == nil {
		// Increment view count
		_, _ = database.DB.Exec(context.Background(), "UPDATE tbl_wiki SET view_count = view_count + 1 WHERE slug = $1", slug)
	}
	return wiki, err
}

func GetWikiCategoriesInternal() ([]dto.WikiCategoryResponse, error) {
	var categories []dto.WikiCategoryResponse
	query := `
		SELECT category, COALESCE(subcategory, '') as subcategory, COUNT(*) as count
		FROM tbl_wiki
		WHERE deleted = 0 AND active = 1 AND is_public = true
		GROUP BY category, subcategory
		ORDER BY category, subcategory`

	err := pgxscan.Select(context.Background(), database.DB, &categories, query)
	return categories, err
}

func CreateWikiInternal(req dto.CreateWikiRequest) (dto.Wiki, error) {
	var wiki dto.Wiki

	priority := 0
	if req.Priority != nil {
		priority = *req.Priority
	}
	isFeatured := false
	if req.IsFeatured != nil {
		isFeatured = *req.IsFeatured
	}
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}
	requiresAuth := false
	if req.RequiresAuth != nil {
		requiresAuth = *req.RequiresAuth
	}

	var fileUrls [5]*string
	var fileInfoEn, fileInfoRu, fileInfoTk [5]*string
	var videoUrls [5]*string
	var videoInfoEn, videoInfoRu, videoInfoTk [5]*string

	for i, file := range req.Files {
		if i >= 5 {
			break
		}
		fileUrls[i] = file.Url
		fileInfoEn[i] = file.InfoEn
		fileInfoRu[i] = file.InfoRu
		fileInfoTk[i] = file.InfoTk
	}

	for i, video := range req.Videos {
		if i >= 5 {
			break
		}
		videoUrls[i] = video.Url
		videoInfoEn[i] = video.InfoEn
		videoInfoRu[i] = video.InfoRu
		videoInfoTk[i] = video.InfoTk
	}

	query := `
		INSERT INTO tbl_wiki (
			title_en, title_ru, title_tk, description_en, description_ru, description_tk, description_type,
			text_md_en, text_md_ru, text_md_tk, text_rich_en, text_rich_ru, text_rich_tk,
			file_url_1, file_url_2, file_url_3, file_url_4, file_url_5,
			file_info_1_en, file_info_1_ru, file_info_1_tk, file_info_2_en, file_info_2_ru, file_info_2_tk,
			file_info_3_en, file_info_3_ru, file_info_3_tk, file_info_4_en, file_info_4_ru, file_info_4_tk,
			file_info_5_en, file_info_5_ru, file_info_5_tk,
			video_url_1, video_url_2, video_url_3, video_url_4, video_url_5,
			video_info_1_en, video_info_1_ru, video_info_1_tk, video_info_2_en, video_info_2_ru, video_info_2_tk,
			video_info_3_en, video_info_3_ru, video_info_3_tk, video_info_4_en, video_info_4_ru, video_info_4_tk,
			video_info_5_en, video_info_5_ru, video_info_5_tk,
			category, subcategory, tags, slug, meta_keywords_en, meta_keywords_ru, meta_keywords_tk,
			priority, is_featured, is_public, requires_auth, content_type, difficulty_level, estimated_read_time
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31,
			$32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49,
			$50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67
		) RETURNING id, uuid, title_en, title_ru, title_tk, description_en, description_ru, description_tk, 
		           description_type, text_md_en, text_md_ru, text_md_tk, text_rich_en, text_rich_ru, text_rich_tk,
		           file_url_1, file_url_2, file_url_3, file_url_4, file_url_5,
		           file_info_1_en, file_info_1_ru, file_info_1_tk, file_info_2_en, file_info_2_ru, file_info_2_tk,
		           file_info_3_en, file_info_3_ru, file_info_3_tk, file_info_4_en, file_info_4_ru, file_info_4_tk,
		           file_info_5_en, file_info_5_ru, file_info_5_tk,
		           video_url_1, video_url_2, video_url_3, video_url_4, video_url_5,
		           video_info_1_en, video_info_1_ru, video_info_1_tk, video_info_2_en, video_info_2_ru, video_info_2_tk,
		           video_info_3_en, video_info_3_ru, video_info_3_tk, video_info_4_en, video_info_4_ru, video_info_4_tk,
		           video_info_5_en, video_info_5_ru, video_info_5_tk,
		           category, subcategory, tags, version, slug, meta_keywords_en, meta_keywords_ru, meta_keywords_tk,
		           priority, view_count, is_featured, is_public, requires_auth, content_type, difficulty_level,
		           estimated_read_time, created_at, updated_at`

	err := pgxscan.Get(context.Background(), database.DB, &wiki, query,
		req.TitleEn, req.TitleRu, req.TitleTk, req.DescriptionEn, req.DescriptionRu, req.DescriptionTk, req.DescriptionType,
		req.TextMdEn, req.TextMdRu, req.TextMdTk, req.TextRichEn, req.TextRichRu, req.TextRichTk,
		fileUrls[0], fileUrls[1], fileUrls[2], fileUrls[3], fileUrls[4],
		fileInfoEn[0], fileInfoRu[0], fileInfoTk[0], fileInfoEn[1], fileInfoRu[1], fileInfoTk[1],
		fileInfoEn[2], fileInfoRu[2], fileInfoTk[2], fileInfoEn[3], fileInfoRu[3], fileInfoTk[3],
		fileInfoEn[4], fileInfoRu[4], fileInfoTk[4],
		videoUrls[0], videoUrls[1], videoUrls[2], videoUrls[3], videoUrls[4],
		videoInfoEn[0], videoInfoRu[0], videoInfoTk[0], videoInfoEn[1], videoInfoRu[1], videoInfoTk[1],
		videoInfoEn[2], videoInfoRu[2], videoInfoTk[2], videoInfoEn[3], videoInfoRu[3], videoInfoTk[3],
		videoInfoEn[4], videoInfoRu[4], videoInfoTk[4],
		req.Category, req.Subcategory, req.Tags, req.Slug, req.MetaKeywordsEn, req.MetaKeywordsRu, req.MetaKeywordsTk,
		priority, isFeatured, isPublic, requiresAuth, req.ContentType, req.DifficultyLevel, req.EstimatedReadTime)
	return wiki, err
}

func UpdateWikiInternal(id int, req dto.UpdateWikiRequest) (dto.Wiki, error) {
	var wiki dto.Wiki
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if req.TitleEn != nil {
		setParts = append(setParts, fmt.Sprintf("title_en = $%d", argIndex))
		args = append(args, *req.TitleEn)
		argIndex++
	}
	if req.TitleRu != nil {
		setParts = append(setParts, fmt.Sprintf("title_ru = $%d", argIndex))
		args = append(args, *req.TitleRu)
		argIndex++
	}
	if req.TitleTk != nil {
		setParts = append(setParts, fmt.Sprintf("title_tk = $%d", argIndex))
		args = append(args, *req.TitleTk)
		argIndex++
	}
	if req.DescriptionEn != nil {
		setParts = append(setParts, fmt.Sprintf("description_en = $%d", argIndex))
		args = append(args, *req.DescriptionEn)
		argIndex++
	}
	if req.DescriptionRu != nil {
		setParts = append(setParts, fmt.Sprintf("description_ru = $%d", argIndex))
		args = append(args, *req.DescriptionRu)
		argIndex++
	}
	if req.DescriptionTk != nil {
		setParts = append(setParts, fmt.Sprintf("description_tk = $%d", argIndex))
		args = append(args, *req.DescriptionTk)
		argIndex++
	}
	if req.DescriptionType != nil {
		setParts = append(setParts, fmt.Sprintf("description_type = $%d", argIndex))
		args = append(args, *req.DescriptionType)
		argIndex++
	}
	if req.TextMdEn != nil {
		setParts = append(setParts, fmt.Sprintf("text_md_en = $%d", argIndex))
		args = append(args, *req.TextMdEn)
		argIndex++
	}
	if req.TextMdRu != nil {
		setParts = append(setParts, fmt.Sprintf("text_md_ru = $%d", argIndex))
		args = append(args, *req.TextMdRu)
		argIndex++
	}
	if req.TextMdTk != nil {
		setParts = append(setParts, fmt.Sprintf("text_md_tk = $%d", argIndex))
		args = append(args, *req.TextMdTk)
		argIndex++
	}
	if req.TextRichEn != nil {
		setParts = append(setParts, fmt.Sprintf("text_rich_en = $%d", argIndex))
		args = append(args, *req.TextRichEn)
		argIndex++
	}
	if req.TextRichRu != nil {
		setParts = append(setParts, fmt.Sprintf("text_rich_ru = $%d", argIndex))
		args = append(args, *req.TextRichRu)
		argIndex++
	}
	if req.TextRichTk != nil {
		setParts = append(setParts, fmt.Sprintf("text_rich_tk = $%d", argIndex))
		args = append(args, *req.TextRichTk)
		argIndex++
	}

	// Metadata fields
	if req.Category != nil {
		setParts = append(setParts, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, *req.Category)
		argIndex++
	}
	if req.Subcategory != nil {
		setParts = append(setParts, fmt.Sprintf("subcategory = $%d", argIndex))
		args = append(args, *req.Subcategory)
		argIndex++
	}
	if req.Tags != nil {
		setParts = append(setParts, fmt.Sprintf("tags = $%d", argIndex))
		args = append(args, *req.Tags)
		argIndex++
	}
	if req.Slug != nil {
		setParts = append(setParts, fmt.Sprintf("slug = $%d", argIndex))
		args = append(args, *req.Slug)
		argIndex++
	}
	if req.MetaKeywordsEn != nil {
		setParts = append(setParts, fmt.Sprintf("meta_keywords_en = $%d", argIndex))
		args = append(args, *req.MetaKeywordsEn)
		argIndex++
	}
	if req.MetaKeywordsRu != nil {
		setParts = append(setParts, fmt.Sprintf("meta_keywords_ru = $%d", argIndex))
		args = append(args, *req.MetaKeywordsRu)
		argIndex++
	}
	if req.MetaKeywordsTk != nil {
		setParts = append(setParts, fmt.Sprintf("meta_keywords_tk = $%d", argIndex))
		args = append(args, *req.MetaKeywordsTk)
		argIndex++
	}
	if req.Priority != nil {
		setParts = append(setParts, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, *req.Priority)
		argIndex++
	}
	if req.IsFeatured != nil {
		setParts = append(setParts, fmt.Sprintf("is_featured = $%d", argIndex))
		args = append(args, *req.IsFeatured)
		argIndex++
	}
	if req.IsPublic != nil {
		setParts = append(setParts, fmt.Sprintf("is_public = $%d", argIndex))
		args = append(args, *req.IsPublic)
		argIndex++
	}
	if req.RequiresAuth != nil {
		setParts = append(setParts, fmt.Sprintf("requires_auth = $%d", argIndex))
		args = append(args, *req.RequiresAuth)
		argIndex++
	}
	if req.ContentType != nil {
		setParts = append(setParts, fmt.Sprintf("content_type = $%d", argIndex))
		args = append(args, *req.ContentType)
		argIndex++
	}
	if req.DifficultyLevel != nil {
		setParts = append(setParts, fmt.Sprintf("difficulty_level = $%d", argIndex))
		args = append(args, *req.DifficultyLevel)
		argIndex++
	}
	if req.EstimatedReadTime != nil {
		setParts = append(setParts, fmt.Sprintf("estimated_read_time = $%d", argIndex))
		args = append(args, *req.EstimatedReadTime)
		argIndex++
	}

	if req.Files != nil {
		fileFields := []string{"file_url_1", "file_url_2", "file_url_3", "file_url_4", "file_url_5"}
		fileInfoEnFields := []string{"file_info_1_en", "file_info_2_en", "file_info_3_en", "file_info_4_en", "file_info_5_en"}
		fileInfoRuFields := []string{"file_info_1_ru", "file_info_2_ru", "file_info_3_ru", "file_info_4_ru", "file_info_5_ru"}
		fileInfoTkFields := []string{"file_info_1_tk", "file_info_2_tk", "file_info_3_tk", "file_info_4_tk", "file_info_5_tk"}

		for i := 0; i < 5; i++ {
			var fileUrl, fileInfoEn, fileInfoRu, fileInfoTk *string
			if i < len(req.Files) {
				fileUrl = req.Files[i].Url
				fileInfoEn = req.Files[i].InfoEn
				fileInfoRu = req.Files[i].InfoRu
				fileInfoTk = req.Files[i].InfoTk
			}

			setParts = append(setParts, fmt.Sprintf("%s = $%d", fileFields[i], argIndex))
			args = append(args, fileUrl)
			argIndex++

			setParts = append(setParts, fmt.Sprintf("%s = $%d", fileInfoEnFields[i], argIndex))
			args = append(args, fileInfoEn)
			argIndex++

			setParts = append(setParts, fmt.Sprintf("%s = $%d", fileInfoRuFields[i], argIndex))
			args = append(args, fileInfoRu)
			argIndex++

			setParts = append(setParts, fmt.Sprintf("%s = $%d", fileInfoTkFields[i], argIndex))
			args = append(args, fileInfoTk)
			argIndex++
		}
	}

	if req.Videos != nil {
		videoFields := []string{"video_url_1", "video_url_2", "video_url_3", "video_url_4", "video_url_5"}
		videoInfoEnFields := []string{"video_info_1_en", "video_info_2_en", "video_info_3_en", "video_info_4_en", "video_info_5_en"}
		videoInfoRuFields := []string{"video_info_1_ru", "video_info_2_ru", "video_info_3_ru", "video_info_4_ru", "video_info_5_ru"}
		videoInfoTkFields := []string{"video_info_1_tk", "video_info_2_tk", "video_info_3_tk", "video_info_4_tk", "video_info_5_tk"}

		for i := 0; i < 5; i++ {
			var videoUrl, videoInfoEn, videoInfoRu, videoInfoTk *string
			if i < len(req.Videos) {
				videoUrl = req.Videos[i].Url
				videoInfoEn = req.Videos[i].InfoEn
				videoInfoRu = req.Videos[i].InfoRu
				videoInfoTk = req.Videos[i].InfoTk
			}

			setParts = append(setParts, fmt.Sprintf("%s = $%d", videoFields[i], argIndex))
			args = append(args, videoUrl)
			argIndex++

			setParts = append(setParts, fmt.Sprintf("%s = $%d", videoInfoEnFields[i], argIndex))
			args = append(args, videoInfoEn)
			argIndex++

			setParts = append(setParts, fmt.Sprintf("%s = $%d", videoInfoRuFields[i], argIndex))
			args = append(args, videoInfoRu)
			argIndex++

			setParts = append(setParts, fmt.Sprintf("%s = $%d", videoInfoTkFields[i], argIndex))
			args = append(args, videoInfoTk)
			argIndex++
		}
	}

	args = append(args, id)
	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf(`
		UPDATE tbl_wiki 
		SET %s 
		WHERE id = $%d AND deleted = 0 
		RETURNING id, uuid, title_en, title_ru, title_tk, description_en, description_ru, description_tk, 
		          description_type, text_md_en, text_md_ru, text_md_tk, text_rich_en, text_rich_ru, text_rich_tk,
		          file_url_1, file_url_2, file_url_3, file_url_4, file_url_5,
		          file_info_1_en, file_info_1_ru, file_info_1_tk, file_info_2_en, file_info_2_ru, file_info_2_tk,
		          file_info_3_en, file_info_3_ru, file_info_3_tk, file_info_4_en, file_info_4_ru, file_info_4_tk,
		          file_info_5_en, file_info_5_ru, file_info_5_tk,
		          video_url_1, video_url_2, video_url_3, video_url_4, video_url_5,
		          video_info_1_en, video_info_1_ru, video_info_1_tk, video_info_2_en, video_info_2_ru, video_info_2_tk,
		          video_info_3_en, video_info_3_ru, video_info_3_tk, video_info_4_en, video_info_4_ru, video_info_4_tk,
		          video_info_5_en, video_info_5_ru, video_info_5_tk,
		          category, subcategory, tags, version, slug, meta_keywords_en, meta_keywords_ru, meta_keywords_tk,
		          priority, view_count, is_featured, is_public, requires_auth, content_type, difficulty_level,
		          estimated_read_time, created_at, updated_at`, setClause, argIndex)

	err := pgxscan.Get(context.Background(), database.DB, &wiki, query, args...)
	return wiki, err
}

func DeleteWikiInternal(id int) error {
	query := `UPDATE tbl_wiki SET deleted = 1, updated_at = NOW() WHERE id = $1 AND deleted = 0`
	_, err := database.DB.Exec(context.Background(), query, id)
	return err
}
