package controllers

import (
	"net/http"
	"strconv"
	"texApi/pkg/schemas/request"
	"texApi/pkg/services"

	"github.com/gin-gonic/gin"
)

// ContentController handles requests related to content.
type ContentController struct {
	service *services.ContentService
}

// NewContentController creates a new ContentController.
func NewContentController(service *services.ContentService) *ContentController {
	return &ContentController{service: service}
}

// GetAll retrieves all content.
func (c *ContentController) GetAll(ctx *gin.Context) {
	contents, err := c.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, contents)
}

// GetByID retrieves content by ID.
func (c *ContentController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	content, err := c.service.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}
	ctx.JSON(http.StatusOK, content)
}

// GetByTitle retrieves content by title.
func (c *ContentController) GetByTitle(ctx *gin.Context) {
	title := ctx.Query("title")
	contents, err := c.service.GetByTitle(ctx, title)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, contents)
}

// GetByUUID retrieves content by UUID.
func (c *ContentController) GetByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	content, err := c.service.GetByUUID(ctx, uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}
	ctx.JSON(http.StatusOK, content)
}

// GetByContentTypeID retrieves content by content type ID.
func (c *ContentController) GetByContentTypeID(ctx *gin.Context) {
	contentTypeID, err := strconv.Atoi(ctx.Param("content_type_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type ID"})
		return
	}
	contents, err := c.service.GetByContentTypeID(ctx, contentTypeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, contents)
}

// Create handles the creation of new content.
func (c *ContentController) Create(ctx *gin.Context) {
	var contentRequest request.ContentRequest
	if err := ctx.ShouldBindJSON(&contentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content, err := c.service.Create(ctx, contentRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, content)
}

// Update handles updating existing content.
func (c *ContentController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var contentRequest request.ContentRequest
	if err := ctx.ShouldBindJSON(&contentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	content, err := c.service.Update(ctx, id, contentRequest)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}
	ctx.JSON(http.StatusOK, content)
}

// Delete handles the deletion of content.
func (c *ContentController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.service.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil) // No content to return
}
