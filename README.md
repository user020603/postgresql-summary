# Tổng quan PostgreSQL

## Overview và So sánh với các SQL Database khác

**Điểm mạnh:**
- Tuân thủ chuẩn SQL và ACID đầy đủ
- Hỗ trợ dữ liệu không cấu trúc (JSON, XML, hstore)
- Extensions phong phú (PostGIS, TimescaleDB, pgVector)
- Khả năng mở rộng cao và xử lý khối lượng dữ liệu lớn
- Mã nguồn mở, cộng đồng lớn

**Điểm yếu:**
- Cấu hình phức tạp hơn MySQL
- Hiệu suất đọc thuần túy thấp hơn MySQL
- Chiếm nhiều tài nguyên hệ thống hơn SQLite

## CRUD Operations

```sql
-- CREATE
INSERT INTO users (name, email) VALUES ('Nguyen Tien Thanh', 'thanhnt@gmail.com');

-- READ
SELECT * FROM users WHERE id = 1;

-- UPDATE
UPDATE users SET email = 'thanhnt_new@gmail.com' WHERE id = 1;

-- DELETE
DELETE FROM users WHERE id = 1;
```

## Foreign Key

```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    order_date DATE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

**Ưu điểm:**
- Đảm bảo tính toàn vẹn của dữ liệu
- Hỗ trợ các hành động: CASCADE, SET NULL, SET DEFAULT, RESTRICT, NO ACTION

## Join

```sql
-- INNER JOIN
SELECT u.name, o.order_date 
FROM users u
INNER JOIN orders o ON u.id = o.user_id;

-- LEFT JOIN
SELECT u.name, o.order_date 
FROM users u
LEFT JOIN orders o ON u.id = o.user_id;

-- Full và Cross Join
SELECT * FROM table1 FULL JOIN table2 ON condition;
SELECT * FROM table1 CROSS JOIN table2;
```

## Subquery

```sql
-- Subquery trong WHERE
SELECT name FROM users 
WHERE id IN (SELECT user_id FROM orders WHERE total > 1000);

-- Subquery trong FROM
SELECT avg_price.average
FROM (SELECT AVG(price) as average FROM products) as avg_price;

-- Subquery trong SELECT
SELECT 
    name,
    (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) as order_count
FROM users;
```

## Index

- **Index** là cấu trúc dữ liệu giúp tăng tốc độ truy vấn dữ liệu trong bảng.
- **Mục đích:** Tăng hiệu suất truy vấn SELECT, UPDATE và DELETE.

1. **B-tree (mặc định):**
   - Tìm kiếm, chèn, cập nhật nhanh.
   - Phù hợp với các phép so sánh `=`, `<`, `>`, `<=`, `>=`.

   ```sql
   CREATE INDEX idx_user_email ON users(email);
   ```

2. **Hash:**
   - Tìm kiếm `=` nhanh.
   - Không hỗ trợ cho các phép so sánh phạm vi.

   ```sql
   CREATE INDEX idx_user_email_hash ON users USING hash (email);
   ```

3. **GiST (Generalized Search Tree):**
   - Dùng cho các kiểu dữ liệu phức tạp như hình học, toàn văn.

   ```sql
   CREATE INDEX idx_geom ON geom_table USING gist (geom);
   ```

4. **GIN (Generalized Inverted Index):**
   - Tìm kiếm toàn văn, mảng, JSONB.

   ```sql
   CREATE INDEX idx_gin ON documents USING gin (content);
   ```

5. **BRIN (Block Range Index):**
   - Cho bảng rất lớn, dữ liệu sắp xếp theo thứ tự tự nhiên.

   ```sql
   CREATE INDEX idx_brin ON large_table USING brin (column);
   ```

6. **Index trên nhiều cột:**

    - Tạo index trên nhiều cột để tối ưu các truy vấn phức tạp.

  ```sql
  CREATE INDEX idx_user_name_email ON users(name, email);
  ```

7. **Unique Index:**

    - Đảm bảo giá trị trong cột là duy nhất.

  ```sql
  CREATE UNIQUE INDEX idx_unique_email ON users(email);
  ```

## Partition

```sql
CREATE TABLE sales (
    id SERIAL,
    sale_date DATE NOT NULL,
    amount DECIMAL(10,2)
) PARTITION BY RANGE (sale_date);

CREATE TABLE sales_2023 PARTITION OF sales
    FOR VALUES FROM ('2023-01-01') TO ('2024-01-01');
    
CREATE TABLE sales_2024 PARTITION OF sales
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
```

**Loại phân vùng:**
- RANGE: phân vùng theo khoảng giá trị
- LIST: phân vùng theo danh sách giá trị
- HASH: phân vùng theo hàm băm

## Transaction

```sql
BEGIN;
    UPDATE accounts SET balance = balance - 100 WHERE id = 1;
    UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;

-- Hoặc sử dụng Savepoint
BEGIN;
    UPDATE accounts SET balance = balance - 100 WHERE id = 1;
    SAVEPOINT my_savepoint;
    UPDATE accounts SET balance = balance + 100 WHERE id = 2;
    -- Nếu có lỗi
    -- ROLLBACK TO my_savepoint;
COMMIT;
```

**Cấp độ Transaction:**
#### READ UNCOMMITTED
- **Mô tả:** Cho phép đọc dữ liệu chưa được commit.
- **Ưu điểm:** Tốc độ cao nhất.
- **Nhược điểm:** Nguy cơ đọc phải dữ liệu chưa hợp lệ (dirty read).

#### READ COMMITTED (mặc định)
- **Mô tả:** Chỉ đọc dữ liệu đã được commit tại thời điểm truy vấn.
- **Ưu điểm:** Tránh được dirty read.
- **Nhược điểm:** Có thể gặp phải non-repeatable read, dữ liệu có thể thay đổi giữa các lần đọc trong cùng một transaction.

#### REPEATABLE READ
- **Mô tả:** Đảm bảo dữ liệu không thay đổi giữa các lần đọc trong cùng một transaction.
- **Ưu điểm:** Tránh được dirty read và non-repeatable read.
- **Nhược điểm:** Có thể gặp phải phantom read, dữ liệu mới có thể xuất hiện trong các lần đọc sau.

#### SERIALIZABLE
- **Mô tả:** Đảm bảo tính tuần tự tuyệt đối, như thể các transaction được thực hiện lần lượt.
- **Ưu điểm:** Tránh được mọi vấn đề về concurrent data (dirty read, non-repeatable read, phantom read).
- **Nhược điểm:** Hiệu suất thấp nhất, dễ gặp deadlock.
