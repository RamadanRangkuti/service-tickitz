create table user_roles(
	id serial unique not null primary key,
	role varchar(40) not null,
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone
);

create table users(
	id serial unique not null primary key,
	email varchar(255) unique not null,
	password varchar(255) not null,
	role_id int not null,
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone,
	foreign key (role_id) references user_roles(id)
);


create table profiles(
	id serial unique not null primary key,
	user_id int not null unique,
	first_name varchar(50),
	last_name varchar(50),
	phone_number varchar(16),
	image varchar(255),
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone,
	foreign key (user_id) references users(id)
);

create table movie_directors(
	id serial unique not null primary key,
	name varchar(100),
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone
);

create table movies(
	id serial unique not null primary key,
	title varchar(255) unique,
	synopsis text,
	director_id int,
	duration int,
	realease_date date,
	image varchar(255),
	banner varchar(255),
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone,
	foreign key(director_id) references movie_directors(id)
);

create table movie_casts(
	id serial unique not null primary key,
	name varchar(100),
	movie_id int,
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone,
	foreign key (movie_id) references movies(id)
);

create table movie_genres(
	id serial unique not null primary key,
	name varchar(100),
	movie_id int,
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone,
	foreign key (movie_id) references movies(id)
);
create table cinemas (
	id serial unique not null primary key,
	name varchar(100),
	image varchar(255),
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone
);

create table locations(
	id serial unique not null primary key,
	city varchar(100),
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone
);

create table cinema_locations(
	id serial unique not null primary key,
	cinema_id int,
	location_id int,
	foreign key(cinema_id) references cinemas(id),
	foreign key(location_id) references locations(id),
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone
);

CREATE TABLE transactions (
    id serial primary key,
    user_id int references users(id) on delete cascade,
    movie_id int references movies(id) on delete cascade,
    date timestamp with time zone not null,
    location varchar(255),
    quantity int not null,
    payment_id int references payments(id) on delete cascade,
    info_fullname varchar(255),
    info_email varchar(255),
    info_phone varchar(16),
    total_price numeric(10, 2),
    status varchar(50) default 'pending',
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);



-- UPDATE
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

CREATE TABLE cinema_seats (
    id serial primary key,
    cinema_id int not null references cinemas(id) on delete cascade,
    seat_number varchar(10) not null,
    is_available boolean default true,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
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
    status varchar(50) default 'pending', -- Status transaksi (pending, completed, canceled)
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
    payment_method varchar(50) not null, -- Metode pembayaran (e-wallet, transfer, dll)
    virtual_account_number varchar(20) not null, -- Nomor virtual account untuk transaksi
    amount numeric(10, 2) not null, -- Total pembayaran
    status varchar(50) default 'pending', -- Status pembayaran (pending, success, failed)
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone
);



INSERT INTO user_roles (role) VALUES
('user'),
('admin');

INSERT INTO users (email, password, role_id) VALUES
('user@example.com', 'hashed_password_admin', 1),
('admin@example.com', 'hashed_password_user', 2);

INSERT INTO profiles (user_id, first_name, last_name, phone_number, image) VALUES
(1, 'User name', 'Last User', '081234567890', 'user_image.jpg'),
(2, 'Admin name', 'Last Admin', '081234567891', 'admin_image.jpg');

INSERT INTO movie_directors (name) VALUES
('Steven Spielberg'),
('Christopher Nolan'),
('Jon Watss');

INSERT INTO movies (title, synopsis, director_id, duration, release_date, image, banner) VALUES
('Inception', 'A thief who steals corporate secrets...', 2, 148, '2010-07-16', 'inception.jpg', 'inception_banner.jpg'),
('Jurassic Park', 'A theme park showcasing cloned dinosaurs...', 1, 127, '1993-06-11', 'jurassic_park.jpg', 'jurassic_park_banner.jpg'),
('Spider-Man: Homecoming', 'Thrilled by his experience with the Avengers, Peter returns home, where he lives with his Aunt May, under the watchful eye of his new mentor Tony Stark, Peter tries to fall back into his normal daily routine - distracted by thoughts of proving himself to be more than just your friendly neighborhood Spider-Man - but when the Vulture emerges as a new villain, everything that Peter holds most important will be threatened', 1, 127, '1993-06-11', 'spiderman.jpg', 'spiderman_banner.jpg');


INSERT INTO movie_casts (name, movie_id) VALUES
('Leonardo DiCaprio', 1),
('Sam Neill', 2),
('Tom Holland', 3),
('Tom Holland', 3),
('Michael Keaton', 3),
('Robert Downey Jr',3);

INSERT INTO movie_genres (name, movie_id) VALUES
('Sci-Fi', 1),
('Adventure', 2),
('Action', 3),
('Adventure', 3);

insert into cinemas(name, image) values
('ebu.id','ebu.png'),
('hiflix','hiflix.png'),
('cineone21','cineone21.png');

INSERT INTO locations (city) VALUES
('Jakarta'),
('Bandung'),
('Medan');

INSERT INTO cinema_locations (cinema_id, location_id) VALUES
(1, 1), -- Cinema XXI di Jakarta
(2, 2), -- CGV di Bandung
(2,3),
(1,3),
(1,2);

INSERT INTO movie_schedules (movie_id, cinema_id, date, time) VALUES
(1, 1, '2025-01-05', '15:00:00'),
(2, 2, '2025-01-06', '18:00:00'),
(3, 3, '2025-01-02', '17:00:00');

INSERT INTO cinema_seats (cinema_id, seat_number, is_available) VALUES
(1, 'A1', true),
(1, 'A2', true),
(1, 'A3', true),
(2, 'B1', true),
(2, 'B2', true);

INSERT INTO transactions (user_id, movie_id, date, location, quantity, total_price, status) VALUES
(2, 1, '2025-01-05 15:00:00+07', 'Jakarta', 2, 100000.00, 'pending');

INSERT INTO transaction_seats (transaction_id, seat_number, cinema_id) VALUES
(1, 'A1', 1),
(1, 'A2', 1);

INSERT INTO payments (transaction_id, payment_method, virtual_account_number, amount, status) VALUES
(1, 'e-wallet', '1234567890', 100000.00, 'pending');
