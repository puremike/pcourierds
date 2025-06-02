-- Drop the CHECK constraint on vehicle_type
ALTER TABLE dispatchers_apply
DROP CONSTRAINT IF EXISTS dispatchers_apply_vehicle_type_check;

-- Drop the CHECK constraint on vehicle_year
ALTER TABLE dispatchers_apply
DROP CONSTRAINT IF EXISTS dispatchers_apply_vehicle_year_check;

-- Drop the CHECK constraint on status
ALTER TABLE dispatchers_apply
DROP CONSTRAINT IF EXISTS check_status;