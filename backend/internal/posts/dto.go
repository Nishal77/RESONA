package posts

type CreatePostRequest struct {
	ContentText *string  `json:"content_text" binding:"omitempty,max=2000"`
	MediaURL    *string  `json:"media_url"`
	MediaType   *string  `json:"media_type" binding:"omitempty,oneof=image video"`
	CommunityID *string  `json:"community_id"`
	TagIDs      []string `json:"tag_ids"`
	// Language can be manually set if confidence < 0.70
	ManualLanguage *string `json:"manual_language" binding:"omitempty,oneof=kannada tamil telugu malayalam hindi english"`
}

type FeedQuery struct {
	Page     int    `form:"page,default=1"`
	Limit    int    `form:"limit,default=20"`
	Language string `form:"language"` // empty = all
}
