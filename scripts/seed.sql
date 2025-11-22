-- Datos de prueba para desarrollo
-- Ejecutar con: make db-seed

-- Limpiar datos existentes (comentar si quieres preservar datos)
TRUNCATE users CASCADE;

-- Insertar usuarios de prueba
INSERT INTO users (email, name, created_at, updated_at) VALUES
    ('alice@example.com', 'Alice Johnson', NOW(), NOW()),
    ('bob@example.com', 'Bob Smith', NOW(), NOW()),
    ('charlie@example.com', 'Charlie Brown', NOW(), NOW()),
    ('diana@example.com', 'Diana Prince', NOW(), NOW()),
    ('eve@example.com', 'Eve Anderson', NOW(), NOW());

-- Verificar inserción
SELECT COUNT(*) as total_users FROM users;
SELECT * FROM users LIMIT 5;
