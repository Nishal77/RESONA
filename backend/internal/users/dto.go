package users

type UpdateUserRequest struct {
	FullName        *string `json:"full_name" binding:"omitempty,min=1,max=100"`
	Bio             *string `json:"bio" binding:"omitempty,max=160"`
	State           *string `json:"state" binding:"omitempty,max=100"`
	City            *string `json:"city" binding:"omitempty,max=100"`
	PrimaryLanguage *string `json:"primary_language" binding:"omitempty,oneof=kannada tamil telugu malayalam hindi english"`
	AvatarURL       *string `json:"avatar_url"`
}

type CompleteOnboardingRequest struct {
	PrimaryLanguage    string   `json:"primary_language" binding:"required,oneof=kannada tamil telugu malayalam hindi english"`
	State              *string  `json:"state"`
	City               *string  `json:"city"`
	CommunityIDs       []string `json:"community_ids"` // communities to join
}

type FeedQuery struct {
	Page     int    `form:"page,default=1"`
	Limit    int    `form:"limit,default=20"`
	Language string `form:"language"` // empty = all
}
