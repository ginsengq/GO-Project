CREATE TABLE orders (
    id serial primary key,
    user_id int references users(id) on delete cascade,
    car_id int references cars(id) on delete restrict,
    status varchar(20) default 'reserved',
    deposit decimal(10, 2) default 0.00,
    total_price decimal(12, 2) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);