-- 书籍
CREATE TABLE m_book
(
    book_id      int8 NOT NULL, -- 书籍 ID
    title        text NOT NULL, -- 书名
    author       text NULL,     -- 作者
    total_page   int4 NULL,     -- 总页数
    publish_date date NULL,     -- 出版日期
    CONSTRAINT book_pkey PRIMARY KEY (book_id)
);
