package models

type Post struct {
	PK         string   `json:"pk" dynamodbav:"pk"`
	Title      string   `json:"title" dynamodbav:"title"`
	UserID     string   `json:"userId" dynamodbav:"userId"`
	Summary    string   `json:"summary" dynamodbav:"summary"`
	Content    string   `json:"content" dynamodbav:"content"`
	PostStatus string   `json:"postStatus" dynamodbav:"postStatus"`
	Categories []string `json:"categories" dynamodbav:"categories"`
	CreatedAt  string   `json:"createdAt" dynamodbav:"createdAt"`
	HeroImage  string   `json:"heroImage" dynamodbav:"heroImage"`
	TotalPosts int      `json:"totalPosts" dynamodbav:"totalPosts"`
}
