-- Matchmaking lobbies

CREATE TABLE lobbies (
    id uuid primary key,
    owner uuid REFERENCES users(id),
    game game not null
);

ALTER TABLE sessions
    ADD COLUMN id uuid primary key DEFAULT gen_random_uuid();

CREATE TABLE lobby_players (
    lobby_id uuid references lobbies(id),
    player_id uuid references users(id),
    session uuid references sessions()
);