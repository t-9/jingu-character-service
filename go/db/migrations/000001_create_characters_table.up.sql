CREATE TABLE IF NOT EXISTS characters (
    id INT(11) AUTO_INCREMENT NOT NULL COMMENT 'キャラクター名',
    surname VARCHAR(30) DEFAULT '' COMMENT '苗字',
    given_name VARCHAR(30) NOT NULL COMMENT '名前',
    PRIMARY KEY (id)
)
CHARACTER SET urf8mb4 
COLLATE utf8mb4_general_ci
COMMENT='キャラクターテーブル';