-- Add vehicle_type column
ALTER TABLE dispatchers_apply
ADD COLUMN vehicle_type TEXT NOT NULL DEFAULT 'car' CHECK (vehicle_type IN ('car', 'motorcycle'));

ALTER TABLE dispatchers_apply
ADD COLUMN vehicle_year INTEGER NOT NULL DEFAULT 2008 CHECK (vehicle_year >= 2008);

-- Rename vehicle to vehicle_plate_number
ALTER TABLE dispatchers_apply
RENAME COLUMN vehicle TO vehicle_plate_number;

-- Add vehicle_model column
ALTER TABLE dispatchers_apply
ADD COLUMN vehicle_model TEXT NOT NULL DEFAULT 'unknown';

-- Rename license to driver_license and make it nullable
ALTER TABLE dispatchers_apply
RENAME COLUMN license TO driver_license;

-- Ensure status has valid values
ALTER TABLE dispatchers_apply
ADD CONSTRAINT check_status CHECK (status IN ('pending', 'approved', 'rejected'));