create database involved;

alter database involved character set utf8 collate utf8_unicode_ci;

use involved;

create table types (
	type varchar(20) not null,
	primary key (type)
);

insert into types VALUES ("host"), ("seeker"), ("center"), ("organization");

create table user (
	id int not null AUTO_INCREMENT,
	name varchar(128) not null,
	password varchar(128) not null,
	email varchar(128) not null unique,
	pic_url varchar(256) default '',
	emailVerified int default 0,
	checkedInWith int default 0,
	belongsTo int default 0,
	description varchar(2048) default '',
	type varchar(20) not null,
	foreign key (type) references types(type) on delete cascade,
	primary key (id)
);

create table token (
	id int not null,
	value varchar(128) not null unique,
	created timestamp default current_timestamp,
	foreign key (id) references user(id) on delete cascade,
	primary key (id, value)
);

create table location (
	id int not null,
	latitude double precision,
	longitude double precision,
	foreign key (id) references user(id) on delete cascade,
	primary key (id)
);

create table room (
	id int not null AUTO_INCREMENT,
	owner int not null,
	numBeds int default 1,
	numDays int default 1,
	offersDinner int default 0,
	offersBreakfast int default 0,
	offersLunch int default 0,
	alcoholPresent int default 0,
	created timestamp default current_timestamp,
	foreign key (owner) references user(id) on delete cascade,
	primary key (id)
);

create table request (
	requester int not null,
	room int not null,
	accepted int default 0,
	foreign key (requester) references user(id) on delete cascade,
	foreign key (room) references room(id) on delete cascade,
	primary key (requester, room)
);

create event KILL_TOKENS
	on schedule every 1 hour
    comment 'removes old login tokens'
do
    delete from token
	where created < date_sub(current_timestamp, interval 180 day);