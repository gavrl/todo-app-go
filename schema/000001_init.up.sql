CREATE TABLE users
(
    listId        serial       not null unique,
    name          varchar(255) not null,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE todo_lists
(
    listId      serial       not null unique,
    title       varchar(255) not null,
    description varchar(255)
);

CREATE TABLE user_lists
(
    listId  serial                                               not null unique,
    user_id int references users (listId) on delete cascade      not null,
    list_id int references todo_lists (listId) on delete cascade not null
);

CREATE TABLE todo_items
(
    listId      serial       not null unique,
    title       varchar(255) not null,
    description varchar(255),
    done        boolean      not null default false
);

CREATE TABLE lists_items
(
    listId  serial                                               not null unique,
    user_id int references todo_items (listId) on delete cascade not null,
    list_id int references todo_lists (listId) on delete cascade not null
);