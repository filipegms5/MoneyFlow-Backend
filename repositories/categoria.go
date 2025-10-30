package repositories

import (
	"github.com/filipegms5/MoneyFlow-Backend/models"
	"gorm.io/gorm"
)

type CategoriaRepository struct {
	db *gorm.DB
}

func NewCategoriaRepository(db *gorm.DB) *CategoriaRepository {
	return &CategoriaRepository{db: db}
}

func (r *CategoriaRepository) GetByID(id uint) (*models.Categoria, error) {
	var c models.Categoria
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoriaRepository) GetByNome(nome string) (*models.Categoria, error) {
	var c models.Categoria
	if err := r.db.Where("nome = ?", nome).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoriaRepository) Create(c *models.Categoria) error {
	return r.db.Create(c).Error
}

func (r *CategoriaRepository) Update(c *models.Categoria) error {
	return r.db.Save(c).Error
}

func (r *CategoriaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Categoria{}, id).Error
}

func (r *CategoriaRepository) GetAll() ([]models.Categoria, error) {
	var list []models.Categoria
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *CategoriaRepository) BulkInsertIfEmpty(categorias []models.Categoria) error {
	var count int64
	if err := r.db.Model(&models.Categoria{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return r.db.Create(&categorias).Error
}
