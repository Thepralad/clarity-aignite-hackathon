package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Thepralad/clarity-aignite-hackathon/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

// DB is the global database connection pool
var DB *sql.DB

// Init initializes the MySQL database connection
func Init(dsn string) error {

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}

	// Set connection pool limits
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Test connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Connected to MySQL database successfully")
	return nil
}

func GetArticles() ([]models.Article, error) {
	rows, err := DB.Query("SELECT id, title, img_url, category FROM articles")
	if err != nil {
		return nil, fmt.Errorf("failed to query articles: %w", err)
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.ImgURL, &article.Category); err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return articles, nil
}
