-- Step 1: Drop the foreign key constraint
ALTER TABLE users DROP CONSTRAINT fk_users_roles;

-- Step 2: Create a temporary column to store the role titles
ALTER TABLE users ADD COLUMN temp_role VARCHAR(100);

-- Step 3: Join with roles table to get the titles
UPDATE users SET temp_role = roles.title
    FROM roles
WHERE users.role_id = roles.id;

-- Step 4: Rename the column back to its original name
ALTER TABLE users RENAME COLUMN role_id TO role;

-- Step 5: Copy the role titles back to the original column
UPDATE users SET role = temp_role;

-- Step 6: If the role column was NOT NULL before, restore that constraint
-- ALTER TABLE users ALTER COLUMN role SET NOT NULL;

-- Step 7: Drop the temporary column
ALTER TABLE users DROP COLUMN temp_role;