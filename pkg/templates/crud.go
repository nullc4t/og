package templates

var CRUDTemplate = `package {{ .Package }}

import (
	"context"
	"github.com/nullc4t/gorm-cruder/crud"
	"gorm.io/gorm"
)

type CRUD interface {
	Create(ctx context.Context, v {{ .Type }}, omit ...string) (*{{ .Type }}, error)
	GetOrCreate(ctx context.Context, v {{ .Type }}, omit ...string) (*{{ .Type }}, error)
	GetByID(ctx context.Context, v {{ .Type }}) (*{{ .Type }}, error)
	Query(ctx context.Context, v {{ .Type }}, omit ...string) ([]*{{ .Type }}, error)
	QueryOne(ctx context.Context, v {{ .Type }}, omit ...string) (*{{ .Type }}, error)
	UpdateField(ctx context.Context, v {{ .Type }}, column string, value any) error
	Update(ctx context.Context, v {{ .Type }}, omit ...string) (err error)
	UpdateMap(ctx context.Context, v map[string]any) error
	Delete(ctx context.Context, v {{ .Type }}) error
}

type impl struct {
	db   *gorm.DB
	crud crud.GenericCRUD[{{ .Type }}]
}

func NewCRUD(db *gorm.DB) CRUD {
	return impl{db, crud.New[{{ .Type }}](db)}
}

func (i impl) Create(ctx context.Context, v {{ .Type }}, omit ...string) (*{{ .Type }}, error) {
	return i.crud.Create(ctx, v, omit...)
}

func (i impl) GetOrCreate(ctx context.Context, v {{ .Type }}, omit ...string) (*{{ .Type }}, error) {
	return i.crud.GetOrCreate(ctx, v, omit...)
}

func (i impl) GetByID(ctx context.Context, v {{ .Type }}) (*{{ .Type }}, error) {
	return i.crud.GetByID(ctx, v)
}

func (i impl) Query(ctx context.Context, v {{ .Type }}, omit ...string) ([]*{{ .Type }}, error) {
	return i.crud.Query(ctx, v, omit...)
}

func (i impl) QueryOne(ctx context.Context, v {{ .Type }}, omit ...string) (*{{ .Type }}, error) {
	return i.crud.QueryOne(ctx, v, omit...)
}

func (i impl) UpdateField(ctx context.Context, v {{ .Type }}, column string, value any) error {
	return i.crud.UpdateField(ctx, v, column, value)
}

func (i impl) Update(ctx context.Context, v {{ .Type }}, omit ...string) (err error) {
	return i.crud.Update(ctx, v, omit...)
}

func (i impl) UpdateMap(ctx context.Context, v map[string]any) error {
	return i.crud.UpdateMap(ctx, v)
}

func (i impl) Delete(ctx context.Context, v {{ .Type }}) error {
	return i.crud.Delete(ctx, v)
}`