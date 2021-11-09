package database

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

func TestRetrieveMap(t *testing.T) {

	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=12 dbname=postgres sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatalln("无法连接数据库")
	}

	// 查询书籍
	sql := "select book_id, title, author, to_char(publish_date, 'YYYY/MM/DD') as publish_date from m_book"
	m, err := RetrieveMap(db, sql)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%v", m)
}

func TestRetrieveBook(t *testing.T) {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=12 dbname=postgres sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatalln("无法连接数据库")
	}

	// 查询书籍
	books, err := RetreiveBook(db)
	log.Printf("%v", books)
}
