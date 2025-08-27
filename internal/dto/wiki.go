package dto

import (
	"github.com/google/uuid"
	"time"
)

type WikiFileInfo struct {
	Url    *string `json:"url"`
	InfoEn *string `json:"info_en"`
	InfoRu *string `json:"info_ru"`
	InfoTk *string `json:"info_tk"`
}

type WikiVideoInfo struct {
	Url    *string `json:"url"`
	InfoEn *string `json:"info_en"`
	InfoRu *string `json:"info_ru"`
	InfoTk *string `json:"info_tk"`
}

type WikiContent struct {
	TitleEn         string  `json:"title_en" binding:"required"`
	TitleRu         string  `json:"title_ru"`
	TitleTk         string  `json:"title_tk"`
	DescriptionEn   *string `json:"description_en"`
	DescriptionRu   *string `json:"description_ru"`
	DescriptionTk   *string `json:"description_tk"`
	DescriptionType string  `json:"description_type" binding:"omitempty,oneof=plain html info"`
	TextMdEn        *string `json:"text_md_en"`
	TextMdRu        *string `json:"text_md_ru"`
	TextMdTk        *string `json:"text_md_tk"`
	TextRichEn      *string `json:"text_rich_en"`
	TextRichRu      *string `json:"text_rich_ru"`
	TextRichTk      *string `json:"text_rich_tk"`
}

type WikiMetadata struct {
	Category          string  `json:"category" binding:"required,oneof=docs wiki guides tutorials api faq changelog troubleshooting"`
	Subcategory       *string `json:"subcategory"`
	Tags              *string `json:"tags"`
	Slug              *string `json:"slug"`
	MetaKeywordsEn    *string `json:"meta_keywords_en"`
	MetaKeywordsRu    *string `json:"meta_keywords_ru"`
	MetaKeywordsTk    *string `json:"meta_keywords_tk"`
	Priority          *int    `json:"priority" binding:"omitempty,min=0"`
	IsFeatured        *bool   `json:"is_featured"`
	IsPublic          *bool   `json:"is_public"`
	RequiresAuth      *bool   `json:"requires_auth"`
	ContentType       string  `json:"content_type" binding:"required,oneof=article tutorial reference guide faq changelog"`
	DifficultyLevel   *string `json:"difficulty_level" binding:"omitempty,oneof=beginner intermediate advanced"`
	EstimatedReadTime *int    `json:"estimated_read_time" binding:"omitempty,min=0"`
}

type Wiki struct {
	ID              int       `json:"id" db:"id"`
	UUID            uuid.UUID `json:"uuid" db:"uuid"`
	TitleEn         string    `json:"title_en" db:"title_en"`
	TitleRu         string    `json:"title_ru" db:"title_ru"`
	TitleTk         string    `json:"title_tk" db:"title_tk"`
	DescriptionEn   *string   `json:"description_en" db:"description_en"`
	DescriptionRu   *string   `json:"description_ru" db:"description_ru"`
	DescriptionTk   *string   `json:"description_tk" db:"description_tk"`
	DescriptionType string    `json:"description_type" db:"description_type"`
	TextMdEn        *string   `json:"text_md_en" db:"text_md_en"`
	TextMdRu        *string   `json:"text_md_ru" db:"text_md_ru"`
	TextMdTk        *string   `json:"text_md_tk" db:"text_md_tk"`
	TextRichEn      *string   `json:"text_rich_en" db:"text_rich_en"`
	TextRichRu      *string   `json:"text_rich_ru" db:"text_rich_ru"`
	TextRichTk      *string   `json:"text_rich_tk" db:"text_rich_tk"`

	FileUrl1 *string `json:"file_url_1" db:"file_url_1"`
	FileUrl2 *string `json:"file_url_2" db:"file_url_2"`
	FileUrl3 *string `json:"file_url_3" db:"file_url_3"`
	FileUrl4 *string `json:"file_url_4" db:"file_url_4"`
	FileUrl5 *string `json:"file_url_5" db:"file_url_5"`

	FileInfo1En *string `json:"file_info_1_en" db:"file_info_1_en"`
	FileInfo2En *string `json:"file_info_2_en" db:"file_info_2_en"`
	FileInfo3En *string `json:"file_info_3_en" db:"file_info_3_en"`
	FileInfo4En *string `json:"file_info_4_en" db:"file_info_4_en"`
	FileInfo5En *string `json:"file_info_5_en" db:"file_info_5_en"`

	FileInfo1Ru *string `json:"file_info_1_ru" db:"file_info_1_ru"`
	FileInfo2Ru *string `json:"file_info_2_ru" db:"file_info_2_ru"`
	FileInfo3Ru *string `json:"file_info_3_ru" db:"file_info_3_ru"`
	FileInfo4Ru *string `json:"file_info_4_ru" db:"file_info_4_ru"`
	FileInfo5Ru *string `json:"file_info_5_ru" db:"file_info_5_ru"`

	FileInfo1Tk *string `json:"file_info_1_tk" db:"file_info_1_tk"`
	FileInfo2Tk *string `json:"file_info_2_tk" db:"file_info_2_tk"`
	FileInfo3Tk *string `json:"file_info_3_tk" db:"file_info_3_tk"`
	FileInfo4Tk *string `json:"file_info_4_tk" db:"file_info_4_tk"`
	FileInfo5Tk *string `json:"file_info_5_tk" db:"file_info_5_tk"`

	VideoUrl1 *string `json:"video_url_1" db:"video_url_1"`
	VideoUrl2 *string `json:"video_url_2" db:"video_url_2"`
	VideoUrl3 *string `json:"video_url_3" db:"video_url_3"`
	VideoUrl4 *string `json:"video_url_4" db:"video_url_4"`
	VideoUrl5 *string `json:"video_url_5" db:"video_url_5"`

	VideoInfo1En *string `json:"video_info_1_en" db:"video_info_1_en"`
	VideoInfo2En *string `json:"video_info_2_en" db:"video_info_2_en"`
	VideoInfo3En *string `json:"video_info_3_en" db:"video_info_3_en"`
	VideoInfo4En *string `json:"video_info_4_en" db:"video_info_4_en"`
	VideoInfo5En *string `json:"video_info_5_en" db:"video_info_5_en"`

	VideoInfo1Ru *string `json:"video_info_1_ru" db:"video_info_1_ru"`
	VideoInfo2Ru *string `json:"video_info_2_ru" db:"video_info_2_ru"`
	VideoInfo3Ru *string `json:"video_info_3_ru" db:"video_info_3_ru"`
	VideoInfo4Ru *string `json:"video_info_4_ru" db:"video_info_4_ru"`
	VideoInfo5Ru *string `json:"video_info_5_ru" db:"video_info_5_ru"`

	VideoInfo1Tk *string `json:"video_info_1_tk" db:"video_info_1_tk"`
	VideoInfo2Tk *string `json:"video_info_2_tk" db:"video_info_2_tk"`
	VideoInfo3Tk *string `json:"video_info_3_tk" db:"video_info_3_tk"`
	VideoInfo4Tk *string `json:"video_info_4_tk" db:"video_info_4_tk"`
	VideoInfo5Tk *string `json:"video_info_5_tk" db:"video_info_5_tk"`

	Category          string    `json:"category" db:"category"`
	Subcategory       *string   `json:"subcategory" db:"subcategory"`
	Tags              *string   `json:"tags" db:"tags"`
	Version           int       `json:"version" db:"version"`
	Slug              *string   `json:"slug" db:"slug"`
	MetaKeywordsEn    *string   `json:"meta_keywords_en" db:"meta_keywords_en"`
	MetaKeywordsRu    *string   `json:"meta_keywords_ru" db:"meta_keywords_ru"`
	MetaKeywordsTk    *string   `json:"meta_keywords_tk" db:"meta_keywords_tk"`
	Priority          int       `json:"priority" db:"priority"`
	ViewCount         int       `json:"view_count" db:"view_count"`
	IsFeatured        bool      `json:"is_featured" db:"is_featured"`
	IsPublic          bool      `json:"is_public" db:"is_public"`
	RequiresAuth      bool      `json:"requires_auth" db:"requires_auth"`
	ContentType       string    `json:"content_type" db:"content_type"`
	DifficultyLevel   *string   `json:"difficulty_level" db:"difficulty_level"`
	EstimatedReadTime *int      `json:"estimated_read_time" db:"estimated_read_time"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	Blocked           *int      `json:"blocked"`
	Active            *int      `json:"active"`
	Deleted           *int      `json:"deleted"`
}

type CreateWikiRequest struct {
	WikiContent
	WikiMetadata
	Files  []WikiFileInfo  `json:"files"`
	Videos []WikiVideoInfo `json:"videos"`
}

type UpdateWikiRequest struct {
	TitleEn         *string `json:"title_en"`
	TitleRu         *string `json:"title_ru"`
	TitleTk         *string `json:"title_tk"`
	DescriptionEn   *string `json:"description_en"`
	DescriptionRu   *string `json:"description_ru"`
	DescriptionTk   *string `json:"description_tk"`
	DescriptionType *string `json:"description_type" binding:"omitempty,oneof=plain html info"`
	TextMdEn        *string `json:"text_md_en"`
	TextMdRu        *string `json:"text_md_ru"`
	TextMdTk        *string `json:"text_md_tk"`
	TextRichEn      *string `json:"text_rich_en"`
	TextRichRu      *string `json:"text_rich_ru"`
	TextRichTk      *string `json:"text_rich_tk"`

	Category          *string `json:"category" binding:"omitempty,oneof=docs wiki guides tutorials api faq changelog troubleshooting"`
	Subcategory       *string `json:"subcategory"`
	Tags              *string `json:"tags"`
	Slug              *string `json:"slug"`
	MetaKeywordsEn    *string `json:"meta_keywords_en"`
	MetaKeywordsRu    *string `json:"meta_keywords_ru"`
	MetaKeywordsTk    *string `json:"meta_keywords_tk"`
	Priority          *int    `json:"priority" binding:"omitempty,min=0"`
	IsFeatured        *bool   `json:"is_featured"`
	IsPublic          *bool   `json:"is_public"`
	RequiresAuth      *bool   `json:"requires_auth"`
	ContentType       *string `json:"content_type" binding:"omitempty,oneof=article tutorial reference guide faq changelog"`
	DifficultyLevel   *string `json:"difficulty_level" binding:"omitempty,oneof=beginner intermediate advanced"`
	EstimatedReadTime *int    `json:"estimated_read_time" binding:"omitempty,min=0"`

	Files  []WikiFileInfo  `json:"files"`
	Videos []WikiVideoInfo `json:"videos"`
}

type WikiFilter struct {
	Category        *string `form:"category" binding:"omitempty,oneof=docs wiki guides tutorials api faq changelog troubleshooting"`
	Subcategory     *string `form:"subcategory"`
	ContentType     *string `form:"content_type" binding:"omitempty,oneof=article tutorial reference guide faq changelog"`
	IsFeatured      *bool   `form:"is_featured"`
	IsPublic        *bool   `form:"is_public"`
	RequiresAuth    *bool   `form:"requires_auth"`
	DifficultyLevel *string `form:"difficulty_level" binding:"omitempty,oneof=beginner intermediate advanced"`
	Language        *string `form:"language" binding:"omitempty,oneof=en ru tk"`
	Tags            *string `form:"tags"`
	Search          *string `form:"search"`
	Active          *bool   `form:"active"`
	Page            int     `form:"page,default=1" binding:"min=1"`
	PerPage         int     `form:"per_page,default=10" binding:"min=1,max=100"`
}

type WikiCategoryResponse struct {
	Category    string `json:"category"`
	Count       int    `json:"count"`
	Subcategory string `json:"subcategory,omitempty"`
}
