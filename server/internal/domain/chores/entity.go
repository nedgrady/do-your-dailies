package chores

type CreateRequest struct {
	Name          string `json:"name"`
	CadenceInDays int    `json:"cadence_in_days"`
}

type UpdateRequest struct {
	Name          *string `json:"name"`
	CadenceInDays *int    `json:"cadence_in_days"`
}
