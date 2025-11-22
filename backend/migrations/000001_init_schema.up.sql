-- Migración inicial: Crear tabla de ejemplo
-- Esta es una migración de ejemplo. Modifícala según tus necesidades.

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Crear índice en email para búsquedas rápidas
CREATE INDEX idx_users_email ON users(email);
