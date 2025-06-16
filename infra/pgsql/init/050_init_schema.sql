CREATE TABLE projects
(
    id         SERIAL PRIMARY KEY,                          -- id записи
    name       TEXT      NOT NULL DEFAULT '',               -- название
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- дата и время
);

CREATE TABLE goods
(
    id          INTEGER,                                    -- id записи
    project_id  INTEGER,                                    -- id компании
    name        TEXT    NOT NULL DEFAULT '',                -- название
    description TEXT    NOT NULL DEFAULT '',                -- описание
    priority    INTEGER NOT NULL DEFAULT 0,                 -- приоритет
    removed     BOOLEAN NOT NULL DEFAULT FALSE,             -- статус удаления
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP, -- дата и время
    PRIMARY KEY (id, project_id)
);

CREATE INDEX goods_name_idx ON goods (name);


CREATE FUNCTION max_goods_priority()
    RETURNS int
    LANGUAGE sql AS
$$
SELECT MAX(priority) + 1
FROM goods
$$;

ALTER TABLE goods
    ALTER COLUMN priority SET DEFAULT max_goods_priority();
