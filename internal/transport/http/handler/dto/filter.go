package dto

type Filter struct {
	Name        *string `form:"name"`
	Surname     *string `form:"surname"`
	MinAge      *int    `form:"min_age"`
	MaxAge      *int    `form:"max_age"`
	Gender      *string `form:"gender"`
	Nationality *string `form:"nationality"`
	Page        int     `form:"page"`
	Size        int     `form:"size"`
}
