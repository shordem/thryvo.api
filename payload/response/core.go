package response

import (
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
	UUID          uuid.UUID          `json:"id"`
	Slug          string             `json:"slug"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Specification string             `json:"specification"`
	Price         float64            `json:"price"`
	SlashPrice    float64            `json:"slash_price"`
	Stock         int                `json:"stock"`
	Sales         int                `json:"sales"`
	Category      CategoryResponse   `json:"category"`
	Images        []ImageResponse    `json:"images"`
	Discounts     []DiscountResponse `json:"discounts"`
	CreatedAt     time.Time          `json:"created_at"`
}

type ImageResponse struct {
	Key string `json:"key"`
}

type DiscountResponse struct {
	Discount float64 `json:"discount"`
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
}

type CategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Slug string    `json:"slug"`
	Name string    `json:"name"`
}

type WishlistResponse struct {
	Product ProductResponse `json:"product"`
}
