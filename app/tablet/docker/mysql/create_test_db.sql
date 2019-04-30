CREATE DATABASE IF NOT EXISTS `claude_test` COLLATE 'utf8_general_ci' ;
GRANT ALL PRIVILEGES ON *.* TO 'claude'@'localhost' IDENTIFIED BY 'password';
FLUSH PRIVILEGES;