ALTER TABLE rings
DROP COLUMN object_id,
    DROP COLUMN shape_leng,
    DROP COLUMN shape_area;



DROP INDEX IF EXISTS idx_segments_object_id;

ALTER TABLE segments
DROP COLUMN IF EXISTS object_id,
    DROP COLUMN IF EXISTS junction,
    DROP COLUMN IF EXISTS shape_leng,
    DROP COLUMN IF EXISTS shape_area;


ALTER TABLE roads
    DROP COLUMN IF EXISTS buffer      BIGINT           NOT NULL,
    DROP COLUMN IF EXISTS shape_leng  DOUBLE PRECISION NOT NULL,
    DROP COLUMN IF EXISTS shape_area  DOUBLE PRECISION NOT NULL;