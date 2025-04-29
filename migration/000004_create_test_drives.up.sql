CREATE TABLE test_drives (
    id serial primary key,
    user_id int references users(id) on delete cascade,
    car_id int references cars(id) on delete restrict,
    order_id int references orders(id) on delete set null,
    date timestamp not null,
    status varchar(20) default 'scheduled',
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
