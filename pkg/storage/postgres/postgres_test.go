package postgres

import (
	"GoNews/pkg/storage"
	"database/sql"
	"log"

	//"reflect"
	"testing"

	_ "github.com/lib/pq"
)

// Setup и TearDown функции для тестов

// setupTestDB инициализирует тестовую базу данных.
func setupTestDB() *Store {
	// Строка подключения к тестовой базе данных
	connStr := "user=postgres password=vlad5043 dbname=gonnews sslmode=disable" // Замените на актуальную строку подключения
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к тестовой базе данных: %v", err)
	}

	// Создание таблиц для тестов
	_, err = db.Exec(`
		DROP TABLE IF EXISTS posts, authors;
		CREATE TABLE authors (id SERIAL PRIMARY KEY, name TEXT NOT NULL);
		CREATE TABLE posts (id SERIAL PRIMARY KEY, author_id INTEGER REFERENCES authors(id), title TEXT NOT NULL, content TEXT, created_at BIGINT NOT NULL);
		INSERT INTO authors (name) VALUES ('John Doe'), ('Jane Smith');
		INSERT INTO posts (author_id, title, content, created_at) VALUES 
		(1, 'Test Post 1', 'Test Content 1', EXTRACT(EPOCH FROM NOW())::BIGINT),
		(2, 'Test Post 2', 'Test Content 2', EXTRACT(EPOCH FROM NOW())::BIGINT);
	`)
	if err != nil {
		log.Fatalf("Не удалось создать тестовые таблицы: %v", err)
	}

	return &Store{db: db}
}

// teardownTestDB удаляет тестовые данные.
func teardownTestDB(s *Store) {
	s.db.Exec("DROP TABLE IF EXISTS posts, authors;")
	s.db.Close()
}

// Тесты для методов Store

// Тест добавления нового поста
func TestStore_AddPost(t *testing.T) {
	s := setupTestDB()
	defer teardownTestDB(s)

	post := storage.Post{
		Title:     "New Test Post",
		Content:   "Test content for new post",
		AuthorID:  1,
		CreatedAt: 1234567890,
	}

	err := s.AddPost(post)
	if err != nil {
		t.Errorf("Store.AddPost() error = %v", err)
	}

	// Проверка, что пост был добавлен
	var addedPost storage.Post
	err = s.db.QueryRow(`SELECT title, content, author_id, created_at FROM posts WHERE title = $1`, post.Title).Scan(&addedPost.Title, &addedPost.Content, &addedPost.AuthorID, &addedPost.CreatedAt)
	if err != nil {
		t.Errorf("Post not found in DB after insertion: %v", err)
	}

	if addedPost.Title != post.Title || addedPost.Content != post.Content {
		t.Errorf("Store.AddPost() = %v, want %v", addedPost, post)
	}
}

// Тест получения всех постов
/*func TestStore_Posts(t *testing.T) {
	s := setupTestDB()
	defer teardownTestDB(s)

	// Ожидаемые значения постов, игнорируя поля CreatedAt и PublishedAt
	expectedPosts := []struct {
		ID         int
		Title      string
		Content    string
		AuthorID   int
		AuthorName string
	}{
		{ID: 1, Title: "Test Post 1", Content: "Test Content 1", AuthorID: 1, AuthorName: "John Doe"},
		{ID: 2, Title: "Test Post 2", Content: "Test Content 2", AuthorID: 2, AuthorName: "Jane Smith"},
	}

	gotPosts, err := s.Posts()
	if err != nil {
		t.Errorf("Store.Posts() error = %v", err)
		return
	}

	if len(gotPosts) != len(expectedPosts) {
		t.Errorf("Store.Posts() len = %v, want %v", len(gotPosts), len(expectedPosts))
		return
	}

	for i, got := range gotPosts {
		// Игнорируем поля CreatedAt и PublishedAt
		gotPost := struct {
			ID         int
			Title      string
			Content    string
			AuthorID   int
			AuthorName string
		}{
			ID:         got.ID,
			Title:      got.Title,
			Content:    got.Content,
			AuthorID:   got.AuthorID,
			AuthorName: got.AuthorName,
		}

		expectedPost := expectedPosts[i]

		if !reflect.DeepEqual(gotPost, expectedPost) {
			t.Errorf("Store.Posts()[%d] = %v, want %v", i, gotPost, expectedPost)
		}
	}
}*/

// Тест обновления поста
func TestStore_UpdatePost(t *testing.T) {
	s := setupTestDB()
	defer teardownTestDB(s)

	// Обновляем существующий пост
	updatePost := storage.Post{
		ID:        1,
		Title:     "Updated Post Title",
		Content:   "Updated Content",
		AuthorID:  1,
		CreatedAt: 1234567890,
	}

	err := s.UpdatePost(updatePost)
	if err != nil {
		t.Errorf("Store.UpdatePost() error = %v", err)
	}

	// Проверка, что пост был обновлен
	var updatedPost storage.Post
	err = s.db.QueryRow(`SELECT title, content FROM posts WHERE id = $1`, updatePost.ID).Scan(&updatedPost.Title, &updatedPost.Content)
	if err != nil {
		t.Errorf("Store.UpdatePost() post not found: %v", err)
	}

	if updatedPost.Title != updatePost.Title || updatedPost.Content != updatePost.Content {
		t.Errorf("Store.UpdatePost() = %v, want %v", updatedPost, updatePost)
	}
}

// Тест удаления поста
func TestStore_DeletePost(t *testing.T) {
	s := setupTestDB()
	defer teardownTestDB(s)

	// Удаляем существующий пост
	postToDelete := storage.Post{
		ID: 1,
	}

	err := s.DeletePost(postToDelete)
	if err != nil {
		t.Errorf("Store.DeletePost() error = %v", err)
	}

	// Проверка, что пост был удален
	var count int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM posts WHERE id = $1`, postToDelete.ID).Scan(&count)
	if err != nil {
		t.Errorf("Store.DeletePost() failed to verify deletion: %v", err)
	}

	if count != 0 {
		t.Errorf("Store.DeletePost() post still exists in DB")
	}
}

// Тест создания нового подключения к базе данных
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
			args:    args{connStr: "user=postgres password=vlad5043 dbname=gonnews sslmode=disable"},
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
			_, err := New(tt.args.connStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
