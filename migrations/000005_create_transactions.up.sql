CREATE TABLE transactions (
    id serial primary key,
    user_id int references users(id) on delete cascade,
    amount decimal(10, 2) not null,
    type varchar(20) not null, 
    created_at timestamp default current_timestamp
);