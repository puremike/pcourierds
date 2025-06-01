-- Rename driver_license back to license
ALTER TABLE dispatchers
RENAME COLUMN driver_license TO license;

-- Rename vehicle_plate_number back to vehicle
ALTER TABLE dispatchers
RENAME COLUMN vehicle_plate_number TO vehicle;

-- Drop vehicle_model column
ALTER TABLE dispatchers
DROP COLUMN vehicle_model;

-- Drop vehicle_year column
ALTER TABLE dispatchers
DROP COLUMN vehicle_year;

-- Drop vehicle_type column
ALTER TABLE dispatchers
DROP COLUMN vehicle_type;