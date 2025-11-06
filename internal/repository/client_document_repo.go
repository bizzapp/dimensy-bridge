package repository

import (
	"dimensy-bridge/internal/model"

	"gorm.io/gorm"
)

type ClientDocumentRepository interface {
	Create(doc *model.ClientDocument) error
	Update(doc *model.ClientDocument) error
	Delete(id int64) error
	FindByID(id int64) (*model.ClientDocument, error)
	FindAll() ([]model.ClientDocument, error)
}

type clientDocumentRepository struct {
	db *gorm.DB
}

func NewClientDocumentRepository(db *gorm.DB) ClientDocumentRepository {
	return &clientDocumentRepository{db: db}
}

func (r *clientDocumentRepository) Create(doc *model.ClientDocument) error {
	return r.db.Create(doc).Error
}

func (r *clientDocumentRepository) Update(doc *model.ClientDocument) error {
	return r.db.Save(doc).Error
}

func (r *clientDocumentRepository) Delete(id int64) error {
	return r.db.Delete(&model.ClientDocument{}, id).Error
}

func (r *clientDocumentRepository) FindByID(id int64) (*model.ClientDocument, error) {
	var doc model.ClientDocument
	if err := r.db.First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *clientDocumentRepository) FindAll() ([]model.ClientDocument, error) {
	var docs []model.ClientDocument
	if err := r.db.Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}
