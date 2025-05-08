package domain

type PersonFilter struct {
	Name        *string
	Surname     *string
	Gender      *string
	Nationality *string
	MinAge      *int
	MaxAge      *int
	Page        int
	Size        int
}
