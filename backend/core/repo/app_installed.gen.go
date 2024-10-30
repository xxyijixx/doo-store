// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package repo

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"doo-store/backend/core/model"
)

func newAppInstalled(db *gorm.DB, opts ...gen.DOOption) appInstalled {
	_appInstalled := appInstalled{}

	_appInstalled.appInstalledDo.UseDB(db, opts...)
	_appInstalled.appInstalledDo.UseModel(&model.AppInstalled{})

	tableName := _appInstalled.appInstalledDo.TableName()
	_appInstalled.ALL = field.NewAsterisk(tableName)
	_appInstalled.ID = field.NewInt64(tableName, "id")
	_appInstalled.CreatedAt = field.NewTime(tableName, "created_at")
	_appInstalled.UpdatedAt = field.NewTime(tableName, "updated_at")
	_appInstalled.Name = field.NewString(tableName, "name")
	_appInstalled.AppID = field.NewInt64(tableName, "app_id")
	_appInstalled.AppDetailID = field.NewInt64(tableName, "app_detail_id")
	_appInstalled.Version = field.NewString(tableName, "version")
	_appInstalled.Params = field.NewString(tableName, "params")
	_appInstalled.Env = field.NewString(tableName, "env")
	_appInstalled.DockerCompose = field.NewString(tableName, "docker_compose")
	_appInstalled.Status = field.NewString(tableName, "status")

	_appInstalled.fillFieldMap()

	return _appInstalled
}

type appInstalled struct {
	appInstalledDo

	ALL           field.Asterisk
	ID            field.Int64
	CreatedAt     field.Time
	UpdatedAt     field.Time
	Name          field.String
	AppID         field.Int64
	AppDetailID   field.Int64
	Version       field.String
	Params        field.String
	Env           field.String
	DockerCompose field.String
	Status        field.String

	fieldMap map[string]field.Expr
}

func (a appInstalled) Table(newTableName string) *appInstalled {
	a.appInstalledDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a appInstalled) As(alias string) *appInstalled {
	a.appInstalledDo.DO = *(a.appInstalledDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *appInstalled) updateTableName(table string) *appInstalled {
	a.ALL = field.NewAsterisk(table)
	a.ID = field.NewInt64(table, "id")
	a.CreatedAt = field.NewTime(table, "created_at")
	a.UpdatedAt = field.NewTime(table, "updated_at")
	a.Name = field.NewString(table, "name")
	a.AppID = field.NewInt64(table, "app_id")
	a.AppDetailID = field.NewInt64(table, "app_detail_id")
	a.Version = field.NewString(table, "version")
	a.Params = field.NewString(table, "params")
	a.Env = field.NewString(table, "env")
	a.DockerCompose = field.NewString(table, "docker_compose")
	a.Status = field.NewString(table, "status")

	a.fillFieldMap()

	return a
}

func (a *appInstalled) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *appInstalled) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 11)
	a.fieldMap["id"] = a.ID
	a.fieldMap["created_at"] = a.CreatedAt
	a.fieldMap["updated_at"] = a.UpdatedAt
	a.fieldMap["name"] = a.Name
	a.fieldMap["app_id"] = a.AppID
	a.fieldMap["app_detail_id"] = a.AppDetailID
	a.fieldMap["version"] = a.Version
	a.fieldMap["params"] = a.Params
	a.fieldMap["env"] = a.Env
	a.fieldMap["docker_compose"] = a.DockerCompose
	a.fieldMap["status"] = a.Status
}

func (a appInstalled) clone(db *gorm.DB) appInstalled {
	a.appInstalledDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a appInstalled) replaceDB(db *gorm.DB) appInstalled {
	a.appInstalledDo.ReplaceDB(db)
	return a
}

type appInstalledDo struct{ gen.DO }

type IAppInstalledDo interface {
	gen.SubQuery
	Debug() IAppInstalledDo
	WithContext(ctx context.Context) IAppInstalledDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IAppInstalledDo
	WriteDB() IAppInstalledDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IAppInstalledDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IAppInstalledDo
	Not(conds ...gen.Condition) IAppInstalledDo
	Or(conds ...gen.Condition) IAppInstalledDo
	Select(conds ...field.Expr) IAppInstalledDo
	Where(conds ...gen.Condition) IAppInstalledDo
	Order(conds ...field.Expr) IAppInstalledDo
	Distinct(cols ...field.Expr) IAppInstalledDo
	Omit(cols ...field.Expr) IAppInstalledDo
	Join(table schema.Tabler, on ...field.Expr) IAppInstalledDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IAppInstalledDo
	RightJoin(table schema.Tabler, on ...field.Expr) IAppInstalledDo
	Group(cols ...field.Expr) IAppInstalledDo
	Having(conds ...gen.Condition) IAppInstalledDo
	Limit(limit int) IAppInstalledDo
	Offset(offset int) IAppInstalledDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IAppInstalledDo
	Unscoped() IAppInstalledDo
	Create(values ...*model.AppInstalled) error
	CreateInBatches(values []*model.AppInstalled, batchSize int) error
	Save(values ...*model.AppInstalled) error
	First() (*model.AppInstalled, error)
	Take() (*model.AppInstalled, error)
	Last() (*model.AppInstalled, error)
	Find() ([]*model.AppInstalled, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.AppInstalled, err error)
	FindInBatches(result *[]*model.AppInstalled, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.AppInstalled) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IAppInstalledDo
	Assign(attrs ...field.AssignExpr) IAppInstalledDo
	Joins(fields ...field.RelationField) IAppInstalledDo
	Preload(fields ...field.RelationField) IAppInstalledDo
	FirstOrInit() (*model.AppInstalled, error)
	FirstOrCreate() (*model.AppInstalled, error)
	FindByPage(offset int, limit int) (result []*model.AppInstalled, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IAppInstalledDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (a appInstalledDo) Debug() IAppInstalledDo {
	return a.withDO(a.DO.Debug())
}

func (a appInstalledDo) WithContext(ctx context.Context) IAppInstalledDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a appInstalledDo) ReadDB() IAppInstalledDo {
	return a.Clauses(dbresolver.Read)
}

func (a appInstalledDo) WriteDB() IAppInstalledDo {
	return a.Clauses(dbresolver.Write)
}

func (a appInstalledDo) Session(config *gorm.Session) IAppInstalledDo {
	return a.withDO(a.DO.Session(config))
}

func (a appInstalledDo) Clauses(conds ...clause.Expression) IAppInstalledDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a appInstalledDo) Returning(value interface{}, columns ...string) IAppInstalledDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a appInstalledDo) Not(conds ...gen.Condition) IAppInstalledDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a appInstalledDo) Or(conds ...gen.Condition) IAppInstalledDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a appInstalledDo) Select(conds ...field.Expr) IAppInstalledDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a appInstalledDo) Where(conds ...gen.Condition) IAppInstalledDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a appInstalledDo) Order(conds ...field.Expr) IAppInstalledDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a appInstalledDo) Distinct(cols ...field.Expr) IAppInstalledDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a appInstalledDo) Omit(cols ...field.Expr) IAppInstalledDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a appInstalledDo) Join(table schema.Tabler, on ...field.Expr) IAppInstalledDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a appInstalledDo) LeftJoin(table schema.Tabler, on ...field.Expr) IAppInstalledDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a appInstalledDo) RightJoin(table schema.Tabler, on ...field.Expr) IAppInstalledDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a appInstalledDo) Group(cols ...field.Expr) IAppInstalledDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a appInstalledDo) Having(conds ...gen.Condition) IAppInstalledDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a appInstalledDo) Limit(limit int) IAppInstalledDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a appInstalledDo) Offset(offset int) IAppInstalledDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a appInstalledDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IAppInstalledDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a appInstalledDo) Unscoped() IAppInstalledDo {
	return a.withDO(a.DO.Unscoped())
}

func (a appInstalledDo) Create(values ...*model.AppInstalled) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a appInstalledDo) CreateInBatches(values []*model.AppInstalled, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a appInstalledDo) Save(values ...*model.AppInstalled) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a appInstalledDo) First() (*model.AppInstalled, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.AppInstalled), nil
	}
}

func (a appInstalledDo) Take() (*model.AppInstalled, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.AppInstalled), nil
	}
}

func (a appInstalledDo) Last() (*model.AppInstalled, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.AppInstalled), nil
	}
}

func (a appInstalledDo) Find() ([]*model.AppInstalled, error) {
	result, err := a.DO.Find()
	return result.([]*model.AppInstalled), err
}

func (a appInstalledDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.AppInstalled, err error) {
	buf := make([]*model.AppInstalled, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a appInstalledDo) FindInBatches(result *[]*model.AppInstalled, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a appInstalledDo) Attrs(attrs ...field.AssignExpr) IAppInstalledDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a appInstalledDo) Assign(attrs ...field.AssignExpr) IAppInstalledDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a appInstalledDo) Joins(fields ...field.RelationField) IAppInstalledDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a appInstalledDo) Preload(fields ...field.RelationField) IAppInstalledDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a appInstalledDo) FirstOrInit() (*model.AppInstalled, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.AppInstalled), nil
	}
}

func (a appInstalledDo) FirstOrCreate() (*model.AppInstalled, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.AppInstalled), nil
	}
}

func (a appInstalledDo) FindByPage(offset int, limit int) (result []*model.AppInstalled, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a appInstalledDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a appInstalledDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a appInstalledDo) Delete(models ...*model.AppInstalled) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *appInstalledDo) withDO(do gen.Dao) *appInstalledDo {
	a.DO = *do.(*gen.DO)
	return a
}
