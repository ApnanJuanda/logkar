package repository

import (
	"gorm.io/gorm"
	"strings"
)

type GeneralRepositoryInterface interface {
	GetDataSequence(tx *gorm.DB, shortName string) (sequenceNumber int64, err error)
	BeginTrans() *gorm.DB
	CommitTrans(tx *gorm.DB) error
	RollbackTrans(tx *gorm.DB) error
}

type generalRepository struct {
	DB *gorm.DB
}

func NewGeneralRepository(DB *gorm.DB) GeneralRepositoryInterface {
	return &generalRepository{
		DB: DB,
	}
}

func (r *generalRepository) GetDataSequence(tx *gorm.DB, shortName string) (sequenceNumber int64, err error) {
	if tx == nil {
		tx = r.DB
	}
	seqName := "bsnack_code_seq_" + strings.ToLower(shortName)

	createSeqQuery := "CREATE SEQUENCE IF NOT EXISTS " + seqName + " START 1"
	err = tx.Exec(createSeqQuery).Error
	if err != nil {
		return 0, err
	}

	var seqNumber float64
	err = tx.Raw("SELECT nextval('" + seqName + "') as seqNumber").Find(&seqNumber).Error
	if err != nil {
		return 0, err
	}

	sequenceNumber = int64(seqNumber)
	return
}

func (r *generalRepository) BeginTrans() *gorm.DB {
	return r.DB.Begin()
}

func (r *generalRepository) CommitTrans(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *generalRepository) RollbackTrans(tx *gorm.DB) error {
	return tx.Rollback().Error
}
