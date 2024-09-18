package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	client   *mongo.Client
	database *mongo.Database
}

func New(uri, dbName string) (*Store, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &Store{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.database.Collection("posts")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find posts: %w", err)
	}
	defer cursor.Close(context.Background())

	var posts []storage.Post
	if err = cursor.All(context.Background(), &posts); err != nil {
		return nil, fmt.Errorf("failed to decode posts: %w", err)
	}

	return posts, nil
}

func (s *Store) AddPost(post storage.Post) error {
	collection := s.database.Collection("posts")
	_, err := collection.InsertOne(context.Background(), post)
	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}
	return nil
}

func (s *Store) UpdatePost(post storage.Post) error {
	collection := s.database.Collection("posts")

	filter := bson.M{"ID": post.ID}

	update := bson.M{
		"$set": bson.M{
			"Title":     post.Title,
			"Content":   post.Content,
			"AuthorID":  post.AuthorID,
			"CreatedAt": post.CreatedAt,
		},
	}

	// Выполнение обновления
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	// Проверка, был ли обновлен хотя бы один документ
	if result.MatchedCount == 0 {
		return fmt.Errorf("no post found with ID: %d", post.ID)
	}

	return nil
}

func (s *Store) DeletePost(post storage.Post) error {
	collection := s.database.Collection("posts")

	filter := bson.M{"ID": post.ID}

	// Выполнение удаления
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	// Проверка, был ли удален хотя бы один документ
	if result.DeletedCount == 0 {
		return fmt.Errorf("no post found with ID: %d", post.ID)
	}

	return nil
}

func (s *Store) Authors() ([]storage.Author, error) {
	collection := s.database.Collection("authors")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find authors: %w", err)
	}
	defer cursor.Close(context.Background())

	var authors []storage.Author
	if err = cursor.All(context.Background(), &authors); err != nil {
		return nil, fmt.Errorf("failed to decode authors: %w", err)
	}

	return authors, nil
}
