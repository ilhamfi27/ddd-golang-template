package model

// ExampleModel represents a license object
// @Description ExampleModel represents a license object
type ExampleModel struct {
	ID          uint   `json:"id" gorm:"primary_key" example:"1"`
	Name        string `json:"name" example:"License Name"`
	Description string `json:"description" example:"License Description"`
}

// TableName overrides the table name used by ExampleModel to `examples`
func (ExampleModel) TableName() string {
	return "examples"
}
