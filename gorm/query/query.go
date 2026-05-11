package query

import (
	"fmt"
	"gorm.io/gorm"
)

func Like(key string, val string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if val != "" {
			return db.Where(fmt.Sprintf("%s %s ?", key, "LIKE"), "%"+val+"%")
		}
		return db
	}
}

func Equal[T comparable](key string, value T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var zero T
		if value == zero {
			return db
		}
		return db.Where(fmt.Sprintf("%s = ?", key), value)
	}
}

func NotEqual[T comparable](key string, value T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var zero T
		if value == zero {
			return db
		}
		return db.Where(fmt.Sprintf("%s != ?", key), value)
	}
}

func In[T comparable](key string, values ...T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(values) == 0 {
			return db
		}
		return db.Where(fmt.Sprintf("%s IN (?)", key), values)
	}
}

func NotIn[T comparable](key string, values ...T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(values) == 0 {
			return db
		}
		return db.Where(fmt.Sprintf("%s NOT IN (?)", key), values)
	}
}

func TimeRange[T comparable](field string, start, end T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var zero T
		if start != zero {
			db = db.Where(fmt.Sprintf("%s >= ?", field), start)
		}
		if end != zero {
			db = db.Where(fmt.Sprintf("%s <= ?", field), end)
		}
		return db
	}
}

// where key > val
func GT(key string, val uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s %s ?", key, ">"), val)
	}
}

func StringListIn(key string, values ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(values) == 0 {
			return db
		}
		return db.Where(fmt.Sprintf("%s IN ?", key), values)
	}
}
