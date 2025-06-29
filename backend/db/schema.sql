CREATE TABLE IF NOT EXISTS user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) 
);

CREATE TABLE IF NOT EXISTS score (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    _value INT NOT NULL,
    CONSTRAINT fk_score_user_id FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);