-- Insert into FacilityAddress
INSERT INTO
    facility_addresses (address, longitude, latitude)
VALUES
    ('Улица Абая, 50, Алматы', 76.9497, 43.2383);

-- Insert into Units
INSERT INTO
    units (name, conversion_factor)
VALUES
    ('Килограмм', 1.0),
    ('Грамм', 0.001),
    ('Литр', 1.0),
    ('Миллилитр', 0.001);

-- Insert into Regions
INSERT INTO
    regions (name)
VALUES
    ('Алматы');

-- Insert into CityWarehouses
INSERT INTO
    warehouses (facility_address_id, name, region_id)
VALUES
    (
        (
            SELECT
                id
            FROM
                facility_addresses
            WHERE
                address = 'Улица Абая, 50, Алматы'
        ),
        'Алматинский склад',
        (
            SELECT
                id
            FROM
                regions
            WHERE
                name = 'Алматы'
        )
    );

-- Insert into Store
INSERT INTO
    stores (
        name,
        facility_address_id,
        franchisee_id,
        warehouse_id,
        is_active,
        contact_phone,
        contact_email,
        store_hours,
        last_inventory_sync_at,
        created_at,
        updated_at
    )
VALUES
    (
        'ZEEP',
        1,
        NULL,
        1,
        true,
        '+79001112233',
        'zeep@example.com',
        '6:00-23:00',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- Insert into Employee
INSERT INTO
    employees (
        first_name,
        last_name,
        phone,
        email,
        is_active,
        hashed_password,
        created_at,
        updated_at
    )
VALUES
    (
        'GoTech',
        'Admin',
        '+79551234567',
        'gotech.admin@example.com',
        true,
        '$2a$10$USFR9OJqElj/xVOht51o1eENreY6aWGu23bQwHBqim3xzix/3kRt6', -- CaramelLatte54# 
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Асет',
        'Акбар',
        '+79667778899',
        'aset.akbar@example.com',
        true,
        '$2a$10$USFR9OJqElj/xVOht51o1eENreY6aWGu23bQwHBqim3xzix/3kRt6', -- CaramelLatte54#
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );


-- Store Employees Table
INSERT INTO
    store_employees (
        employee_id,
        store_id,
        role,
        created_at,
        updated_at
    )
VALUES
    (
        2,
        1,
        'STORE_MANAGER',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- Admin Employees Table
INSERT INTO
    admin_employees (employee_id, role, created_at, updated_at)
VALUES
    (1, 'ADMIN', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);