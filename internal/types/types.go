package types

type Student struct {
	Id       int    `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password"`
}

type StudentPatch struct {
	Id       *int
	Name     *string
	Email    *string
	Password *string
}
