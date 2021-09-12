DROP DATABASE IF EXISTS test_db;
CREATE DATABASE test_db;
USE test_db;

CREATE TABLE statuses
(
    id     INT         NOT NULL AUTO_INCREMENT,
    status VARCHAR(10) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE priorities
(
    id       INT     NOT NULL AUTO_INCREMENT,
    priority CHAR(1) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE todos
(
  id                  INT         NOT NULL AUTO_INCREMENT,
  title               VARCHAR(50) NOT NULL,
  implementation_date DATE        NOT NULL,
  due_date            DATE        NOT NULL,
  status_id           INT         NOT NULL,
  priority_id         INT  DEFAULT NULL,
  memo                TEXT DEFAULT NULL,
  PRIMARY KEY (id),

  FOREIGN KEY fk_status_id (status_id)
    REFERENCES statuses (id)
    ON DELETE RESTRICT ON UPDATE CASCADE,

  FOREIGN KEY fk_priority_id (priority_id)
    REFERENCES priorities (id)
    ON DELETE RESTRICT ON UPDATE CASCADE
);
