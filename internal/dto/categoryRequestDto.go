package dto

type CreateCategoryRequestDto struct {
	Name         string `json:"name"`
	ParentId     uint   `json:"parent_id"`
	ImageURL     string `json:"image_url"`
	DisplayOrder int    `json:"display_order"`
}
