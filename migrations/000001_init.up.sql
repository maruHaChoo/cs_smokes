CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    username TEXT,
    first_name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_sessions (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    chat_id BIGINT NOT NULL,
    menu_message_id INTEGER NOT NULL DEFAULT 0,
    content_message_id INTEGER,
    current_screen TEXT NOT NULL DEFAULT 'main',
    current_map_slug TEXT NOT NULL DEFAULT '',
    current_mode TEXT NOT NULL DEFAULT '',
    current_zone_slug TEXT NOT NULL DEFAULT '',
    current_target_slug TEXT NOT NULL DEFAULT '',
    current_smoke_slug TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE maps (
    id BIGSERIAL PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE smokes (
    id BIGSERIAL PRIMARY KEY,
    map_id BIGINT NOT NULL REFERENCES maps(id) ON DELETE CASCADE,
    zone_slug TEXT NOT NULL,
    target_slug TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    from_position TEXT NOT NULL,
    to_position TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    video_file_id TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_smokes_map_id ON smokes(map_id);
CREATE INDEX idx_smokes_map_zone ON smokes(map_id, zone_slug) WHERE is_active = TRUE;
CREATE INDEX idx_smokes_map_target ON smokes(map_id, target_slug) WHERE is_active = TRUE;
CREATE INDEX idx_smokes_slug ON smokes(slug);
CREATE INDEX idx_maps_sort_order ON maps(sort_order);

INSERT INTO maps (slug, title, sort_order) VALUES
    ('mirage', 'Mirage', 10),
    ('inferno', 'Inferno', 20);

INSERT INTO smokes (map_id, zone_slug, target_slug, slug, title, from_position, to_position, description, video_file_id)
SELECT m.id, x.zone_slug, x.target_slug, x.slug, x.title, x.from_position, x.to_position, x.description, x.video_file_id
FROM maps m
JOIN (
    VALUES
        ('mirage', 'mid', 'window', 'window_from_t_spawn', 'Window from T-spawn', 'T-spawn', 'Window', 'Быстрый стандартный смок в окно.', 'REPLACE_WITH_REAL_FILE_ID'),
        ('mirage', 'a', 'jungle', 'jungle_from_ramp', 'Jungle from Ramp', 'Ramp', 'Jungle', 'Классический выходной смок в джангл.', 'REPLACE_WITH_REAL_FILE_ID'),
        ('mirage', 'a', 'ct', 'ct_from_tetris', 'CT from Tetris', 'Tetris', 'CT', 'Выходной смок в КТ на A site.', 'REPLACE_WITH_REAL_FILE_ID'),
        ('inferno', 'b', 'coffins', 'coffins_from_banana', 'Coffins from Banana', 'Banana', 'Coffins', 'Смок в гробы с банана.', 'REPLACE_WITH_REAL_FILE_ID'),
        ('inferno', 'a', 'library', 'library_from_second_mid', 'Library from Second Mid', 'Second Mid', 'Library', 'Смок в библиотеку для выхода на A.', 'REPLACE_WITH_REAL_FILE_ID')
) AS x(map_slug, zone_slug, target_slug, slug, title, from_position, to_position, description, video_file_id)
    ON x.map_slug = m.slug;
