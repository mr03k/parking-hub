ALTER TABLE rings
    ADD COLUMN object_id BIGINT,
    ADD COLUMN shape_leng BIGINT,
    ADD COLUMN shape_area BIGINT;


ALTER TABLE segments
    ADD COLUMN object_id BIGINT,
    ADD COLUMN junction SMALLINT,
    ADD COLUMN shape_leng SMALLINT,
    ADD COLUMN shape_area SMALLINT;

CREATE INDEX idx_segments_object_id ON segments(object_id);


ALTER TABLE roads
    ADD COLUMN buffer      BIGINT           NOT NULL,
    ADD COLUMN shape_leng  DOUBLE PRECISION NOT NULL,
    ADD COLUMN shape_area  DOUBLE PRECISION NOT NULL;