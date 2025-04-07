CREATE TABLE IF NOT EXISTS metrics(
                                    id VARCHAR (50) UNIQUE NOT NULL,
                                    type VARCHAR (50)  NOT NULL,
                                    value double precision ,
                                    delta bigint
);