CREATE SCHEMA IF NOT EXISTS stats;

-- Статистика по квизам (QuizStat)
CREATE TABLE IF NOT EXISTS stats.quizzes (
    quiz_id VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) NOT NULL,
    num_sessions INTEGER NOT NULL DEFAULT 0,
    avg_rate FLOAT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (quiz_id, author_id)
);

-- Статистика игроков (PlayerStat)
CREATE TABLE IF NOT EXISTS stats.players (
    user_id VARCHAR(255) PRIMARY KEY,
    total_score FLOAT NOT NULL DEFAULT 0,
    best_score FLOAT NOT NULL DEFAULT 0,
    avg_score FLOAT NOT NULL DEFAULT 0,
    num_sessions INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Статистика авторов (AuthorStat)
CREATE TABLE IF NOT EXISTS stats.authors (
    user_id VARCHAR(255) PRIMARY KEY,
    num_quizzes INTEGER NOT NULL DEFAULT 0,
    avg_quiz_rate FLOAT NOT NULL DEFAULT 0,
    best_quiz_rate FLOAT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для ускорения запросов
CREATE INDEX IF NOT EXISTS idx_quiz_stats_author ON stats.quiz_stats(author_id);
CREATE INDEX IF NOT EXISTS idx_quiz_stats_quiz ON stats.quiz_stats(quiz_id);

-- Функция и триггеры для обновления временных меток
CREATE OR REPLACE FUNCTION stats.update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_quiz_stats_timestamp
BEFORE UPDATE ON stats.quiz_stats
FOR EACH ROW EXECUTE FUNCTION stats.update_timestamp();

CREATE TRIGGER update_player_stats_timestamp
BEFORE UPDATE ON stats.player_stats
FOR EACH ROW EXECUTE FUNCTION stats.update_timestamp();

CREATE TRIGGER update_author_stats_timestamp
BEFORE UPDATE ON stats.author_stats
FOR EACH ROW EXECUTE FUNCTION stats.update_timestamp();