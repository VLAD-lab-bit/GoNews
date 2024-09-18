package mongo

// setupTestDB инициализирует тестовую базу данных.
/*func setupTestDB() *Store {
	connStr := "mongodb://localhost:27017" // Замените на актуальную строку подключения
	clientOptions := options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Не удалось подключиться к тестовой базе данных: %v", err)
	}

	store := &Store{
		client: client,
		posts:  client.Database("gonnews").Collection("posts"),
	}

	_, err = store.posts.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		log.Fatalf("Не удалось очистить коллекцию: %v", err)
	}

	_, err = store.posts.InsertMany(context.Background(), []interface{}{
		storage.Post{ID: 1, Title: "Test Post 1", Content: "Test Content 1", AuthorID: 1, CreatedAt: 1234567890},
		storage.Post{ID: 2, Title: "Test Post 2", Content: "Test Content 2", AuthorID: 2, CreatedAt: 1234567890},
	})
	if err != nil {
		log.Fatalf("Не удалось вставить тестовые данные: %v", err)
	}

	return store
}

// teardownTestDB удаляет тестовые данные.
func teardownTestDB(s *Store) {
	s.posts.Drop(context.Background())
	s.client.Disconnect(context.Background())
}

func TestNew(t *testing.T) {
	type args struct {
		connStr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Valid connection",
			args:    args{connStr: "mongodb://localhost:27017"},
			wantErr: false,
		},
		{
			name:    "Invalid connection",
			args:    args{connStr: "invalid_connection_string"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.connStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Errorf("New() = nil, want non-nil")
				return
			}
			if err == nil && (got.client == nil || got.posts == nil) {
				t.Errorf("New() = %v, want non-nil client and posts", got)
			}
		})
	}
}

func TestStore_Posts(t *testing.T) {
	s := setupTestDB()
	defer teardownTestDB(s)

	expectedPosts := []storage.Post{
		{ID: 1, Title: "Test Post 1", Content: "Test Content 1", AuthorID: 1, CreatedAt: 1234567890},
		{ID: 2, Title: "Test Post 2", Content: "Test Content 2", AuthorID: 2, CreatedAt: 1234567890},
	}

	tests := []struct {
		name    string
		s       *Store
		want    []storage.Post
		wantErr bool
	}{
		{
			name:    "Get all posts",
			s:       s,
			want:    expectedPosts,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Posts()
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.Posts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Store.Posts() len = %v, want %v", len(got), len(tt.want))
			}
			for i, gotPost := range got {
				if !reflect.DeepEqual(gotPost, tt.want[i]) {
					t.Errorf("Store.Posts()[%d] = %v, want %v", i, gotPost, tt.want[i])
				}
			}
		})
	}
}

func TestStore_AddPost(t *testing.T) {
	s := setupTestDB()
	defer teardownTestDB(s)

	newPost := storage.Post{
		ID:        3,
		Title:     "New Test Post",
		Content:   "Test content for new post",
		AuthorID:  1,
		CreatedAt: 1234567890,
	}

	tests := []struct {
		name    string
		s       *Store
		args    storage.Post
		wantErr bool
	}{
		{
			name:    "Add new post",
			s:       s,
			args:    newPost,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddPost(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Store.AddPost() error = %v, wantErr %v", err, tt.wantErr)
			}

			var addedPost storage.Post
			err := s.posts.FindOne(context.Background(), bson.D{{"id", tt.args.ID}}).Decode(&addedPost)
			if err != nil {
				t.Errorf("Post not found in DB after insertion: %v", err)
			}
			if !reflect.DeepEqual(addedPost, tt.args) {
				t.Errorf("Store.AddPost() = %v, want %v", addedPost, tt.args)
			}
		})
	}
}

func TestStore_UpdatePost(t *testing.T) {
	s := setupTestDB()
	defer teardownTestDB(s)

	updatedPost := storage.Post{
		ID:        1,
		Title:     "Updated Post Title",
		Content:   "Updated Content",
		AuthorID:  1,
		CreatedAt: 1234567890,
	}

	tests := []struct {
		name    string
		s       *Store
		args    storage.Post
		wantErr bool
	}{
		{
			name:    "Update existing post",
			s:       s,
			args:    updatedPost,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.UpdatePost(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Store.UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
			}

			var updated storage.Post
			err := s.posts.FindOne(context.Background(), bson.D{{"id", tt.args.ID}}).Decode(&updated)
			if err != nil {
				t.Errorf("Post not found in DB after update: %v", err)
			}
			if !reflect.DeepEqual(updated, tt.args) {
				t.Errorf("Store.UpdatePost() = %v, want %v", updated, tt.args)
			}
		})
	}
}

func TestStore_DeletePost(t *testing.T) {
	s := setupTestDB()
	defer teardownTestDB(s)

	postToDelete := storage.Post{
		ID: 1,
	}

	tests := []struct {
		name    string
		s       *Store
		args    storage.Post
		wantErr bool
	}{
		{
			name:    "Delete existing post",
			s:       s,
			args:    postToDelete,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.DeletePost(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Store.DeletePost() error = %v, wantErr %v", err, tt.wantErr)
			}

			count, err := s.posts.CountDocuments(context.Background(), bson.D{{"id", tt.args.ID}})
			if err != nil {
				t.Errorf("Failed to count documents after deletion: %v", err)
			}
			if count != 0 {
				t.Errorf("Post still exists in DB after deletion")
			}
		})
	}
}
*/
