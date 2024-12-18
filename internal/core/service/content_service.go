package service

import (
	"bwanews/config"
	"bwanews/internal/adapter/cloudflare"
	"bwanews/internal/adapter/repository"
	"bwanews/internal/core/domain/entity"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type ContentService interface {
	GetContents(ctx context.Context) ([]entity.ContentEntity, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
	UploadImageR2(ctx context.Context, req entity.FileUploadEntity) (string, error)
}

type contentService struct {
	contentRepo repository.ContentRepository
	cfg         *config.Config
	r2          cloudflare.CloudflareR2Adapter
}

// CreateContent implements ContentService.
func (c *contentService) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	panic("unimplemented")
}

// DeleteContext implements ContentService.
func (c *contentService) DeleteContent(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetContentByID implements ContentService.
func (c *contentService) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	panic("unimplemented")
}

// GetContents implements ContentService.
func (c *contentService) GetContents(ctx context.Context) ([]entity.ContentEntity, error) {
	result, err := c.contentRepo.GetContents(ctx)
	if err!= nil {
		code = "[SERVICE] GetContents - 1"
		log.Errorw(code, err)
		return nil, err	
	}

	return result, nil
}

// UpdateContext implements ContentService.
func (c *contentService) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	panic("unimplemented")
}

// UploadImageR2 implements ContentService.
func (c *contentService) UploadImageR2(ctx context.Context, req entity.FileUploadEntity) (string, error) {
	panic("unimplemented")
}

func NewContentService(repo repository.ContentRepository, cfg *config.Config, r2 cloudflare.CloudflareR2Adapter) ContentService {
	return &contentService{
		contentRepo: repo,
		cfg:         cfg,
		r2:          r2,
	}
}
