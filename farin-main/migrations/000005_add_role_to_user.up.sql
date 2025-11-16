-- Step 1: Rename the column
ALTER TABLE users ADD COLUMN role_id UUID;

-- Step 6: Add a foreign key constraint to reference the roles table
ALTER TABLE users ADD CONSTRAINT fk_users_roles
    FOREIGN KEY (role_id) REFERENCES roles(id);

-- Step 7: Update users with Admin role to reference the Admin role UUID
UPDATE users SET role_id = (SELECT id FROM roles WHERE title = 'Admin')
WHERE role = 'admin';
-- -- Step 8: Update users with Driver role to reference the Driver role UUID
UPDATE users SET role_id = (SELECT id FROM roles WHERE title = 'Driver')
WHERE role = 'driver';

-- -- Step 9: Drop the temporary column
ALTER TABLE users DROP COLUMN role;