INSERT INTO maps (slug, title)
VALUES
    ('mirage', 'Mirage'),
    ('inferno', 'Inferno')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO smokes (
    map_slug,
    zone_slug,
    target_slug,
    slug,
    title,
    from_position,
    to_position,
    description,
    video_file_id,
    is_active
)
VALUES
    (
        'mirage',
        'mid',
        'window',
        'window_from_t_spawn',
        'Window from T-spawn',
        'T-spawn',
        'Window',
        'Быстрый стандартный смок в окно.',
        'REPLACE_WITH_REAL_FILE_ID',
        TRUE
    ),
    (
        'mirage',
        'a',
        'jungle',
        'jungle_from_ramp',
        'Jungle from Ramp',
        'Ramp',
        'Jungle',
        'Классический выходной смок в джангл.',
        'REPLACE_WITH_REAL_FILE_ID',
        TRUE
    ),
    (
        'mirage',
        'a',
        'ct',
        'ct_from_tetris',
        'CT from Tetris',
        'Tetris',
        'CT',
        'Выходной смок в КТ на A site.',
        'REPLACE_WITH_REAL_FILE_ID',
        TRUE
    ),
    (
        'inferno',
        'b',
        'coffins',
        'coffins_from_banana',
        'Coffins from Banana',
        'Banana',
        'Coffins',
        'Смок в гробы с банана.',
        'REPLACE_WITH_REAL_FILE_ID',
        TRUE
    ),
    (
        'inferno',
        'a',
        'library',
        'library_from_second_mid',
        'Library from Second Mid',
        'Second Mid',
        'Library',
        'Смок в библиотеку для выхода на A.',
        'REPLACE_WITH_REAL_FILE_ID',
        TRUE
    );