create table events (
  id int,
  event_type text not null,
  created_at timestamptz
);

CREATE MATERIALIZED VIEW something AS
select * from events
where event_type = 'sale'
order by created_at desc;

create schema computed_tables;
alter materialized view something set schema computed_tables;
