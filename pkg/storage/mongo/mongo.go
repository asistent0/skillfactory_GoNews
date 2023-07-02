package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Post struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	Content     string             `bson:"content"`
	AuthorID    int                `bson:"author_id"`
	AuthorName  string             `bson:"author_name"`
	CreatedAt   int64              `bson:"created_at"`
	PublishedAt int64              `bson:"published_at"`
}

// Storage Хранилище данных.
type Storage struct {
	db *mongo.Collection
}

// New Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	db, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db.Database("goNews").Collection("post"),
	}
	return &s, nil
}

func (s *Storage) Posts() ([]storage.Post, error) {
	filter := bson.D{}
	cur, err := s.db.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			return
		}
	}(cur, context.Background())
	var posts []storage.Post
	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, cur.Err()
}

func (s *Storage) AddPost(p storage.Post) error {
	newPost := Post{
		Title:       p.Title,
		Content:     p.Content,
		AuthorID:    p.AuthorID,
		AuthorName:  p.AuthorName,
		CreatedAt:   time.Now().Unix(),
		PublishedAt: time.Now().Unix(),
	}
	_, err := s.db.InsertOne(context.TODO(), newPost)

	return err
}
func (s *Storage) UpdatePost(p storage.Post) error {
	filter := bson.D{{"_id", p.ID}}
	update := bson.D{{"$set", bson.D{
		{"title", p.Title},
		{"content", p.Content},
	}}}
	_, err := s.db.UpdateOne(context.TODO(), filter, update)

	return err
}
func (s *Storage) DeletePost(p storage.Post) error {
	filter := bson.D{{"_id", p.ID}}
	_, err := s.db.DeleteOne(context.TODO(), filter)

	return err
}
