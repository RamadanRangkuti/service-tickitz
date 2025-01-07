--DDL=--
--Users
CREATE TABLE user_roles (
    id serial primary key,
    role varchar(40) not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

CREATE TABLE users (
    id serial primary key,
    email varchar(255) unique not null,
    password varchar(255) not null,
    role_id int not null references user_roles(id) on delete cascade,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

CREATE TABLE profiles (
    id serial primary key,
    user_id int unique not null references users(id) on delete cascade,
    first_name varchar(50),
    last_name varchar(50),
    phone_number varchar(16),
    image varchar(255),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

-- Movies
CREATE TABLE movie_directors (
    id serial primary key,
    name varchar(100),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

CREATE TABLE movies (
    id serial primary key,
    title varchar(255) unique not null,
    synopsis text,
    director_id int references movie_directors(id) on delete cascade,
    duration int, -- Durasi dalam menit
    release_date date,
    image varchar(255),
    banner varchar(255),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

CREATE TABLE movie_casts (
    id serial primary key,
    name varchar(100),
    movie_id int not null references movies(id) on delete cascade,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

CREATE TABLE movie_genres (
    id serial primary key,
    name varchar(100),
    movie_id int not null references movies(id) on delete cascade,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

-- Cinemas
CREATE TABLE cinemas (
    id serial primary key,
    name varchar(100),
    image varchar(255),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

CREATE TABLE locations (
    id serial primary key,
    city varchar(100),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

CREATE TABLE cinema_locations (
    id serial primary key,
    cinema_id int not null references cinemas(id) on delete cascade,
    location_id int not null references locations(id) on delete cascade,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);


--Jadwal tayang dan kursi
CREATE TABLE movie_schedules (
    id serial primary key,
    movie_id int not null references movies(id) on delete cascade,
    cinema_id int not null references cinemas(id) on delete cascade,
    date date not null,
    time time not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);


-- Seat Prices
CREATE TABLE seat_prices (
	id SERIAL PRIMARY KEY,
    category VARCHAR(20), -- regular, vip, love
    price NUMERIC(10, 2) NOT NULL
);

-- Cinema Seats
CREATE TABLE cinema_seats (
    id SERIAL PRIMARY KEY,
    cinema_id INT NOT NULL REFERENCES cinemas(id) ON DELETE CASCADE,
    seat_number VARCHAR(10) NOT NULL,
    seat_prices_id int NOT NULL REFERENCES seat_prices(id),
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE
);


-- Trx
CREATE TABLE transactions (
    id serial primary key,
    user_id int not null references users(id) on delete cascade,
    movie_id int not null references movies(id) on delete cascade,
    date timestamp with time zone not null,
    location varchar(255),
    quantity int not null,
    total_price numeric(10, 2),
    status varchar(50) default 'pending',
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

CREATE TABLE transaction_seats (
    id serial primary key,
    transaction_id int not null references transactions(id) on delete cascade,
    seat_number varchar(10) not null,
    cinema_id int not null references cinemas(id) on delete cascade,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);

CREATE TABLE payments (
    id serial primary key,
    transaction_id int not null references transactions(id) on delete cascade,
    payment_method varchar(50) not null, 
    virtual_account_number varchar(20) not null,
    amount numeric(10, 2) not null,
    status varchar(50) default 'pending',
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);


INSERT INTO user_roles (role) VALUES
('user'),
('admin');