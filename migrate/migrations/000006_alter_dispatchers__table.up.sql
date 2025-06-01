-- Add vehicle_type column with default and CHECK constraint
ALTER TABLE dispatchers
ADD COLUMN vehicle_type TEXT NOT NULL DEFAULT 'car' CHECK (vehicle_type IN ('car', 'motorcycle'));

-- Add vehicle_year column with default and CHECK constraint
ALTER TABLE dispatchers
ADD COLUMN vehicle_year INTEGER NOT NULL DEFAULT 2008 CHECK (vehicle_year >= 2008);

-- Add vehicle_model column with default
ALTER TABLE dispatchers
ADD COLUMN vehicle_model TEXT NOT NULL DEFAULT 'unknown';

-- Rename vehicle to vehicle_plate_number
ALTER TABLE dispatchers
RENAME COLUMN vehicle TO vehicle_plate_number;

-- Rename license to driver_license (remains NOT NULL)
ALTER TABLE dispatchers
RENAME COLUMN license TO driver_license;