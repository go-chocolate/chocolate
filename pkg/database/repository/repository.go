package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

type Repository[T any] interface {
	FindOne(ctx context.Context, where any) (*T, error)
	List(ctx context.Context, where any, offset, limit int, order ...any) ([]*T, int64, error)
	Count(ctx context.Context, where any) (int64, error)
	Update(ctx context.Context, where any, update any) (int64, error)
	Delete(ctx context.Context, where any) (int64, error)
	Insert(ctx context.Context, data any) (int64, error)
	Iterate(ctx context.Context, column string, where any, orders ...any) (Iterator[T], error)
}

var ErrUnimplemented = errors.New("not implemented")

type UnimplementedRepository[T any] struct{}

func (rep *UnimplementedRepository[T]) FindOne(ctx context.Context, where any) (*T, error) {
	return nil, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) List(ctx context.Context, where any, offset, limit int, order ...any) ([]*T, int64, error) {
	return nil, 0, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) Count(ctx context.Context, where any) (int64, error) {
	return 0, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) Update(ctx context.Context, where any, update any) (int64, error) {
	return 0, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) Delete(ctx context.Context, where any) (int64, error) {
	return 0, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) Insert(ctx context.Context, data any) (int64, error) {
	return 0, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) Iterate(ctx context.Context, column string, where any, orders ...any) (Iterator[T], error) {
	return nil, ErrUnimplemented
}

type SimpleRepository[T any] struct {
	db *gorm.DB
}

func NewSimpleRepository[T any](db *gorm.DB) *SimpleRepository[T] {
	return &SimpleRepository[T]{db: db}
}

func (rep *SimpleRepository[T]) SetDB(db *gorm.DB) *SimpleRepository[T] {
	rep.db = db
	return rep
}

func (rep *SimpleRepository[T]) FindOne(ctx context.Context, where any) (*T, error) {
	query := Where(rep.db, where)
	var one = new(T)
	err := query.Take(one).Error
	if err != nil {
		return nil, err
	}
	return one, err
}

func (rep *SimpleRepository[T]) List(ctx context.Context, where any, offset, limit int, orders ...any) ([]*T, int64, error) {
	query := rep.db.Model(new(T))
	query = Where(query, where)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	for _, order := range orders {
		query = query.Order(order)
	}
	query = query.Offset(offset).Limit(limit)
	var results []*T
	if err := query.Find(&results).Error; err != nil {
		return nil, 0, err
	}
	return results, count, nil
}

func (rep *SimpleRepository[T]) Count(ctx context.Context, where any) (int64, error) {
	query := rep.db.Model(new(T))
	query = Where(query, where)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (rep *SimpleRepository[T]) Update(ctx context.Context, where any, update any) (int64, error) {
	query := rep.db.Model(new(T))
	query = Where(query, where)
	query = query.Updates(update)
	return query.RowsAffected, query.Error
}

func (rep *SimpleRepository[T]) Delete(ctx context.Context, where any) (int64, error) {
	query := Where(rep.db, where)
	query.Delete(new(T))
	return query.RowsAffected, query.Error
}

func (rep *SimpleRepository[T]) Insert(ctx context.Context, data any) (int64, error) {
	exp := rep.db.Model(new(T))
	switch reflect.TypeOf(data).Kind() {
	case reflect.Array, reflect.Slice:
		exp = exp.CreateInBatches(data, 100)
	default:
		exp = exp.Create(data)
	}
	return exp.RowsAffected, exp.Error
}

func (rep *SimpleRepository[T]) Iterate(ctx context.Context, column string, where any, orders ...any) (Iterator[T], error) {
	return NewOffsetIterator[T](rep.db, where, 1000, orders...)
}

// Where
// 根据传参类型自动构建where语句
// 传参类型为 int64 时直接查询 id = ?
// 传参类型为 clause.Expression 或 []clause.Expression 时，使用 db.Clause 查询
// 其他类型使用 db.Where 查询
func Where(query *gorm.DB, where any) *gorm.DB {
	if where == nil {
		return query
	}
	if exp, ok := where.(clause.Expression); ok {
		query = query.Clauses(exp)
	} else if exps, ok := where.([]clause.Expression); ok {
		query = query.Clauses(exps...)
	} else if val, ok := where.(int64); ok {
		query = query.Where("id = ?", val)
	} else {
		query = query.Where(where)
	}
	return query
}
