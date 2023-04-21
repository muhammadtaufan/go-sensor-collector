CREATE TABLE IF NOT EXISTS sensor (
  id CHAR(36) PRIMARY KEY,
  sensor_value FLOAT,
  sensor_type VARCHAR(50),
  id1 VARCHAR(50),
  id2 INT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_created_at (created_at),
  INDEX idx_id1_created_at (id1, created_at),
  INDEX idx_id1_id2 (id1, id2)
);
