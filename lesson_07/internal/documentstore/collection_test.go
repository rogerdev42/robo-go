package documentstore

import (
	"reflect"
	"testing"
)

func TestCollection_Put(t *testing.T) {
	type fields struct {
		cfg       CollectionConfig
		documents map[string]Document
	}
	type args struct {
		doc Document
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{

		{
			name: "success",
			fields: fields{
				cfg:       CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"id": {Type: DocumentFieldTypeString, Value: "1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error_no_primary_key",
			fields: fields{
				cfg:       CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{},
				},
			},
			wantErr: true,
		},
		{
			name: "error_invalid_key_type",
			fields: fields{
				cfg:       CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"id": {Type: DocumentFieldTypeNumber, Value: 1},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error_empty_key",
			fields: fields{
				cfg:       CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"id": {Type: DocumentFieldTypeString, Value: ""},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Collection{
				cfg:       tt.fields.cfg,
				documents: tt.fields.documents,
			}
			if err := s.Put(tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollection_Get(t *testing.T) {
	type fields struct {
		cfg       CollectionConfig
		documents map[string]Document
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Document
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				cfg: CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{
					"1": {
						Fields: map[string]DocumentField{
							"id": {Type: DocumentFieldTypeString, Value: "1"},
						},
					},
				},
			},
			args: args{
				key: "1",
			},
			want: &Document{
				Fields: map[string]DocumentField{
					"id": {Type: DocumentFieldTypeString, Value: "1"},
				},
			},
			wantErr: false,
		},
		{
			name: "error_not_found",
			fields: fields{
				cfg:       CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{},
			},
			args: args{
				key: "1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error_empty_key",
			fields: fields{
				cfg:       CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{},
			},
			args: args{
				key: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Collection{
				cfg:       tt.fields.cfg,
				documents: tt.fields.documents,
			}
			got, err := s.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_Delete(t *testing.T) {
	type fields struct {
		cfg       CollectionConfig
		documents map[string]Document
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				cfg: CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{
					"1": {
						Fields: map[string]DocumentField{
							"id": {Type: DocumentFieldTypeString, Value: "1"},
						},
					},
				},
			},
			args: args{
				key: "1",
			},
			wantErr: false,
		},
		{
			name: "error_not_found",
			fields: fields{
				cfg:       CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{},
			},
			args: args{
				key: "1",
			},
			wantErr: true,
		},
		{
			name: "error_empty_key",
			fields: fields{
				cfg:       CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{},
			},
			args: args{
				key: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Collection{
				cfg:       tt.fields.cfg,
				documents: tt.fields.documents,
			}
			if err := s.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollection_List(t *testing.T) {
	type fields struct {
		cfg       CollectionConfig
		documents map[string]Document
	}
	tests := []struct {
		name   string
		fields fields
		want   []Document
	}{
		{
			name: "empty",
			fields: fields{
				cfg:       CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{},
			},
			want: []Document{},
		},
		{
			name: "single_document",
			fields: fields{
				cfg: CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{
					"1": {
						Fields: map[string]DocumentField{
							"id": {Type: DocumentFieldTypeString, Value: "1"},
						},
					},
				},
			},
			want: []Document{
				{
					Fields: map[string]DocumentField{
						"id": {Type: DocumentFieldTypeString, Value: "1"},
					},
				},
			},
		},
		{
			name: "multiple_documents",
			fields: fields{
				cfg: CollectionConfig{PrimaryKey: "id"},
				documents: map[string]Document{
					"1": {
						Fields: map[string]DocumentField{
							"id": {Type: DocumentFieldTypeString, Value: "1"},
						},
					},
					"2": {
						Fields: map[string]DocumentField{
							"id": {Type: DocumentFieldTypeString, Value: "2"},
						},
					},
				},
			},
			want: []Document{
				{
					Fields: map[string]DocumentField{
						"id": {Type: DocumentFieldTypeString, Value: "1"},
					},
				},
				{
					Fields: map[string]DocumentField{
						"id": {Type: DocumentFieldTypeString, Value: "2"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Collection{
				cfg:       tt.fields.cfg,
				documents: tt.fields.documents,
			}
			if got := s.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}
