-- 创建database----->选择database---->创建具体账目表
-- 创建database
CREATE DATABASE IF NOT EXISTS ledger
CHARACTER SET utf8mb4;

-- 选择database
USE ledger;

-- 创建具体账目表
-- 名称id 类型kind（支出还是收入） 金额amount_cents 分类category（餐饮/交通） 发生时间occurred_on 备注note
-- 更新时间updated_at
CREATE TABLE IF NOT EXISTS transactions(
        id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
        kind VARCHAR(10) NOT NULL,
        amount_cents BIGINT NOT NULL,
        category VARCHAR(50) NOT NULL,
        occurred_on DATE NOT NULL,
        note VARCHAR(255) NOT NULL DEFAULT '',
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP 
            ON UPDATE CURRENT_TIMESTAMP,

        CONSTRAINT chk_kind
            CHECK (kind IN ('income', 'expense')),
            
        CONSTRAINT cuk_amounts
            CHECK (amount_cents > 0),

        CONSTRAINT cuk_category
            CHECK (CHAR_LENGTH(TRIM(category))>0)     
);