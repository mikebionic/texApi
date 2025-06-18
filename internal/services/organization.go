package services

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

func GetOrganization(ctx *gin.Context) {
	wishlist, err := GetOrg()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve organization", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Organization retrieved", wishlist))
}

func CreateOrganization(ctx *gin.Context) {
	var req dto.CreateOrganizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	organization, err := CreateOrg(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to create organization", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Organization created successfully", organization))
}

func UpdateOrganization(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid organization ID", err.Error()))
		return
	}

	var req dto.UpdateOrganizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	organization, err := UpdateOrg(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to update organization", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Organization updated successfully", organization))
}

func DeleteOrganization(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid organization ID", err.Error()))
		return
	}

	err = DeleteOrg(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to delete organization", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Organization deleted successfully", nil))
}

func GetOrg() (organization dto.Organization, err error) {
	query := `SELECT * FROM tbl_organization WHERE active = 1 AND deleted = 0 LIMIT 1`
	err = pgxscan.Get(context.Background(), db.DB, &organization, query)
	return
}

func CreateOrg(req dto.CreateOrganizationRequest) (dto.Organization, error) {
	var organization dto.Organization
	query := `
		INSERT INTO tbl_organization (
			name,
			description_en,
			description_ru,
			description_tk,
			email,
			image_url,
			logo_url,
			icon_url,
			banner_url,
			website_url,
			about_text,
			refund_text,
			delivery_text,
			contact_text,
			terms_conditions,
			privacy_policy,
			address1,
			address2,
			address3,
			address4,
			address_title1,
			address_title2,
			address_title3,
			address_title4,
			contact_phone1,
			contact_phone2,
			contact_phone3,
			contact_phone4,
			contact_title1,
			contact_title2,
			contact_title3,
			contact_title4,
			meta,
			meta2,
			meta3,
			active,
			deleted
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30, $31, $32, $33, $34, $35, 1, 0
		) RETURNING *`

	err := pgxscan.Get(context.Background(), db.DB, &organization, query,
		req.Name, req.DescriptionEN, req.DescriptionRU, req.DescriptionTK,
		req.Email, req.ImageUrl, req.LogoUrl, req.IconUrl, req.BannerUrl,
		req.WebsiteUrl, req.AboutText, req.RefundText, req.DeliveryText,
		req.ContactText, req.TermsConditions, req.PrivacyPolicy,
		req.Address1, req.Address2, req.Address3, req.Address4,
		req.AddressTitle1, req.AddressTitle2, req.AddressTitle3, req.AddressTitle4,
		req.ContactPhone1, req.ContactPhone2, req.ContactPhone3, req.ContactPhone4,
		req.ContactTitle1, req.ContactTitle2, req.ContactTitle3, req.ContactTitle4,
		req.Meta, req.Meta2, req.Meta3,
	)

	return organization, err
}

func UpdateOrg(id int, req dto.UpdateOrganizationRequest) (dto.Organization, error) {
	var organization dto.Organization

	existsQuery := `SELECT id FROM tbl_organization WHERE id = $1 AND deleted = 0`
	var existingID int
	err := pgxscan.Get(context.Background(), db.DB, &existingID, existsQuery, id)
	if err != nil {
		return organization, fmt.Errorf("organization not found")
	}

	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.DescriptionEN != nil {
		setParts = append(setParts, fmt.Sprintf("description_en = $%d", argIndex))
		args = append(args, req.DescriptionEN)
		argIndex++
	}
	if req.DescriptionRU != nil {
		setParts = append(setParts, fmt.Sprintf("description_ru = $%d", argIndex))
		args = append(args, req.DescriptionRU)
		argIndex++
	}
	if req.DescriptionTK != nil {
		setParts = append(setParts, fmt.Sprintf("description_tk = $%d", argIndex))
		args = append(args, req.DescriptionTK)
		argIndex++
	}
	if req.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, req.Email)
		argIndex++
	}
	if req.ImageUrl != nil {
		setParts = append(setParts, fmt.Sprintf("image_url = $%d", argIndex))
		args = append(args, req.ImageUrl)
		argIndex++
	}
	if req.LogoUrl != nil {
		setParts = append(setParts, fmt.Sprintf("logo_url = $%d", argIndex))
		args = append(args, req.LogoUrl)
		argIndex++
	}
	if req.IconUrl != nil {
		setParts = append(setParts, fmt.Sprintf("icon_url = $%d", argIndex))
		args = append(args, req.IconUrl)
		argIndex++
	}
	if req.BannerUrl != nil {
		setParts = append(setParts, fmt.Sprintf("banner_url = $%d", argIndex))
		args = append(args, req.BannerUrl)
		argIndex++
	}
	if req.WebsiteUrl != nil {
		setParts = append(setParts, fmt.Sprintf("website_url = $%d", argIndex))
		args = append(args, req.WebsiteUrl)
		argIndex++
	}
	if req.AboutText != nil {
		setParts = append(setParts, fmt.Sprintf("about_text = $%d", argIndex))
		args = append(args, req.AboutText)
		argIndex++
	}
	if req.RefundText != nil {
		setParts = append(setParts, fmt.Sprintf("refund_text = $%d", argIndex))
		args = append(args, req.RefundText)
		argIndex++
	}
	if req.DeliveryText != nil {
		setParts = append(setParts, fmt.Sprintf("delivery_text = $%d", argIndex))
		args = append(args, req.DeliveryText)
		argIndex++
	}
	if req.ContactText != nil {
		setParts = append(setParts, fmt.Sprintf("contact_text = $%d", argIndex))
		args = append(args, req.ContactText)
		argIndex++
	}
	if req.TermsConditions != nil {
		setParts = append(setParts, fmt.Sprintf("terms_conditions = $%d", argIndex))
		args = append(args, req.TermsConditions)
		argIndex++
	}
	if req.PrivacyPolicy != nil {
		setParts = append(setParts, fmt.Sprintf("privacy_policy = $%d", argIndex))
		args = append(args, req.PrivacyPolicy)
		argIndex++
	}
	if req.Address1 != nil {
		setParts = append(setParts, fmt.Sprintf("address1 = $%d", argIndex))
		args = append(args, req.Address1)
		argIndex++
	}
	if req.Address2 != nil {
		setParts = append(setParts, fmt.Sprintf("address2 = $%d", argIndex))
		args = append(args, req.Address2)
		argIndex++
	}
	if req.Address3 != nil {
		setParts = append(setParts, fmt.Sprintf("address3 = $%d", argIndex))
		args = append(args, req.Address3)
		argIndex++
	}
	if req.Address4 != nil {
		setParts = append(setParts, fmt.Sprintf("address4 = $%d", argIndex))
		args = append(args, req.Address4)
		argIndex++
	}
	if req.AddressTitle1 != nil {
		setParts = append(setParts, fmt.Sprintf("address_title1 = $%d", argIndex))
		args = append(args, req.AddressTitle1)
		argIndex++
	}
	if req.AddressTitle2 != nil {
		setParts = append(setParts, fmt.Sprintf("address_title2 = $%d", argIndex))
		args = append(args, req.AddressTitle2)
		argIndex++
	}
	if req.AddressTitle3 != nil {
		setParts = append(setParts, fmt.Sprintf("address_title3 = $%d", argIndex))
		args = append(args, req.AddressTitle3)
		argIndex++
	}
	if req.AddressTitle4 != nil {
		setParts = append(setParts, fmt.Sprintf("address_title4 = $%d", argIndex))
		args = append(args, req.AddressTitle4)
		argIndex++
	}
	if req.ContactPhone1 != nil {
		setParts = append(setParts, fmt.Sprintf("contact_phone1 = $%d", argIndex))
		args = append(args, req.ContactPhone1)
		argIndex++
	}
	if req.ContactPhone2 != nil {
		setParts = append(setParts, fmt.Sprintf("contact_phone2 = $%d", argIndex))
		args = append(args, req.ContactPhone2)
		argIndex++
	}
	if req.ContactPhone3 != nil {
		setParts = append(setParts, fmt.Sprintf("contact_phone3 = $%d", argIndex))
		args = append(args, req.ContactPhone3)
		argIndex++
	}
	if req.ContactPhone4 != nil {
		setParts = append(setParts, fmt.Sprintf("contact_phone4 = $%d", argIndex))
		args = append(args, req.ContactPhone4)
		argIndex++
	}
	if req.ContactTitle1 != nil {
		setParts = append(setParts, fmt.Sprintf("contact_title1 = $%d", argIndex))
		args = append(args, req.ContactTitle1)
		argIndex++
	}
	if req.ContactTitle2 != nil {
		setParts = append(setParts, fmt.Sprintf("contact_title2 = $%d", argIndex))
		args = append(args, req.ContactTitle2)
		argIndex++
	}
	if req.ContactTitle3 != nil {
		setParts = append(setParts, fmt.Sprintf("contact_title3 = $%d", argIndex))
		args = append(args, req.ContactTitle3)
		argIndex++
	}
	if req.ContactTitle4 != nil {
		setParts = append(setParts, fmt.Sprintf("contact_title4 = $%d", argIndex))
		args = append(args, req.ContactTitle4)
		argIndex++
	}
	if req.Meta != nil {
		setParts = append(setParts, fmt.Sprintf("meta = $%d", argIndex))
		args = append(args, req.Meta)
		argIndex++
	}
	if req.Meta2 != nil {
		setParts = append(setParts, fmt.Sprintf("meta2 = $%d", argIndex))
		args = append(args, req.Meta2)
		argIndex++
	}
	if req.Meta3 != nil {
		setParts = append(setParts, fmt.Sprintf("meta3 = $%d", argIndex))
		args = append(args, req.Meta3)
		argIndex++
	}
	if req.Active != nil {
		setParts = append(setParts, fmt.Sprintf("active = $%d", argIndex))
		args = append(args, *req.Active)
		argIndex++
	}

	args = append(args, id)

	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf(`
	UPDATE tbl_organization 
	SET %s 
	WHERE id = $%d AND deleted = 0 
	RETURNING *`, setClause, argIndex)

	err = pgxscan.Get(context.Background(), db.DB, &organization, query, args...)
	return organization, err
}

func DeleteOrg(id int) error {
	query := `UPDATE tbl_organization SET deleted = 1, updated_at = NOW() WHERE id = $1 AND deleted = 0`
	_, err := db.DB.Exec(context.Background(), query, id)
	return err
}
