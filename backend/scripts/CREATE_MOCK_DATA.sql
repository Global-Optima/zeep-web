-- Insert into Customer
INSERT INTO
    customers (
        first_name,
        last_name,
        password,
        phone,
        is_verified,
        is_banned
    )
VALUES
    (
        'Иван',
        'Иванов',
        'hashed_password_123',
        '+79031234567',
        true,
        false
    ),
    (
        'Мария',
        'Смирнова',
        'hashed_password_456',
        '+79876543210',
        false,
        false
    ),
    (
        'Алексей',
        'Петров',
        'hashed_password_789',
        '+79998887766',
        true,
        false
    );

-- Insert into FacilityAddress
INSERT INTO
    facility_addresses (address, longitude, latitude)
VALUES
    ('Улица Абая, 50, Алматы', 76.9497, 43.2383),
    (
        'Проспект Республики, 25, Астана',
        73.0948,
        49.8028
    ),
    (
        'Улица Байтурсынова, 15, Алматы',
        76.9422,
        43.2528
    ),
    ('Улица Сарыарка, 12, Астана', 71.4100, 51.1690),
    ('Улица Панфилова, 98, Алматы', 76.9271, 43.2575),
    ('Улица Ш. Уалиханова, 7, Астана', 71.4451, 51.1),
    ('Улица Улы Дала, 10, Астана', 75.52, 45.20);

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
    ('Алматы'),
    ('Астана');

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
    ),
    (
        (
            SELECT
                id
            FROM
                facility_addresses
            WHERE
                address = 'Проспект Республики, 25, Астана'
        ),
        'Астанинский склад',
        (
            SELECT
                id
            FROM
                regions
            WHERE
                name = 'Астана'
        )
    );

-- Insert into ProductCategory
INSERT INTO
    product_categories (name, description)
VALUES
    ('Напитки', 'Различные виды напитков'),
    ('Кофе', 'Горячие кофейные напитки'),
    ('Чай', 'Различные виды чая'),
    ('Смузи', 'Фруктовые и овощные смузи'),
    ('Соки', 'Свежевыжатые соки и напитки'),
    (
        'Коктейли',
        'Алкогольные и безалкогольные коктейли'
    ),
    (
        'Газированные напитки',
        'Прохладительные газированные напитки'
    ),
    (
        'Энергетические напитки',
        'Напитки для повышения энергии'
    ),
    (
        'Молочные коктейли',
        'Коктейли на основе молока и сливок'
    ),
    (
        'Минеральная вода',
        'Природная и газированная минеральная вода'
    ),
    ('Фраппе', 'Кофейные напитки со льдом'),
    ('Травяные чаи', 'Настои и чаи на травах'),
    (
        'Круассаны',
        'Свежая выпечка с хрустящей корочкой и разнообразной начинкой — идеально к кофе'
    ),
    (
        'Тестовые продукты',
        'Продукты для тестирования функционала QR-кодов'
    );


-- Insert into AdditiveCategory
INSERT INTO
    additive_categories (name, description, is_multiple_select)
VALUES
    (
        'Ароматизаторы',
        'Дополнительные вкусы для усиления аромата',
        true
    ),
    ('Подсластители', 'Добавление сладости', false),
    (
        'Топпинги',
        'Украшения для десертов и напитков',
        true
    ),
    ('Сиропы', 'Сиропы для напитков и десертов', false),
    ('Специи', 'Ароматные специи для напитков', true),
    (
        'Молочные добавки',
        'Добавление молока и сливок',
        false
    ),
    ('Фрукты', 'Свежие и сушеные фрукты', true),
    (
        'Орехи',
        'Измельченные орехи для украшения',
        false
    ),
    ('Шоколад', 'Шоколадная стружка и какао', true),
    (
        'Мед',
        'Естественный подсластитель на основе меда',
        false
    ),
    ('Сахарные добавки', 'Различные виды сахара', true),
    (
        'Кубики льда',
        'Лед для охлаждения напитков',
        false
    );

-- Insert into IngredientCategories
INSERT INTO
    ingredient_categories (name, description)
VALUES
    (
        'Молочные продукты',
        'Категория для молока, сливок, и других молочных ингредиентов'
    ),
    (
        'Фрукты',
        'Категория для фруктов, таких как яблоки, бананы, апельсины'
    ),
    (
        'Подсластители',
        'Категория для сахара, мёда, сиропов'
    ),
    (
        'Специи',
        'Категория для специй, таких как корица, ваниль, мята, имбирь'
    ),
    (
        'Орехи и семена',
        'Категория для орехов, фисташек, и других семян'
    ),
    (
        'Шоколад и какао',
        'Категория для шоколада, какао-порошка, и других шоколадных продуктов'
    ),
    (
        'Охлаждающие добавки',
        'Категория для льда, мороженого, и других охлаждающих ингредиентов'
    );

-- Insert into Products
INSERT INTO
    products (
        name,
        description,
        image_key,
        video_key,
        category_id
    )
VALUES
    (
        'Латте',
        'Нежный кофе с молоком',
        'https://png.pngtree.com/png-vector/20240416/ourmid/pngtree-coffee-latte-seen-up-close-png-image_12286454.png',
        'https://example.com/videos/latte.mp4',
        2
    ),
    (
        'Капучино',
        'Кофе с густой пенкой',
        'https://www.pngplay.com/wp-content/uploads/7/Cappuccino-Coffee-PNG-Clipart-Background.png',
        'https://example.com/videos/cappuccino.mp4',
        2
    ),
    (
        'Эспрессо',
        'Крепкий кофе, подается в маленькой чашке',
        'https://static.vecteezy.com/system/resources/previews/023/438/448/non_2x/espresso-coffee-cutout-free-png.png',
        NULL,
        2
    ),
    (
        'Зеленый чай',
        'Свежий зеленый чай из лучших листьев',
        'https://static.vecteezy.com/system/resources/thumbnails/024/108/075/small_2x/fresh-herbal-tea-cup-with-green-leaves-isolated-on-transparent-background-png.png',
        NULL,
        3
    ),
    (
        'Черный чай',
        'Классический черный чай с насыщенным вкусом',
        'https://static.vecteezy.com/system/resources/thumbnails/037/280/520/small_2x/ai-generated-black-tea-held-in-a-glass-isolated-on-transparent-background-free-png.png',
        NULL,
        3
    ),
    (
        'Клубничный смузи',
        'Смузи с клубникой и бананом',
        'https://png.pngtree.com/png-vector/20240807/ourmid/pngtree-smooth-and-sweet-strawberry-smoothie-ready-to-serve-in-a-cool-png-image_13402854.png',
        NULL,
        4
    ),
    (
        'Манговый смузи',
        'Тропический смузи с манго',
        'https://static.vecteezy.com/system/resources/thumbnails/033/321/478/small_2x/mango-smoothie-in-a-glass-isolated-png.png',
        NULL,
        4
    ),
    (
        'Апельсиновый сок',
        'Свежевыжатый апельсиновый сок',
        'https://www.pngplay.com/wp-content/uploads/6/Orange-Flavor-Juice-PNG.png',
        NULL,
        5
    ),
    (
        'Яблочный сок',
        'Сок из спелых яблок',
        'https://static.vecteezy.com/system/resources/previews/027/145/688/non_2x/apple-juice-ice-surrounded-by-apples-and-leaves-ai-generated-png.png',
        NULL,
        5
    ),
    (
        'Мохито',
        'Освежающий коктейль с мятой и лаймом',
        'https://png.pngtree.com/png-clipart/20231119/original/pngtree-mojito-cocktail-herb-picture-image_13268789.png',
        NULL,
        6
    ),
    (
        'Кола',
        'Классический газированный напиток',
        'https://freepngimg.com/save/2438-coca-cola-bottle-png-image/906x906',
        NULL,
        7
    ),
    (
        'Энергетический напиток',
        'Энергетический напиток для бодрости',
        'https://pngimg.com/uploads/red_bull/red_bull_PNG29.png',
        NULL,
        8
    ),
    (
        'Шоколадный молочный коктейль',
        'Сладкий коктейль с шоколадным вкусом',
        'https://png.pngtree.com/png-vector/20240603/ourmid/pngtree-hyper-realistic-image-of-chocolate-milkshake-png-image_12614596.png',
        NULL,
        9
    ),
    (
        'Минеральная вода',
        'Природная минеральная вода без газа',
        'https://static.vecteezy.com/system/resources/previews/036/573/060/non_2x/cold-mineral-water-bottle-borjomi-free-png.png',
        NULL,
        10
    ),
    (
        'Лимонный фраппе',
        'Кофейный напиток с лимоном и льдом',
        'https://png.pngtree.com/png-vector/20240801/ourmid/pngtree-mocha-coffee-frappe-in-glass-png-image_13321780.png',
        NULL,
        11
    ),
    (
        'Имбирный чай',
        'Травяной чай с имбирем для бодрости',
        'https://static.vecteezy.com/system/resources/thumbnails/039/336/072/small_2x/ai-generated-a-glass-of-ginger-tea-isolated-on-transparent-background-png.png',
        NULL,
        12
    ),
    (
        'Чай с мятой',
        'Травяной чай с ароматом свежей мяты',
        'https://www.freepnglogos.com/uploads/tea-png/tea-latte-mix-cafe-17.png',
        NULL,
        12
    ),
    (
        'Малиновый лимонад',
        'Освежающий лимонад с малиной',
        'https://static.vecteezy.com/system/resources/thumbnails/036/257/929/small_2x/ai-generated-a-glass-of-strawberry-juice-isolated-on-transparent-background-free-png.png',
        NULL,
        7
    ),
    (
        'Клубничный чай',
        'Чай с добавлением клубники',
        'https://static.vecteezy.com/system/resources/previews/048/894/947/non_2x/a-strawberry-iced-tea-in-a-plastic-cup-transparent-background-png.png',
        NULL,
        3
    ),
    (
        'Карамельный латте',
        'Кофе с молоком и карамельным вкусом',
        'https://static.vecteezy.com/system/resources/thumbnails/027/145/750/small_2x/iced-caramel-latte-topped-with-whipped-cream-and-caramel-sauce-perfect-for-drink-catalog-ai-generated-png.png',
        'https://example.com/videos/caramel-latte.mp4',
        2
    ),
    (
        'Круассан с шоколадом',
        'Нежный хрустящий круассан с шоколадной начинкой',
        'https://www.pngplay.com/wp-content/uploads/15/Pain-Au-Chocolat-Transparent-PNG.png',
        NULL,
        13
    ),
    (
        'Круассан с миндалем',
        'Ароматный круассан с миндальной начинкой и посыпкой',
        'https://www.pngplay.com/wp-content/uploads/15/Croissants-Transparent-Image.png',
        NULL,
        13
    );

-- Insert into RecipeStep
INSERT INTO
    recipe_steps (product_id, step, name, description, image_key)
VALUES
    (
        1,
        1,
        'Приготовление эспрессо',
        'Сварить эспрессо',
        'https://example.com/images/espresso.jpg'
    ),
    (
        1,
        2,
        'Вспенивание молока',
        'Вспенить молоко до нужной текстуры',
        'https://example.com/images/steamed-milk.jpg'
    ),
    (
        2,
        1,
        'Приготовление пены',
        'Сделать густую пену для капучино',
        'https://example.com/images/foam.jpg'
    );

-- Insert into ProductSize
INSERT INTO
    product_sizes (
        name,
        unit_id,
        base_price,
        size,
        product_id,
        machine_id
    )
VALUES
    -- Все unit_id заменены на 4 (миллилитры)
    (
        'S',
        4,
        1000.00,
        250,
        1,
        'ZEEP0000000000000000000000000001'
    ),
    (
        'M',
        4,
        1250.00,
        350,
        1,
        'ZEEP0000000000000000000000000002'
    ),
    (
        'L',
        4,
        1500.00,
        450,
        1,
        'ZEEP0000000000000000000000000003'
    ),
    (
        'S',
        4,
        1100.00,
        300,
        2,
        'ZEEP0000000000000000000000000004'
    ),
    (
        'S',
        4,
        750.00,
        200,
        3,
        'ZEEP0000000000000000000000000005'
    ),
    (
        'M',
        4,
        900.00,
        300,
        3,
        'ZEEP0000000000000000000000000006'
    ),
    (
        'L',
        4,
        1150.00,
        400,
        3,
        'ZEEP0000000000000000000000000007'
    ),
    (
        'S',
        4,
        600.00,
        250,
        4,
        'ZEEP0000000000000000000000000008'
    ),
    (
        'M',
        4,
        850.00,
        350,
        4,
        'ZEEP0000000000000000000000000009'
    ),
    (
        'L',
        4,
        1100.00,
        450,
        4,
        'ZEEP0000000000000000000000000010'
    ),
    (
        'S',
        4,
        900.00,
        250,
        5,
        'ZEEP0000000000000000000000000011'
    ),
    (
        'M',
        4,
        1100.00,
        350,
        5,
        'ZEEP0000000000000000000000000012'
    ),
    (
        'L',
        4,
        1300.00,
        450,
        5,
        'ZEEP0000000000000000000000000013'
    ),
    (
        'S',
        4,
        1500.00,
        300,
        6,
        'ZEEP0000000000000000000000000014'
    ),
    (
        'L',
        4,
        1750.00,
        500,
        6,
        'ZEEP0000000000000000000000000015'
    ),
    (
        'S',
        4,
        500.00,
        200,
        7,
        'ZEEP0000000000000000000000000016'
    ),
    (
        'M',
        4,
        750.00,
        300,
        7,
        'ZEEP0000000000000000000000000017'
    ),
    (
        'L',
        4,
        900.00,
        400,
        7,
        'ZEEP0000000000000000000000000018'
    ),
    (
        'S',
        4,
        600.00,
        300,
        8,
        'ZEEP0000000000000000000000000019'
    ),
    (
        'L',
        4,
        750.00,
        500,
        8,
        'ZEEP0000000000000000000000000020'
    ),
    (
        'M',
        4,
        800.00,
        350,
        9,
        'ZEEP0000000000000000000000000021'
    ),
    (
        'L',
        4,
        1000.00,
        550,
        9,
        'ZEEP0000000000000000000000000022'
    ),
    (
        'S',
        4,
        500.00,
        200,
        10,
        'ZEEP0000000000000000000000000023'
    ),
    (
        'M',
        4,
        700.00,
        350,
        10,
        'ZEEP0000000000000000000000000024'
    ),
    (
        'S',
        4,
        600.00,
        500,
        11,
        'ZEEP0000000000000000000000000025'
    ),
    (
        'S',
        4,
        900.00,
        500,
        12,
        'ZEEP0000000000000000000000000026'
    ),
    (
        'S',
        4,
        400.00,
        500,
        13,
        'ZEEP0000000000000000000000000027'
    ),
    (
        'S',
        4,
        900.00,
        300,
        14,
        'ZEEP0000000000000000000000000028'
    ),
    (
        'S',
        4,
        1200.00,
        300,
        15,
        'ZEEP0000000000000000000000000029'
    ),
    (
        'S',
        4,
        800.00,
        300,
        16,
        'ZEEP0000000000000000000000000030'
    ),
    (
        'S',
        4,
        800.00,
        300,
        17,
        'ZEEP0000000000000000000000000031'
    ),
    (
        'S',
        4,
        1400.00,
        300,
        18,
        'ZEEP0000000000000000000000000032'
    ),
    (
        'S',
        4,
        1100.00,
        300,
        19,
        'ZEEP0000000000000000000000000033'
    ),
    (
        'S',
        4,
        1250.00,
        500,
        20,
        'ZEEP0000000000000000000000000034'
    ),
    (
        'S',
        4,
        1250.00,
        500,
        21,
        'ZEEP0000000000000000000000000035'
    ),
    (
        'S',
        4,
        1500.00,
        500,
        22,
        'ZEEP0000000000000000000000000036'
    );

-- Insert into Additives
INSERT INTO
    additives (
        name,
        description,
        base_price,
        size,
        unit_id,
        additive_category_id,
        image_key,
        machine_id
    )
VALUES
    (
        'Ванильный сироп',
        'Сладкий аромат ванили',
        50.00,
        20,
        4, -- Миллилитр (исправлено)
        1,
        'https://monin.us/cdn/shop/files/750mL-Vanilla.png?v=1724939521&width=1946',
        'ZEEP0000111122223333000001'
    ),
    (
        'Карамельный сироп',
        'Сладкий вкус карамели',
        70.00,
        40,
        4, -- Миллилитр (исправлено)
        1,
        'https://www.giffard.com/327-large_default/karamell-sirup.jpg',
        'ZEEP0000111122223333000002'
    ),
    (
        'Взбитые сливки',
        'Сливки для украшения напитков',
        80.00,
        500,
        2, -- Грамм (исправлено)
        3,
        'https://static.vecteezy.com/system/resources/previews/033/888/680/non_2x/white-whipped-cream-with-ai-generated-free-png.png',
        'ZEEP0000111122223333000003'
    ),
    (
        'Корица',
        'Пряная корица для добавления аромата',
        30.00,
        20,
        2, -- Грамм (исправлено)
        2,
        'https://png.pngtree.com/png-vector/20240803/ourmid/pngtree-cinnamon-isolated-on-white-background-png-image_13361859.png',
        'ZEEP0000111122223333000004'
    ),
    (
        'Шоколадная крошка',
        'Крошка из шоколада для украшения',
        90.00,
        20,
        2, -- Грамм (исправлено)
        3,
        'https://mlfwajaoc5gr.i.optimole.com/w:1080/h:1080/q:mauto/ig:avif/http://icecreambakery.in/wp-content/uploads/2024/03/Chocolate-Chips-Manufacturer-RPG-Industries-2.png',
        'ZEEP0000111122223333000005'
    ),
    (
        'Мед',
        'Натуральный мед для подслащивания',
        60.00,
        20,
        2, -- Грамм (исправлено)
        2,
        'https://png.pngtree.com/png-vector/20240801/ourmid/pngtree-3d-bowl-of-honey-with-a-driper-on-transparent-background-png-image_13330438.png',
        'ZEEP0000111122223333000006'
    ),
    (
        'Кокосовое молоко',
        'Молоко с ароматом кокоса',
        100.00,
        20,
        3, -- Литр (исправлено)
        4,
        'https://static.vecteezy.com/system/resources/previews/049/642/396/non_2x/coconut-milk-in-a-coconut-shell-with-a-transparent-background-png.png',
        'ZEEP0000111122223333000007'
    ),
    (
        'Клубничный сироп',
        'Сироп с ароматом клубники',
        65.00,
        20,
        4, -- Миллилитр (исправлено)
        1,
        'https://monin.ca/cdn/shop/files/750mL-Strawberry.png?v=1727882931&width=1445',
        'ZEEP0000111122223333000008'
    ),
    (
        'Мята',
        'Свежие листья мяты для аромата',
        40.00,
        20,
        2, -- Грамм (исправлено)
        5,
        'https://static.vecteezy.com/system/resources/thumbnails/047/732/281/small/basil-leaves-in-a-glass-bowl-transparent-background-png.png',
        'ZEEP0000111122223333000009'
    ),
    (
        'Сахар',
        'Белый сахар для подслащивания',
        10.00,
        20,
        1, -- Килограмм (корректно)
        6,
        'https://static.vecteezy.com/system/resources/thumbnails/044/902/343/small_2x/sweet-delight-enhance-your-recipes-with-sugar-in-bowl-free-png.png',
        'ZEEP0000111122223333000010'
    ),
    (
        'Кубики льда',
        'Лед для охлаждения напитков',
        5.00,
        20,
        4, -- Миллилитр (исправлено, если лед меряется в объеме)
        7,
        'https://img.pikbest.com/origin/09/14/36/70CpIkbEsTeTk.png!sw800',
        'ZEEP0000111122223333000011'
    ),
    (
        'Сироп из клена',
        'Ароматный сироп из клена',
        80.00,
        20,
        4, -- Миллилитр (исправлено)
        1,
        'https://png.pngtree.com/png-vector/20240812/ourmid/pngtree-sweet-syrup-drizzling-over-breakfast-treats-png-image_13462282.png',
        'ZEEP0000111122223333000012'
    ),
    (
        'Лимонный сок',
        'Свежий лимонный сок для аромата',
        45.00,
        20,
        4, -- Миллилитр (исправлено)
        2,
        'https://www.bestlobster.com/test/wp-content/uploads/2021/02/Lemon-PNG-Image.png',
        'ZEEP0000111122223333000013'
    ),
    (
        'Какао-порошок',
        'Порошок какао для украшения напитков',
        35.00,
        15,
        2, -- Грамм (корректно)
        3,
        'https://static.vecteezy.com/system/resources/previews/041/042/862/non_2x/ai-generated-heaping-spoonful-of-cocoa-powder-free-png.png',
        'ZEEP0000111122223333000014'
    );

INSERT INTO
    franchisees (name, description)
VALUES
    (
        'Кофейня "Астана"',
        'Франчайзинговая сеть кофеен в городе Астана'
    ),
    (
        'Сеть кофеен "Юг"',
        'Сеть кофеен на юге Казахстана'
    ),
    (
        'Кофейня "Шымкент"',
        'Кофейня в центре города Шымкент'
    );

-- Insert into Store
INSERT INTO
    stores (
        name,
        facility_address_id,
        franchisee_id,
        warehouse_id, -- Directly linking store to warehouse
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
        'Центральное Кафе',
        3,
        NULL,
        1,
        true,
        '+79001112233',
        'central@example.com',
        '8:00-20:00',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Assigned to 'Алматинский склад'
    (
        'Кафе на Углу',
        4,
        1,
        2,
        true,
        '+79002223344',
        'corner@example.com',
        '9:00-22:00',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Assigned to 'Астанинский склад'
    (
        'Маленькое Кафе',
        5,
        2,
        1,
        true,
        '+79003334455',
        'smallstore@example.com',
        '8:00-18:00',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Assigned to 'Алматинский склад'
    (
        'Городское Кафе',
        6,
        NULL,
        2,
        true,
        '+79004445566',
        'citycoffee@example.com',
        '7:00-23:00',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'New Cafe',
        7,
        NULL,
        2,
        true,
        '+79004445577',
        'newcafe@example.com',
        '7:00-23:00',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- Assigned to 'Астанинский склад'
-- Insert into StoreAdditives store 1 additives are loaded later in the script
INSERT INTO
    store_additives (store_id, additive_id)
VALUES
    (2, 1),
    (2, 3),
    (2, 4),
    (2, 5),
    (2, 6),
    (2, 8),
    (2, 9),
    (2, 10),
    (2, 11),
    (2, 12),
    (2, 14),
    (3, 1),
    (3, 13),
    (3, 14),
    (4, 2),
    (4, 4),
    (4, 5),
    (4, 6),
    (4, 7),
    (4, 9),
    (4, 11),
    (4, 12);

INSERT INTO
    store_products (
        store_id,
        product_id,
        is_available,
        created_at,
        updated_at
    )
SELECT
    1,
    p.id,
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM
    products p
WHERE
    p.deleted_at IS NULL;

-- Store 2
INSERT INTO
    store_products (store_id, product_id, is_available)
VALUES
    (2, 3, true), -- Espresso
    (2, 5, true), -- Black Tea
    (2, 6, false), -- Strawberry Smoothie
    (2, 8, true), -- Orange Juice
    (2, 14, false), -- Mineral Water
    (2, 18, true), -- Raspberry Lemonade
    (2, 17, false), -- Mint Tea
    (2, 16, true);

-- Ginger Tea
-- Store 3
INSERT INTO
    store_products (store_id, product_id, is_available)
VALUES
    (3, 12, true), -- Energy Drink
    (3, 19, true), -- Strawberry Tea
    (3, 21, true), -- Croissant w/ Chocolate
    (3, 10, true), -- Mojito
    (3, 15, true), -- Lemon Frappe
    (3, 7, true), -- Mango Smoothie
    (3, 20, true), -- Caramel Latte
    (3, 13, true);

-- Chocolate Milkshake
-- Store 4
INSERT INTO
    store_products (store_id, product_id, is_available)
VALUES
    (4, 20, true), -- (Assume re-using product 20 for "Maple Syrup Latte"?)
    (4, 16, true), -- Ginger Tea
    (4, 4, true), -- Green Tea
    (4, 6, false), -- Strawberry Smoothie
    (4, 10, true), -- Mojito
    (4, 18, true), -- Raspberry Lemonade
    (4, 17, true), -- Mint Tea
    (4, 12, true);

INSERT INTO
    store_product_sizes (
        store_product_id,
        product_size_id,
        store_price,
        created_at,
        updated_at
    )
SELECT
    sp.id,
    ps.id,
    ps.base_price,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM
    store_products sp
    JOIN product_sizes ps ON ps.product_id = sp.product_id
WHERE
    sp.store_id IS NOT NULL;

-- L size for Product 11;
-- Insert into ProductAdditive
INSERT INTO
    product_size_additives (product_size_id, additive_id, is_default)
VALUES
    (1, 1, true),
    (1, 2, true),
    (2, 3, true),
    (3, 1, true),
    (4, 2, true),
    (5, 4, false),
    (6, 5, false),
    (7, 3, false),
    (8, 6, false),
    (9, 7, true),
    (10, 2, true),
    (11, 8, false),
    (12, 9, true),
    (13, 10, true),
    (14, 11, false),
    (15, 12, false),
    (16, 13, true),
    (17, 14, false),
    (18, 1, false),
    (19, 3, true),
    (20, 6, false),
    (3, 8, false),
    (4, 10, false),
    (5, 12, true),
    (6, 14, false),
    (7, 1, false),
    (8, 4, false),
    (9, 5, false),
    (10, 9, false);

INSERT INTO
    store_additives (
        store_id,
        additive_id,
        store_price,
        created_at,
        updated_at
    )
SELECT DISTINCT
    1,
    psa.additive_id,
    a.base_price,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM
    store_product_sizes sps
    JOIN product_size_additives psa ON psa.product_size_id = sps.product_size_id
    JOIN additives a ON a.id = psa.additive_id
    JOIN store_products sp ON sp.id = sps.store_product_id
WHERE
    sp.store_id = 1;

-- Insert into Ingredients
INSERT INTO
    ingredients (
        name,
        calories,
        fat,
        carbs,
        proteins,
        expiration_in_days,
        unit_id,
        category_id,
        is_allergen
    )
VALUES
    ('Сахар', 387, 0, 100, 0, 365, 1, 3, false), -- Килограмм
    ('Молоко', 42, 1, 5, 3, 7, 3, 1, false), -- Литр
    ('Шоколад', 546, 30, 61, 7, 180, 2, 6, true), -- Грамм
    ('Корица', 247, 1.2, 81, 4, 365, 2, 4, true), -- Грамм
    ('Мед', 304, 0, 82, 0, 365, 2, 3, true), -- Грамм
    ('Ваниль', 288, 12, 55, 0, 730, 2, 4, true), -- Грамм
    ('Орехи', 607, 54, 18, 20, 365, 2, 5, true), -- Грамм
    ('Кокосовое молоко', 230, 23, 6, 2, 120, 3, 1, false), -- Литр
    ('Яблоки', 52, 0.2, 14, 0.3, 14, 1, 2, false), -- Килограмм
    ('Бананы', 96, 0.3, 27, 1.3, 7, 1, 2, false), -- Килограмм
    ('Сливки', 195, 20, 3, 2, 10, 3, 1, false), -- Литр (исправлено)
    ('Апельсины', 47, 0.1, 12, 0.9, 14, 1, 2, true), -- Килограмм
    ('Мята', 44, 0.7, 8, 3.3, 180, 2, 4, true), -- Грамм
    ('Лимонный сок', 123, 0.2, 6, 0.3, 60, 3, 2, true), -- Литр
    ('Какао-порошок', 228, 13, 58, 19, 730, 2, 6, true), -- Грамм
    ('Кленовый сироп', 261, 0, 67, 0, 365, 4, 3, true), -- Миллилитр
    ('Клубника', 33, 0.3, 8, 0.7, 10, 1, 2, true), -- Килограмм
    ('Имбирь', 80, 0.8, 18, 1.8, 180, 2, 4, true), -- Грамм (исправлено)
    ('Соль', 0, 0, 0, 0, 1095, 2, 4, false), -- Грамм (исправлено)
    ('Фисташки', 562, 45, 28, 20, 365, 2, 5, false);

-- Грамм (исправлено)
-- Орехи и семена
-- Insert into AdditiveIngredients
INSERT INTO
    additive_ingredients (additive_id, ingredient_id, quantity)
VALUES
    (1, 1, 1.00),
    (2, 1, 1.00),
    (3, 2, 1.00),
    (1, 3, 1.00),
    (4, 1, 1.00),
    (5, 4, 1.00),
    (6, 5, 1.00),
    (3, 4, 1.00),
    (6, 6, 1.00),
    (5, 1, 1.00),
    (2, 2, 1.00);

INSERT INTO
    additive_ingredients (
        additive_id,
        ingredient_id,
        quantity,
        created_at,
        updated_at
    )
SELECT
    a.id AS additive_id,
    1 AS ingredient_id, -- e.g. "Sugar" or some real ingredient
    1.0 AS quantity,
    NOW (),
    NOW ()
FROM
    additives a
    LEFT JOIN additive_ingredients ai ON a.id = ai.additive_id
WHERE
    ai.id IS NULL;


-- =============================================
-- 1) INSERT THE 21 NEW TEST PRODUCTS
--    (All in category_id=14 = "Тестовые продукты")
-- =============================================
INSERT INTO products (
    name,
    description,
    image_key,
    video_key,
    category_id
)
VALUES
    ('Zeep Grape',
     'Тестовый продукт: виноградный вкус',
     'https://example.com/img/zeep_grape.png',
     NULL,
     14
    ),
    ('Zeep Light Grape',
     'Тестовый продукт: легкий виноградный вкус',
     'https://example.com/img/zeep_light_grape.png',
     NULL,
     14
    ),
    ('Zeep Shine Muskat',
     'Тестовый продукт: мускатный аромат',
     'https://example.com/img/zeep_shine_muskat.png',
     NULL,
     14
    ),
    ('Zeep Mango',
     'Тестовый продукт: манговый вкус',
     'https://example.com/img/zeep_mango.png',
     NULL,
     14
    ),
    ('Pulpy Grape Tea',
     'Тестовый продукт: чай с виноградной мякотью',
     'https://example.com/img/pulpy_grape_tea.png',
     NULL,
     14
    ),
    ('Pulpy Coco-Mango Tea',
     'Тестовый продукт: чай с кокосом и манго',
     'https://example.com/img/pulpy_coco_mango.png',
     NULL,
     14
    ),
    ('Pulpy Mango Saga Tea',
     'Тестовый продукт: манговый чай с мякотью',
     'https://example.com/img/pulpy_mango_saga.png',
     NULL,
     14
    ),
    ('Pulpy Super Mango Tea',
     'Тестовый продукт: интенсивный манговый вкус',
     'https://example.com/img/pulpy_super_mango.png',
     NULL,
     14
    ),
    ('Pulpy Shine Muscat Tea',
     'Тестовый продукт: Shine Muscat чай с мякотью',
     'https://example.com/img/pulpy_shine_muscat.png',
     NULL,
     14
    ),
    ('Passion Green Tea',
     'Тестовый продукт: зеленый чай с маракуйей',
     'https://example.com/img/passion_green_tea.png',
     NULL,
     14
    ),
    ('Mango Passion Tea',
     'Тестовый продукт: чай с манго и маракуйей',
     'https://example.com/img/mango_passion_tea.png',
     NULL,
     14
    ),
    ('Ruby Grapefruit Tea',
     'Тестовый продукт: чай с красным грейпфрутом',
     'https://example.com/img/ruby_grapefruit_tea.png',
     NULL,
     14
    ),
    ('Duck Citrus Lemon Tea',
     'Тестовый продукт: чай с лимоном и цитрусами',
     'https://example.com/img/duck_citrus_lemon.png',
     NULL,
     14
    ),
    ('Lemon Tea',
     'Тестовый продукт: классический чай с лимоном',
     'https://example.com/img/lemon_tea.png',
     NULL,
     14
    ),
    ('Cheesy Matcha',
     'Тестовый продукт: матча с сырной шапкой',
     'https://example.com/img/cheesy_matcha.png',
     NULL,
     14
    ),
    ('Matcha Milk Pudding',
     'Тестовый продукт: матча с молочным пудингом',
     'https://example.com/img/matcha_milk_pudding.png',
     NULL,
     14
    ),
    ('Milk Ceylon Tea',
     'Тестовый продукт: Цейлонский чай с молоком',
     'https://example.com/img/milk_ceylon.png',
     NULL,
     14
    ),
    ('Milk Tapioca Tea',
     'Тестовый продукт: молочный чай с тапиокой',
     'https://example.com/img/milk_tapioca.png',
     NULL,
     14
    ),
    ('Tapioca&Milk Pudding',
     'Тестовый продукт: чай с тапиокой и пудингом',
     'https://example.com/img/tapioca_milk_pudding.png',
     NULL,
     14
    ),
    ('Black Sugar Milk',
     'Тестовый продукт: молоко с чёрным сахаром',
     'https://example.com/img/black_sugar_milk.png',
     NULL,
     14
    ),
    ('OREO CREME BRULEE',
     'Тестовый продукт: крем-брюле со вкусом Oreo',
     'https://example.com/img/oreo_creme_brulee.png',
     NULL,
     14
    );


-- =============================================
-- 2) INSERT PRODUCT_SIZES FOR THOSE 21 NEW ITEMS
--    All "S", unit_id = 4 (ml), pick your base_price
--    NOTE: product_id values below assume these 21
--          new products are assigned IDs 23..43.
--          Adjust if needed based on your DB.
-- =============================================

INSERT INTO product_sizes (
    name,
    unit_id,
    base_price,
    size,
    product_id,
    machine_id
)
VALUES
    -- Zeep Grape (product_id = 23?)
    ('S', 4, 1300.00, 300, 23, 'ZG0001'),

    -- Zeep Light Grape (product_id = 24?)
    ('S', 4, 1400.00, 300, 24, 'ZLG002'),

    -- Zeep Shine Muskat (25?)
    ('S', 4, 1500.00, 300, 25, 'ZSM003'),

    -- Zeep Mango (26?)
    ('S', 4, 1400.00, 300, 26, 'ZM0004'),

    -- Pulpy Grape Tea (27?)
    ('S', 4, 1300.00, 300, 27, 'PGT005'),

    -- Pulpy Coco-Mango Tea (28?)
    ('S', 4, 1350.00, 300, 28, 'PCMT06'),

    -- Pulpy Mango Saga Tea (29?)
    ('S', 4, 1400.00, 300, 29, 'PMST07'),

    -- Pulpy Super Mango Tea (30?)
    ('S', 4, 1450.00, 300, 30, 'PSMT08'),

    -- Pulpy Shine Muscat Tea (31?)
    ('S', 4, 1450.00, 300, 31, 'PSMT09'),

    -- Passion Green Tea (32?)
    ('S', 4, 1300.00, 300, 32, 'PGT010'),

    -- Mango Passion Tea (33?)
    ('S', 4, 1400.00, 300, 33, 'MPT011'),

    -- Ruby Grapefruit Tea (34?)
    ('S', 4, 1500.00, 300, 34, 'RGT012'),

    -- Duck Citrus Lemon Tea (35?)
    ('S', 4, 1300.00, 300, 35, 'DCLT13'),

    -- Lemon Tea (36?)
    ('S', 4, 1200.00, 300, 36, 'LT0014'),

    -- Cheesy Matcha (37?)
    ('S', 4, 1800.00, 300, 37, 'CM0015'),

    -- Matcha Milk Pudding (38?)
    ('S', 4, 1800.00, 300, 38, 'MMP016'),

    -- Milk Ceylon Tea (39?)
    ('S', 4, 1300.00, 300, 39, 'MCT017'),

    -- Milk Tapioca Tea (40?)
    ('S', 4, 1600.00, 300, 40, 'MTT018'),

    -- Tapioca&Milk Pudding (41?)
    ('S', 4, 1700.00, 300, 41, 'TMP019'),

    -- Black Sugar Milk (42?)
    ('S', 4, 1800.00, 300, 42, 'BSM020'),

    -- OREO CREME BRULEE (43?)
    ('S', 4, 2000.00, 300, 43, 'OCM021');


-- Орехи и семена
-- Insert into ProductIngredients
INSERT INTO
    product_size_ingredients (ingredient_id, product_size_id, quantity)
VALUES
    -- Product Size 1 (S, Product 1)
    (1, 1, 1), -- Сахар
    (2, 1, 1), -- Молоко
    (4, 1, 1), -- Корица
    -- Product Size 2 (M, Product 1)
    (1, 2, 1), -- Сахар
    (3, 2, 1), -- Шоколад
    (5, 2, 1), -- Мед
    -- Product Size 3 (L, Product 1)
    (3, 3, 1), -- Шоколад
    (6, 3, 1), -- Ваниль
    (7, 3, 1), -- Орехи
    -- Product Size 4 (S, Product 2)
    (4, 4, 1), -- Корица
    (2, 4, 1), -- Молоко
    (8, 4, 1), -- Кокосовое молоко
    -- Product Size 5 (M, Product 2)
    (5, 5, 1), -- Мед
    (1, 5, 1), -- Сахар
    (9, 5, 1), -- Яблоки
    -- Product Size 6 (L, Product 2)
    (6, 6, 1), -- Ваниль
    (3, 6, 1), -- Шоколад
    (10, 6, 1), -- Бананы
    -- Product Size 7 (S, Product 3)
    (7, 7, 1), -- Орехи
    (8, 7, 1), -- Кокосовое молоко
    (13, 7, 1), -- Мята
    -- Product Size 8 (M, Product 3)
    (9, 8, 1), -- Яблоки
    (14, 8, 1), -- Лимонный сок
    (4, 8, 1), -- Корица
    -- Product Size 9 (L, Product 3)
    (10, 9, 1), -- Бананы
    (15, 9, 1), -- Какао-порошок
    (11, 9, 1), -- Сливки
    -- Product Size 10 (S, Product 4)
    (11, 10, 1), -- Сливки
    (16, 10, 1), -- Кленовый сироп
    (2, 10, 1), -- Молоко
    -- Product Size 11 (M, Product 4)
    (12, 11, 1), -- Апельсины
    (17, 11, 1), -- Клубника
    (8, 11, 1), -- Кокосовое молоко
    -- Product Size 12 (L, Product 4)
    (13, 12, 1), -- Мята
    (18, 12, 1), -- Имбирь
    (5, 12, 1), -- Мед
    -- Product Size 13 (S, Product 5)
    (14, 13, 1), -- Лимонный сок
    (19, 13, 1), -- Соль
    (15, 13, 1), -- Какао-порошок
    -- Product Size 14 (M, Product 5)
    (20, 14, 1), -- Фисташки
    (3, 14, 1), -- Шоколад
    (6, 14, 1), -- Ваниль
    -- Product Size 15 (L, Product 5)
    (7, 15, 1), -- Орехи
    (1, 15, 1), -- Сахар
    (11, 15, 1), -- Сливки
    -- Product Size 16 (S, Product 6)
    (9, 16, 1), -- Яблоки
    (4, 16, 1), -- Корица
    (14, 16, 1), -- Лимонный сок
    -- Product Size 17 (M, Product 6)
    (10, 17, 1), -- Бананы
    (17, 17, 1), -- Клубника
    (13, 17, 1), -- Мята
    -- Product Size 18 (L, Product 6)
    (18, 18, 1), -- Имбирь
    (12, 18, 1), -- Апельсины
    (20, 18, 1);

INSERT INTO
    product_size_ingredients (
        product_size_id,
        ingredient_id,
        quantity,
        created_at,
        updated_at
    )
SELECT
    ps.id AS product_size_id,
    1 AS ingredient_id, -- e.g. "Sugar" or some default
    1.0 AS quantity,
    NOW (),
    NOW ()
FROM
    product_sizes ps
    LEFT JOIN product_size_ingredients psi ON ps.id = psi.product_size_id
WHERE
    psi.id IS NULL;

-- Фисташки
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
        'Елена',
        'Соколова',
        '+79551234567',
        'elena@example.com',
        true,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Сергей',
        'Павлов',
        '+79667778899',
        'sergey@example.com',
        true,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Анна',
        'Федорова',
        '+79223334455',
        'anna@example.com',
        true,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Иван',
        'Иванов',
        '+79161234567',
        'ivan@example.com',
        true,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Мария',
        'Смирнова',
        '+79345566778',
        'maria@example.com',
        true,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Олег',
        'Кузнецов',
        '+79991234567',
        'oleg@example.com',
        false,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Татьяна',
        'Орлова',
        '+79882233445',
        'tatiana@example.com',
        true,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Алексей',
        'Попов',
        '+79002221133',
        'alexei@example.com',
        true,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Юлия',
        'Петрова',
        '+79115555666',
        'yulia@example.com',
        true,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Дмитрий',
        'Фролов',
        '+79553334456',
        'dmitry@example.com',
        false,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Елена',
        'Орлова',
        '+79151236817',
        'elenaOrlova@example.com',
        true,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Диар',
        'Бегисбаев',
        '+79191236817',
        'begisbayev@example.com',
        true,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'New',
        'Manager',
        '+77771236507',
        'newmanager@example.com',
        true,
        '$2a$10$TpQjaWD1c2cj8Omkb6l36.tVrR8dl0EtuNwcrD09THT9dL7bo5aQy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- Franchisee Employees Table
INSERT INTO
    franchisee_employees (
        franchisee_id,
        employee_id,
        role,
        created_at,
        updated_at
    )
VALUES
    (
        1,
        8,
        'FRANCHISEE_OWNER',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        1,
        9,
        'FRANCHISEE_MANAGER',
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
        2,
        'STORE_MANAGER',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        3,
        3,
        'STORE_MANAGER',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        4,
        1,
        'BARISTA',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        10,
        2,
        'BARISTA',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        12,
        4,
        'STORE_MANAGER',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        13,
        5,
        'STORE_MANAGER',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- Warehouse Employees Table
INSERT INTO
    warehouse_employees (
        employee_id,
        warehouse_id,
        role,
        created_at,
        updated_at
    )
VALUES
    (
        5,
        1,
        'WAREHOUSE_EMPLOYEE',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        6,
        2,
        'WAREHOUSE_MANAGER',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- Region Managers Table
INSERT INTO
    region_employees (
        employee_id,
        region_id,
        role,
        created_at,
        updated_at
    )
VALUES
    (
        7,
        1,
        'REGION_WAREHOUSE_MANAGER',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- Admin Employees Table
INSERT INTO
    admin_employees (employee_id, role, created_at, updated_at)
VALUES
    (1, 'ADMIN', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (11, 'OWNER', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert into EmployeeAudit
INSERT INTO
    employee_work_tracks (start_work_at, end_work_at, employee_id)
VALUES
    (
        '2024-10-01 09:00:00+00',
        '2024-10-01 17:00:00+00',
        1
    ),
    (
        '2024-10-02 09:00:00+00',
        '2024-10-02 17:00:00+00',
        1
    );

-- Insert into EmployeeWorkday
INSERT INTO
    employee_workdays (day, start_at, end_at, employee_id)
VALUES
    ('Понедельник', '08:00:00', '16:00:00', 1),
    ('Вторник', '08:00:00', '16:00:00', 2),
    ('Среда', '08:00:00', '16:00:00', 3);

-- Insert into Referral
INSERT INTO
    referrals (customer_id, referee_id)
VALUES
    (1, 2),
    (2, 3);

-- Insert into VerificationCode
INSERT INTO
    verification_codes (customer_id, code, expires_at)
VALUES
    (1, '123456', '2024-12-31 23:59:59+00'),
    (2, '654321', '2024-12-31 23:59:59+00');

-- Insert into CustomerAddress
INSERT INTO
    customer_addresses (customer_id, address, longitude, latitude)
VALUES
    (1, 'Улица Ленина, дом 34', 37.6173, 55.7558),
    (2, 'Проспект Мира, дом 45', 30.3158, 59.9343),
    (3, 'Улица Советская, дом 89', 60.6094, 56.8389);

-- Insert into Bonus
INSERT INTO
    bonuses (bonuses, customer_id, expires_at)
VALUES
    (100.00, 1, '2024-12-31 23:59:59+00'),
    (50.00, 2, '2024-06-30 23:59:59+00'),
    (25.00, 3, '2024-04-30 23:59:59+00');

INSERT INTO
    store_stocks (
        store_id,
        ingredient_id,
        quantity,
        low_stock_threshold,
        created_at,
        updated_at
    )
VALUES
    -- Store 1 Stocks (Центральное Кафе)
    (
        1,
        1,
        10000,
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Молоко (литр)
    (
        1,
        2,
        10000,
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Сахар (кг)
    (
        1,
        3,
        100000,
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Шоколад (грамм)
    (
        1,
        5,
        100000,
        20,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Мед (грамм)
    (
        1,
        6,
        100000,
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Ваниль (грамм)
    -- Store 2 Stocks (Кафе на Углу)
    (
        2,
        1,
        10000,
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Молоко (литр)
    (
        2,
        2,
        10000,
        5,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Сахар (кг)
    (
        2,
        7,
        100000,
        15,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Мята (грамм)
    (
        2,
        8,
        10000,
        20,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Кокосовое молоко (литр)
    (
        2,
        10,
        10000,
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Лимонный сок (литр)
    -- Store 3 Stocks (Маленькое Кафе)
    (
        3,
        1,
        10000,
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Молоко (литр)
    (
        3,
        3,
        100000,
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Шоколад (грамм)
    (
        3,
        5,
        100000,
        20,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Мед (грамм)
    (
        3,
        12,
        100000,
        5,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Имбирь (грамм)
    (
        3,
        14,
        100000,
        15,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Какао-порошок (грамм)
    -- Store 4 Stocks (Городское Кафе)
    (
        4,
        1,
        10000,
        20,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Молоко (литр)
    (
        4,
        2,
        10000,
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Сахар (кг)
    (
        4,
        7,
        100000,
        15,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Мята (грамм)
    (
        4,
        9,
        100000,
        5,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Кубики льда (мл)
    (
        4,
        14,
        100000,
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- Какао-порошок (грамм)
WITH
    store_list AS (
        SELECT
            id AS store_id
        FROM
            stores
    )
INSERT INTO
    store_stocks (
        store_id,
        ingredient_id,
        quantity,
        low_stock_threshold,
        created_at,
        updated_at
    )
SELECT
    sl.store_id,
    i.id AS ingredient_id,
    10000 AS quantity, -- Example initial stock amount
    100 AS low_stock_threshold,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM
    (
        SELECT
            sa.store_id,
            ai.ingredient_id
        FROM
            store_additives sa
            JOIN additive_ingredients ai ON ai.additive_id = sa.additive_id
        UNION
        SELECT
            sp.store_id,
            psi.ingredient_id
        FROM
            store_products sp
            JOIN store_product_sizes sps ON sps.store_product_id = sp.id
            JOIN product_size_ingredients psi ON psi.product_size_id = sps.product_size_id
    ) AS needed
    JOIN store_list sl ON sl.store_id = needed.store_id
    JOIN ingredients i ON i.id = needed.ingredient_id
WHERE
    NOT EXISTS (
        SELECT
            1
        FROM
            store_stocks s
        WHERE
            s.store_id = sl.store_id
            AND s.ingredient_id = i.id
            AND s.deleted_at IS NULL
    );

-- Insert stock material categories
INSERT INTO
    stock_material_categories (name, description)
VALUES
    (
        'Молочные продукты',
        'Молоко, сливки, йогурты и другие молочные продукты'
    ),
    (
        'Подсластители',
        'Сахар, мед и другие подсластители'
    ),
    (
        'Кондитерские изделия',
        'Шоколад, какао и другие кондитерские ингредиенты'
    ),
    ('Специи', 'Различные специи и пряности'),
    (
        'Ароматизаторы',
        'Ванильный экстракт и другие ароматизаторы'
    );

-- Insert stock materials with Russian names and category references
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
        'Простоквашино Молоко 3.2%',
        'Молоко пастеризованное 1л',
        2,
        50,
        3,
        1,
        1,
        '111111111111',
        1095,
        TRUE
    ),
    (
        'Русский сахар Экстра',
        'Сахар песок высший сорт 1кг',
        1,
        20,
        1,
        1,
        2,
        '222222222222',
        1095,
        TRUE
    ),
    (
        'Бабаевский горький шоколад 75%',
        'Темный шоколад 500г',
        3,
        15,
        2,
        500,
        3,
        '333333333333',
        730,
        TRUE
    ),
    (
        'Индийская корица молотая премиум',
        'Корица молотая 200г',
        4,
        10,
        2,
        200,
        4,
        '444444444444',
        1460,
        TRUE
    ),
    (
        'Dr.Oetker Ванильный экстракт',
        'Натуральный экстракт ванили 50мл',
        5,
        25,
        4,
        50,
        5,
        '555555555555',
        1460,
        TRUE
    );

-- Insert into Suppliers
INSERT INTO
    suppliers (name, contact_email, contact_phone, city, address)
VALUES
    (
        'ООО "Нестле Россия"',
        'contact@nestle.ru',
        '+79005556677',
        'Москва',
        'Павелецкая площадь, д. 2, стр. 1'
    ),
    (
        'АО "Кока-Кола ЭйчБиСи Евразия"',
        'info@coca-cola.ru',
        '+78002002222',
        'Москва',
        'ул. Новоорловская, д. 7'
    ),
    (
        'ООО "ПепсиКо Холдингс"',
        'support@pepsico.ru',
        '+78007001000',
        'Московская обл.',
        'г. Солнечногорск, территория свободной экономической зоны "Шерризон", стр. 1'
    ),
    (
        'ООО "Юнилевер Русь"',
        'info@unilever.ru',
        '+78002001200',
        'Москва',
        'ул. Сергея Макеева, д. 13'
    ),
    (
        'ООО "Штарбакс"',
        'help@starbucks.ru',
        '+78001008333',
        'Москва',
        'ул. Большая Новодмитровская, д. 23, стр. 1'
    ),
    (
        'ООО "Мон`дэлис Русь"',
        'support@mdlz.ru',
        '+74959602424',
        'Владимирская обл.',
        'г. Покров, ул. Франца Штольверка, д. 10'
    ),
    (
        'АО "ДАНОН РОССИЯ"',
        'contact@danone.ru',
        '+78002000201',
        'Москва',
        'ул. Вятская, д. 27, корп. 13-14'
    ),
    (
        'ООО "Марс"',
        'support@mars.ru',
        '+74957212100',
        'Московская обл.',
        'г. Ступино, ул. Ситенка, д. 12'
    ),
    (
        'ООО "Юнилевер Русь"',
        'contact@unilever.ru',
        '+78002001201',
        'Омск',
        'ул. 10 лет Октября, д. 205'
    ),
    (
        'ООО "Дженерал Миллс Рус"',
        'support@generalmills.ru',
        '+74959373400',
        'Москва',
        'ул. Большая Новодмитровская, д. 14, стр. 2'
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
    (
        1,
        1,
        '2024-09-01',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Delivery 1
    (
        2,
        1,
        '2024-11-01',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ), -- Delivery 2
    (
        1,
        2,
        '2024-10-01',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- Delivery 3
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
    (1, 1, '111111111111', 50, '2026-12-01'), -- Milk Delivery
    (1, 2, '222222222222', 30, '2025-06-05'), -- Sugar Delivery
    (2, 3, '333333333333', 40, '2025-11-20'), -- Chocolate Delivery
    (3, 4, '444444444444', 20, '2026-06-10'), -- Cinnamon Delivery
    (3, 5, '555555555555', 15, '2027-12-15');

-- Vanilla Delivery
-- Insert mock data into supplier_materials
INSERT INTO
    supplier_materials (
        supplier_id,
        stock_material_id,
        created_at,
        updated_at
    )
VALUES
    -- Nestlé supplies Milk and Sugar
    (1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Milk
    (1, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Sugar
    -- Coca-Cola supplies Sugar
    (2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Sugar
    -- PepsiCo supplies Chocolate and Cinnamon
    (3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Chocolate
    (3, 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Cinnamon
    -- Unilever supplies Vanilla
    (4, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Vanilla
    (4, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Chocolate
-- Insert mock data into supplier_prices
INSERT INTO
    supplier_prices (supplier_material_id, base_price)
VALUES
    -- Prices for supplier_material_id 1
    (1, 50.00),
    (1, 48.00),
    -- Prices for supplier_material_id 2
    (2, 25.00),
    (2, 24.50),
    -- Prices for supplier_material_id 3
    (3, 100.00),
    (3, 98.00),
    -- Prices for supplier_material_id 4
    (4, 30.00),
    (4, 28.00),
    -- Prices for supplier_material_id 5
    (5, 75.00),
    (5, 72.50);

INSERT INTO
    warehouse_stocks (warehouse_id, stock_material_id, quantity)
VALUES
    (1, 1, 50), -- Milk in Warehouse 1
    (1, 2, 30), -- Sugar in Warehouse 1
    (1, 3, 40), -- Chocolate in Warehouse 1
    (2, 4, 20), -- Cinnamon in Warehouse 2
    (2, 5, 15);

-- Vanilla in Warehouse 2
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
        2,
        2,
        'PROCESSED',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        3,
        2,
        'IN_DELIVERY',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        4,
        1,
        'COMPLETED',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

INSERT INTO
    stock_request_ingredients (
        stock_request_id,
        stock_material_id,
        quantity,
        created_at,
        updated_at
    )
VALUES
    -- StockRequest 1 (Store 1 -> Warehouse 1)
    (1, 2, 10.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Sugar
    (1, 1, 20.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Milk
    -- StockRequest 2 (Store 2 -> Warehouse 2)
    (2, 3, 5.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Chocolate
    (2, 4, 2.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Cinnamon
    -- StockRequest 3 (Store 3 -> Warehouse 3)
    (3, 5, 1.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Vanilla
    (3, 2, 15.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Sugar
    -- StockRequest 4 (Store 4 -> Warehouse 4)
    (4, 1, 10.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Milk
    (4, 3, 8.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Chocolate
