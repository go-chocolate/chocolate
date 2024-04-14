package repository

import (
	"errors"
	"gorm.io/gorm"
)

var ErrIterationEOF = errors.New("iteration eof")

type Iterator[T any] interface {
	Count() (int64, error)
	Next() ([]*T, error)
}

type OffsetIterator[T any] struct {
	limit  int
	offset int
	count  int64
	query  *gorm.DB
}

func NewOffsetIterator[T any](db *gorm.DB, where any, limit int, orders ...any) (Iterator[T], error) {
	iter := &OffsetIterator[T]{}
	iter.limit = limit
	iter.query = Where(db.Model(new(T)), where)
	if err := iter.query.Count(&iter.count).Error; err != nil {
		return nil, err
	}
	for _, order := range orders {
		iter.query = iter.query.Order(order)
	}
	return iter, nil
}

func (iter *OffsetIterator[T]) Count() (int64, error) {
	return iter.count, nil
}

func (iter *OffsetIterator[T]) Next() ([]*T, error) {
	defer func() { iter.offset += iter.limit }()
	var data []*T
	err := iter.query.Offset(iter.offset).Limit(iter.limit).Find(&data).Error
	if err != nil {
		return nil, err
	}
	if len(data) == 0 || len(data) < iter.limit {
		return data, ErrIterationEOF
	}
	return data, nil
}

func (iter *OffsetIterator[T]) Reset() {
	iter.offset = 0
}
