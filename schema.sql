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
