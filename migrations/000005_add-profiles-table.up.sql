CREATE TYPE PROFILESTATUS AS ENUM ('ONLINE','BUSY','IDLE','OFFLINE');

CREATE TABLE IF NOT EXISTS profiles(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(30) NOT NULL,
    image_url TEXT,
    banner_url TEXT,
    bio Text,
    status PROFILESTATUS NOT NULL DEFAULT 'OFFLINE',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);