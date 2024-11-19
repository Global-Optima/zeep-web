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
  stores (name, facility_address_id, is_franchise, admin_id)
VALUES
  ('Центральное кафе', 1, false, NULL),
  ('Кофейня на углу', 2, true, NULL),
  ('Маленький магазин на Советской', 3, true, NULL),
  ('Кофейня у вокзала', 4, false, NULL),
  ('Городской кофе', 5, true, NULL),
  ('Летняя терраса', 6, true, NULL),
  ('Кафе на проспекте', 7, false, NULL),
  ('Заведение у реки', 8, true, NULL),
  ('Чайный дом', 9, false, NULL),
  ('Кофе и компании', 10, true, NULL),
  ('Парк-кафе', 11, false, NULL),
  ('Восточный уголок', 12, true, NULL),
  ('Семейная кофейня', 13, false, NULL);

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
  ingredients (name, calories, fat, carbs, proteins, expires_at)
VALUES
  ('Сахар', 387, 0, 100, 0, '2024-12-31 00:00:00+00'),
  ('Молоко', 42, 1, 5, 3, '2024-01-15 00:00:00+00'),
  (
    'Шоколад',
    546,
    30,
    61,
    7,
    '2024-06-30 00:00:00+00'
  ),
  (
    'Корица',
    247,
    1.2,
    81,
    4,
    '2024-09-15 00:00:00+00'
  ),
  ('Мед', 304, 0, 82, 0, '2024-10-20 00:00:00+00'),
  (
    'Ваниль',
    288,
    12,
    55,
    0,
    '2025-01-30 00:00:00+00'
  ),
  (
    'Орехи',
    607,
    54,
    18,
    20,
    '2024-08-15 00:00:00+00'
  ),
  (
    'Кокосовое молоко',
    230,
    23,
    6,
    2,
    '2024-05-01 00:00:00+00'
  ),
  (
    'Яблоки',
    52,
    0.2,
    14,
    0.3,
    '2024-02-28 00:00:00+00'
  ),
  (
    'Бананы',
    96,
    0.3,
    27,
    1.3,
    '2024-03-15 00:00:00+00'
  ),
  ('Сливки', 195, 20, 3, 2, '2024-04-10 00:00:00+00'),
  (
    'Апельсины',
    47,
    0.1,
    12,
    0.9,
    '2024-02-05 00:00:00+00'
  ),
  ('Мята', 44, 0.7, 8, 3.3, '2024-06-01 00:00:00+00'),
  (
    'Лимонный сок',
    22,
    0.2,
    6,
    0.3,
    '2024-03-20 00:00:00+00'
  ),
  (
    'Какао-порошок',
    228,
    13,
    58,
    19,
    '2025-06-30 00:00:00+00'
  ),
  (
    'Кленовый сироп',
    261,
    0,
    67,
    0,
    '2024-12-15 00:00:00+00'
  ),
  (
    'Клубника',
    33,
    0.3,
    8,
    0.7,
    '2024-05-05 00:00:00+00'
  ),
  (
    'Имбирь',
    80,
    0.8,
    18,
    1.8,
    '2024-07-01 00:00:00+00'
  ),
  ('Соль', 0, 0, 0, 0, '2026-12-31 00:00:00+00'),
  (
    'Фисташки',
    562,
    45,
    28,
    20,
    '2024-09-15 00:00:00+00'
  );

-- Insert into ProductIngredients
INSERT INTO
  product_ingredients (item_ingredient_id, product_id)
VALUES
  (1, 1),
  (2, 1),
  (1, 2),
  (3, 3),
  (4, 1),
  (5, 2),
  (6, 2),
  (7, 3),
  (8, 4),
  (9, 5),
  (10, 5),
  (11, 6),
  (12, 7),
  (13, 8),
  (14, 8),
  (15, 9),
  (16, 10),
  (17, 11),
  (18, 12),
  (19, 13),
  (20, 14),
  (3, 4),
  (4, 5),
  (6, 6),
  (7, 7),
  (8, 8),
  (9, 9),
  (10, 10),
  (11, 11),
  (12, 12),
  (13, 13);

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

-- Insert into EmployeeRole
INSERT INTO
  employee_roles (name)
VALUES
  ('Менеджер'),
  ('Бариста'),
  ('Кассир'),
  ('Администратор'),
  ('Официант'),
  ('Уборщик'),
  ('Шеф-повар'),
  ('Помощник повара'),
  ('Старший кассир'),
  ('Складской работник'),
  ('Заведующий'),
  ('Маркетолог'),
  ('Менеджер по персоналу');

-- Insert into Employee
INSERT INTO
  employees (name, phone, email, role_id, store_id, is_active)
VALUES
  (
    'Елена Соколова',
    '79551234567',
    'elena@example.com',
    1,
    1,
    true
  ),
  (
    'Сергей Павлов',
    '79667778899',
    'sergey@example.com',
    2,
    2,
    true
  ),
  (
    'Анна Федорова',
    '79223334455',
    'anna@example.com',
    3,
    1,
    true
  ),
  (
    'Иван Иванов',
    '79161234567',
    'ivan@example.com',
    4,
    3,
    true
  ),
  (
    'Мария Смирнова',
    '79345566778',
    'maria@example.com',
    5,
    2,
    true
  ),
  (
    'Олег Кузнецов',
    '79991234567',
    'oleg@example.com',
    6,
    1,
    false
  ),
  (
    'Татьяна Орлова',
    '79882233445',
    'tatiana@example.com',
    7,
    4,
    true
  ),
  (
    'Алексей Попов',
    '79002221133',
    'alexei@example.com',
    8,
    3,
    true
  ),
  (
    'Юлия Петрова',
    '79115555666',
    'yulia@example.com',
    9,
    2,
    true
  ),
  (
    'Дмитрий Фролов',
    '79553334456',
    'dmitry@example.com',
    10,
    1,
    false
  ),
  (
    'Наталья Волкова',
    '79225556677',
    'natalya@example.com',
    11,
    5,
    true
  ),
  (
    'Светлана Крылова',
    '79331112223',
    'svetlana@example.com',
    12,
    3,
    true
  ),
  (
    'Виктор Соколов',
    '79442233445',
    'victor@example.com',
    13,
    4,
    true
  ),
  (
    'Андрей Лебедев',
    '79118887799',
    'andrei@example.com',
    1,
    1,
    true
  ),
  (
    'Ольга Николаева',
    '79664445566',
    'olga@example.com',
    2,
    5,
    true
  ),
  (
    'Максим Морозов',
    '79552233456',
    'maksim@example.com',
    3,
    2,
    true
  ),
  (
    'Людмила Громова',
    '79774445511',
    'lyudmila@example.com',
    4,
    3,
    false
  ),
  (
    'София Зайцева',
    '79883332211',
    'sofia@example.com',
    5,
    4,
    true
  ),
  (
    'Владимир Козлов',
    '79338887799',
    'vladimir@example.com',
    6,
    5,
    true
  ),
  (
    'Алена Белова',
    '79229998877',
    'alena@example.com',
    7,
    1,
    true
  );

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