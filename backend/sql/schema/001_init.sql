CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,
    company_name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE facilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    coordinates geography(POINT, 4326) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE waste_streams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID REFERENCES facilities(id) ON DELETE CASCADE,
    primary_chemical TEXT NOT NULL,
    purity_percentage DECIMAL(5,2) NOT NULL,
    tonnage_available INTEGER NOT NULL,
    local_landfill_fee_per_ton DECIMAL(10,2) NOT NULL,
    lab_report_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE buyer_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID REFERENCES facilities(id) ON DELETE CASCADE,
    required_chemical TEXT NOT NULL,
    minimum_purity DECIMAL(5,2) NOT NULL,
    max_acceptable_distance_meters INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    waste_stream_id UUID NOT NULL REFERENCES waste_streams(id),
    buyer_requirement_id UUID NOT NULL REFERENCES buyer_requirements(id),
    tonnage_exchanged INT NOT NULL,
    freight_cost_estimated DECIMAL(10,2) NOT NULL,
    net_savings_estimated DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'locked',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
