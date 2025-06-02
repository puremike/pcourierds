-- Re-add the CHECK constraint on vehicle_type
ALTER TABLE dispatchers
ADD CONSTRAINT dispatchers_vehicle_type_check CHECK (vehicle_type IN ('car', 'motorcycle'));

-- Re-add the CHECK constraint on vehicle_year
ALTER TABLE dispatchers
ADD CONSTRAINT dispatchers_vehicle_year_check CHECK (vehicle_year >= 2008);

-- Re-add DEFAULT constraints
ALTER TABLE dispatchers
ALTER COLUMN vehicle_type SET DEFAULT 'car';

ALTER TABLE dispatchers
ALTER COLUMN vehicle_year SET DEFAULT 2008;

ALTER TABLE dispatchers
ALTER COLUMN vehicle_model SET DEFAULT 'unknown';