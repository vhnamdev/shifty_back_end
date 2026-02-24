package dto

type UserFilter struct {
	Search   *string
	Position *string
}

type UserRank struct {
	UserID string
	Rank   int
}
