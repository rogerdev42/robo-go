package documentstore

import (
	"fmt"
	"reflect"
)

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

func MarshalDocument(input any) (*Document, error) {
	document := Document{Fields: make(map[string]DocumentField)}

	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		key := t.Field(i).Name
		fieldType := t.Field(i).Type
		fieldValue := v.Field(i).Interface()

		var docFieldType DocumentFieldType
		switch fieldType.Kind() {
		case reflect.String:
			docFieldType = DocumentFieldTypeString
		case reflect.Int, reflect.Float32, reflect.Float64:
			docFieldType = DocumentFieldTypeNumber
		case reflect.Bool:
			docFieldType = DocumentFieldTypeBool
		case reflect.Slice, reflect.Array:
			docFieldType = DocumentFieldTypeArray
		case reflect.Struct, reflect.Map:
			docFieldType = DocumentFieldTypeObject
		default:
			return nil, fmt.Errorf("unsupported field type: %s", fieldType.Kind())
		}

		document.Fields[key] = DocumentField{
			Type:  docFieldType,
			Value: fieldValue,
		}
	}

	return &document, nil
}

func UnmarshalDocument(doc *Document, output any) error {
	v := reflect.ValueOf(output)
	t := v.Elem().Type()

	for i := 0; i < v.Elem().NumField(); i++ {
		field := v.Elem().Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name

		docField, ok := doc.Fields[fieldName]
		if !ok {
			continue
		}

		switch docField.Type {
		case DocumentFieldTypeString:
			if field.Kind() == reflect.String {
				field.SetString(docField.Value.(string))
			} else {
				return fmt.Errorf("field %s type mismatch: expected string", fieldName)
			}
		case DocumentFieldTypeNumber:
			switch field.Kind() {
			case reflect.Int:
				field.SetInt(int64(docField.Value.(float64)))
			case reflect.Float32, reflect.Float64:
				field.SetFloat(docField.Value.(float64))
			default:
				return fmt.Errorf("field %s type mismatch: expected number", fieldName)
			}
		case DocumentFieldTypeBool:
			if field.Kind() == reflect.Bool {
				field.SetBool(docField.Value.(bool))
			} else {
				return fmt.Errorf("field %s type mismatch: expected bool", fieldName)
			}
		case DocumentFieldTypeArray:
			if field.Kind() == reflect.Slice {
				values, ok := docField.Value.([]any)
				if !ok {
					return fmt.Errorf("field %s contains invalid array value", fieldName)
				}

				slice := reflect.MakeSlice(field.Type(), len(values), len(values))
				for i, v := range values {
					slice.Index(i).Set(reflect.ValueOf(v))
				}
				field.Set(slice)
			} else {
				return fmt.Errorf("field %s type mismatch: expected array", fieldName)
			}
		case DocumentFieldTypeObject:
			if field.Kind() == reflect.Struct {
				// Рекурсивная обработка вложенных структур
				nestedDoc := docField.Value.(Document)
				err := UnmarshalDocument(&nestedDoc, field.Addr().Interface())
				if err != nil {
					return fmt.Errorf("failed to unmarshal nested object for field %s: %w", fieldName, err)
				}
			} else {
				return fmt.Errorf("field %s type mismatch: expected object", fieldName)
			}
		default:
			return fmt.Errorf("unsupported field type: %s", docField.Type)
		}
	}

	return nil
}
