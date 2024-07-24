package request

import "github.com/google/uuid"

type CreateCategoryRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
}

type UpdateCategoryRequest struct {
	CreateCategoryRequest
}

type CreateProductRequest struct {
	CategoryUUID  string            `json:"category_id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Specification string            `json:"specification"`
	Price         int               `json:"price"`
	Stock         int               `json:"stock"`
	SlashPrice    int               `json:"slash_price"`
	Images        []string          `json:"images"`
	Discounts     []DiscountRequest `json:"discounts"`
	// ProductVariations []struct {
	// 	Name  string `json:"name"`
	// 	Price int    `json:"price"`
	// 	Stock int    `json:"stock"`
	// } `json:"product_variations"`
}

type DiscountRequest struct {
	Discount int `json:"discount"`
	Min      int `json:"min"`
	Max      int `json:"max"`
}

type CreateWishlistRequest struct {
	ProductUUID string `json:"product_id"`
}

type CreateReviewRequest struct {
	ProductID string `json:"product_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}
