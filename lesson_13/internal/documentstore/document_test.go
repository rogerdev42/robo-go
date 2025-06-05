package documentstore

import (
	"reflect"
	"testing"
)

func TestMarshalDocument(t *testing.T) {
	type args struct {
		input any
	}
	tests := []struct {
		name    string
		args    args
		want    *Document
		wantErr bool
	}{
		{
			name: "string field",
			args: args{
				input: struct {
					Name string `json:"name"`
				}{
					Name: "test",
				},
			},
			want: &Document{
				Fields: map[string]DocumentField{
					"Name": {Type: DocumentFieldTypeString, Value: "test"},
				},
			},
			wantErr: false,
		},
		{
			name: "number field",
			args: args{
				input: struct {
					Age int `json:"age"`
				}{
					Age: 25,
				},
			},
			want: &Document{
				Fields: map[string]DocumentField{
					"Age": {Type: DocumentFieldTypeNumber, Value: 25},
				},
			},
			wantErr: false,
		},
		{
			name: "boolean field",
			args: args{
				input: struct {
					Active bool `json:"active"`
				}{
					Active: true,
				},
			},
			want: &Document{
				Fields: map[string]DocumentField{
					"Active": {Type: DocumentFieldTypeBool, Value: true},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalDocument(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalDocument() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalDocument(t *testing.T) {
	type testStruct struct {
		Value int
	}

	tests := []struct {
		name    string
		doc     *Document
		output  interface{}
		wantErr bool
		wantVal int
	}{
		{
			name: "valid number",
			doc: &Document{
				Fields: map[string]DocumentField{
					"Value": {
						Type:  DocumentFieldTypeNumber,
						Value: float64(42),
					},
				},
			},
			output:  &testStruct{},
			wantErr: false,
			wantVal: 42,
		},
		{
			name: "wrong type",
			doc: &Document{
				Fields: map[string]DocumentField{
					"Value": {
						Type:  DocumentFieldTypeString,
						Value: "not a number",
					},
				},
			},
			output:  &testStruct{},
			wantErr: true,
			wantVal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalDocument(tt.doc, tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				result := tt.output.(*testStruct)
				if result.Value != tt.wantVal {
					t.Errorf("UnmarshalDocument() got Value = %v, want %v", result.Value, tt.wantVal)
				}
			}
		})
	}
}
