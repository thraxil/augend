CREATE TABLE users (
			 username varchar(256) primary key,
			 password varchar(256)
);

CREATE TABLE facts (
			 id uuid primary key,
			 title text,
			 details text,
			 source_name text,
			 source_url text,
			 added timestamp default current_timestamp
);

CREATE TABLE tags (
			 slug varchar(256) primary key,
			 tagname varchar(256)
);

CREATE TABLE fact_tags (
			 fact_uuid uuid references facts (id) on delete cascade,
			 tag_slug varchar(256) references tags (slug) on delete cascade
);

CREATE INDEX fact_tags_fact_idx on fact_tags (fact_uuid);
CREATE INDEX fact_tags_tag_idx on fact_tags (tag_slug);
