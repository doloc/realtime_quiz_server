package storage

import (
	"context"
	"realtime_quiz_server/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *storage {
	return &storage{
		db: db,
	}
}

type BaseStorage interface {
	BeginTx() *gorm.DB
	CloseTx(tx *gorm.DB, err error)
	Create(tx *gorm.DB, data interface{}, result interface{}) error
	Get(ctx context.Context, cond map[string]interface{}, result interface{}) error
	GetForUpdate(ctx context.Context, cond map[string]interface{}, result interface{}) error
	Find(ctx context.Context, cond map[string]interface{}, result interface{}) error
	Update(tx *gorm.DB, cond map[string]interface{}, data interface{}, result interface{}) error
	Delete(tx *gorm.DB, cond map[string]interface{}, result interface{}) error
	GetAll(ctx context.Context, result interface{}) error
	CountTotal(ctx context.Context, cond map[string]interface{}, result interface{}) (int64, error)
	CommitTx(tx *gorm.DB) error
	RollbackTx(tx *gorm.DB)
}

func (s *storage) CommitTx(tx *gorm.DB) error {
	return tx.Commit().Error
}
func (s *storage) RollbackTx(tx *gorm.DB) {
	tx.Rollback()
}
func (s *storage) BeginTx() *gorm.DB {
	return s.db.Begin()
}

func (s *storage) CloseTx(tx *gorm.DB, err error) {
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (s *storage) Create(tx *gorm.DB, data interface{}, result interface{}) error {
	if err := utils.ConvertStruct(data, result); err != nil {
		return err
	}
	if err := tx.Create(result).Error; err != nil {
		return err
	}
	return nil
}

func (s *storage) Update(tx *gorm.DB, cond map[string]interface{}, data interface{}, result interface{}) error {
	if err := utils.ConvertStruct(data, result); err != nil {
		return err
	}
	db := tx.Where(cond).Updates(result)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (s *storage) Delete(tx *gorm.DB, cond map[string]interface{}, result interface{}) error {
	db := tx.Where(cond).Delete(result)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (s *storage) Get(ctx context.Context, cond map[string]interface{}, result interface{}) error {
	if err := s.db.Where(cond).First(result).Error; err != nil {
		return err
	}
	return nil
}

func (s *storage) GetForUpdate(ctx context.Context, cond map[string]interface{}, result interface{}) error {
	if err := s.db.Clauses(clause.Locking{Strength: "NO KEY UPDATE"}).Where(cond).First(result).Error; err != nil {
		return err
	}
	return nil
}

func (s *storage) Find(ctx context.Context, cond map[string]interface{}, result interface{}) error {
	if err := s.db.Where(cond).Find(result).Error; err != nil {
		return err
	}
	return nil
}

func (s *storage) GetAll(ctx context.Context, result interface{}) error {
	if err := s.db.Find(result).Error; err != nil {
		return err
	}
	return nil
}
func (s *storage) CountTotal(ctx context.Context, cond map[string]interface{}, result interface{}) (int64, error) {
	var count int64
	if err := s.db.Model(result).Where(cond).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
