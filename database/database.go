package database

import (
	"database/sql"
	"log"
)

// Book 书籍
type Book struct {
	BookId      int    // 书籍 ID
	Title       string // 书名
	Author      string // 作者
	PublishDate string // 出版日期
}

// RetreiveBook 查询书籍
func RetreiveBook(db *sql.DB) ([]Book, error) {
	var books []Book
	// 查询
	sql := `select book_id, title, author, to_char(publish_date, 'YYYY/MM/DD') as publish_date from m_book`
	rows, err := db.Query(sql)
	if err != nil {
		log.Println(err)
		return books, err
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		// 获取各列的值，放到对应的地址中
		rows.Scan(&book.BookId, &book.Title, &book.Author, &book.PublishDate)
		books = append(books, book)
	}

	return books, nil
}

// RetrieveMap SQL查询结果输出为Map
func RetrieveMap(db *sql.DB, sSql string) ([]map[string]interface{}, error) {

	// 准备查询语句
	stmt, err := db.Prepare(sSql)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	// 查询
	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	// 数据列
	columns, err := rows.Columns()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 列的个数
	count := len(columns)

	// 返回值 Map切片
	mData := make([]map[string]interface{}, 0)
	// 一条数据的各列的值（需要指定长度为列的个数，以便获取地址）
	values := make([]interface{}, count)
	// 一条数据的各列的值的地址
	valPointers := make([]interface{}, count)
	for rows.Next() {

		// 获取各列的值的地址
		for i := 0; i < count; i++ {
			valPointers[i] = &values[i]
		}

		// 获取各列的值，放到对应的地址中
		rows.Scan(valPointers...)

		// 一条数据的Map (列名和值的键值对)
		entry := make(map[string]interface{})

		// Map 赋值
		for i, col := range columns {
			var v interface{}

			// 值复制给val(所以Scan时指定的地址可重复使用)
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				// 字符切片转为字符串
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}

		mData = append(mData, entry)
	}

	return mData, nil
}
