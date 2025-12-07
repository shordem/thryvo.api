package request

type CreatePlan struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Currency    string  `json:"currency" validate:"required,len=3"`
	Duration    int     `json:"duration" validate:"required,gt=0"` // days
}

type UpdatePlan struct {
	Name        string  `json:"name" validate:"omitempty,min=3,max=100"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	Duration    int     `json:"duration" validate:"omitempty,gt=0"` // days
	IsActive    *bool   `json:"is_active"`
}
