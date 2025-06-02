-- Drop the CHECK constraint on vehicle_type
ALTER TABLE dispatchers
DROP CONSTRAINT IF EXISTS dispatchers_vehicle_type_check;

-- Drop the CHECK constraint on vehicle_year
ALTER TABLE dispatchers
DROP CONSTRAINT IF EXISTS dispatchers_vehicle_year_check;

-- Drop DEFAULT constraints
ALTER TABLE dispatchers
ALTER COLUMN vehicle_type DROP DEFAULT;

ALTER TABLE dispatchers
ALTER COLUMN vehicle_year DROP DEFAULT;

ALTER TABLE dispatchers
ALTER COLUMN vehicle_model DROP DEFAULT;