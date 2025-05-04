CREATE TABLE users (
    id serial primary key,
    name varchar(100) not null,
    email varchar(100) unique not null,
    password_hash text not null,
    balance decimal(10, 2) default 0.00,
    role varchar(50) default 'user',
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

