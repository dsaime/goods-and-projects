CREATE TABLE goods
(
    id Int, -- идентификатор
    projectId Int, -- идентификатор
    name String, --  название
    description String, -- описание
    priority Int, -- приоритет
    removed bool, -- статус удаления
    eventTime DateTime -- дата и время
)
ENGINE = MergeTree
PRIMARY KEY (id, projectId, name)