CREATE TABLE IF NOT EXISTS examples (
                                        example_id  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
                                        value       numeric NOT NULL,
                                        create_time timestamp with time zone NOT NULL DEFAULT now(),
                                        update_time timestamp with time zone NOT NULL DEFAULT now()
);