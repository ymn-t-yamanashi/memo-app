AIで生成したメモ


初学者の方向けに、**「これだけで動く・データが消えない・構成がシンプル」**な教材用セットを作成しました。

データベースには設定不要で使える **SQLite** を採用し、Goの標準ライブラリに近い `database/sql` を使って、データの流れが直感的にわかるようにしています。

---

## 1. プロジェクト構造

まず、以下のようなフォルダ構成にします。

```text
memo-app/
├── backend/
│   ├── main.go
│   ├── go.mod
│   └── memos.db (自動生成されます)
└── frontend/
    ├── src/
    │   └── App.js
    └── (その他のReactファイル)

```

---

## 2. バックエンド (Go + Gin + SQLite)

### セットアップ

```bash
mkdir backend
cd backend
go mod init simple-memo-backend
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go get github.com/mattn/go-sqlite3

```

### `main.go` の作成

```go
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

```

---

## 3. フロントエンド (React)

### セットアップ

```bash
cd ..
npx create-react-app frontend
cd frontend
npm install axios

```

### `src/App.js` の書き換え

```jsx
import React, { useState, useEffect } from 'react';
import axios from 'axios';

const API_URL = 'http://localhost:8080/memos';

function App() {
  const [memos, setMemos] = useState([]);
  const [text, setText] = useState('');

  // 1. データ取得 (Read)
  const fetchMemos = async () => {
    const res = await axios.get(API_URL);
    setMemos(res.data || []);
  };

  useEffect(() => { fetchMemos(); }, []);

  // 2. データ保存 (Create)
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!text.trim()) return;
    await axios.post(API_URL, { body: text });
    setText('');
    fetchMemos();
  };

  return (
    <div style={{ maxWidth: '500px', margin: '40px auto', fontFamily: 'sans-serif' }}>
      <h1>Memo App</h1>
      <form onSubmit={handleSubmit} style={{ display: 'flex', gap: '10px' }}>
        <input 
          style={{ flex: 1, padding: '10px' }}
          value={text} 
          onChange={(e) => setText(e.target.value)} 
          placeholder="今日の一言を記録..."
        />
        <button style={{ padding: '10px 20px' }} type="submit">保存</button>
      </form>

      <div style={{ marginTop: '20px' }}>
        {memos.map(m => (
          <div key={m.id} style={{ padding: '10px', borderBottom: '1px solid #eee' }}>
            {m.body}
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;

```

---

## 4. 学習のポイント（解説のヒント）

教材として教える際は、以下の3点に注目させると理解が深まります。

1. **JSONの橋渡し**:
* フロントが `axios.post` で送る **JSON** を、Go側が `ShouldBindJSON` で受け取る流れ。


2. **SQLの基本**:
* `CREATE TABLE` (器を作る), `INSERT` (入れる), `SELECT` (出す) というDB操作の基本。


3. **ステート管理**:
* Reactの `useState` が更新されると、HTMLが自動的に書き換わる仕組み。



---

## 5. 実行手順まとめ

1. **バックエンド起動**: `backend` フォルダで `go run main.go`
2. **フロントエンド起動**: `frontend` フォルダで `npm start`
3. **確認**: ブラウザで `localhost:3000` を開き、文字を入力して保存。
4. **永続化確認**: バックエンドを止めて（Ctrl+C）再度起動しても、データが残っていることを確認。

このサンプルをベースに、「削除ボタン」を追加する課題などを出してみるのはいかがでしょうか？必要であれば、削除機能（DELETE）のコードも作成します。

実行方法

```
cd backend
go run main.go

cd backend
npm start
```

# Goでホットリード

```
go install github.com/air-verse/air@latest
air init 
air 
```