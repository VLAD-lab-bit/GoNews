package postgres

import (
	"GoNews/pkg/storage"
	"database/sql"

	_ "github.com/lib/pq"
)

// Store - структура для работы с PostgreSQL.
type Store struct {
	db *sql.DB
}

// New создаёт новое соединение с базой данных PostgreSQL.
func New(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Posts возвращает все публикации.
func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(`SELECT posts.id, posts.title, posts.content, posts.author_id, authors.name, posts.created_at FROM posts JOIN authors ON posts.author_id = authors.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var post storage.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AuthorName, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// AddPost добавляет новую публикацию.
func (s *Store) AddPost(post storage.Post) error {
	query := `INSERT INTO posts (title, content, author_id, created_at) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, post.Title, post.Content, post.AuthorID, post.CreatedAt)
	return err
}

// UpdatePost обновляет существующую публикацию.
func (s *Store) UpdatePost(post storage.Post) error {
	query := `UPDATE posts SET title = $1, content = $2, author_id = $3, created_at = $4 WHERE id = $5`
	_, err := s.db.Exec(query, post.Title, post.Content, post.AuthorID, post.CreatedAt, post.ID)
	return err
}

// DeletePost удаляет публикацию по ID.
func (s *Store) DeletePost(post storage.Post) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := s.db.Exec(query, post.ID)
	return err
}

// Authors возвращает всех авторов.
func (s *Store) Authors() ([]storage.Author, error) {
	rows, err := s.db.Query(`SELECT id, name FROM authors`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []storage.Author
	for rows.Next() {
		var author storage.Author
		err := rows.Scan(&author.ID, &author.Name)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}
