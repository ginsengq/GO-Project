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