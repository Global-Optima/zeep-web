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
    ('Килограмм', 1.0),
    ('Грамм', 0.001),
    ('Литр', 1.0),
    ('Миллилитр', 0.001);

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
  ('Коктейли','Алкогольные и безалкогольные коктейли'),
  ('Газированные напитки','Прохладительные газированные напитки'),
  ('Энергетические напитки','Напитки для повышения энергии'),
  ('Молочные коктейли','Коктейли на основе молока и сливок'),
  ('Минеральная вода','Природная и газированная минеральная вода'),
  ('Фраппе', 'Кофейные напитки со льдом'),
  ('Травяные чаи', 'Настои и чаи на травах'),
  ('Круассаны', 'Свежая выпечка с хрустящей корочкой и разнообразной начинкой — идеально к кофе');

-- Insert into AdditiveCategory
INSERT INTO
  additive_categories (name, description, is_multiple_select)
VALUES
  ('Ароматизаторы','Дополнительные вкусы для усиления аромата',true),
  ('Подсластители', 'Добавление сладости', false),
  ('Топпинги','Украшения для десертов и напитков',true),
  ('Сиропы', 'Сиропы для напитков и десертов', false),
  ('Специи', 'Ароматные специи для напитков', true),
  ('Молочные добавки','Добавление молока и сливок',false),
  ('Фрукты', 'Свежие и сушеные фрукты', true),
  ('Орехи','Измельченные орехи для украшения',false),
  ('Шоколад', 'Шоколадная стружка и какао', true),
  ('Мед','Естественный подсластитель на основе меда',false),
  ('Сахарные добавки', 'Различные виды сахара', true),
  ('Кубики льда','Лед для охлаждения напитков',false);

  -- Insert into IngredientCategories
INSERT INTO ingredient_categories (name, description)
VALUES
    ('Молочные продукты', 'Категория для молока, сливок, и других молочных ингредиентов'),
    ('Фрукты', 'Категория для фруктов, таких как яблоки, бананы, апельсины'),
    ('Подсластители', 'Категория для сахара, мёда, сиропов'),
    ('Специи', 'Категория для специй, таких как корица, ваниль, мята, имбирь'),
    ('Орехи и семена', 'Категория для орехов, фисташек, и других семян'),
    ('Шоколад и какао', 'Категория для шоколада, какао-порошка, и других шоколадных продуктов');


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

-- Insert into Additives
INSERT INTO
  additives (
    name,
    description,
    base_price,
    size,
    unit_id,
    additive_category_id,
    image_url
  )
VALUES
  (
    'Ванильный сироп',
    'Сладкий аромат ванили',
    50.00,
    20,
    1,
    1,
    'https://monin.us/cdn/shop/files/750mL-Vanilla.png?v=1724939521&width=1946'
  ),
  (
    'Карамельный сироп',
    'Сладкий вкус карамели',
    70.00,
    40,
    3,
    1,
    'https://www.giffard.com/327-large_default/karamell-sirup.jpg'
  ),
  (
    'Взбитые сливки',
    'Сливки для украшения напитков',
    80.00,
    500,
    2,
    3,
    'https://static.vecteezy.com/system/resources/previews/033/888/680/non_2x/white-whipped-cream-with-ai-generated-free-png.png'
  ),
  (
    'Корица',
    'Пряная корица для добавления аромата',
    30.00,
    20,
    1,
    2,
    'https://png.pngtree.com/png-vector/20240803/ourmid/pngtree-cinnamon-isolated-on-white-background-png-image_13361859.png'
  ),
  (
    'Шоколадная крошка',
    'Крошка из шоколада для украшения',
    90.00,
    20,
    1,
    3,
    'https://mlfwajaoc5gr.i.optimole.com/w:1080/h:1080/q:mauto/ig:avif/http://icecreambakery.in/wp-content/uploads/2024/03/Chocolate-Chips-Manufacturer-RPG-Industries-2.png'
  ),
  (
    'Мед',
    'Натуральный мед для подслащивания',
    60.00,
    20,
    1,
    2,
    'https://png.pngtree.com/png-vector/20240801/ourmid/pngtree-3d-bowl-of-honey-with-a-driper-on-transparent-background-png-image_13330438.png'
  ),
  (
    'Кокосовое молоко',
    'Молоко с ароматом кокоса',
    100.00,
    20,
    1,
    4,
    'https://static.vecteezy.com/system/resources/previews/049/642/396/non_2x/coconut-milk-in-a-coconut-shell-with-a-transparent-background-png.png'
  ),
  (
    'Клубничный сироп',
    'Сироп с ароматом клубники',
    65.00,
    20,
    1,
    1,
    'https://monin.ca/cdn/shop/files/750mL-Strawberry.png?v=1727882931&width=1445'
  ),
  (
    'Мята',
    'Свежие листья мяты для аромата',
    40.00,
    20,
    1,
    5,
    'https://static.vecteezy.com/system/resources/thumbnails/047/732/281/small/basil-leaves-in-a-glass-bowl-transparent-background-png.png'
  ),
  (
    'Сахар',
    'Белый сахар для подслащивания',
    10.00,
    20,
    1,
    6,
    'https://static.vecteezy.com/system/resources/thumbnails/044/902/343/small_2x/sweet-delight-enhance-your-recipes-with-sugar-in-bowl-free-png.png'
  ),
  (
    'Кубики льда',
    'Лед для охлаждения напитков',
    5.00,
    20,
    1,
    7,
    'https://img.pikbest.com/origin/09/14/36/70CpIkbEsTeTk.png!sw800'
  ),
  (
    'Сироп из клена',
    'Ароматный сироп из клена',
    80.00,
    20,
    1,
    1,
    'https://png.pngtree.com/png-vector/20240812/ourmid/pngtree-sweet-syrup-drizzling-over-breakfast-treats-png-image_13462282.png'
  ),
  (
    'Лимонный сок',
    'Свежий лимонный сок для аромата',
    45.00,
    20,
    2,
    2,
    'https://www.bestlobster.com/test/wp-content/uploads/2021/02/Lemon-PNG-Image.png'
  ),
  (
    'Какао-порошок',
    'Порошок какао для украшения напитков',
    35.00,
    15,
    2,
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
        '+79001112233',
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
        '+79002223344',
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
        '+79003334455',
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
        '+79004445566',
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
        '+79005556677',
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
        '+79006667788',
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
        '+79007778899',
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
        '+79008889900',
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
        '+79009990011',
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
        '+79010001122',
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
        '+79011002233',
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
        '+79012003344',
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
        '+79013004455',
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


-- Insert into StoreProductSizes
INSERT INTO
    store_product_sizes (product_size_id, store_product_id, price)
VALUES
    -- Store 1 Products (All sizes included)
    (1, 1, 1000.00), -- S size for Product 1
    (2, 1, 1250.00), -- M size for Product 1
    (3, 1, 1500.00), -- L size for Product 1

    -- Store 2 Products (Only M and L sizes included)
    -- S size excluded intentionally
    (2, 2, 1100.00), -- M size for Product 2
    (3, 2, 1350.00), -- L size for Product 2

    -- Store 3 Products (Only S and M sizes included)
    -- L size excluded intentionally
    (4, 3, 750.00),  -- S size for Product 3
    (5, 3, 900.00),  -- M size for Product 3

    -- Store 4 Products (All sizes included)
    (6, 4, 600.00),  -- S size for Product 4
    (7, 4, 850.00),  -- M size for Product 4
    (8, 4, 1100.00), -- L size for Product 4

    -- Store 5 Products (Only S size included)
    -- M and L sizes excluded intentionally
    (9, 5, 900.00),  -- S size for Product 5

    -- Store 6 Products (All sizes included)
    (10, 6, 1500.00), -- S size for Product 6
    (11, 6, 1750.00), -- L size for Product 6

    -- Store 7 Products (Only M size included)
    -- S and L sizes excluded intentionally
    (12, 7, 750.00),  -- M size for Product 7

    -- Store 8 Products (All sizes included)
    (13, 8, 600.00),  -- S size for Product 8
    (14, 8, 750.00),  -- L size for Product 8

    -- Store 9 Products (All sizes included)
    (15, 9, 800.00),  -- M size for Product 9
    (16, 9, 1000.00), -- L size for Product 9

    -- Store 10 Products (Only S size included)
    -- M and L sizes excluded intentionally
    (17, 10, 500.00), -- S size for Product 10

    -- Store 11 Products (Only L size included)
    -- S and M sizes excluded intentionally
    (18, 11, 1200.00); -- L size for Product 11;

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

-- Insert into Ingredients
INSERT INTO
  ingredients (name, calories, fat, carbs, proteins, expiration_in_days, unit_id, category_id)
VALUES
  ('Сахар', 387, 0, 100, 0, 365, 1, 3), -- Подсластители
  ('Молоко', 42, 1, 5, 3, 7, 3, 1), -- Молочные продукты
  ('Шоколад', 546, 30, 61, 7, 180, 2, 6), -- Шоколад и какао
  ('Корица', 247, 1.2, 81, 4, 365, 2, 4), -- Специи
  ('Мед', 304, 0, 82, 0, 365, 2, 3), -- Подсластители
  ('Ваниль', 288, 12, 55, 0, 730, 2, 4), -- Специи
  ('Орехи', 607, 54, 18, 20, 365, 2, 5), -- Орехи и семена
  ('Кокосовое молоко', 230, 23, 6, 2, 120, 3, 1), -- Молочные продукты
  ('Яблоки', 52, 0.2, 14, 0.3, 14, 1, 2), -- Фрукты
  ('Бананы', 96, 0.3, 27, 1.3, 7, 1, 2), -- Фрукты
  ('Сливки', 195, 20, 3, 2, 10, 1, 1), -- Молочные продукты
  ('Апельсины', 47, 0.1, 12, 0.9, 14, 1, 2), -- Фрукты
  ('Мята', 44, 0.7, 8, 3.3, 180, 2, 4), -- Специи
  ('Лимонный сок', 123, 0.2, 6, 0.3, 60, 3, 2), -- Фрукты
  ('Какао-порошок', 228, 13, 58, 19, 730, 2, 6), -- Шоколад и какао
  ('Кленовый сироп', 261, 0, 67, 0, 365, 4, 3), -- Подсластители
  ('Клубника', 33, 0.3, 8, 0.7, 10, 1, 2), -- Фрукты
  ('Имбирь', 80, 0.8, 18, 1.8, 180, 1, 4), -- Специи
  ('Соль', 0, 0, 0, 0, 1095, 1, 4), -- Специи
  ('Фисташки', 562, 45, 28, 20, 365, 1, 5); -- Орехи и семена


-- Insert into AdditiveIngredients
INSERT INTO
    additive_ingredients (additive_id, ingredient_id, quantity)
VALUES
    (1, 1, 60.00),
    (2, 1, 75.00),
    (3, 2, 85.00),
    (1, 3, 65.00),
    (4, 1, 50.00),
    (5, 4, 90.00),
    (6, 5, 70.00),
    (3, 4, 75.00),
    (6, 6, 85.00),
    (5, 1, 95.00),
    (2, 2, 70.00);

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

-- Фисташки
-- Insert into Customer
INSERT INTO
    customers (first_name, last_name, password, phone, is_verified, is_banned)
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


-- Insert into Employee
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
        '+79551234567',
        'elena@example.com',
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
        '+79667778899',
        'sergey@example.com',
        'DIRECTOR',
        'STORE',
        true,
        '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Анна',
        'Федорова',
        '+79223334455',
        'anna@example.com',
        'MANAGER',
        'STORE',
        true,
        '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Иван',
        'Иванов',
        '+79161234567',
        'ivan@example.com',
        'BARISTA',
        'STORE',
        true,
        '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Мария',
        'Смирнова',
        '+79345566778',
        'maria@example.com',
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
        '+79991234567',
        'oleg@example.com',
        'WAREHOUSE_EMPLOYEE',
        'WAREHOUSE',
        false,
        '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Татьяна',
        'Орлова',
        '+79882233445',
        'tatiana@example.com',
        'WAREHOUSE_EMPLOYEE',
        'WAREHOUSE',
        true,
        '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Алексей',
        'Попов',
        '+79002221133',
        'alexei@example.com',
        'WAREHOUSE_EMPLOYEE',
        'WAREHOUSE',
        true,
        '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Юлия',
        'Петрова',
        '+79115555666',
        'yulia@example.com',
        'WAREHOUSE_EMPLOYEE',
        'WAREHOUSE',
        true,
        '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'Дмитрий',
        'Фролов',
        '+79553334456',
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
  ('2024-10-01 09:00:00+00','2024-10-01 17:00:00+00',1),
  ('2024-10-02 09:00:00+00','2024-10-02 17:00:00+00',2);

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
  (5, 5),
  (6, 6);

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
    (1, 2, 50, 500, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 1, 20, 30, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 2, 30, 100, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, 1, 10000, 1000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, 3, 120, 500, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);


-- Insert stock material categories
INSERT INTO stock_material_categories (name, description)
VALUES 
    ('Молочные продукты', 'Молоко, сливки, йогурты и другие молочные продукты'),
    ('Подсластители', 'Сахар, мед и другие подсластители'),
    ('Кондитерские изделия', 'Шоколад, какао и другие кондитерские ингредиенты'),
    ('Специи', 'Различные специи и пряности'),
    ('Ароматизаторы', 'Ванильный экстракт и другие ароматизаторы');

-- Insert stock materials with Russian names and category references
INSERT INTO stock_materials (name, description, ingredient_id, safety_stock, unit_id, category_id, barcode, expiration_period_in_days, is_active) 
VALUES 
    ('Простоквашино Молоко 3.2%', 'Молоко пастеризованное 1л', 2, 50, 3, 1, '111111111111', 1095, TRUE),
    ('Русский сахар Экстра', 'Сахар песок высший сорт 1кг', 1, 20, 1, 2, '222222222222', 1095, TRUE),
    ('Бабаевский горький шоколад 75%', 'Темный шоколад 500г', 3, 15, 2, 3, '333333333333', 730, TRUE),
    ('Индийская корица молотая премиум', 'Корица молотая 200г', 4, 10, 2, 4, '444444444444', 1460, TRUE),
    ('Dr.Oetker Ванильный экстракт', 'Натуральный экстракт ванили 50мл', 5, 25, 4, 5, '555555555555', 1460, TRUE);

-- Insert into Suppliers
INSERT INTO suppliers (name, contact_email, contact_phone, city, address)
VALUES
    ('ООО "Нестле Россия"', 'contact@nestle.ru', '+79005556677', 'Москва', 'Павелецкая площадь, д. 2, стр. 1'),
    ('АО "Кока-Кола ЭйчБиСи Евразия"', 'info@coca-cola.ru', '+78002002222', 'Москва', 'ул. Новоорловская, д. 7'),
    ('ООО "ПепсиКо Холдингс"', 'support@pepsico.ru', '+78007001000', 'Московская обл.', 'г. Солнечногорск, территория свободной экономической зоны "Шерризон", стр. 1'),
    ('ООО "Юнилевер Русь"', 'info@unilever.ru', '+78002001200', 'Москва', 'ул. Сергея Макеева, д. 13'),
    ('ООО "Штарбакс"', 'help@starbucks.ru', '+78001008333', 'Москва', 'ул. Большая Новодмитровская, д. 23, стр. 1'),
    ('ООО "Мон`дэлис Русь"', 'support@mdlz.ru', '+74959602424', 'Владимирская обл.', 'г. Покров, ул. Франца Штольверка, д. 10'),
    ('АО "ДАНОН РОССИЯ"', 'contact@danone.ru', '+78002000201', 'Москва', 'ул. Вятская, д. 27, корп. 13-14'),
    ('ООО "Марс"', 'support@mars.ru', '+74957212100', 'Московская обл.', 'г. Ступино, ул. Ситенка, д. 12'),
    ('ООО "Юнилевер Русь"', 'contact@unilever.ru', '+78002001201', 'Омск', 'ул. 10 лет Октября, д. 205'),
    ('ООО "Дженерал Миллс Рус"', 'support@generalmills.ru', '+74959373400', 'Москва', 'ул. Большая Новодмитровская, д. 14, стр. 2');

INSERT INTO stock_material_packages (stock_material_id, size, unit_id)
VALUES
    (1, 1.0, 3), -- 1L Milk
    (2, 1.0, 1), -- 1kg Sugar
    (3, 500, 2), -- 500g Chocolate
    (4, 200, 2), -- 200g Cinnamon
    (5, 50, 4); -- 50ml Vanilla


-- Insert into supplier_warehouse_deliveries
INSERT INTO supplier_warehouse_deliveries (supplier_id, warehouse_id, delivery_date, created_at, updated_at)
VALUES
    (1, 1, '2024-09-01', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Delivery 1
    (2, 1,'2024-11-01',  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Delivery 2
    (1, 2, '2024-10-01', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP); -- Delivery 3

-- Insert into supplier_warehouse_delivery_materials
INSERT INTO supplier_warehouse_delivery_materials (delivery_id, stock_material_id, package_id, barcode, quantity, expiration_date)
VALUES
    (1, 1, 1, '111111111111', 50,  '2026-12-01'), -- Milk Delivery
    (1, 2, 2, '222222222222', 30,  '2025-06-05'), -- Sugar Delivery
    (2, 3, 3, '333333333333', 40,  '2025-11-20'), -- Chocolate Delivery
    (3, 4, 4, '444444444444', 20,  '2026-06-10'), -- Cinnamon Delivery
    (3, 5, 5, '555555555555', 15,  '2027-12-15'); -- Vanilla Delivery



-- Insert mock data into supplier_materials
INSERT INTO supplier_materials (supplier_id, stock_material_id, created_at, updated_at)
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
    (4, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP); -- Chocolate


-- Insert mock data into supplier_prices
INSERT INTO supplier_prices (supplier_material_id, base_price)
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


INSERT INTO warehouse_stocks (warehouse_id, stock_material_id, quantity)
VALUES
    (1, 1, 50), -- Milk in Warehouse 1
    (1, 2, 30), -- Sugar in Warehouse 1
    (1, 3, 40), -- Chocolate in Warehouse 1
    (2, 4, 20), -- Cinnamon in Warehouse 2
    (2, 5, 15); -- Vanilla in Warehouse 2


-- Insert into StockRequests (Initial Requests)
INSERT INTO stock_requests (store_id, warehouse_id, status, request_date, created_at, updated_at)
VALUES
    (1, 1, 'CREATED', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 2, 'PROCESSED', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, 3, 'IN_DELIVERY', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (4, 4, 'COMPLETED', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO stock_request_ingredients (stock_request_id, ingredient_id, stock_material_id, quantity, created_at, updated_at)
VALUES
    -- StockRequest 1 (Store 1 -> Warehouse 1)
    (1, 1, 2, 10.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Sugar
    (1, 2, 1, 20.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Milk

    -- StockRequest 2 (Store 2 -> Warehouse 2)
    (2, 3, 3, 5.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),  -- Chocolate
    (2, 4, 4, 2.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),  -- Cinnamon

    -- StockRequest 3 (Store 3 -> Warehouse 3)
    (3, 5, 5, 1.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),  -- Vanilla
    (3, 1, 2, 15.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Sugar

    -- StockRequest 4 (Store 4 -> Warehouse 4)
    (4, 2, 1, 10.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Milk
    (4, 3, 3, 8.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);  -- Chocolate

