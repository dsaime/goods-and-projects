CREATE TABLE goods
(
    Id          Int,      -- идентификатор
    ProjectId   Int,      -- идентификатор
    Name        String,   --  название
    Description String,   -- описание
    Priority    Int,      -- приоритет
    Removed     bool,     -- статус удаления
    EventTime   DateTime, -- дата и время
    INDEX id_index Id TYPE set(0) GRANULARITY 1,
    INDEX project_id_index ProjectId TYPE set(0) GRANULARITY 1,
    INDEX name_index Name TYPE set(0) GRANULARITY 1
)
    ENGINE = MergeTree
        ORDER BY (Id, ProjectId, EventTime, Removed);

