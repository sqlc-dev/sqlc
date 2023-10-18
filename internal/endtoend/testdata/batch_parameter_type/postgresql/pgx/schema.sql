CREATE TABLE public.solar_commcard_mapping (
    id
        INT8 NOT NULL,
    "deviceId"
        INT8 NOT NULL,
    version
        VARCHAR(32) DEFAULT ''::VARCHAR NOT NULL,
    sn
        VARCHAR(32) DEFAULT ''::VARCHAR NOT NULL,
    "createdAt"
        TIMESTAMPTZ DEFAULT now(),
    "updatedAt"
        TIMESTAMPTZ DEFAULT now()
);
