package orm

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/go-chocolate/chocolate/pkg/database/repository"
)

type GormRepository[T any] struct {
	db *gorm.DB
}

func (rep *GormRepository[T]) GetDB(ctx context.Context) *gorm.DB {
	return rep.db.WithContext(ctx)
}

func (rep *GormRepository[T]) where(ctx context.Context, where any) *gorm.DB {
	db := rep.GetDB(ctx).Model(new(T))
	switch condition := where.(type) {
	case int64:
		return db.Where("id = ?", condition)
	case clause.Expression:
		return db.Clauses(condition)
	case []clause.Expression:
		return db.Clauses(condition...)
	default:
		return db.Where(where)
	}
}

func (rep *GormRepository[T]) FindOne(ctx context.Context, where any) (one *T, err error) {
	one = new(T)
	err = rep.where(ctx, where).Take(one).Error
	return
}

func (rep *GormRepository[T]) List(ctx context.Context, where any, offset, limit int, order ...any) (list []*T, count int64, err error) {
	query := rep.where(ctx, where)
	if err = query.Count(&count).Error; err != nil {
		return
	}
	list = make([]*T, 0)
	if count == 0 {
		return
	}
	if len(order) > 0 {
		query = query.Order(order[0])
	}
	if err = query.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return
	}
	return
}

func (rep *GormRepository[T]) Count(ctx context.Context, where any) (count int64, err error) {
	err = rep.where(ctx, where).Count(&count).Error
	return
}

func (rep *GormRepository[T]) Update(ctx context.Context, where any, update any) (rowsAffected int64, err error) {
	db := rep.where(ctx, where).Updates(update)
	rowsAffected = db.RowsAffected
	err = db.Error
	return
}

func (rep *GormRepository[T]) Delete(ctx context.Context, where any) (rowsAffected int64, err error) {
	var model T
	db := rep.where(ctx, where).Delete(&model)
	rowsAffected = db.RowsAffected
	err = db.Error
	return
}

func (rep *GormRepository[T]) Insert(ctx context.Context, data any) (rowsAffected int64, err error) {
	db := rep.GetDB(ctx).Create(data)
	rowsAffected = db.RowsAffected
	err = db.Error
	return
}

func (rep *GormRepository[T]) Iterate(ctx context.Context, column string, where any) (repository.Iterator[T], error) {
	return &iterator[T]{column: column, db: rep.GetDB(ctx), where: where}, nil
}

type iterator[T any] struct {
	lastID       any
	lastIDColumn string
	column       string
	db           *gorm.DB
	where        any
}

func (i *iterator[T]) Count() (count int64, err error) {
	query := i.db
	if i.where != nil {
		switch cond := i.where.(type) {
		case clause.Expression:
			query = query.Clauses(cond)
		default:
			query = query.Where(cond)
		}
	}
	err = query.Count(&count).Error
	return
}

var ErrIterationEoF = errors.New("iteration end of file")

func (i *iterator[T]) Next() ([]*T, error) {
	query := i.db
	if i.where != nil {
		switch cond := i.where.(type) {
		case clause.Expression:
			if i.lastID != nil {
				query = query.Clauses(clause.And(cond, clause.Gt{Column: "`" + i.column + "`", Value: i.lastID}))
			} else {
				query = query.Clauses(cond)
			}
		default:
			query = query.Where(cond)
			if i.lastID != nil {
				query = query.Where("`"+i.column+"` > ?", i.lastID)
			}
		}
	} else if i.lastID != nil {
		query = query.Where("`"+i.column+"` > ?", i.lastID)
	}
	query = query.Order("`" + i.column + "`")

	var result []*T
	err := query.Find(&result).Error
	if err != nil {
		return result, err
	}
	if len(result) == 0 {
		return nil, ErrIterationEoF
	}
	return result, i.extractLastID(result[len(result)-1])
}

func (i *iterator[T]) extractLastID(item any) error {
	val := reflect.ValueOf(item)

	if i.lastIDColumn == "" {
		typ := reflect.TypeOf(item)
		for n := 0; n < val.NumField(); n++ {
			field := typ.Field(n)
			//decode gorm tag and find column name
			if tag := field.Tag.Get("gorm"); tag != "" {
				var column string
				for _, v := range strings.Split(tag, ";") {
					if len(v) > 7 && strings.ToLower(v[:7]) == "column:" {
						column = v[7:]
						break
					}
				}
				if column != "" && column == i.column {
					i.lastIDColumn = field.Name
					break
				}
			}

			// decode table column name to struct field name
			var column []byte
			var toUpper = true //make first character to upper
			for _, v := range column {
				if v == '_' { // skip underline and make next character to upper
					toUpper = true
					continue
				}
				if toUpper && (v >= 'a' && v <= 'z') {
					v = v - 32
				}
				column = append(column, v)
			}

			if string(column) == field.Name {
				i.lastIDColumn = field.Name
				break
			}
		}
	}
	if i.lastIDColumn == "" {
		typ := reflect.TypeOf(item)
		return fmt.Errorf("struct %s does not contain field %s", typ.String(), i.column)
	}
	i.lastID = val.FieldByName(i.lastIDColumn).Interface()
	return nil
}
