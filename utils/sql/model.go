package sql

import (
	"encoding/json"
	"reflect"

	form "movies/utils/form"
	pg "movies/utils/pg"

	pgtype "github.com/jackc/pgtype"
)

type Table interface {
	TableName() string
}

/*============================================================================*/
/*=====*                              Base                              *=====*/
/*============================================================================*/

type Model struct {
	ID pgtype.UUID `json:"id" db:"id"`
}

// Conform struct fields
func (obj *Model) StructConform(T Table) error { return form.ConformStruct(T) }

// Validate struct fields
func (obj *Model) StructValidate(T Table) error { return form.ValidateStruct(T, false) }

// Conform and Validate struct fields
func (obj *Model) StructPrepare(T Table) error {
	if err := obj.StructConform(T); err != nil {
		return err
	}
	return obj.StructValidate(T)
}

func (obj Model) GetPK() pgtype.UUID { return obj.ID }

/*============================================================================*/
/*=====*                            Extended                            *=====*/
/*============================================================================*/

type Extended struct {
	Model
	Created
	Updated
	Deleted
}

// Conform struct fields
func (obj *Extended) StructConform(T Table) error {
	obj.Created.StructConform()
	obj.Updated.StructConform()
	obj.Deleted.StructConform()
	return obj.Model.StructConform(T)
}

// Conform and Validate struct fields
func (obj *Extended) StructPrepare(T Table) error {
	if err := obj.StructConform(T); err != nil {
		return err
	}
	return obj.Model.StructValidate(T)
}

/*============================================================================*/
/*=====*                            Created                             *=====*/
/*============================================================================*/

type Created struct {
	CreatedAt pgtype.Timestamptz `json:"created_at" db:"created_at"`
	CreatedBy pgtype.UUID        `json:"created_by" db:"created_by"`
}

func (obj *Created) StructConform() {
	form.RemoveUndefined(&obj.CreatedAt.Status)
	form.RemoveUndefined(&obj.CreatedBy.Status)
}

func (obj *Created) SetCreatedID(collaboratorID pgtype.UUID) {
	obj.CreatedBy = collaboratorID
}

/*============================================================================*/
/*=====*                            Updated                             *=====*/
/*============================================================================*/

type Updated struct {
	UpdatedAt pgtype.Timestamptz `json:"updated_at" db:"updated_at"`
	UpdatedBy pgtype.UUID        `json:"updated_by" db:"updated_by"`
}

func (obj Updated) GetUpdated() Updated { return obj }

func (obj *Updated) StructConform() {
	form.RemoveUndefined(&obj.UpdatedAt.Status)
	form.RemoveUndefined(&obj.UpdatedBy.Status)
}

func (obj *Updated) SetUpdatedID(collaboratorID pgtype.UUID) {
	obj.UpdatedBy = collaboratorID
}

/*============================================================================*/
/*=====*                            Deleted                             *=====*/
/*============================================================================*/

type Deleted struct {
	DeletedAt pgtype.Timestamptz `json:"deleted_at" db:"deleted_at"`
	DeletedBy pgtype.UUID        `json:"deleted_by" db:"deleted_by"`
}

func (obj Deleted) GetDeleted() Deleted { return obj }

func (obj *Deleted) StructConform() {
	form.RemoveUndefined(&obj.DeletedAt.Status)
	form.RemoveUndefined(&obj.DeletedBy.Status)
}

func (obj *Deleted) SetDeletedID(collaboratorID pgtype.UUID) {
	obj.DeletedBy = collaboratorID
}

/*============================================================================*/
/*=====*                           Archived                             *=====*/
/*============================================================================*/

type Archived struct {
	ArchivedAt pgtype.Timestamptz `json:"archived_at" db:"archived_at"`
	ArchivedBy pgtype.UUID        `json:"archived_by" db:"archived_by"`
}

func (obj Archived) GetArchived() Archived { return obj }

func (obj *Archived) StructConform() {
	form.RemoveUndefined(&obj.ArchivedAt.Status)
	form.RemoveUndefined(&obj.ArchivedBy.Status)
}

func (obj *Archived) SetArchivedID(collaboratorID pgtype.UUID) {
	obj.ArchivedBy = collaboratorID
}

/*============================================================================*/
/*=====*                            Setting                             *=====*/
/*============================================================================*/

type Setting struct {
	Settings pgtype.JSONB `json:"settings" db:"settings"`
}

func (obj *Setting) SetSetting(key string, value any) error {
	settings, err := obj.GetSettings()
	if err == nil {
		if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
			delete(settings, key)
		} else {
			settings[key] = value
		}
		return obj.Settings.Set(settings)
	}
	return err
}

func (obj Setting) GetSettings() (map[string]any, error) {
	if obj.Settings.Status != pgtype.Present {
		return map[string]any{}, nil
	}
	data := make(map[string]any, 0)
	return data, json.Unmarshal(obj.Settings.Bytes, &data)
}

func (obj Setting) HaveSetting(key string) (bool, error) {
	if obj.Settings.Status != pgtype.Present {
		return false, nil
	}
	settings, err := obj.GetSettings()
	if err != nil {
		return false, err
	}
	_, ok := settings[key]
	return ok, nil
}

func (obj Setting) GetSetting(key string, empty any) (any, error) {
	settings, err := obj.GetSettings()
	if err != nil {
		return empty, err
	} else if value, ok := settings[key]; ok {
		return value, nil
	}
	return empty, nil
}

func (obj Setting) GetSettingBool(key string, empty bool) (bool, error) {
	value, err := obj.GetSetting(key, empty)
	return value.(bool), err
}

func (obj Setting) GetSettingInt(key string, empty int) (int, error) {
	value, err := obj.GetSetting(key, empty)
	if v, ok := value.(int); ok {
		return v, nil
	} else if v, ok := value.(float64); ok {
		return int(v), nil
	} else if v, ok := value.(float32); ok {
		return int(v), nil
	}
	return 0, err
}

func (obj Setting) GetSettingFloat(key string, empty float64) (float64, error) {
	value, err := obj.GetSetting(key, empty)
	return value.(float64), err
}

func (obj Setting) GetSettingUUID(key string, empty pgtype.UUID) (pgtype.UUID, error) {
	value, err := obj.GetSetting(key, empty)
	if err == nil {
		return pg.ParseUUID(value.(string))
	}
	return value.(pgtype.UUID), err
}
