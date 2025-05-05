package domain

type PersonFilter struct {
    Name        *string
    Surname     *string
    Age         *int
    Gender      *string
    Nationality *string
    MinAge      *int
    MaxAge      *int
}