ALTER TABLE characters 
MODIFY id INT(11)
AUTO_INCREMENT NOT NULL COMMENT 'キャラクター名',
DROP created_at,
DROP updated_at;