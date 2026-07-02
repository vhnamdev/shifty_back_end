package constants

const (
	ReactionLike  = "LIKE"
	ReactionLove  = "LOVE"
	ReactionCare  = "CARE"
	ReactionHaha  = "HAHA"
	ReactionWow   = "WOW"
	ReactionSad   = "SAD"
	ReactionAngry = "ANGRY"
)

var ReactionTypes = []string{
	ReactionLike,
	ReactionLove,
	ReactionCare,
	ReactionHaha,
	ReactionWow,
	ReactionSad,
	ReactionAngry,
}

var validReactionTypes = map[string]struct{}{
	ReactionLike:  {},
	ReactionLove:  {},
	ReactionCare:  {},
	ReactionHaha:  {},
	ReactionWow:   {},
	ReactionSad:   {},
	ReactionAngry: {},
}

func IsValidReactionType(reactionType string) bool {
	_, ok := validReactionTypes[reactionType]
	return ok
}
