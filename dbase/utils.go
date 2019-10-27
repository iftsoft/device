package dbase

import (
	"database/sql"
	"github.com/iftsoft/device/core"
	"reflect"
	"time"
)


const (
	strErrMakeBuilder	= "Can't create query builder"
	strErrNullDaoPtr	= "Invalid pointer to DAO object"
	strErrNullBldPtr	= "Invalid pointer to SQL builder"
	strErrNullDataPtr	= "Attempt to load into an invalid pointer"
	strErrDataIsSlice	= "Unit target is slice"
	strErrDataIsUnit	= "List target is not slice"
)

// Fetch columns data from selected record and fill object structure
func fetchReturnedRow(rows *sql.Rows, m map[string]reflect.Value) (int64, error) {
	defer rows.Close()

	column, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	// Make column data holders
	colsNum := len(column)
	refs := getDummyColumnValues(colsNum)

	var count int64 = 0
	if rows.Next() {
		// Scan column data from record
		err = rows.Scan(refs...)
		if err != nil {
			return 0, err
		}
		// Fill value map with column data
		err = iterateColumnValues(refs, column, m)
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

// Fetch columns data from selected record and fill object structure
func fetchSelectedRow(rows *sql.Rows, value interface{}) (int64, error) {
	defer rows.Close()

	column, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		panic(strErrNullDataPtr)
	}
	obj := v.Elem()
	if obj.Kind() == reflect.Slice {
		panic(strErrDataIsSlice)
	}

	// Make column data holders
	colsNum := len(column)
	refs := getDummyColumnValues(colsNum)
	m := DefineColumnValues(obj)

	var count int64 = 0
	if rows.Next() {
		// Scan column data from record
		err = rows.Scan(refs...)
		if err != nil {
			return 0, err
		}
		// Fill Object with column data
		err = iterateColumnValues(refs, column, m)
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

// Fetch columns data from record set and fill objects slice
func fetchSearchedRows(rows *sql.Rows, value interface{}) (int64, error) {
	defer rows.Close()

	column, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		panic(strErrNullDataPtr)
	}
	list := v.Elem()
	if list.Kind() != reflect.Slice {
		panic(strErrDataIsUnit)
	}
	var count int64 = 0
	colsNum := len(column)

	for rows.Next() {
		// Make column data holders and scan data from current row
		refs := getDummyColumnValues(colsNum)
		err = rows.Scan(refs...)
		if err != nil {
			return 0, err
		}
		// Create new Object for slice
		elem := reflect.New(list.Type().Elem()).Elem()
		m := DefineColumnValues(elem)
		// Fill Object with column data
		err = iterateColumnValues(refs, column, m)
		if err != nil {
			return 0, err
		}
		count++
		list.Set(reflect.Append(list, elem))
	}
	return count, nil
}

// Return slice of empty interfaces for column data holders
func getDummyColumnValues(colsNum int) []interface{} {
	refs := make([]interface{}, colsNum)
	for i := range refs {
		var ref interface{}
		refs[i] = &ref
	}
	return refs
}


func DefineColumnValues(value reflect.Value) map[string]reflect.Value {
	// Init and fill map of Value holders
	m := make(map[string]reflect.Value)
	iterateObjectValues(value, m)
	return m
}

// Recursive iterate through object structure and fill column map with value holders
func iterateObjectValues(value reflect.Value, m map[string]reflect.Value) {
	switch value.Kind() {
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		if !value.IsNil() {
			// Recursive call for pointer to struct
			iterateObjectValues(value.Elem(), m)
		}
	case reflect.Struct:
		t := value.Type()
		// Iterate through struct fields
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" && !field.Anonymous {
				// Skip unexported field
				continue
			}
			if field.Anonymous && field.Type.Kind() == reflect.Struct {
				// Recursive call for anonymous include struct
				iterateObjectValues(value.Field(i), m)
				continue
			}
			colName := core.FormatSnakeString(field.Name)
			tag := field.Tag.Get("col")
			if len(tag) > 0 {
				colName = tag
			}
			// Add struct field to column map
			if _, ok := m[colName]; !ok {
				m[colName] = value.Field(i)
			}
		}
	}
}


// Fill Object structure with column data
func iterateColumnValues(refs []interface{}, column []string, m map[string]reflect.Value) error {
	// Iterate through columns in record
	for i, col := range column {
		if dst, ok := m[col]; ok {
			// Get column data interface
			src := reflect.Indirect(reflect.ValueOf(refs[i]))
			// Fill field value in object
			err := setUnitColumnValue(src, dst)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Copy data from record column to object field
func setUnitColumnValue(val reflect.Value, dst reflect.Value) error {
	dstKind := dst.Kind()
	src := val.Interface()

	var fld reflect.Value
	// Find or create actual field
	if dstKind == reflect.Ptr {
		if src == nil {
			dst.Set(reflect.Zero(dst.Type()))
			return nil
		}
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		fld = dst.Elem()
	} else {
		fld = dst
	}

	var str *Value
	// Check for string type of source data
	switch v := src.(type) {
	case []byte:
		s := Value(string(v))
		str = &s
	case string:
		s := Value(v)
		str = &s
	}

	fldKind := fld.Kind()
	// Check field type
	switch fldKind {
	// String
	case reflect.String:
		if str != nil {
			fld.SetString(str.String())
		} else {
			fld.SetString(ToValue(src))
		}
	// Boolean
	case reflect.Bool:
		if str != nil {
			b, err := str.Bool()
			if err != nil {
				return err
			}
			fld.SetBool(b)
		} else {
			fld.SetBool(src.(bool))
		}
	// Integer
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if str != nil {
			i, err := str.Int64()
			if err != nil {
				return err
			}
			fld.SetInt(i)
		} else {
			fld.SetInt(src.(int64))
		}
	// Unsigned
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if str != nil {
			u, err := str.Uint64()
			if err != nil {
				return err
			}
			fld.SetUint(u)
		} else {
			fld.SetUint(src.(uint64))
		}
	// Float
	case reflect.Float32, reflect.Float64:
		if str != nil {
			v, err := str.Float64()
			if err != nil {
				return err
			}
			fld.SetFloat(v)
		} else {
			fld.SetFloat(src.(float64))
		}
	// Structure
	case reflect.Struct:
		switch fld.Interface().(type) {
		case time.Time:
			if tm, ok := src.(time.Time); ok {
				fld.Set(reflect.ValueOf(tm))
			} else {
				fld.Set(reflect.Zero(fld.Type()))
			}
		// Undefined structure
		default:
			fld.Set(reflect.Zero(fld.Type()))
		}
	// Slice or map
	case reflect.Slice, reflect.Map :
		fld.Set(reflect.Zero(fld.Type()))
	}
	return nil
}


