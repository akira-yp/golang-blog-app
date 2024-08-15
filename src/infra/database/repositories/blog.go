package repository

import (
	"context"
	"gin-todo-app/src/domain/models"
	"gin-todo-app/src/domain/repositories"

	"gorm.io/gorm"
)

// GormによるBlogリポジトリの実装
type BlogRepository struct {
	DB *gorm.DB
}

func NewBlogRepository(db *gorm.DB) repositories.BlogRepository {
	return &BlogRepository{DB: db}
}

func (r *BlogRepository) GetByID(ctx context.Context, id uint) (*models.Blog, error) {
	var blog models.Blog
	result := r.DB.First(&blog, id)
	return &blog, result.Error
}

func (r *BlogRepository) Create(ctx context.Context, blog *models.Blog) error {
	return r.DB.Create(blog).Error
}

func (r *BlogRepository) Update(ctx context.Context, blog *models.Blog) error {
	return r.DB.Save(blog).Error
}

func (r *BlogRepository) Delete(ctx context.Context, id uint) error {
	return r.DB.Delete(&models.Blog{}, id).Error
}

func (r *BlogRepository) List(ctx context.Context) ([]*models.Blog, error) {
	var blogs []*models.Blog
	result := r.DB.Find(&blogs)
	return blogs, result.Error
}
