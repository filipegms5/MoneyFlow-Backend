package repositories

import (
	"github.com/filipegms5/MoneyFlow-Backend/models"
	"gorm.io/gorm"
)

type UsuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

func (r *UsuarioRepository) Create(usuario *models.Usuario) error {
	return r.db.Create(usuario).Error
}

func (r *UsuarioRepository) Update(usuario *models.Usuario) error {
	return r.db.Save(usuario).Error
}

func (r *UsuarioRepository) Delete(id uint) error {
	return r.db.Delete(&models.Usuario{}, id).Error
}

func (r *UsuarioRepository) GetAll() ([]models.Usuario, error) {
	var usuarios []models.Usuario
	err := r.db.Find(&usuarios).Error
	return usuarios, err
}

func (r *UsuarioRepository) GetByID(id uint) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.First(&usuario, id).Error
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (r *UsuarioRepository) GetByEmail(email string) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.Where("email = ?", email).First(&usuario).Error
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}
