ALTER TABLE characters 
MODIFY id INT(11)
UNSIGNED AUTO_INCREMENT NOT NULL COMMENT 'キャラクターID',
ADD created_at TIMESTAMP
DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '作成日時',
ADD updated_at TIMESTAMP
DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
NOT NULL COMMENT '更新日時';