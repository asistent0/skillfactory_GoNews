package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Storage Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// New Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

func (s *Storage) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			post.id,
			post.title,
			post.content,
			author.id as author_id,
			author.name,
			post.created_at,
			post.created_at
		FROM post
		LEFT JOIN author ON post.author_id = author.id
		ORDER BY post.id;
	`,
	)
	if err != nil {
		return nil, err
	}
	var posts []storage.Post
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
			&p.PublishedAt,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, p)
	}
	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}

func (s *Storage) AddPost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO post (title, content, author_id, created_at, published_at)
		VALUES ($1, $2, $3, $4, $4);
		`,
		p.Title,
		p.Content,
		p.AuthorID,
		time.Now().Unix(),
	)
	return err
}
func (s *Storage) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE post SET title = $1, content = $2 WHERE id = $3;
		`,
		p.Title,
		p.Content,
		p.ID,
	)
	return err
}
func (s *Storage) DeletePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM post WHERE id = $1;
		`,
		p.ID,
	)
	return err
}
