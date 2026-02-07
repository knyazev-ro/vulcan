-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    last_name VARCHAR(255)
);

-- Таблица профилей (Has-One к пользователю)
CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    bio TEXT,
    avatar VARCHAR(255)
);

-- Таблица постов (Has-Many к пользователю)
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    user_id INTEGER REFERENCES users(id)
);

-- Таблица тегов
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

-- Связующая таблица (Pivot) для Post-Tags
CREATE TABLE post_tags (
    post_id INTEGER REFERENCES posts(id),
    tag_id INTEGER REFERENCES tags(id),
    PRIMARY KEY (post_id, tag_id)
);

-- 1. Пользователи (разные кейсы: с профилем, без профиля, с кучей постов)
INSERT INTO users (name, last_name) VALUES 
('Luka', 'Original'),    -- Имеет профиль, 3 поста
('Luka', 'Clone'),       -- Имеет профиль, 1 пост
('Zachary', 'Istrian'),  -- Нет профиля, 5 постов
('Prince', 'None'),      -- Имеет профиль, 0 постов
('Ghost', 'User');       -- Нет профиля, 0 постов (крайний случай)

-- 2. Профили
INSERT INTO profiles (user_id, bio, avatar) VALUES 
(1, 'The real brother who died in the ruins.', 'luka_orig.png'),
(2, 'A soul in a mechanical shell. High-level Apollo user.', 'luka_clone.png'),
(4, 'A competitor for the throne. Dangerous.', 'prince.png');

-- 3. Теги (категории магии и технологий)
INSERT INTO tags (name) VALUES 
('Apollo'), ('Low-level'), ('Magic'), ('Hardware'), 
('Imperial'), ('Memory-Leak'), ('Refraction'), ('Kernel');

-- 4. Посты
INSERT INTO posts (name, user_id) VALUES 
('Memory Management in Apollo', 1),
('The Day the School Fell', 1),
('Unauthorized Access to Throne Registry', 1),
('Mechanical Body Maintenance', 2),
('Zachary Terminals: Budget Edition', 3),
('Quartz vs Magnetic Discs', 3),
('Compiling Fireballs for Beginners', 3),
('Optimization of Soul Runtimes', 3),
('Inter-world State Management', 3);

-- 5. Связи Post-Tags (разное количество тегов на пост)
INSERT INTO post_tags (post_id, tag_id) VALUES 
(1, 1), (1, 2), (1, 6), 
(2, 5),                 
(3, 5), (3, 8),         
(4, 4), (4, 1),         
(5, 4),                 
(6, 4), (6, 3),         
(7, 3), (7, 1), (7, 2), 
(8, 1), (8, 6), (8, 7), 
(9, 5);                 