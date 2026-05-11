package fn

import (
	"base/gorm/query"
	"context"
	"errors"
	"gorm.io/gorm"
)

type IRepo interface {
	GetDB() *gorm.DB
	Create(ctx context.Context, tx *gorm.DB, im IModel, fs ...func(*gorm.DB) *gorm.DB) error
	Update(ctx context.Context, tx *gorm.DB, im IModel, fs ...func(*gorm.DB) *gorm.DB) error
	FirstByCond(ctx context.Context, tx *gorm.DB, im IModel, fs ...func(*gorm.DB) *gorm.DB) error
}

type IModel interface {
	TableName() string
}

type Model struct {
	db        *gorm.DB
	isHardDel bool
}

func NewModel(db *gorm.DB, opts ...ModelOption) *Model {
	model := &Model{
		db: db,
	}
	if opts != nil {
		for _, opt := range opts {
			opt(model)
		}
	}
	return model
}

func (m *Model) GetDB() *gorm.DB {
	return m.db
}

func (m *Model) checkDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return m.db
}

func (m *Model) Create(ctx context.Context, tx *gorm.DB, im IModel, fs ...func(*gorm.DB) *gorm.DB) error {
	return func(db *gorm.DB) error {
		return db.Scopes(fs...).Create(im).Error
	}(m.checkDB(tx))
}

func (m *Model) Update(ctx context.Context, tx *gorm.DB, im IModel, fs ...func(*gorm.DB) *gorm.DB) error {
	if len(fs) == 0 {
		return errors.New("update method must pass conditions")
	}
	return func(db *gorm.DB) error {
		return db.Scopes(fs...).Updates(im).Error
	}(m.checkDB(tx))
}

func (m *Model) FirstByCond(ctx context.Context, tx *gorm.DB, im IModel, fs ...func(*gorm.DB) *gorm.DB) error {
	if len(fs) == 0 {
		return errors.New("first method must pass conditions")
	}
	if !m.isHardDel {
		fs = append(fs, query.Equal("delete_user_id", 0))
	}
	return func(db *gorm.DB) error {
		return db.Scopes(fs...).First(im).Error
	}(m.checkDB(tx))
}
