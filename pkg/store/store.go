package store

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DB Database
type DB struct {
	pool *pgxpool.Pool
}

// Post structura got from rss
type Post struct {
	ID      int
	Title   string
	Content string
	PubTime int64
	Link    string
}

func New() (*DB, error) {
	connstr := os.Getenv("newsdb")
	if connstr == "" {
		return nil, errors.New("не ууказано подключение к БД")
	}
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	db := DB{
		pool: pool,
	}
	return &db, nil
}

// RecordNews recording news in DB
func (db *DB) RecordNews(news []Post) error {
	for _, post := range news {
		_, err := db.pool.Exec(context.Background(), `
		INSERT INTO news(title, content, pub_time, link)
		VALUES ($1, $2, $3, $4)`,
			post.Title,
			post.Content,
			post.PubTime,
			post.Link,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReturnNews return news from BD
func (db *DB) ReturnNews(count int) ([]Post, error) {
	if count == 0 {
		count = 10
	}
	rows, err := db.pool.Query(context.Background(), `
	SELECT id, title, content, pub_time, link FROM news
	ORDER BY pub_time DESC
	LIMIT $1
	`,
		count,
	)
	if err != nil {
		return nil, err
	}
	var news []Post
	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, p)
	}
	return news, rows.Err()
}
