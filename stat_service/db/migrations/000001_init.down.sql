-- -- Удаление триггеров
-- DROP TRIGGER IF EXISTS update_authors_timestamp ON stats.authors;
-- DROP TRIGGER IF EXISTS update_players_timestamp ON stats.players;
-- DROP TRIGGER IF EXISTS update_quizzes_timestamp ON stats.quizzes;

-- -- Удаление функции обновления временных меток
-- DROP FUNCTION IF EXISTS stats.update_timestamp;

-- -- Удаление индексов
-- DROP INDEX IF EXISTS stats.idx_quizzes_quiz;

-- -- Удаление таблиц (в обратном порядке создания)
-- DROP TABLE IF EXISTS stats.authors;
-- DROP TABLE IF EXISTS stats.players;
-- DROP TABLE IF EXISTS stats.quizzes;

-- -- Удаление схемы
-- DROP SCHEMA IF EXISTS stats CASCADE;