package dto

type PackagingTypeResponse struct {
	ID            int     `json:"id"`
	NameRu        string  `json:"name_ru"`
	NameEn        string  `json:"name_en"`
	NameTk        string  `json:"name_tk"`
	CategoryRu    string  `json:"category_ru"`
	CategoryEn    string  `json:"category_en"`
	CategoryTk    string  `json:"category_tk"`
	Material      string  `json:"material"`
	Dimensions    string  `json:"dimensions"`
	Weight        float64 `json:"weight"`
	DescriptionRu string  `json:"description_ru"`
	DescriptionEn string  `json:"description_en"`
	DescriptionTk string  `json:"description_tk"`
	Active        int     `json:"active"`
	Deleted       int     `json:"deleted"`
}

type CreatePackagingType struct {
	NameRu        *string  `json:"name_ru,omitempty"`
	NameEn        *string  `json:"name_en,omitempty"`
	NameTk        *string  `json:"name_tk,omitempty"`
	CategoryRu    *string  `json:"category_ru,omitempty"`
	CategoryEn    *string  `json:"category_en,omitempty"`
	CategoryTk    *string  `json:"category_tk,omitempty"`
	Material      *string  `json:"material,omitempty"`
	Dimensions    *string  `json:"dimensions,omitempty"`
	Weight        *float64 `json:"weight,omitempty"`
	DescriptionRu *string  `json:"description_ru,omitempty"`
	DescriptionEn *string  `json:"description_en,omitempty"`
	DescriptionTk *string  `json:"description_tk,omitempty"`
	Active        *int     `json:"active,omitempty"`
}
