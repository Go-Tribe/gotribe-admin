-- 初始化数据库脚本
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS gotribe CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE gotribe;

-- 创建用户（如果不存在）
CREATE USER IF NOT EXISTS 'gotribe'@'%' IDENTIFIED BY 'gotribe123';

-- 授权
GRANT ALL PRIVILEGES ON gotribe.* TO 'gotribe'@'%';

-- 刷新权限
FLUSH PRIVILEGES;
