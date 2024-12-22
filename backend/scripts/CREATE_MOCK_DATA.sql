-- Insert into FacilityAddress
INSERT INTO
  facility_addresses (address, longitude, latitude)
VALUES
  ('Улица Ленина, 12, Москва', 37.6173, 55.7558),
  (
    'Проспект Мира, 45, Санкт-Петербург',
    30.3158,
    59.9343
  ),
  (
    'Улица Советская, 89, Екатеринбург',
    60.6094,
    56.8389
  ),
  (
    'Улица Куйбышева, 101, Новосибирск',
    82.9204,
    55.0415
  ),
  (
    'Площадь Революции, 17, Нижний Новгород',
    44.0059,
    56.3269
  ),
  ('Проспект Гагарина, 27, Казань', 49.1233, 55.8304),
  ('Улица Ленина, 64, Пермь', 56.2516, 58.0105),
  ('Проспект Победы, 5, Самара', 50.1834, 53.2038),
  (
    'Улица Большая Садовая, 101, Ростов-на-Дону',
    39.7015,
    47.2225
  ),
  (
    'Невский проспект, 88, Санкт-Петербург',
    30.3543,
    59.9311
  ),
  (
    'Улица Советская, 18, Волгоград',
    44.5018,
    48.7080
  ),
  (
    'Улица Октябрьская, 5, Челябинск',
    61.4026,
    55.1600
  ),
  ('Улица Кирова, 2, Уфа', 56.0367, 54.7352);

-- Insert into Units
INSERT INTO
  units (name, conversion_factor)
VALUES
  ('kg', 1.0),
  ('g', 0.001),
  ('L', 1.0),
  ('ml', 0.001);

-- Insert into CityWarehouses
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
        address = 'Улица Ленина, 12, Москва'
    ),
    'Московский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Проспект Мира, 45, Санкт-Петербург'
    ),
    'Санкт-Петербургский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Улица Советская, 89, Екатеринбург'
    ),
    'Екатеринбургский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Улица Куйбышева, 101, Новосибирск'
    ),
    'Новосибирский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Площадь Революции, 17, Нижний Новгород'
    ),
    'Нижегородский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Проспект Гагарина, 27, Казань'
    ),
    'Казанский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Улица Ленина, 64, Пермь'
    ),
    'Пермский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Проспект Победы, 5, Самара'
    ),
    'Самарский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Улица Большая Садовая, 101, Ростов-на-Дону'
    ),
    'Ростовский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Невский проспект, 88, Санкт-Петербург'
    ),
    'Второй Санкт-Петербургский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Улица Советская, 18, Волгоград'
    ),
    'Волгоградский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Улица Октябрьская, 5, Челябинск'
    ),
    'Челябинский склад'
  ),
  (
    (
      SELECT
        id
      FROM
        facility_addresses
      WHERE
        address = 'Улица Кирова, 2, Уфа'
    ),
    'Уфимский склад'
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
  ('Травяные чаи', 'Настои и чаи на травах');

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
    'https://gallerypng.com/wp-content/uploads/2024/05/mango-shake-png-image.png',
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
  );

-- Insert into RecipeStep
INSERT INTO
  recipe_steps (product_id, step, name, description, image_url)
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
    measure,
    base_price,
    size,
    is_default,
    product_id
  )
VALUES
  ('S', 'мл', 1000.00, 250, true, 1),
  ('M', 'мл', 1250.00, 350, false, 1),
  ('L', 'мл', 1500.00, 450, false, 1),
  ('S', 'мл', 1100.00, 300, true, 2),
  ('S', 'мл', 750.00, 200, true, 3),
  ('M', 'мл', 900.00, 300, false, 3),
  ('L', 'мл', 1150.00, 400, false, 3),
  ('S', 'мл', 600.00, 250, true, 4),
  ('M', 'мл', 850.00, 350, false, 4),
  ('L', 'мл', 1100.00, 450, false, 4),
  ('S', 'мл', 900.00, 250, true, 5),
  ('M', 'мл', 1100.00, 350, false, 5),
  ('L', 'мл', 1300.00, 450, false, 5),
  ('S', 'мл', 1500.00, 300, true, 6),
  ('L', 'мл', 1750.00, 500, false, 6),
  ('S', 'мл', 500.00, 200, true, 7),
  ('M', 'мл', 750.00, 300, false, 7),
  ('L', 'мл', 900.00, 400, false, 7),
  ('S', 'мл', 600.00, 300, true, 8),
  ('L', 'мл', 750.00, 500, false, 8),
  ('M', 'мл', 800.00, 350, true, 9),
  ('L', 'мл', 1000.00, 550, false, 9),
  ('S', 'мл', 500.00, 200, true, 10),
  ('M', 'мл', 700.00, 350, false, 10),
  ('S', 'мл', 600.00, 500, true, 11),
  ('S', 'мл', 900.00, 500, true, 18),
  ('S', 'мл', 400.00, 500, true, 13),
  ('S', 'мл', 900.00, 300, true, 20),
  ('S', 'мл', 1200.00, 300, true, 14),
  ('S', 'мл', 800.00, 300, true, 16),
  ('S', 'мл', 800.00, 300, true, 17),
  ('S', 'мл', 1400.00, 300, true, 15),
  ('S', 'мл', 1100.00, 300, true, 19),
  ('S', 'мл', 1250.00, 500, true, 12);

-- Insert into Additives
INSERT INTO
  additives (
    name,
    description,
    base_price,
    size,
    additive_category_id,
    image_url
  )
VALUES
  (
    'Ванильный сироп',
    'Сладкий аромат ванили',
    50.00,
    'мл',
    1,
    'https://monin.us/cdn/shop/files/750mL-Vanilla.png?v=1724939521&width=1946'
  ),
  (
    'Карамельный сироп',
    'Сладкий вкус карамели',
    70.00,
    'мл',
    1,
    'https://www.giffard.com/327-large_default/karamell-sirup.jpg'
  ),
  (
    'Взбитые сливки',
    'Сливки для украшения напитков',
    80.00,
    'г',
    3,
    'https://static.vecteezy.com/system/resources/previews/033/888/680/non_2x/white-whipped-cream-with-ai-generated-free-png.png'
  ),
  (
    'Корица',
    'Пряная корица для добавления аромата',
    30.00,
    'г',
    2,
    'https://png.pngtree.com/png-vector/20240803/ourmid/pngtree-cinnamon-isolated-on-white-background-png-image_13361859.png'
  ),
  (
    'Шоколадная крошка',
    'Крошка из шоколада для украшения',
    90.00,
    'г',
    3,
    'https://mlfwajaoc5gr.i.optimole.com/w:1080/h:1080/q:mauto/ig:avif/http://icecreambakery.in/wp-content/uploads/2024/03/Chocolate-Chips-Manufacturer-RPG-Industries-2.png'
  ),
  (
    'Мед',
    'Натуральный мед для подслащивания',
    60.00,
    'мл',
    2,
    'https://png.pngtree.com/png-vector/20240801/ourmid/pngtree-3d-bowl-of-honey-with-a-driper-on-transparent-background-png-image_13330438.png'
  ),
  (
    'Кокосовое молоко',
    'Молоко с ароматом кокоса',
    100.00,
    'мл',
    4,
    'https://static.vecteezy.com/system/resources/previews/049/642/396/non_2x/coconut-milk-in-a-coconut-shell-with-a-transparent-background-png.png'
  ),
  (
    'Клубничный сироп',
    'Сироп с ароматом клубники',
    65.00,
    'мл',
    1,
    'https://monin.ca/cdn/shop/files/750mL-Strawberry.png?v=1727882931&width=1445'
  ),
  (
    'Мята',
    'Свежие листья мяты для аромата',
    40.00,
    'г',
    5,
    'https://static.vecteezy.com/system/resources/thumbnails/047/732/281/small/basil-leaves-in-a-glass-bowl-transparent-background-png.png'
  ),
  (
    'Сахар',
    'Белый сахар для подслащивания',
    10.00,
    'г',
    6,
    'https://static.vecteezy.com/system/resources/thumbnails/044/902/343/small_2x/sweet-delight-enhance-your-recipes-with-sugar-in-bowl-free-png.png'
  ),
  (
    'Кубики льда',
    'Лед для охлаждения напитков',
    5.00,
    'г',
    7,
    'https://img.pikbest.com/origin/09/14/36/70CpIkbEsTeTk.png!sw800'
  ),
  (
    'Сироп из клена',
    'Ароматный сироп из клена',
    80.00,
    'мл',
    1,
    'https://png.pngtree.com/png-vector/20240812/ourmid/pngtree-sweet-syrup-drizzling-over-breakfast-treats-png-image_13462282.png'
  ),
  (
    'Лимонный сок',
    'Свежий лимонный сок для аромата',
    45.00,
    'мл',
    2,
    'https://www.bestlobster.com/test/wp-content/uploads/2021/02/Lemon-PNG-Image.png'
  ),
  (
    'Какао-порошок',
    'Порошок какао для украшения напитков',
    35.00,
    'г',
    3,
    'https://static.vecteezy.com/system/resources/previews/041/042/862/non_2x/ai-generated-heaping-spoonful-of-cocoa-powder-free-png.png'
  );

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
    'Центральное кафе',
    1,
    false,
    'ACTIVE',
    '79001112233',
    'central@example.com',
    '8:00-20:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Кофейня на углу',
    2,
    true,
    'ACTIVE',
    '79002223344',
    'corner@example.com',
    '9:00-22:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Маленький магазин на Советской',
    3,
    true,
    'ACTIVE',
    '79003334455',
    'smallstore@example.com',
    '8:00-18:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Кофейня у вокзала',
    4,
    false,
    'DISABLED',
    '79004445566',
    'station@example.com',
    '10:00-22:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Городской кофе',
    5,
    true,
    'ACTIVE',
    '79005556677',
    'citycoffee@example.com',
    '7:00-23:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Летняя терраса',
    6,
    true,
    'ACTIVE',
    '79006667788',
    'terrace@example.com',
    '10:00-22:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Кафе на проспекте',
    7,
    false,
    'ACTIVE',
    '79007778899',
    'avenuecafe@example.com',
    '9:00-21:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Заведение у реки',
    8,
    true,
    'ACTIVE',
    '79008889900',
    'riverside@example.com',
    '10:00-22:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Чайный дом',
    9,
    false,
    'DISABLED',
    '79009990011',
    'teahouse@example.com',
    '8:00-20:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Кофе и компании',
    10,
    true,
    'ACTIVE',
    '79010001122',
    'coffeeandco@example.com',
    '8:00-22:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Парк-кафе',
    11,
    false,
    'ACTIVE',
    '79011002233',
    'parkcafe@example.com',
    '10:00-21:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Восточный уголок',
    12,
    true,
    'ACTIVE',
    '79012003344',
    'easterncorner@example.com',
    '10:00-23:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Семейная кофейня',
    13,
    false,
    'ACTIVE',
    '79013004455',
    'familycafe@example.com',
    '9:00-22:00',
    NULL,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

-- Insert into StoreAdditives
INSERT INTO
  store_additives (additive_id, store_id, price)
VALUES
  (1, 1, 60.00),
  (2, 1, 75.00),
  (3, 2, 85.00),
  (1, 3, 65.00),
  (4, 1, 50.00),
  (5, 4, 90.00),
  (6, 5, 70.00),
  (7, 6, 100.00),
  (8, 7, 55.00),
  (9, 8, 45.00),
  (10, 9, 15.00),
  (11, 10, 5.00),
  (12, 11, 80.00),
  (13, 12, 40.00),
  (14, 13, 30.00),
  (3, 4, 75.00),
  (6, 6, 85.00),
  (9, 7, 65.00),
  (5, 8, 95.00),
  (2, 9, 70.00);

-- Insert into StoreProductSizes
INSERT INTO
  store_product_sizes (product_size_id, store_id, price)
VALUES
  (1, 1, 210.00),
  (2, 1, 260.00),
  (3, 2, 310.00),
  (1, 3, 220.00),
  (2, 4, 270.00),
  (3, 5, 330.00),
  (1, 6, 215.00),
  (2, 7, 265.00),
  (3, 8, 320.00),
  (1, 9, 225.00),
  (2, 10, 275.00),
  (3, 11, 315.00),
  (1, 12, 230.00),
  (2, 13, 280.00),
  (3, 1, 300.00),
  (1, 2, 210.00),
  (2, 3, 270.00),
  (3, 4, 325.00),
  (1, 5, 250.00),
  (2, 6, 290.00);

-- Insert into StoreProduct
INSERT INTO
  store_products (product_id, store_id, is_available)
VALUES
  -- All products available in the first store
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
  -- Products in other stores with varied availability
  (1, 2, true),
  (2, 2, false),
  (3, 2, true),
  (4, 2, true),
  (5, 3, false),
  (6, 3, true),
  (7, 3, false),
  (8, 3, true),
  (9, 4, true),
  (10, 4, true),
  (11, 4, false),
  (12, 4, true),
  (13, 5, true),
  (14, 5, true),
  (15, 5, false),
  (16, 5, true),
  (17, 6, true),
  (18, 6, false),
  (19, 6, true),
  (20, 6, true);

-- Insert into ProductAdditive
INSERT INTO
  product_additives (product_size_id, additive_id)
VALUES
  (1, 1),
  (1, 2),
  (2, 3),
  (3, 1),
  (4, 2),
  (5, 4),
  (6, 5),
  (7, 3),
  (8, 6),
  (9, 7),
  (10, 2),
  (11, 8),
  (12, 9),
  (13, 10),
  (14, 11),
  (15, 12),
  (16, 13),
  (17, 14),
  (18, 1),
  (19, 3),
  (20, 6),
  (3, 8),
  (4, 10),
  (5, 12),
  (6, 14),
  (7, 1),
  (8, 4),
  (9, 5),
  (10, 9);

-- Insert into DefaultProductAdditive
INSERT INTO
  default_product_additives (product_id, additive_id)
VALUES
  (1, 1),
  (1, 2),
  (2, 3),
  (2, 4),
  (3, 5),
  (3, 1),
  (4, 2),
  (4, 6),
  (5, 7),
  (5, 8),
  (6, 9),
  (6, 3),
  (7, 2),
  (8, 4),
  (9, 10),
  (10, 11),
  (11, 5),
  (12, 12),
  (13, 6),
  (14, 13),
  (15, 7),
  (16, 1),
  (17, 9),
  (18, 8),
  (19, 14),
  (20, 10);

-- Insert into Ingredients
INSERT INTO
  ingredients (
    name,
    calories,
    fat,
    carbs,
    proteins,
    expires_at,
    unit_id
  )
VALUES
  (
    'Сахар',
    387,
    0,
    100,
    0,
    '2024-12-31 00:00:00+00',
    2
  ),
  (
    'Молоко',
    42,
    1,
    5,
    3,
    '2024-01-15 00:00:00+00',
    3
  ),
  (
    'Шоколад',
    546,
    30,
    61,
    7,
    '2024-06-30 00:00:00+00',
    1
  ),
  (
    'Корица',
    247,
    1.2,
    81,
    4,
    '2024-09-15 00:00:00+00',
    2
  ),
  ('Мед', 304, 0, 82, 0, '2024-10-20 00:00:00+00', 2),
  (
    'Ваниль',
    288,
    12,
    55,
    0,
    '2025-01-30 00:00:00+00',
    2
  ),
  (
    'Орехи',
    607,
    54,
    18,
    20,
    '2024-08-15 00:00:00+00',
    2
  ),
  (
    'Кокосовое молоко',
    230,
    23,
    6,
    2,
    '2024-05-01 00:00:00+00',
    3
  ),
  (
    'Яблоки',
    52,
    0.2,
    14,
    0.3,
    '2024-02-28 00:00:00+00',
    1
  ),
  (
    'Бананы',
    96,
    0.3,
    27,
    1.3,
    '2024-03-15 00:00:00+00',
    1
  ),
  (
    'Сливки',
    195,
    20,
    3,
    2,
    '2024-04-10 00:00:00+00',
    1
  ),
  (
    'Апельсины',
    47,
    0.1,
    12,
    0.9,
    '2024-02-05 00:00:00+00',
    1
  ),
  (
    'Мята',
    44,
    0.7,
    8,
    3.3,
    '2024-06-01 00:00:00+00',
    3
  ),
  (
    'Лимонный сок',
    123,
    0.2,
    6,
    0.3,
    '2024-03-20 00:00:00+00',
    2
  ),
  (
    'Какао-порошок',
    228,
    13,
    58,
    19,
    '2025-06-30 00:00:00+00',
    2
  ),
  (
    'Кленовый сироп',
    261,
    0,
    67,
    0,
    '2024-12-15 00:00:00+00',
    2
  ),
  (
    'Клубника',
    33,
    0.3,
    8,
    0.7,
    '2024-05-05 00:00:00+00',
    1
  ),
  (
    'Имбирь',
    80,
    0.8,
    18,
    1.8,
    '2024-07-01 00:00:00+00',
    2
  ),
  ('Соль', 0, 0, 0, 0, '2026-12-31 00:00:00+00', 2),
  (
    'Фисташки',
    562,
    45,
    28,
    20,
    '2024-09-15 00:00:00+00',
    1
  );

-- Insert into ProductIngredients
INSERT INTO
  product_ingredients (ingredient_id, product_size_id, quantity)
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

-- Фисташки
-- Insert into Customer
INSERT INTO
  customers (name, password, phone, is_verified, is_banned)
VALUES
  (
    'Иван Иванов',
    'hashed_password_123',
    '79031234567',
    true,
    false
  ),
  (
    'Мария Смирнова',
    'hashed_password_456',
    '79876543210',
    false,
    false
  ),
  (
    'Алексей Петров',
    'hashed_password_789',
    '79998887766',
    true,
    false
  );

-- Insert into Employee
INSERT INTO
  employees (
    name,
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
    'Елена Соколова',
    '79551234567',
    'elena@example.com',
    'ADMIN',
    'STORE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Сергей Павлов',
    '79667778899',
    'sergey@example.com',
    'DIRECTOR',
    'STORE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Анна Федорова',
    '79223334455',
    'anna@example.com',
    'MANAGER',
    'STORE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Иван Иванов',
    '79161234567',
    'ivan@example.com',
    'MANAGER',
    'STORE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Мария Смирнова',
    '79345566778',
    'maria@example.com',
    'WAREHOUSE_EMPLOYEE',
    'WAREHOUSE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Олег Кузнецов',
    '79991234567',
    'oleg@example.com',
    'WAREHOUSE_EMPLOYEE',
    'WAREHOUSE',
    false,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Татьяна Орлова',
    '79882233445',
    'tatiana@example.com',
    'WAREHOUSE_EMPLOYEE',
    'WAREHOUSE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Алексей Попов',
    '79002221133',
    'alexei@example.com',
    'WAREHOUSE_EMPLOYEE',
    'WAREHOUSE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Юлия Петрова',
    '79115555666',
    'yulia@example.com',
    'WAREHOUSE_EMPLOYEE',
    'WAREHOUSE',
    true,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'Дмитрий Фролов',
    '79553334456',
    'dmitry@example.com',
    'WAREHOUSE_EMPLOYEE',
    'WAREHOUSE',
    false,
    '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

-- Insert into StoreEmployee
INSERT INTO
  store_employees (
    employee_id,
    store_id,
    is_franchise,
    created_at,
    updated_at
  )
VALUES
  (1, 1, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (2, 2, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (3, 3, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (4, 1, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (5, 2, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert into WarehouseEmployee
INSERT INTO
  warehouse_employees (employee_id, warehouse_id, created_at, updated_at)
VALUES
  (6, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (7, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (8, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (9, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (10, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert into EmployeeAudit
INSERT INTO
  employee_audits (start_work_at, end_work_at, employee_id)
VALUES
  (
    '2024-10-01 09:00:00+00',
    '2024-10-01 17:00:00+00',
    1
  ),
  (
    '2024-10-02 09:00:00+00',
    '2024-10-02 17:00:00+00',
    2
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
  store_warehouses (store_id, warehouse_id)
VALUES
  (1, 1), -- Store 1 linked to Central Warehouse in Moscow
  (2, 2), -- Store 2 linked to Central Warehouse in St. Petersburg
  (3, 3), -- Store 3 linked to Central Warehouse in Ekaterinburg
  (4, 4), -- Store 4 linked to Central Warehouse in Novosibirsk
  (5, 5);

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
  (1, 1, 20, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (
    1,
    2,
    50,
    500,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    2,
    1,
    20,
    30,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    2,
    2,
    30,
    100,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    2,
    1,
    40,
    80,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    3,
    1,
    10000,
    1000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    3,
    3,
    120,
    500,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO
  warehouses (facility_address_id, name)
VALUES
  (1, 'Central Warehouse - Moscow'),
  (2, 'Central Warehouse - St. Petersburg'),
  (3, 'Central Warehouse - Ekaterinburg'),
  (4, 'Central Warehouse - Novosibirsk'),
  (5, 'Central Warehouse - Kazan');

-- Store 5 linked to Central Warehouse in Kazan
-- Insert into Suppliers
INSERT INTO
  suppliers (name, contact_email, contact_phone, address)
VALUES
  (
    'Nestlé',
    'contact@nestle.com',
    '+1 800 225 2270',
    'Avenue Nestlé 55, 1800 Vevey, Switzerland'
  ),
  (
    'Coca-Cola',
    'info@coca-cola.com',
    '+1 800 438 2653',
    '1 Coca-Cola Plaza, Atlanta, GA 30313, USA'
  ),
  (
    'PepsiCo',
    'support@pepsico.com',
    '+1 914 253 2000',
    '700 Anderson Hill Rd, Purchase, NY 10577, USA'
  ),
  (
    'Lipton',
    'info@lipton.com',
    '+44 800 776 647',
    'Unilever House, Springfield Dr, Leatherhead KT22 7GR, UK'
  ),
  (
    'Starbucks',
    'help@starbucks.com',
    '+1 800 782 7282',
    '2401 Utah Ave S, Seattle, WA 98134, USA'
  ),
  (
    'Mondelez',
    'support@mondelez.com',
    '+1 855 535 5648',
    '100 Deforest Ave, East Hanover, NJ 07936, USA'
  ),
  (
    'Danone',
    'contact@danone.com',
    '+33 1 44 35 20 20',
    '17 Boulevard Haussmann, 75009 Paris, France'
  ),
  (
    'Mars',
    'support@mars.com',
    '+1 703 821 4900',
    '6885 Elm St, McLean, VA 22101, USA'
  ),
  (
    'Unilever',
    'contact@unilever.com',
    '+44 20 7822 5252',
    '100 Victoria Embankment, London EC4Y 0DY, UK'
  ),
  (
    'General Mills',
    'support@generalmills.com',
    '+1 800 248 7310',
    '1 General Mills Blvd, Minneapolis, MN 55426, USA'
  );

INSERT INTO
  stock_materials (
    name,
    description,
    safety_stock,
    expiration_flag,
    unit_id,
    category,
    barcode,
    expiration_period_in_days,
    is_active
  )
VALUES
  (
    'Молоко',
    '1L pack of milk',
    50,
    TRUE,
    3,
    'Dairy',
    '111111111111',
    1095,
    TRUE
  ),
  (
    'Сахар',
    '1kg pack of sugar',
    20,
    TRUE,
    2,
    'Sweeteners',
    '222222222222',
    1095,
    TRUE
  ),
  (
    'Шоколад',
    '500g pack of chocolate',
    15,
    TRUE,
    2,
    'Confectionery',
    '333333333333',
    730,
    TRUE
  ),
  (
    'Корица',
    '200g pack of cinnamon',
    10,
    TRUE,
    2,
    'Spices',
    '444444444444',
    1460,
    TRUE
  ),
  (
    'Ваниль',
    '50ml vanilla extract bottle',
    25,
    TRUE,
    4,
    'Flavorings',
    '555555555555',
    1460,
    TRUE
  );

-- Vanilla linked to stock material
INSERT INTO
  stock_material_packages (stock_material_id, package_size, package_unit_id)
VALUES
  (1, 1.0, 3), -- 1L Milk
  (2, 1.0, 2), -- 1kg Sugar
  (3, 0.5, 2), -- 500g Chocolate
  (4, 0.2, 2), -- 200g Cinnamon
  (5, 0.05, 4);

-- 50ml Vanilla
INSERT INTO
  supplier_warehouse_deliveries (
    stock_material_id,
    supplier_id,
    warehouse_id,
    barcode,
    quantity,
    delivery_date,
    expiration_date
  )
VALUES
  (
    1,
    1,
    1,
    '111111111111',
    50,
    '2024-12-01',
    '2026-12-01'
  ), -- Milk Delivery
  (
    2,
    2,
    1,
    '222222222222',
    30,
    '2024-12-05',
    '2025-06-05'
  ), -- Sugar Delivery
  (
    3,
    1,
    1,
    '333333333333',
    40,
    '2024-11-20',
    '2025-11-20'
  ), -- Chocolate Delivery
  (
    4,
    2,
    2,
    '444444444444',
    20,
    '2024-12-10',
    '2026-06-10'
  ), -- Cinnamon Delivery
  (
    5,
    1,
    2,
    '555555555555',
    15,
    '2024-12-15',
    '2027-12-15'
  );

-- Vanilla Delivery
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
    request_date,
    created_at,
    updated_at
  )
VALUES
  (
    1,
    1,
    'CREATED',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    2,
    2,
    'PROCESSED',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    3,
    3,
    'IN_DELIVERY',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    4,
    4,
    'COMPLETED',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

-- Insert into StockRequestIngredients (Items for Stock Requests)
INSERT INTO
  stock_request_ingredients (
    stock_request_id,
    ingredient_id,
    quantity,
    created_at,
    updated_at
  )
VALUES
  -- StockRequest 1 (Store 1 -> Warehouse 1)
  (1, 1, 10.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Sugar
  (1, 2, 20.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Milk
  -- StockRequest 2 (Store 2 -> Warehouse 2)
  (2, 3, 5.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Chocolate
  (2, 4, 2.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Cinnamon
  -- StockRequest 3 (Store 3 -> Warehouse 3)
  (3, 5, 1.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Vanilla
  (3, 1, 15.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Sugar
  -- StockRequest 4 (Store 4 -> Warehouse 4)
  (4, 2, 10.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Milk
  (4, 3, 8.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO
  ingredient_stock_material_mappings (ingredient_id, stock_material_id)
VALUES
  (2, 1), -- Milk linked to stock material
  (1, 2), -- Sugar linked to stock material
  (3, 3), -- Chocolate linked to stock material
  (4, 4), -- Cinnamon linked to stock material
  (6, 5);