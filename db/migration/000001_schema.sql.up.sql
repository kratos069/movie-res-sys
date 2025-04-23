CREATE TABLE "users" (
  "user_id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "name" text NOT NULL,
  "email" text UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "role" varchar NOT NULL DEFAULT 'customer',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "genres" (
  "genre_id" serial PRIMARY KEY,
  "name" text UNIQUE NOT NULL
);

CREATE TABLE "movies" (
  "movie_id" serial PRIMARY KEY,
  "title" text NOT NULL,
  "description" text NOT NULL,
  "poster_url" text NOT NULL,
  "genre_id" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "showtimes" (
  "showtime_id" serial PRIMARY KEY,
  "movie_id" int NOT NULL,
  "start_time" timestamp,
  "price" numeric(5,2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "seats" (
  "seat_id" serial PRIMARY KEY,
  "row" int NOT NULL,
  "number" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "reservations" (
  "reservation_id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "showtime_id" int NOT NULL,
  "seat_id" int NOT NULL,
  "reserved_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "reservations" ("showtime_id", "seat_id");
CREATE INDEX idx_users_email ON users(email);

COMMENT ON TABLE "seats" IS 'This table represents the fixed seat layout';

ALTER TABLE "movies" ADD FOREIGN KEY ("genre_id") REFERENCES "genres" ("genre_id");

ALTER TABLE "showtimes" ADD FOREIGN KEY ("movie_id") REFERENCES "movies" ("movie_id") ON DELETE CASCADE;

ALTER TABLE "reservations" ADD FOREIGN KEY ("showtime_id") REFERENCES "showtimes" ("showtime_id") ON DELETE CASCADE;

ALTER TABLE "reservations" ADD FOREIGN KEY ("seat_id") REFERENCES "seats" ("seat_id") ON DELETE CASCADE;

ALTER TABLE "reservations" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

-- Insert fixed genres
INSERT INTO genres (name) VALUES
  ('Action'),
  ('Drama'),
  ('Comedy'),
  ('Horror'),
  ('Sci-Fi'),
  ('Romance'),
  ('Thriller'),
  ('Documentary'),
  ('Animation'),
  ('Fantasy')
ON CONFLICT DO NOTHING;

-- Insert fixed seat layout (e.g., 5 rows × 10 seats each = 50 seats)
DO $$
DECLARE
  row_num INT;
  seat_num INT;
BEGIN
  FOR row_num IN 1..5 LOOP
    FOR seat_num IN 1..10 LOOP
      INSERT INTO seats (row, number)
      VALUES (row_num, seat_num);
    END LOOP;
  END LOOP;
END $$;