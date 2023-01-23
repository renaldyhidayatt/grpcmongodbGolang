package services

import "examplegrpcgin/models"

type PostService interface {
	CreatePost(post *models.CreatePostRequest) (*models.DBPost, error)
	UpdatePost(string, *models.UpdatePost) (*models.DBPost, error)
	FindPostById(string) (*models.DBPost, error)
	FindPosts(page int, limit int) ([]*models.DBPost, error)
	DeletePost(string) error
}
