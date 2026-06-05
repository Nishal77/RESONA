package communities

type CreateCommunityRequest struct {
	Name               string   `json:"name" binding:"required,min=3,max=100"`
	Description        *string  `json:"description"`
	PrimaryLanguage    string   `json:"primary_language" binding:"required,oneof=kannada tamil telugu malayalam hindi english"`
	SecondaryLanguages []string `json:"secondary_languages"`
	AvatarURL          *string  `json:"avatar_url"`
	BannerURL          *string  `json:"banner_url"`
}

type UpdateCommunityRequest struct {
	Name               *string  `json:"name" binding:"omitempty,min=3,max=100"`
	Description        *string  `json:"description"`
	SecondaryLanguages []string `json:"secondary_languages"`
	AvatarURL          *string  `json:"avatar_url"`
	BannerURL          *string  `json:"banner_url"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=member moderator admin"`
}
