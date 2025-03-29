DO $$
DECLARE
tables TEXT[] := ARRAY[
        'supplier_materials',
        'warehouse_stocks',
        'supplier_warehouse_deliveries',
        'stock_material_packages',
        'suppliers',
        'stock_request_ingredients',
        'stock_materials',
        'stock_requests',
        'suborder_additives',
        'suborders',
        'orders',
        'bonuses',
        'customer_addresses',
        'verification_codes',
        'referrals',
        'employee_workdays',
        'employee_audits',
        'employee_work_tracks',
        'warehouse_employees',
        'store_employees',
        'employees',
        'customers',
        'store_warehouse_stocks',
        'store_warehouses',
        'warehouses',
        'additive_ingredients',
        'product_size_ingredients',
        'ingredients',
        'product_size_additives',
        'store_product_sizes',
        'store_products',
        'store_additives',
        'stores',
        'additives',
        'product_sizes',
        'recipe_steps',
        'products',
        'stock_material_categories',
        'ingredient_categories',
        'additive_categories',
        'product_categories',
        'units',
        'facility_addresses',
        'regions',
        'region_managers',
        'franchisees',
        'franchisee_employees',
        'admin_employees'
    ];
BEGIN
    -- Drop all tables in reverse order
    FOREACH table_name IN ARRAY tables
    LOOP
        EXECUTE format('DROP TABLE IF EXISTS %I CASCADE', table_name);
END LOOP;

    -- Drop custom domain
DROP DOMAIN IF EXISTS valid_phone;
END $$;
