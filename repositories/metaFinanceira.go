package repositories

import (
	"github.com/filipegms5/MoneyFlow-Backend/models"
	"gorm.io/gorm"
)

type MetaFinanceiraRepository struct {
	db *gorm.DB
}

func NewMetaFinanceiraRepository(db *gorm.DB) *MetaFinanceiraRepository {
	return &MetaFinanceiraRepository{db: db}
}

func (r *MetaFinanceiraRepository) Create(m *models.MetaFinanceira) error {
	return r.db.Create(m).Error
}

func (r *MetaFinanceiraRepository) Update(m *models.MetaFinanceira) error {
	return r.db.Save(m).Error
}

func (r *MetaFinanceiraRepository) Delete(id uint) error {
	return r.db.Delete(&models.MetaFinanceira{}, id).Error
}

func (r *MetaFinanceiraRepository) GetByID(id uint) (*models.MetaFinanceira, error) {
	var m models.MetaFinanceira
	if err := r.db.Preload("Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id", "email") }).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MetaFinanceiraRepository) GetAll() ([]models.MetaFinanceira, error) {
	var metas []models.MetaFinanceira
	if err := r.db.Preload("Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id", "email") }).Find(&metas).Error; err != nil {
		return nil, err
	}
	return metas, nil
}

func (r *MetaFinanceiraRepository) GetByUsuarioID(usuarioID uint) ([]models.MetaFinanceira, error) {
	var metas []models.MetaFinanceira
	if err := r.db.Preload("Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id", "email") }).Where("usuario_id = ?", usuarioID).Find(&metas).Error; err != nil {
		return nil, err
	}
	return metas, nil
}

func (r *MetaFinanceiraRepository) DeactivateAllByUsuarioID(usuarioID uint) error {
	return r.db.Model(&models.MetaFinanceira{}).Where("usuario_id = ?", usuarioID).Update("ativa", false).Error
}

func (r *MetaFinanceiraRepository) GetActiveByUsuarioID(usuarioID uint) (*models.MetaFinanceira, error) {
	var meta models.MetaFinanceira
	if err := r.db.Preload("Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id", "email") }).
		Where("usuario_id = ? AND ativa = true", usuarioID).Order("id DESC").First(&meta).Error; err != nil {
		return nil, err
	}
	return &meta, nil
}
