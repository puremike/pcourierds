-- Remove the CHECK constraint on status
ALTER TABLE dispatchers_apply
DROP CONSTRAINT check_status_valid_values;

-- Rename driver_license back to license
ALTER TABLE dispatchers_apply
RENAME COLUMN driver_license TO license;

-- Drop vehicle_model column
ALTER TABLE dispatchers_apply
DROP COLUMN vehicle_model;

-- Rename vehicle_plate_number back to vehicle
ALTER TABLE dispatchers_apply
RENAME COLUMN vehicle_plate_number TO vehicle;

-- Drop vehicle_year column
ALTER TABLE dispatchers_apply
DROP COLUMN vehicle_year;

-- Drop vehicle_type column
ALTER TABLE dispatchers_apply
DROP COLUMN vehicle_type;