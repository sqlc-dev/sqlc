CREATE TABLE vampires (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid ()
);

CREATE TABLE memories (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    vampire_id uuid REFERENCES vampires (id) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp
);
