-- Создание таблицы сессий
CREATE TABLE sessions (
    session_id VARCHAR(255) PRIMARY KEY,
    quiz_id VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);

-- Создание таблицы ответов пользователей
CREATE TABLE user_answers (
    session_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    answer_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (session_id, user_id, answer_id),
    FOREIGN KEY (session_id) REFERENCES sessions(session_id)
);

-- Создание индексов для ускорения поиска
CREATE INDEX idx_sessions_quiz ON sessions(quiz_id);
CREATE INDEX idx_answers_user ON user_answers(user_id);