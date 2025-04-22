CREATE TABLE favorite_quizzes (
    user_id BIGINT NOT NULL,
    quiz_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, quiz_id),
    CONSTRAINT fk_favorite_quizzes_user
        FOREIGN KEY (user_id) 
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_favorite_quizzes_user_id ON favorite_quizzes(user_id);
CREATE INDEX idx_favorite_quizzes_quiz_id ON favorite_quizzes(quiz_id);