create schema metrics collate utf8_general_ci;

create table services
(
	date datetime not null,
	container_id varchar(255) not null,
	container_name varchar(255) not null,
	cpu_percentage float not null,
	memory_usage_mib float not null,
	memory_limit_mib float not null,
	memory_percentage float not null,
	primary key (date, container_id)
);
