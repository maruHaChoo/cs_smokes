DELETE FROM smokes
WHERE slug IN (
    'window_from_t_spawn',
    'jungle_from_ramp',
    'ct_from_tetris',
    'coffins_from_banana',
    'library_from_second_mid'
);

DELETE FROM maps
WHERE slug IN ('mirage', 'inferno');