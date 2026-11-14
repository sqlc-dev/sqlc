create table books (id int primary key, title text not null);
create table reviews (id int primary key, book_id int not null);
create view book_view as
  select b.id, (select r.id from reviews r where r.book_id = b.id limit 1) review_id,
    (select count(*) from reviews r where r.book_id = b.id) reviews
  from books b;
