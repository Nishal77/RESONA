package comments

type CreateCommentRequest struct {
	Content         string  `json:"content" binding:"required,min=1,max=2000"`
	ParentCommentID *string `json:"parent_comment_id"`
}
