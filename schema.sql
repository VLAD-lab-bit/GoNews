-- Удаление таблиц, если они существуют, чтобы избежать конфликтов при повторном создании
DROP TABLE IF EXISTS posts, authors;

-- Создание таблицы авторов
CREATE TABLE authors (
    id SERIAL PRIMARY KEY,  -- Идентификатор автора, автоинкрементируемое поле
    name TEXT NOT NULL      -- Имя автора
);

-- Создание таблицы постов
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,        -- Идентификатор поста, автоинкрементируемое поле
    author_id INTEGER REFERENCES authors(id) NOT NULL,  -- Внешний ключ на таблицу authors
    title TEXT NOT NULL,          -- Заголовок поста
    content TEXT,                 -- Содержание поста
    created_at BIGINT NOT NULL     -- Время создания поста (UNIX timestamp)
);

-- Вставка примера данных в таблицу authors
INSERT INTO authors (name) VALUES
('John Doe'),
('Jane Smith');

-- Вставка примера данных в таблицу posts
INSERT INTO posts (author_id, title, content, created_at) VALUES
(1, 'First Post', 'This is the content of the first post', EXTRACT(EPOCH FROM NOW())::BIGINT),
(2, 'Second Post', 'Here is another interesting post', EXTRACT(EPOCH FROM NOW())::BIGINT);
