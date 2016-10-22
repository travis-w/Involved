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
	user_id int not null,
	value varchar(128) not null unique,
	created timestamp default current_timestamp,
	foreign key (user_id) references user(id) on delete cascade,
	primary key (user_id, value)
);

create table user_location (
	user_id int not null,
	latitude double precision,
	longitude double precision,
	foreign key (user_id) references user(id) on delete cascade,
	primary key (user_id)
);

create table seeker_dependent (
	user_id int not null,
	sub_id int default 1,
	name varchar(128) not null,
	foreign key (user_id) references user(id),
	primary key (user_id, sub_id)
);

create table user_meta (
	user_id int not null,
	sub_id int default 0,
	meta_key varchar(32) not null,
	value varchar(256) default '',
	foreign key (user_id) references user(id) on delete cascade,
	primary key (user_id, meta_key)
);

create table event (
	event_id int not null AUTO_INCREMENT,
	user_id int not null,
	availableSlots int default 1,
	maximumDivisions int default 1,
	description varchar(2048) default '',
	created timestamp default current_timestamp,
	type varchar(20) default '',
	foreign key (user_id) references user(id) on delete cascade,
	primary key (event_id)
);

create table event_location (
	event_id int not null,
	latitude double precision,
	longitude double precision,
	foreign key (event_id) references event(event_id) on delete cascade,
	primary key (event_id)
);

create table event_meta (
	event_id int not null,
	meta_key varchar(32) not null,
	value varchar(256) default '',
	isNeed tinyint(1) default 0,
	foreign key (event_id) references event(event_id) on delete cascade,
	primary key (event_id, meta_key)
);

create table seeker_event_response (
	event_id int not null,
	user_id int not null,
	accepted tinyint(1) default 0,
	count int default -1,
	foreign key (event_id) references event(event_id),
	foreign key (user_id) references user(id),
	primary key (event_id, user_id)
);

create table host_event_response (
	event_id int not null,
	user_id int not null,
	accepted tinyint(1) default 0,
	meta_key varchar(32) not null,
	value varchar(256) default '',
	foreign key (event_id, meta_key) references event_meta(event_id, meta_key),
	foreign key (user_id) references user(id),
	primary key (event_id, user_id, meta_key)
);

create event KILL_TOKENS
	on schedule every 1 hour
    comment 'removes old login tokens'
do
    delete from token
	where created < date_sub(current_timestamp, interval 180 day);