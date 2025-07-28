package storage

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

type News struct {
	ID			int
	Title		string
	Content		string
    CreatedAt	int64 //дата добавления в базу
    PublishedAt	int64 //дата публикации новости
    Link		string //ссылка на оригинальную новость
}


func New(connstr string) (*DB, error) {
	conn, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	// Проверка подключения
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Ping failed:", err)
	}
	db := DB{
		pool: conn,
	}
	return &db, nil
}


//запрос новостей из базы
func (db *DB) News(n int) ([]News, error) {
	var news []News
	rows, err := db.pool.Query(context.Background(), `
	SELECT
		id,
		title,
		content,
		created_at,
		published_at,
		link
	FROM news
	ORDER BY published_at DESC
	LIMIT $1
	`, n)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var nw News
		err = rows.Scan(
			&nw.ID,
			&nw.Title,
			&nw.Content,
			&nw.CreatedAt,
			&nw.PublishedAt,
			&nw.Link,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, nw)
	}
	return news, rows.Err()
}

//добавление новостей в БД
func (db *DB) AddNews(n []News) error {
	for _, nw := range n {
		_, err := db.pool.Exec(context.Background(), `
		INSERT INTO news(title, content, published_at, link, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
			nw.Title,
			nw.Content,
			nw.PublishedAt,
			nw.Link,
			time.Now().Unix(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}