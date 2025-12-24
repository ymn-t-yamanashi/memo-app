package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // SQLiteドライバ
)

type Memo struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

var db *sql.DB

func initDB() {
	var err error
	// データベース接続（ファイルがなければ作成される）
	db, err = sql.Open("sqlite3", "./memos.db")
	if err != nil {
		log.Fatal(err)
	}

	// テーブル作成
	query := `
	CREATE TABLE IF NOT EXISTS memos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		body TEXT NOT NULL
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.Use(cors.Default()) // Reactからの接続許可

	// 全件取得 API
	r.GET("/memos", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, body FROM memos ORDER BY id DESC")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var memos []Memo
		for rows.Next() {
			var m Memo
			rows.Scan(&m.ID, &m.Body)
			memos = append(memos, m)
		}
		c.JSON(http.StatusOK, memos)
	})

	// 保存 API
	r.POST("/memos", func(c *gin.Context) {
		var input Memo
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		result, err := db.Exec("INSERT INTO memos (body) VALUES (?)", input.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		id, _ := result.LastInsertId()
		input.ID = int(id)
		c.JSON(http.StatusCreated, input)
	})

	r.Run(":8080")
}
