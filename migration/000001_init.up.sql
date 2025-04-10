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

CREATE TABLE cars (
    id serial primary key,
    brand varchar(50) not null,
    model varchar(50) not null,
    year integer not null,
    price decimal(12, 2) not null,
    mileage integer default 0,
    color varchar(50),
    status varchar(20) default 'available',
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

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

CREATE TABLE transactions (
    id serial primary key,
    user_id int references users(id) on delete cascade,
    amount decimal(10, 2) not null,
    type varchar(20) not null, 
    created_at timestamp default current_timestamp,
);

CREATE INDEX idx_cars_brand_model ON cars(brand, model);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
