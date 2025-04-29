package documentstore

import (
	"lesson_09/pkg/bst"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewStore(t *testing.T) {
	tests := []struct {
		name string
		want *Store
	}{
		{
			name: "successfully create new store",
			want: &Store{
				collections: make(map[string]*Collection),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_CreateCollection(t *testing.T) {
	type fields struct {
		collections map[string]*Collection
	}
	type args struct {
		name string
		cfg  *CollectionConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Collection
		wantErr bool
	}{
		{
			name: "successfully create new collection",
			fields: fields{
				collections: make(map[string]*Collection),
			},
			args: args{
				name: "test",
				cfg: &CollectionConfig{
					PrimaryKey: "id",
				},
			},
			want: &Collection{
				cfg: CollectionConfig{
					PrimaryKey: "id",
				},
				documents: make(map[string]Document),
				indexes:   make(map[string]*bst.BinarySearchTree),
			},
			wantErr: false,
		},
		{
			name: "error when collection already exists",
			fields: fields{
				collections: map[string]*Collection{
					"test": {
						cfg: CollectionConfig{
							PrimaryKey: "id",
						},
						documents: make(map[string]Document),
					},
				},
			},
			args: args{
				name: "test",
				cfg: &CollectionConfig{
					PrimaryKey: "id",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				collections: tt.fields.collections,
			}
			got, err := s.CreateCollection(tt.args.name, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateCollection() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_GetCollection(t *testing.T) {
	type fields struct {
		collections map[string]*Collection
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Collection
		wantErr bool
	}{
		{
			name: "successfully get existing collection",
			fields: fields{
				collections: map[string]*Collection{
					"test": {
						cfg: CollectionConfig{
							PrimaryKey: "id",
						},
						documents: make(map[string]Document),
					},
				},
			},
			args: args{
				name: "test",
			},
			want: &Collection{
				cfg: CollectionConfig{
					PrimaryKey: "id",
				},
				documents: make(map[string]Document),
			},
			wantErr: false,
		},
		{
			name: "error when collection not found",
			fields: fields{
				collections: make(map[string]*Collection),
			},
			args: args{
				name: "test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				collections: tt.fields.collections,
			}
			got, err := s.GetCollection(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCollection() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_DeleteCollection(t *testing.T) {
	type fields struct {
		collections map[string]*Collection
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successfully delete existing collection",
			fields: fields{
				collections: map[string]*Collection{
					"test": {
						cfg: CollectionConfig{
							PrimaryKey: "id",
						},
						documents: make(map[string]Document),
					},
				},
			},
			args: args{
				name: "test",
			},
			wantErr: false,
		},
		{
			name: "error when collection not found",
			fields: fields{
				collections: make(map[string]*Collection),
			},
			args: args{
				name: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				collections: tt.fields.collections,
			}
			if err := s.DeleteCollection(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_DumpToFile(t *testing.T) {
	// Создаем временную директорию для тестовых файлов
	tmpDir, err := os.MkdirTemp("", "store_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		store   *Store
		wantErr bool
	}{
		{
			name:    "dump empty store",
			store:   NewStore(),
			wantErr: false,
		},
		{
			name: "dump store with collections",
			store: func() *Store {
				s := NewStore()
				cfg := &CollectionConfig{PrimaryKey: "ID"}
				col, _ := s.CreateCollection("users", cfg)

				user, _ := MarshalDocument(struct {
					ID    string
					Name  string
					Email string
				}{
					ID:    "1",
					Name:  "Test User",
					Email: "test@example.com",
				})
				col.Put(*user)
				return s
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join(tmpDir, tt.name+".dump")

			if err := tt.store.DumpToFile(filename); (err != nil) != tt.wantErr {
				t.Errorf("Store.DumpToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if _, err := os.Stat(filename); err != nil {
				t.Errorf("File was not created: %v", err)
				return
			}

			restored, err := NewStoreFromFile(filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStoreFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			originalDump, err := tt.store.Dump()
			if err != nil {
				t.Errorf("Failed to dump original store: %v", err)
				return
			}

			restoredDump, err := restored.Dump()
			if err != nil {
				t.Errorf("Failed to dump restored store: %v", err)
				return
			}

			if string(originalDump) != string(restoredDump) {
				t.Errorf("Restored store doesn't match original.\nOriginal: %s\nRestored: %s",
					string(originalDump), string(restoredDump))
			}
		})
	}
}

func TestStore_NewStoreFromFile_Errors(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "non-existent file",
			filename: "non_existent_file.dump",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			filename: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStoreFromFile(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStoreFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_DumpToFile_Errors(t *testing.T) {
	store := NewStore()

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "empty filename",
			filename: "",
			wantErr:  true,
		},
		{
			name:     "invalid path",
			filename: "/non/existent/path/file.dump",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.DumpToFile(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.DumpToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
