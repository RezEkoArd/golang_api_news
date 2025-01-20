package repository

import (
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/domain/model"
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContentRepository interface {
	GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, *entity.Page, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
}

type contentRepository struct {
	db *gorm.DB
}

// CreateContent implements ContentRepository.
func (c *contentRepository) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	// join dari array string menjadi string dengan menggunakan Join
	tags := strings.Join(req.Tags,",")
	modelContent := model.Content{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: req.CreatedByID,
	}

	err = c.db.Create(&modelContent).Error
	if err != nil {
		code = "[REPOSITORY] createContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// DeleteContext implements ContentRepository.
func (c *contentRepository) DeleteContent(ctx context.Context, id int64) error {

	err = c.db.Where("id = ?", id).Delete(&model.Content{}).Error
	if err!=nil{
		code := "[REPOSITORY] DeleteContent - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// GetContentByID implements ContentRepository.
func (c *contentRepository) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	var modelContent model.Content

	err = c.db.Where("id = ?", id).Preload("Category").Preload("User").First(&modelContent).Error
	if err != nil {
		code = "[REPOSITORY] GetContentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	tags := strings.Split(modelContent.Tags, ",")
	resp := entity.ContentEntity{
		ID:          modelContent.ID,
		Title:       modelContent.Title,
		Excerpt:     modelContent.Excerpt,
		Description: modelContent.Description,
		Image:       modelContent.Image,
		Tags:        tags,
		Status:      modelContent.Status,
		CategoryID:  modelContent.CategoryID,
		CreatedByID: modelContent.CreatedByID,
		CreatedAt:   modelContent.CreatedAt,
		Category: entity.CategoryEntity{
			ID:    modelContent.Category.ID,
			Title: modelContent.Category.Title,
			Slug:  modelContent.Category.Slug,
		},
		User: entity.UserEntity{
			ID:   modelContent.User.ID,
			Name: modelContent.User.Name,
		},
	}

	return &resp, nil
}

// GetContents implements ContentRepository.
func (c *contentRepository) GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, *entity.Page, error) {
	var modelContents []model.Content
	var totalCount int64

	orderQuery := ""
	if query.OrderBy != "" && query.OrderType != "" {
		validOrderColumns := map[string]bool{
			"created_at" : true,
			"title" : true,
		}
		validOrderTypes := map[string]bool{
            "ASC":  true,
            "DESC": true,
        }

		if validOrderColumns[query.OrderBy] && validOrderTypes[strings.ToUpper(query.OrderType)] {
			orderQuery = fmt.Sprintf("%s %s", query.OrderBy, strings.ToUpper(query.OrderType))
		}
	} else {
		orderQuery = "created_at DESC"
	}


	// Validation Pagination
	if query.Page < 1 {
		query.Page = 1
	}

	if query.Limit < 1{
		query.Limit = 10 // default
	}

	offset := (query.Page - 1) * query.Limit
	
	// Build Query
	baseQuery := c.db.Model(&model.Content{})

	//Tambahkan kondisi pencarian
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		baseQuery = baseQuery.Where("title ILIKE ? OR excerpt ILIKE ? OR description ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	//Tambah filter Status
	if query.Status != "" {
		baseQuery = baseQuery.Where("status ILIKE ?", "%"+query.Status+"%")
	}

	if query.CategoryID > 0 {
		baseQuery = baseQuery.Where("category_id = ?", query.CategoryID)
	}

	//Hitung Total Records
	err := baseQuery.Count(&totalCount).Error
	if err != nil {
		log.Errorw("[Repository] GetContents - Count", "error", err)
		return nil, nil, err
	}

	//Hitung Total Halaman
	totalPages := int(math.Ceil(float64(totalCount) / float64(query.Limit)))

	//EXecute
	err = baseQuery.
		Preload(clause.Associations).
		Order(orderQuery).
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Find(&modelContents).Error
		
	if err != nil {
		code = "[REPOSITORY] GetContents - Find"
		log.Errorw(code, "error", err)
		return nil, nil, err
	}

	resps := []entity.ContentEntity{}
	for _, val := range modelContents {
		tags := strings.Split(val.Tags, ".")
		resp := entity.ContentEntity{
			ID:           val.ID,
			Title:        val.Title,
			Excerpt:      val.Excerpt,
			Description:  val.Description,
			Image:        val.Image,
			Tags:         tags,
			Status:       val.Status,
			CategoryID:   val.CategoryID,
			CreatedByID: val.CreatedByID,
			CreatedAt:    val.CreatedAt,
			Category:     entity.CategoryEntity{
				ID:    val.CategoryID,
				Title: val.Category.Title,
				Slug:  val.Category.Slug, 
			},
			User:         entity.UserEntity{
				ID:       val.User.ID,
				Name:     val.User.Name,
			},
		}
		resps = append(resps, resp)
	}

	//Info Pagination
	pagination := &entity.Page{
		Page:       query.Page,
		Perpage:    query.Limit,
		PageCount:  totalPages,
		TotalCount: int(totalCount),
		First:      1,
		Last:       totalPages,
	}
	return resps, pagination, nil
}


// UpdateContent implements ContentRepository.
func (c *contentRepository) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: req.CreatedByID,
	}

	err = c.db.Where("id = ?", req.ID).Updates(&modelContent).Error
	if err != nil {
		code = "[REPOSITORY] UpdateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db: db}
}
