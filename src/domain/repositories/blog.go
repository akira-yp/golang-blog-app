package repositories

import (
	"context"
	"gin-todo-app/src/domain/models"
)

type BlogRepository interface {
	GetByID(ctx context.Context, id uint) (*models.Blog, error)
	Create(ctx context.Context, blog *models.Blog) error
	Update(ctx context.Context, blog *models.Blog) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*models.Blog, error)
}
