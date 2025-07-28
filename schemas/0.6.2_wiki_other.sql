CREATE TABLE tbl_wiki
(
    id                          SERIAL PRIMARY KEY,
    uuid                        UUID          NOT NULL DEFAULT gen_random_uuid(),

    title_en                    VARCHAR(500)  NOT NULL DEFAULT '',
    title_ru                    VARCHAR(500)  NOT NULL DEFAULT '',
    title_tk                    VARCHAR(500)  NOT NULL DEFAULT '',

    description_en              TEXT,
    description_ru              TEXT,
    description_tk              TEXT,
    description_type     VARCHAR(50)   NOT NULL DEFAULT 'plain', -- 'plain', 'html', 'info'

    text_md_en                  TEXT,         -- Markdown content in English
    text_md_ru                  TEXT,         -- Markdown content in Russian
    text_md_tk                  TEXT,         -- Markdown content in Turkmen

    text_rich_en                 TEXT,         -- Rich text/HTML content in English
    text_rich_ru                 TEXT,         -- Rich text/HTML content in Russian
    text_rich_tk                 TEXT,         -- Rich text/HTML content in Turkmen

    file_url_1                  VARCHAR(1000),
    file_url_2                  VARCHAR(1000),
    file_url_3                  VARCHAR(1000),
    file_url_4                  VARCHAR(1000),
    file_url_5                  VARCHAR(1000),

    file_info_1_en              VARCHAR(500), -- File description in English
    file_info_1_ru              VARCHAR(500), -- File description in Russian
    file_info_1_tk              VARCHAR(500), -- File description in Turkmen
    file_info_2_en              VARCHAR(500),
    file_info_2_ru              VARCHAR(500),
    file_info_2_tk              VARCHAR(500),
    file_info_3_en              VARCHAR(500),
    file_info_3_ru              VARCHAR(500),
    file_info_3_tk              VARCHAR(500),
    file_info_4_en              VARCHAR(500),
    file_info_4_ru              VARCHAR(500),
    file_info_4_tk              VARCHAR(500),
    file_info_5_en              VARCHAR(500),
    file_info_5_ru              VARCHAR(500),
    file_info_5_tk              VARCHAR(500),

    video_url_1                 VARCHAR(1000),
    video_url_2                 VARCHAR(1000),
    video_url_3                 VARCHAR(1000),
    video_url_4                 VARCHAR(1000),
    video_url_5                 VARCHAR(1000),

    video_info_1_en             VARCHAR(500), -- Video description in English
    video_info_1_ru             VARCHAR(500), -- Video description in Russian
    video_info_1_tk             VARCHAR(500), -- Video description in Turkmen
    video_info_2_en             VARCHAR(500),
    video_info_2_ru             VARCHAR(500),
    video_info_2_tk             VARCHAR(500),
    video_info_3_en             VARCHAR(500),
    video_info_3_ru             VARCHAR(500),
    video_info_3_tk             VARCHAR(500),
    video_info_4_en             VARCHAR(500),
    video_info_4_ru             VARCHAR(500),
    video_info_4_tk             VARCHAR(500),
    video_info_5_en             VARCHAR(500),
    video_info_5_ru             VARCHAR(500),
    video_info_5_tk             VARCHAR(500),

    category                    VARCHAR(50)   NOT NULL DEFAULT 'docs', -- 'docs', 'wiki', 'guides', 'tutorials', 'api', 'faq'
    subcategory                 VARCHAR(50),  -- Additional categorization
    tags                        VARCHAR(500), -- Comma-separated tags for search
    version                     INT           NOT NULL DEFAULT 1,

    slug                        VARCHAR(200), -- URL-friendly identifier
    meta_keywords_en            VARCHAR(500), -- SEO keywords in English
    meta_keywords_ru            VARCHAR(500), -- SEO keywords in Russian
    meta_keywords_tk            VARCHAR(500), -- SEO keywords in Turkmen
    priority                    INT           NOT NULL DEFAULT 0, -- Display priority (higher = more important)
    view_count                  INT           NOT NULL DEFAULT 0, -- Track popularity
    is_featured                 BOOLEAN       NOT NULL DEFAULT FALSE, -- Featured content
    is_public                   BOOLEAN       NOT NULL DEFAULT TRUE,  -- Public visibility
    requires_auth               BOOLEAN       NOT NULL DEFAULT FALSE, -- Authentication required
    content_type                VARCHAR(20)   NOT NULL DEFAULT 'article', -- 'article', 'tutorial', 'reference', 'guide'
    difficulty_level            VARCHAR(20)   DEFAULT 'beginner', -- 'beginner', 'intermediate', 'advanced'
    estimated_read_time         INT,          -- in minutes

    created_at                  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active                      INT           NOT NULL DEFAULT 1,
    deleted                     INT           NOT NULL DEFAULT 0,

    CONSTRAINT unique_slug UNIQUE (slug, deleted),
    CONSTRAINT valid_priority CHECK (priority >= 0),
    CONSTRAINT valid_version CHECK (version >= 1),
    CONSTRAINT valid_view_count CHECK (view_count >= 0),
    CONSTRAINT valid_read_time CHECK (estimated_read_time >= 0)
--     CONSTRAINT valid_category CHECK (category IN ('docs', 'wiki', 'guides', 'tutorials', 'api', 'faq', 'changelog', 'troubleshooting')),
--     CONSTRAINT valid_content_type CHECK (content_type IN ('article', 'tutorial', 'reference', 'guide', 'faq', 'changelog')),
--     CONSTRAINT valid_difficulty CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced'))
);

-- Indexes for better performance
CREATE INDEX idx_wiki_category_active ON tbl_wiki(category, active, deleted);
CREATE INDEX idx_wiki_subcategory ON tbl_wiki(subcategory);
CREATE INDEX idx_wiki_content_type ON tbl_wiki(content_type);
CREATE INDEX idx_wiki_priority ON tbl_wiki(priority DESC);
CREATE INDEX idx_wiki_uuid ON tbl_wiki(uuid);
CREATE INDEX idx_wiki_slug ON tbl_wiki(slug);
CREATE INDEX idx_wiki_tags ON tbl_wiki USING gin (to_tsvector('english', tags));
CREATE INDEX idx_wiki_title_search ON tbl_wiki USING gin (
                                                          to_tsvector('english', coalesce(title_en, '') || ' ' || coalesce(title_ru, '') || ' ' || coalesce(title_tk, ''))
    );
CREATE INDEX idx_wiki_content_search ON tbl_wiki USING gin (
                                                            to_tsvector('english', coalesce(text_md_en, '') || ' ' || coalesce(description_en, ''))
    );
CREATE INDEX idx_wiki_featured_public ON tbl_wiki(is_featured, is_public, active, deleted);

CREATE OR REPLACE FUNCTION update_wiki_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_wiki_updated_at
    BEFORE UPDATE ON tbl_wiki
    FOR EACH ROW
EXECUTE FUNCTION update_wiki_updated_at();


INSERT INTO tbl_wiki (
    uuid,
    title_en, title_ru, title_tk,
    description_en, description_ru, description_tk,
    description_type,
    text_md_en, text_md_ru, text_md_tk,
    text_rich_en, text_rich_ru, text_rich_tk,
    file_url_1,
    file_info_1_en, file_info_1_ru, file_info_1_tk,
    video_url_1,
    video_info_1_en, video_info_1_ru, video_info_1_tk,
    category, tags, version, slug, meta_keywords_en,
    meta_keywords_ru, meta_keywords_tk, priority, view_count, is_featured,
    is_public, requires_auth, content_type, difficulty_level, estimated_read_time
)
VALUES (
       gen_random_uuid(),
       'How to Use TEX Express Web and Mobile Application', -- English Title
       'Как использовать веб- и мобильное приложение TEX Express', -- Russian Title
       'TEX Express web we mobil programmasyny nähili ulanmalydygyny', -- Turkmen Title
       'TEX Express is a revolutionary logistics platform that allows users to create offers for cargo or transport, track shipments, and more. Learn how to navigate the platform and create offers using the web and mobile applications.',
       'TEX Express - это революционная логистическая платформа, которая позволяет пользователям создавать предложения для грузов или транспорта, отслеживать отправления и многое другое. Узнайте, как навигировать по платформе и создавать предложения с помощью веб- и мобильных приложений.',
       'TEX Express, bu ulanyjylara gurnama ýa-da ulag üçin teklip döretmäge, ýükleri yzarlamaga we başga-da köp zada mümkinçilik berýän inqilaby logistik platformadyr. Web we mobil programma arkaly platformany nädip ulanmalydygyny we teklipleri nädip döretmeli bilen tanyş boluň.',
       'plain', -- Description type
       '### Getting Started\nTo start using TEX Express, you need to sign up and log in.\n\n**Web Application**\n1. Visit the official website: [TEX Express](https://texexpress.pro/)\n2. Click on the Sign Up button at the top right of the page.\n\n### Creating Offers\nAs a **shipper**, you can create an offer to find carriers for your cargo.\n\n1. Log in to your account.\n2. Go to the Find Cargo section.\n3. Click on Create Offer.', -- Markdown content in English
       '### Начало работы\nЧтобы начать использовать TEX Express, вам нужно зарегистрироваться и войти в систему.\n\n**Веб-приложение**\n1. Перейдите на официальный сайт: [TEX Express](https://texexpress.pro/)\n2. Нажмите на кнопку Регистрация в верхней правой части страницы.\n\n### Создание предложений\nКак **грузоотправитель**, вы можете создать предложение для поиска перевозчиков для вашего груза.\n\n1. Войдите в свой аккаунт.\n2. Перейдите в раздел Найти груз.\n3. Нажмите Создать предложение.', -- Markdown content in Russian
       '### Başlangyç\nTEX Express-i ulanmak üçin hasabyňyza girmeli we hasapdan geçmeli.\n\n**Web Programmasy**\n1. Resmi web sahypasyna giriň: [TEX Express](https://texexpress.pro/)\n2. Sahypanyň sag üst bölegindäki "Hasaba alyş" düwmesini basyň.\n\n### Teklip döretmek\n**Ýükdaşary** hökmünde, ýüküňiz üçin taşujylary tapmak üçin teklip döretmek mümkinçiligine eýe bolarsyňyz.\n\n1. Hasabyňyza giriň.\n2. "Ýük Tapmak" bölümini açyň.\n3. "Teklip Döretmek" düwmesini basyň.', -- Markdown content in Turkmen
       '#### How TEX Express Works\nTEX Express optimizes your logistics operations and connects you with carriers or shippers. You can easily create offers, track shipments, and manage logistics operations from the platform.',
       '#### Как работает TEX Express\nTEX Express оптимизирует ваши логистические операции и связывает вас с перевозчиками или грузоотправителями. Вы можете легко создавать предложения, отслеживать отправления и управлять логистическими операциями прямо с платформы.',
       '#### TEX Express nähili işleýär\nTEX Express logistika amallaryňyzy optimizirleýär we sizi taşujylary ýa-da ýükdaşarylar bilen baglanyşdyrýar. Platformadan aňsatlyk bilen teklipler döretmek, ýükleriňizden yzarlama geçirmek we logistika amallaryny dolandyrmak mümkin.',
       'https://www.itf-oecd.org/sites/default/files/docs/02logisticse.pdf',
    'Transport Logistics - SHARED SOLUTIONS TO COMMON CHALLENGES', 'Транспортная логистика', NULL, 'https://office.belentlik.tm/index.php/s/HkA25rx6cZmGDEw/download/logistics_AD.mp4', 'TEX Express showcase', 'TEX Express демонстрация', NULL, 'How to use', -- No files or videos for this example
       'logistics, transport, cargo, shipping, TEX Express, platform', -- Tags
       1, -- Version
       'how-to-use-tex-express-web-and-mobile-application', -- Slug
       'TEX Express, logistics, cargo, transport, mobile application, web application', -- SEO Keywords in English
       'TEX Express, логистика, груз, транспорт, мобильное приложение, веб-приложение', -- SEO Keywords in Russian
       'TEX Express, logistika, ýük, ulag, mobil programma, web programma', -- SEO Keywords in Turkmen
       0,
       0,
       FALSE,
       TRUE,
       FALSE,
       'article',
       'beginner',
       15
    );

