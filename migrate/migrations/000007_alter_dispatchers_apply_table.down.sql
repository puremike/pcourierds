-- Re-add the CHECK constraint on vehicle_type
ALTER TABLE dispatchers_apply
ADD CONSTRAINT dispatchers_apply_vehicle_type_check CHECK (vehicle_type IN ('car', 'motorcycle'));

-- Re-add the CHECK constraint on vehicle_year
ALTER TABLE dispatchers_apply
ADD CONSTRAINT dispatchers_apply_vehicle_year_check CHECK (vehicle_year >= 2008);

-- Re-add the CHECK constraint on status
ALTER TABLE dispatchers_apply
ADD CONSTRAINT check_status CHECK (status IN ('pending', 'approved', 'rejected'));