SELECT sqlc.embed(users), sqlc.embed(pets) FROM users JOIN pets ON true;
