package domain



type CreateOrUpdateVariableRequest struct {
	Name         string `json:"name" validate:"required,min=1,max=100"`
	Value        string `json:"value" validate:"required,min=1,max=255"`
	VariableType string `json:"variableType" validate:"required,min=1,max=50"`
}
