-- Insert into FacilityAddress
-- Only 4 addresses in Kazakhstan
INSERT INTO
  facility_addresses (address, longitude, latitude)
VALUES
  ('Улица Абая, 50, Алматы', 76.9497, 43.2383),
  (
    'Проспект Назарбаева, 10, Астана',
    71.4278,
    51.1694
  ),
  ('Улица Ауэзова, 12, Шымкент', 69.5983, 42.3417),
  (
    'Проспект Республики, 25, Караганда',
    73.0948,
    49.8028
  );

-- Insert into Units
-- Unit measures remain unchanged
INSERT INTO
  units (name, conversion_factor)
VALUES
  ('Килограмм', 1.0),
  ('Грамм', 0.001),
  ('Штука', 1.0),
  ('Комплект', 1.0);

-- Insert into Warehouses
-- Only 4 warehouses linked to the new facility addresses
INSERT INTO
  warehouses (facility_address_id, name)
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
    'Алматинский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Проспект Назарбаева, 10, Астана'
    ),
    'Астанинский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Улица Ауэзова, 12, Шымкент'
    ),
    'Шымкентский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Проспект Республики, 25, Караганда'
    ),
    'Карагандинский склад'
  );

INSERT INTO
  ingredient_categories (name, description)
VALUES
  (
    'Электронные компоненты',
    'Микросхемы, транзисторы, резисторы и т.д.'
  ),
  (
    'Крепёж',
    'Болты, гайки, дюбели, саморезы и другие крепёжные изделия'
  ),
  (
    'Кабели и провода',
    'Электрические кабели и провода различных сечений'
  ),
  (
    'Листовой металл',
    'Металлические листы, пластины и другие материалы'
  ),
  (
    'Инструменты',
    'Ручные и электрические инструменты для ремонта и сборки'
  ),
  (
    'Расходные материалы',
    'Изоленты, масла, расходные материалы для производства'
  );

INSERT INTO
  ingredients (
    name,
    calories,
    fat,
    carbs,
    proteins,
    expiration_in_days,
    unit_id,
    category_id
  )
VALUES
  -- Technical items connected to updated categories
  ('Микросхема ATmega328', 0, 0, 0, 0, 3650, 3, 1), -- Электронные компоненты
  ('Болты М8', 0, 0, 0, 0, 1825, 3, 2), -- Крепёж
  ('Кабель ВВГ 3x2.5', 0, 0, 0, 0, 3650, 1, 3), -- Кабели и провода
  ('Листовая сталь 1.5 мм', 0, 0, 0, 0, 730, 1, 4), -- Листовой металл
  ('Набор сверл (10 шт.)', 0, 0, 0, 0, 0, 4, 5), -- Инструменты
  ('Резистор 10кОм (пачка)', 0, 0, 0, 0, 3650, 3, 1), -- Электронные компоненты
  ('Транзистор IRFZ44N', 0, 0, 0, 0, 3650, 3, 1), -- Электронные компоненты
  ('Разъём RJ-45 (упаковка)', 0, 0, 0, 0, 1095, 3, 1), -- Электронные компоненты
  ('Изолента ПВХ', 0, 0, 0, 0, 1095, 3, 6), -- Расходные материалы
  ('Масло для резки', 0, 0, 0, 0, 730, 1, 6);

-- Расходные материалы
-- Categories remain focused on technical goods, localized for Kazakhstan
INSERT INTO
  product_categories (name, description)
VALUES
  (
    'Электроника',
    'Электронные устройства и компоненты'
  ),
  (
    'Строительные материалы',
    'Материалы для строительства и ремонта'
  ),
  (
    'Кабели и провода',
    'Кабельная продукция для электроснабжения и связи'
  ),
  (
    'Компьютерное оборудование',
    'Компьютеры, мониторы и периферия'
  ),
  (
    'Офисная техника',
    'Принтеры, сканеры, копировальные аппараты'
  ),
  (
    'Запасные части',
    'Запчасти и аксессуары для различного оборудования'
  ),
  (
    'Промышленное оборудование',
    'Станки и производственное оборудование'
  ),
  (
    'Ручной инструмент',
    'Инструменты для ремонта и монтажа'
  ),
  (
    'Медицинское оборудование',
    'Аппаратура для медицинских учреждений'
  ),
  (
    'Химические реагенты',
    'Химические вещества для производства'
  ),
  (
    'Электродвигатели',
    'Двигатели и моторы для техники'
  ),
  (
    'Автозапчасти',
    'Комплектующие для автомобилей и спецтехники'
  ),
  (
    'Энергетическое оборудование',
    'Генераторы и трансформаторы'
  );

-- Insert into Products
INSERT INTO
  products (
    name,
    description,
    image_url,
    video_url,
    category_id
  )
VALUES
  -- 1 (category_id=2: "Строительные материалы")
  (
    'Цемент М500',
    'Высококачественный цемент для строительных работ',
    'https://example.com/images/cement_m500.png',
    'https://example.com/videos/cement_m500.mp4',
    2
  ),
  -- 2 (category_id=2: "Строительные материалы")
  (
    'Песок карьерный',
    'Крупнозернистый песок для изготовления бетонных смесей',
    'https://example.com/images/sand.png',
    NULL,
    2
  ),
  -- 3 (category_id=2: "Строительные материалы")
  (
    'Арматура 12 мм',
    'Стальная арматура для укрепления железобетонных конструкций',
    'https://example.com/images/rebar_12mm.png',
    NULL,
    2
  ),
  -- 4 (category_id=3: "Кабели и провода")
  (
    'Кабель ВВГнг 3x2.5',
    'Негорючий силовой кабель для внутренней проводки',
    'https://example.com/images/cable_vvgn.png',
    NULL,
    3
  ),
  -- 5 (category_id=3: "Кабели и провода")
  (
    'Провод ПВС 3x1.5',
    'Гибкий провод для подключения бытовых приборов',
    'https://example.com/images/pvs_3x1.5.png',
    NULL,
    3
  ),
  -- 6 (category_id=4: "Компьютерное оборудование")
  (
    'Ноутбук Lenovo ThinkPad',
    'Надёжный бизнес-ноутбук с длительным временем работы',
    'https://example.com/images/lenovo_thinkpad.png',
    'https://example.com/videos/thinkpad_review.mp4',
    4
  ),
  -- 7 (category_id=4: "Компьютерное оборудование")
  (
    'Монитор Samsung 24"',
    'Монитор с диагональю 24 дюйма и высоким разрешением',
    'https://example.com/images/samsung_24.png',
    NULL,
    4
  ),
  -- 8 (category_id=5: "Офисная техника")
  (
    'МФУ HP LaserJet',
    'Многофункциональное устройство для печати, сканирования и копирования',
    'https://example.com/images/hp_laserjet_mfp.png',
    NULL,
    5
  ),
  -- 9 (category_id=5: "Офисная техника")
  (
    'Ламинатор OfficeKit',
    'Компактный ламинатор для офисного использования',
    'https://example.com/images/laminator_officekit.png',
    NULL,
    5
  ),
  -- 10 (category_id=6: "Запасные части")
  (
    'Комплект ремня ГРМ',
    'Набор деталей для замены ремня ГРМ в двигателе',
    'https://example.com/images/timing_belt_kit.png',
    NULL,
    6
  ),
  -- 11 (category_id=7: "Промышленное оборудование")
  (
    'Токарный станок 16К20',
    'Универсальный токарный станок для металлообработки',
    'https://example.com/images/lathe_16k20.png',
    NULL,
    7
  ),
  -- 12 (category_id=8: "Ручной инструмент")
  (
    'Набор отверток',
    'Универсальный набор отверток для разных типов крепежа',
    'https://example.com/images/screwdriver_set.png',
    NULL,
    8
  ),
  -- 13 (category_id=9: "Медицинское оборудование")
  (
    'Аппарат УЗИ',
    'Ультразвуковой диагностический аппарат для медицинских учреждений',
    'https://example.com/images/ultrasound_machine.png',
    NULL,
    9
  ),
  -- 14 (category_id=10: "Химические реагенты")
  (
    'Соляная кислота (HCl)',
    'Химический реагент для различных производственных процессов',
    'https://example.com/images/hydrochloric_acid.png',
    NULL,
    10
  ),
  -- 15 (category_id=11: "Электродвигатели")
  (
    'Электродвигатель АИР80B2',
    'Надёжный трёхфазный двигатель мощностью 1.5 кВт',
    'https://example.com/images/electric_motor.png',
    NULL,
    11
  ),
  -- 16 (category_id=12: "Автозапчасти")
  (
    'Автомобильный аккумулятор 75Ач',
    'Аккумулятор повышенной ёмкости для легкового автомобиля',
    'https://example.com/images/car_battery_75ah.png',
    NULL,
    12
  ),
  -- 17 (category_id=12: "Автозапчасти")
  (
    'Масляный фильтр MANN',
    'Фильтр для очистки моторного масла',
    'https://example.com/images/mann_oil_filter.png',
    NULL,
    12
  ),
  -- 18 (category_id=7: "Промышленное оборудование")
  (
    'Пресс гидравлический',
    'Гидравлический пресс для обработки металла и сборки деталей',
    'https://example.com/images/hydraulic_press.png',
    NULL,
    7
  ),
  -- 19 (category_id=3: "Кабели и провода")
  (
    'Кабель КГ 4x10',
    'Гибкий кабель силовой для сварочного оборудования',
    'https://example.com/images/kg_cable.png',
    NULL,
    3
  ),
  -- 20 (category_id=2: "Строительные материалы")
  (
    'Гипсокартон 2500x1200x9.5',
    'Листовой материал для отделочных работ',
    'https://example.com/images/drywall.png',
    'https://example.com/videos/drywall_installation.mp4',
    2
  ),
  -- 21 (category_id=13: "Энергетическое оборудование")
  (
    'Генератор дизельный 5кВт',
    'Портативный дизельный генератор для резервного питания',
    'https://example.com/images/diesel_generator_5kw.png',
    NULL,
    13
  ),
  -- 22 (category_id=13: "Энергетическое оборудование")
  (
    'ИБП APC 1500VA',
    'Источник бесперебойного питания для защиты серверов и ПК',
    'https://example.com/images/apc_ups_1500.png',
    NULL,
    13
  );

-- Insert into product_sizes
INSERT INTO
  product_sizes (
    name,
    unit_id,
    base_price,
    size,
    is_default,
    product_id
  )
VALUES
  ('S', 1, 1000.00, 250, true, 1),
  ('M', 1, 1250.00, 350, false, 1),
  ('L', 1, 1500.00, 450, false, 1),
  ('S', 1, 1100.00, 300, true, 2),
  ('S', 1, 750.00, 200, true, 3),
  ('M', 1, 900.00, 300, false, 3),
  ('L', 2, 1150.00, 400, false, 3),
  ('S', 2, 600.00, 250, true, 4),
  ('M', 2, 850.00, 350, false, 4),
  ('L', 2, 1100.00, 450, false, 4),
  ('S', 2, 900.00, 250, true, 5),
  ('M', 2, 1100.00, 350, false, 5),
  ('L', 2, 1300.00, 450, false, 5),
  ('S', 2, 1500.00, 300, true, 6),
  ('L', 2, 1750.00, 500, false, 6),
  ('S', 2, 500.00, 200, true, 7),
  ('M', 2, 750.00, 300, false, 7),
  ('L', 3, 900.00, 400, false, 7),
  ('S', 3, 600.00, 300, true, 8),
  ('L', 3, 750.00, 500, false, 8),
  ('M', 3, 800.00, 350, true, 9),
  ('L', 3, 1000.00, 550, false, 9),
  ('S', 3, 500.00, 200, true, 10),
  ('M', 3, 700.00, 350, false, 10),
  ('S', 3, 600.00, 500, true, 11),
  ('S', 3, 900.00, 500, true, 18),
  ('S', 3, 400.00, 500, true, 13),
  ('S', 4, 900.00, 300, true, 20),
  ('S', 4, 1200.00, 300, true, 14),
  ('S', 4, 800.00, 300, true, 16),
  ('S', 4, 800.00, 300, true, 17),
  ('S', 4, 1400.00, 300, true, 15),
  ('S', 4, 1100.00, 300, true, 19),
  ('S', 4, 1250.00, 500, true, 12);

-- Insert into Store
INSERT INTO
  stores (
    name,
    facility_address_id,
    is_franchise,
    status,
    contact_phone,
    contact_email,
    store_hours,
    admin_id,
    created_at,
    updated_at
  )
VALUES
  (
    'Центральный склад',
    1, -- Ссылается на первую запись в facility_addresses
    false,
    'ACTIVE',
    '+77001112233', -- Телефонный код Казахстана
    'central@warehouse.kz',
    '8:00-20:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

-- Insert into Store Products
INSERT INTO
  store_products (product_id, store_id, is_available)
VALUES
  -- Все товары прикреплены к центральному складу
  (1, 1, true),
  (2, 1, true),
  (3, 1, true),
  (4, 1, true),
  (5, 1, true),
  (6, 1, true),
  (7, 1, true),
  (8, 1, true),
  (9, 1, true),
  (10, 1, true),
  (11, 1, true),
  (12, 1, true),
  (13, 1, true),
  (14, 1, true),
  (15, 1, true),
  (16, 1, true),
  (17, 1, true),
  (18, 1, true),
  (19, 1, true),
  (20, 1, true),
  (21, 1, true),
  (22, 1, true);

-- Insert into Store Product Sizes
INSERT INTO
  store_product_sizes (product_size_id, store_product_id, price)
VALUES
  -- Пример данных для единственного магазина
  (1, 1, 1000.00), -- Продукт 1, размер S
  (2, 1, 1250.00), -- Продукт 1, размер M
  (3, 1, 1500.00), -- Продукт 1, размер L
  (4, 2, 1100.00), -- Продукт 2, размер S
  (5, 3, 750.00), -- Продукт 3, размер S
  (6, 3, 900.00), -- Продукт 3, размер M
  (7, 3, 1150.00), -- Продукт 3, размер L
  (8, 4, 600.00), -- Продукт 4, размер S
  (9, 4, 850.00), -- Продукт 4, размер M
  (10, 4, 1100.00), -- Продукт 4, размер L
  (11, 5, 900.00), -- Продукт 5, размер S
  (12, 5, 1100.00), -- Продукт 5, размер M
  (13, 5, 1300.00), -- Продукт 5, размер L
  (14, 6, 1500.00), -- Продукт 6, размер S
  (15, 6, 1750.00), -- Продукт 6, размер L
  (16, 7, 500.00), -- Продукт 7, размер S
  (17, 7, 750.00), -- Продукт 7, размер M
  (18, 7, 900.00), -- Продукт 7, размер L
  (19, 8, 600.00), -- Продукт 8, размер S
  (20, 8, 750.00), -- Продукт 8, размер L
  (21, 9, 800.00), -- Продукт 9, размер M
  (22, 9, 1000.00), -- Продукт 9, размер L
  (23, 10, 500.00), -- Продукт 10, размер S
  (24, 10, 700.00), -- Продукт 10, размер M
  (25, 11, 600.00), -- Продукт 11, размер S
  (26, 12, 900.00), -- Продукт 12, размер S
  (27, 13, 400.00), -- Продукт 13, размер S
  (28, 14, 900.00), -- Продукт 14, размер S
  (29, 15, 1200.00), -- Продукт 15, размер S
  (30, 16, 800.00), -- Продукт 16, размер S
  (31, 17, 800.00), -- Продукт 17, размер S
  (32, 18, 1400.00), -- Продукт 18, размер S
  (33, 19, 1100.00), -- Продукт 19, размер S
  (34, 20, 1250.00);

-- Insert into Employees
INSERT INTO
  employees (
    first_name,
    last_name,
    phone,
    email,
    role,
    type,
    is_active,
    hashed_password,
    created_at,
    updated_at
  )
VALUES
  (
    'Елена',
    'Соколова',
    '+77001112233',
    'elena@company.kz',
    'ADMIN',
    'STORE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Сергей',
    'Павлов',
    '+77002223344',
    'sergey@company.kz',
    'MANAGER',
    'STORE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Мария',
    'Смирнова',
    '+77003334455',
    'maria@company.kz',
    'WAREHOUSE_EMPLOYEE',
    'WAREHOUSE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Олег',
    'Кузнецов',
    '+77004445566',
    'oleg@company.kz',
    'WAREHOUSE_EMPLOYEE',
    'WAREHOUSE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Анна',
    'Федорова',
    '+77005556677',
    'anna@company.kz',
    'MANAGER',
    'STORE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

-- Insert into Store Employees
INSERT INTO
  store_employees (
    employee_id,
    store_id,
    is_franchise,
    created_at,
    updated_at
  )
VALUES
  (1, 1, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Елена (ADMIN)
  (2, 1, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Сергей (MANAGER)
  (5, 1, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert into Warehouse Employees
INSERT INTO
  warehouse_employees (employee_id, warehouse_id, created_at, updated_at)
VALUES
  (3, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Мария (WAREHOUSE_EMPLOYEE)
  (4, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert store warehouses
INSERT INTO
  store_warehouses (store_id, warehouse_id)
VALUES
  (1, 1), -- Связь магазина с 1-м складом (Алматы)
  (1, 2), -- Связь магазина со 2-м складом (Астана)
  (1, 3), -- Связь магазина с 3-м складом (Шымкент)
  (1, 4);

-- Insert store warehouse stock
INSERT INTO
  store_warehouse_stocks (
    store_warehouse_id,
    ingredient_id,
    quantity,
    low_stock_threshold,
    created_at,
    updated_at
  )
VALUES
  -- Склад 1 (id=1), храним 2 вида материалов
  (
    1,
    1,
    100,
    10,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ), -- Микросхема XYZ
  (
    1,
    2,
    300,
    50,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ), -- Болты М8
  -- Склад 2 (id=2), храним ещё 2 вида материалов
  (
    2,
    3,
    500,
    20,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ), -- Кабель ВВГ 3x2.5
  (2, 4, 20, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Стальной лист
  -- Склад 3 (id=3), храним 1 материал
  (
    3,
    2,
    150,
    30,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ), -- Болты М8 (пример)
  -- Склад 4 (id=4), храним оставшийся материал
  (4, 5, 40, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert stock material categories
INSERT INTO
  stock_material_categories (name, description)
VALUES
  (
    'Электронные компоненты',
    'Микросхемы, резисторы, конденсаторы и т.д.'
  ),
  ('Крепёж', 'Винты, гайки, саморезы, дюбели и т.п.'),
  (
    'Кабели и провода',
    'Различные типы и сечения проводов и кабелей'
  ),
  (
    'Листовой металл',
    'Листовая сталь, алюминий и другие металлы'
  ),
  (
    'Инструменты',
    'Сверла, пилки, диски и другие расходные материалы для инструмента'
  );

-- Insert stock materials
INSERT INTO
  stock_materials (
    name,
    description,
    ingredient_id,
    safety_stock,
    unit_id,
    size,
    category_id,
    barcode,
    expiration_period_in_days,
    is_active
  )
VALUES
  (
    'Микросхема XYZ',
    'Высокочастотный микропроцессорный модуль',
    1, -- условный ID "ингредиента"
    30, -- минимальный страховой запас
    3, -- допустим, 3 = "Штука"
    1, -- размер (количество штук)
    1, -- категория 1 ("Электронные компоненты")
    '1111111111111',
    3650, -- ~10 лет условного срока
    TRUE
  ),
  (
    'Болты М8',
    'Набор болтов М8 с гайками',
    2,
    50,
    3, -- "Штука"
    100, -- в комплекте 100 штук
    2, -- категория 2 ("Крепёж")
    '2222222222222',
    0, -- нет срока хранения, ставим 0
    TRUE
  ),
  (
    'Кабель ВВГ 3x2.5',
    'Силовой кабель для внутренней электропроводки',
    3,
    10,
    1, -- 1 = "Килограмм" (или можно "Метр", если бы была такая ед.изм.)
    50, -- условная длина 50 м
    3, -- категория 3 ("Кабели и провода")
    '3333333333333',
    0,
    TRUE
  ),
  (
    'Стальной лист 1.5 мм',
    'Листовая сталь толщиной 1.5 мм',
    4,
    5,
    1, -- "Килограмм" (или другое)
    10, -- условно 10 кг/лист
    4, -- категория 4 ("Листовой металл")
    '4444444444444',
    0,
    TRUE
  ),
  (
    'Набор сверл по металлу',
    'Набор из 10 сверл разного диаметра',
    5,
    20,
    4, -- 4 = "Комплект"
    1, -- 1 комплект
    5, -- категория 5 ("Инструменты")
    '5555555555555',
    0,
    TRUE
  );

-- Insert into Suppliers
INSERT INTO
  suppliers (name, contact_email, contact_phone, city, address)
VALUES
  (
    'ТОО "KazTech Solutions"',
    'contact@kaztech.kz',
    '+77001112233',
    'Алматы',
    'пр. Абая, д. 15'
  ),
  (
    'ТОО "Almaty Krepyozh"',
    'info@krepyozh.kz',
    '+77002223344',
    'Алматы',
    'ул. Толе би, д. 45'
  ),
  (
    'ТОО "QazaqCable"',
    'sales@qazaqcable.kz',
    '+77003334455',
    'Астана',
    'пр. Туран, д. 10, оф. 101'
  ),
  (
    'ТОО "Astana Metal"',
    'info@ametall.kz',
    '+77004445566',
    'Астана',
    'ул. Кабанбай батыра, д. 2'
  ),
  (
    'ТОО "MegaInstrument"',
    'support@megainst.kz',
    '+77005556677',
    'Шымкент',
    'ул. Рыскулова, д. 87'
  ),
  (
    'ТОО "AutoParts KZ"',
    'parts@autokz.kz',
    '+77006667788',
    'Караганда',
    'пр. Республики, д. 40'
  ),
  (
    'ТОО "Resistor Co."',
    'contact@resistor.kz',
    '+77007778899',
    'Алматы',
    'ул. Сейфуллина, д. 120'
  ),
  (
    'ТОО "CondenserGroup"',
    'info@condenser.kz',
    '+77008889900',
    'Астана',
    'ул. Достык, д. 14'
  ),
  (
    'ТОО "StanKaz"',
    'stan@kzindustry.kz',
    '+77009990011',
    'Актобе',
    'пр. Молдагуловой, д. 30'
  ),
  (
    'ТОО "GenIndustrial"',
    'support@genindus.kz',
    '+77010001122',
    'Алматы',
    'ул. Желтоксан, д. 56'
  );

-- Insert into supplier_warehouse_deliveries
INSERT INTO
  supplier_warehouse_deliveries (
    supplier_id,
    warehouse_id,
    delivery_date,
    created_at,
    updated_at
  )
VALUES
  -- Поставка №1 от KazTech Solutions на склад №1 (Алматы)
  (
    1,
    1,
    '2024-09-01',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  -- Поставка №2 от QazaqCable на склад №2 (Астана)
  (
    3,
    2,
    '2024-10-15',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  -- Поставка №3 от MegaInstrument на склад №3 (Шымкент)
  (
    5,
    3,
    '2024-11-01',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  -- Поставка №4 от Astana Metal на склад №4 (Караганда)
  (
    4,
    4,
    '2024-12-01',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

-- Insert into supplier_warehouse_delivery_materials
INSERT INTO
  supplier_warehouse_delivery_materials (
    delivery_id,
    stock_material_id,
    barcode,
    quantity,
    expiration_date
  )
VALUES
  -- Поставка #1 (delivery_id=1) на склад #1
  (1, 1, '1111111111111', 100, '2030-12-31'), -- Микросхема XYZ
  (1, 2, '2222222222222', 200, '2030-12-31'), -- Болты М8
  -- Поставка #2 (delivery_id=2) на склад #2
  (2, 3, '3333333333333', 50, '2030-12-31'), -- Кабель ВВГ 3x2.5
  -- Поставка #3 (delivery_id=3) на склад #3
  (3, 5, '5555555555555', 30, '2030-12-31'), -- Набор сверл
  -- Поставка #4 (delivery_id=4) на склад #4
  (4, 4, '4444444444444', 25, '2030-12-31');

-- Insert into supplier_materials
INSERT INTO
  supplier_materials (
    supplier_id,
    stock_material_id,
    created_at,
    updated_at
  )
VALUES
  -- KazTech Solutions (id=1) поставляет Микросхему и Болты
  (1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Микросхема XYZ
  (1, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Болты М8
  -- Almaty Krepyozh (id=2) поставляет Болты М8 и Набор сверл
  (2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Болты
  (2, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Сверла
  -- QazaqCable (id=3) поставляет Кабель ВВГ
  (3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Кабель
  -- Astana Metal (id=4) поставляет Стальной лист
  (4, 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Лист
  -- MegaInstrument (id=5) поставляет Набор сверл и Микросхему
  (5, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (5, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert mock data into supplier_prices
INSERT INTO
  supplier_prices (supplier_material_id, base_price)
VALUES
  -- К примеру, поставщик #1 (KazTech Solutions) может поставлять
  -- микросхему (supplier_material_id = 1) и болты (supplier_material_id = 2)
  (1, 1500.00),
  (1, 1490.00),
  (2, 300.00),
  (2, 295.00),
  -- Допустим, поставщик #2 (Almaty Krepyozh) может поставлять
  -- кабель (supplier_material_id = 3) и сверла (supplier_material_id = 4)
  (3, 2000.00),
  (3, 1950.00),
  (4, 500.00),
  (4, 480.00),
  -- Поставщик #3 (QazaqCable), например, сталь (supplier_material_id = 5)
  (5, 1000.00),
  (5, 990.00);

-- Insert mock data into warehouse_stocks
INSERT INTO
  warehouse_stocks (warehouse_id, stock_material_id, quantity)
VALUES
  (1, 1, 100), -- Микросхема ATmega328 на Склада №1
  (1, 2, 200), -- Болты М8 на Склада №1
  (1, 3, 320), -- Микросхема ATmega328 на Склада №1
  (1, 4, 120), -- Болты М8 на Склада №1
  (1, 5, 30), -- Микросхема ATmega328 на Склада №1
  (2, 3, 50), -- Кабель ВВГ на Склада №2
  (3, 4, 20), -- Листовая сталь на Склада №3
  (4, 5, 15);

-- Insert into StockRequests (Initial Requests)
INSERT INTO
  stock_requests (
    store_id,
    warehouse_id,
    status,
    created_at,
    updated_at
  )
VALUES
  (
    1,
    1,
    'CREATED',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    1,
    1,
    'PROCESSED',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    1,
    1,
    'IN_DELIVERY',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    1,
    1,
    'COMPLETED',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO
  stock_request_ingredients (
    stock_request_id,
    ingredient_id,
    stock_material_id,
    quantity,
    created_at,
    updated_at
  )
VALUES
  -- Запрос #1 (store_id=1 -> warehouse_id=1)
  (
    1,
    1,
    1,
    10.0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ), -- запросить 10 шт. (Микросхем)
  (
    1,
    2,
    2,
    50.0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ), -- запросить 50 шт. (Болтов)
  -- Запрос #2 (store_id=1 -> warehouse_id=2)
  (
    2,
    3,
    3,
    20.0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ), -- запросить 20 м кабеля
  (
    2,
    2,
    2,
    100.0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ), -- ещё болты
  -- Запрос #3 (store_id=1 -> warehouse_id=3)
  (
    3,
    4,
    4,
    5.0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ), -- немного листовой стали
  (
    3,
    5,
    5,
    2.0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ), -- несколько наборов сверл
  -- Запрос #4 (store_id=1 -> warehouse_id=4)
  (
    4,
    1,
    1,
    30.0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ), -- микросхемы
  (
    4,
    3,
    3,
    10.0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );