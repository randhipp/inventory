package models

type User struct {
	BaseModel
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `gorm:"varchar(60)" json:"-"`
}

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Status string `json:"status"`
	User   *User  `json:"user"`
}

type UsersResponse struct {
	Status      string `json:"status"`
	Users       []User `json:"users"`
	PerPage     int    `json:"perPage,omitempty"`
	TotalPage   int    `json:"totalPage,omitempty"`
	CurrentPage int    `json:"currentPage"`
	NextPage    int    `json:"nextPage,omitempty"`
	PrevPage    int    `json:"prevPage,omitempty"`
}

type UserDeleteResponse struct {
	Status string `json:"status"`
}
