package documentstore

import (
	"github.com/stretchr/testify/assert"
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

func setupTestCollection() *Collection {
	coll := &Collection{
		cfg: CollectionConfig{PrimaryKey: "id"},
		documents: map[string]Document{
			"1": {Fields: map[string]DocumentField{
				"id":    {Type: DocumentFieldTypeString, Value: "1"},
				"name":  {Type: DocumentFieldTypeString, Value: "Alice"},
				"email": {Type: DocumentFieldTypeString, Value: "alice@example.com"},
			}},
			"2": {Fields: map[string]DocumentField{
				"id":    {Type: DocumentFieldTypeString, Value: "2"},
				"name":  {Type: DocumentFieldTypeString, Value: "Bob"},
				"email": {Type: DocumentFieldTypeString, Value: "bob@example.com"},
			}},
			"3": {Fields: map[string]DocumentField{
				"id":    {Type: DocumentFieldTypeString, Value: "3"},
				"name":  {Type: DocumentFieldTypeString, Value: "Charlie"},
				"email": {Type: DocumentFieldTypeString, Value: "charlie@example.com"},
			}},
		},
	}
	_ = coll.CreateIndex("name")
	_ = coll.CreateIndex("email")
	return coll
}

func TestQuery_AscendFullRange(t *testing.T) {
	coll := setupTestCollection()
	docs, err := coll.Query("name", QueryParams{})
	assert.NoError(t, err)
	gotNames := []string{}
	for _, d := range docs {
		gotNames = append(gotNames, d.Fields["name"].Value.(string))
	}
	assert.Equal(t, []string{"Alice", "Bob", "Charlie"}, gotNames)
}

func TestQuery_DescendFullRange(t *testing.T) {
	coll := setupTestCollection()
	docs, err := coll.Query("name", QueryParams{Desc: true})
	assert.NoError(t, err)
	gotNames := []string{}
	for _, d := range docs {
		gotNames = append(gotNames, d.Fields["name"].Value.(string))
	}
	assert.Equal(t, []string{"Charlie", "Bob", "Alice"}, gotNames)
}

func TestQuery_MinValue(t *testing.T) {
	coll := setupTestCollection()
	minValue := "Bob"
	docs, err := coll.Query("name", QueryParams{MinValue: &minValue})
	assert.NoError(t, err)
	gotNames := []string{}
	for _, d := range docs {
		gotNames = append(gotNames, d.Fields["name"].Value.(string))
	}
	assert.Equal(t, []string{"Bob", "Charlie"}, gotNames)
}

func TestQuery_MaxValue(t *testing.T) {
	coll := setupTestCollection()
	maxValue := "Bob"
	docs, err := coll.Query("name", QueryParams{MaxValue: &maxValue})
	assert.NoError(t, err)
	gotNames := []string{}
	for _, d := range docs {
		gotNames = append(gotNames, d.Fields["name"].Value.(string))
	}
	assert.Equal(t, []string{"Alice", "Bob"}, gotNames)
}

func TestQuery_MinAndMaxValue(t *testing.T) {
	coll := setupTestCollection()
	minValue := "Bob"
	maxValue := "Charlie"
	docs, err := coll.Query("name", QueryParams{MinValue: &minValue, MaxValue: &maxValue})
	assert.NoError(t, err)
	gotNames := []string{}
	for _, d := range docs {
		gotNames = append(gotNames, d.Fields["name"].Value.(string))
	}
	assert.Equal(t, []string{"Bob", "Charlie"}, gotNames)
}

func TestQuery_NoIndex(t *testing.T) {
	coll := setupTestCollection()
	_ = coll.DeleteIndex("name")
	_, err := coll.Query("name", QueryParams{})
	assert.ErrorIs(t, err, ErrIndexNotFound)
}

func TestQuery_EmptyResult(t *testing.T) {
	coll := setupTestCollection()
	minValue := "D"
	docs, err := coll.Query("name", QueryParams{MinValue: &minValue})
	assert.NoError(t, err)
	assert.Empty(t, docs)
}

func TestCreateIndex_Success(t *testing.T) {
	coll := setupTestCollection()
	// Попробуем создать новый индекс по несуществующему полю
	err := coll.CreateIndex("nonexistent")
	assert.NoError(t, err)
	tree, ok := coll.indexes["nonexistent"]
	assert.True(t, ok)
	count := 0
	tree.Ascend(func(item *indexItem) bool {
		count++
		return true
	})
	assert.Equal(t, 0, count)
}

func TestCreateIndex_AlreadyExists(t *testing.T) {
	coll := setupTestCollection()
	err := coll.CreateIndex("name")
	assert.ErrorIs(t, err, ErrIndexExists)
}

func TestCreateIndex_NonStringField(t *testing.T) {
	coll := setupTestCollection()
	coll.documents["4"] = Document{Fields: map[string]DocumentField{
		"id":  {Type: DocumentFieldTypeString, Value: "4"},
		"age": {Type: DocumentFieldTypeNumber, Value: 42},
	}}
	err := coll.CreateIndex("age")
	assert.NoError(t, err)
	tree, ok := coll.indexes["age"]
	assert.True(t, ok)
	count := 0
	tree.Ascend(func(item *indexItem) bool {
		count++
		return true
	})
	assert.Equal(t, 0, count)
}

func TestDeleteIndex_Success(t *testing.T) {
	coll := setupTestCollection()
	err := coll.DeleteIndex("name")
	assert.NoError(t, err)
	_, ok := coll.indexes["name"]
	assert.False(t, ok)
}

func TestDeleteIndex_NotFound(t *testing.T) {
	coll := setupTestCollection()
	err := coll.DeleteIndex("notfound")
	assert.ErrorIs(t, err, ErrIndexNotFound)
}
