-- SpotSync initial schema migration
-- Run manually if not using GORM AutoMigrate

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'driver',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

CREATE TABLE IF NOT EXISTS parking_zones (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    location VARCHAR(255) NOT NULL,
    description VARCHAR(500),
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    hourly_rate DECIMAL(10,2) NOT NULL CHECK (hourly_rate >= 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_parking_zones_name ON parking_zones(name);
CREATE INDEX IF NOT EXISTS idx_parking_zones_location ON parking_zones(location);
CREATE INDEX IF NOT EXISTS idx_parking_zones_is_active ON parking_zones(is_active);
CREATE INDEX IF NOT EXISTS idx_parking_zones_deleted_at ON parking_zones(deleted_at);

CREATE TABLE IF NOT EXISTS reservations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    parking_zone_id INTEGER NOT NULL REFERENCES parking_zones(id),
    vehicle_number VARCHAR(20) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    total_cost DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT chk_reservation_time CHECK (end_time > start_time)
);

CREATE INDEX IF NOT EXISTS idx_reservations_user_id ON reservations(user_id);
CREATE INDEX IF NOT EXISTS idx_reservations_parking_zone_id ON reservations(parking_zone_id);
CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations(status);
CREATE INDEX IF NOT EXISTS idx_reservations_start_time ON reservations(start_time);
CREATE INDEX IF NOT EXISTS idx_reservations_end_time ON reservations(end_time);
CREATE INDEX IF NOT EXISTS idx_reservations_deleted_at ON reservations(deleted_at);

-- Composite index for active reservation overlap queries
CREATE INDEX IF NOT EXISTS idx_reservations_zone_status_time
    ON reservations(parking_zone_id, status, start_time, end_time)
    WHERE deleted_at IS NULL;
