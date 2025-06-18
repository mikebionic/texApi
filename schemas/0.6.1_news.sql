-- News api db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "unaccent";

CREATE TYPE article_status AS ENUM ('draft', 'published', 'archived', 'deleted');
CREATE TYPE content_type AS ENUM ('article', 'breaking_news', 'opinion', 'analysis', 'interview', 'review');
CREATE TYPE priority_level AS ENUM ('low', 'medium', 'high', 'urgent');

CREATE TABLE tbl_article (
    -- Primary identification
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(255), -- For imported articles
    slug VARCHAR(500) UNIQUE NOT NULL,

    -- Content fields
    title VARCHAR(500) NOT NULL,
    subtitle VARCHAR(1000),
    excerpt TEXT, -- Short summary
    content TEXT NOT NULL, -- Main article content (rich text/HTML)
    content_plain TEXT, -- Plain text version for search

    -- Media and images
    featured_image_url VARCHAR(500),
    featured_image_alt VARCHAR(500),
    featured_image_caption TEXT,
    media_urls TEXT[], -- Array of additional media URLs
    media_captions TEXT[], -- Corresponding captions for media
    media_types TEXT[], -- 'image', 'video', 'audio', 'document'

    -- Author information (embedded)
    author_name VARCHAR(200) NOT NULL,
    author_email VARCHAR(255),
    author_bio TEXT,
    author_avatar_url VARCHAR(500),
    author_social_links TEXT[], -- Array of social media URLs
    author_is_verified BOOLEAN DEFAULT FALSE,

    -- Categorization
    category_primary VARCHAR(100) NOT NULL,
    category_secondary VARCHAR(100),
    categories TEXT[], -- Array of all categories
    tags TEXT[], -- Array of tags

    -- Source information
    source_name VARCHAR(200),
    source_domain VARCHAR(255),
    source_url VARCHAR(500),
    source_logo_url VARCHAR(500),
    external_url VARCHAR(500), -- Original source URL if syndicated

    -- Content metadata
    content_type content_type DEFAULT 'article',
    status article_status DEFAULT 'draft',
    priority priority_level DEFAULT 'medium',
    language_code VARCHAR(5) DEFAULT 'en', -- ISO language code
    country_code VARCHAR(2), -- ISO country code for local news

    -- Content metrics
    word_count INTEGER,
    reading_time_minutes INTEGER,

    -- Publishing and scheduling
    published_at TIMESTAMP WITH TIME ZONE,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- SEO fields
    meta_title VARCHAR(200),
    meta_description VARCHAR(500),
    meta_keywords TEXT[],
    canonical_url VARCHAR(500),

    -- Social media optimization
    social_title VARCHAR(200),
    social_description VARCHAR(500),
    social_image_url VARCHAR(500),

    -- Engagement metrics
    view_count BIGINT DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    dislike_count INTEGER DEFAULT 0,
    share_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,

    -- Article flags and settings
    is_featured BOOLEAN DEFAULT FALSE,
    is_breaking BOOLEAN DEFAULT FALSE,
    is_trending BOOLEAN DEFAULT FALSE,
    is_premium BOOLEAN DEFAULT FALSE,
    is_sponsored BOOLEAN DEFAULT FALSE,
    allow_comments BOOLEAN DEFAULT TRUE,

    -- Location data (for local news)
    location_city VARCHAR(100),
    location_state VARCHAR(100),
    location_country VARCHAR(100),
    location_coordinates POINT, -- PostgreSQL point type for lat/lng

    -- Comments (embedded as JSONB array)
    comments JSONB DEFAULT '[]'::jsonb, -- [{author, email, content, created_at, approved, replies: []}]

    -- Version history (embedded as JSONB array)
    versions JSONB DEFAULT '[]'::jsonb, -- [{version, title, content, changed_by, changed_at, summary}]

    -- View tracking (recent views as JSONB)
    recent_views JSONB DEFAULT '[]'::jsonb, -- [{ip, user_agent, viewed_at, country}] - keep last 100

    -- Newsletter/distribution
    newsletter_sent BOOLEAN DEFAULT FALSE,
    newsletter_sent_at TIMESTAMP WITH TIME ZONE,
    distribution_channels TEXT[], -- 'website', 'newsletter', 'social', 'rss'

    -- Analytics and performance
    bounce_rate DECIMAL(5,2),
    avg_time_on_page INTEGER, -- seconds
    search_keywords TEXT[], -- Keywords this article ranks for
    referrer_domains TEXT[], -- Top referrer domains

    -- Moderation and quality
    is_fact_checked BOOLEAN DEFAULT FALSE,
    fact_checker_name VARCHAR(200),
    fact_check_date TIMESTAMP WITH TIME ZONE,
    quality_score DECIMAL(3,2), -- 0.00 to 10.00
    editorial_notes TEXT,

    -- Related content
    related_article_ids UUID[],
    similar_tags TEXT[],
    recommended_reads TEXT[], -- Article titles or IDs

    -- Syndication and rights
    syndication_rights VARCHAR(50), -- 'exclusive', 'shared', 'public'
    copyright_holder VARCHAR(200),
    license_type VARCHAR(100), -- 'all_rights_reserved', 'creative_commons', etc.

    -- Custom fields (flexible storage)
    custom_fields JSONB DEFAULT '{}'::jsonb,

    meta TEXT,
    meta2 TEXT,
    meta3 TEXT,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    active           INT          NOT NULL DEFAULT 1,
    deleted          INT          NOT NULL DEFAULT 0
);

-- Create comprehensive indexes for performance
CREATE INDEX idx_tbl_article_status ON tbl_article(status);
CREATE INDEX idx_tbl_article_published_at ON tbl_article(published_at DESC);
CREATE INDEX idx_tbl_article_category ON tbl_article(category_primary);
CREATE INDEX idx_tbl_article_author ON tbl_article(author_name);
CREATE INDEX idx_tbl_article_featured ON tbl_article(is_featured) WHERE is_featured = TRUE;
CREATE INDEX idx_tbl_article_breaking ON tbl_article(is_breaking) WHERE is_breaking = TRUE;
CREATE INDEX idx_tbl_article_trending ON tbl_article(is_trending) WHERE is_trending = TRUE;
CREATE INDEX idx_tbl_article_slug ON tbl_article(slug);
CREATE INDEX idx_tbl_article_external_id ON tbl_article(external_id);

-- Full-text search index
-- Add column manually (not GENERATED)
ALTER TABLE tbl_article ADD COLUMN search_vector tsvector;

-- Create a trigger to update it
CREATE FUNCTION tbl_article_search_vector_trigger() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector(
            'english',
            coalesce(NEW.title, '') || ' ' ||
            coalesce(NEW.subtitle, '') || ' ' ||
            coalesce(NEW.content_plain, '') || ' ' ||
            coalesce(array_to_string(NEW.tags, ' '))
                         );
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
    ON tbl_article FOR EACH ROW EXECUTE FUNCTION tbl_article_search_vector_trigger();

-- Then index it:
CREATE INDEX idx_tbl_article_search_vector ON tbl_article USING GIN (search_vector);



-- Array indexes for tags and categories
CREATE INDEX idx_tbl_article_tags ON tbl_article USING gin(tags);
CREATE INDEX idx_tbl_article_categories ON tbl_article USING gin(categories);
CREATE INDEX idx_tbl_article_keywords ON tbl_article USING gin(meta_keywords);

-- JSONB indexes for flexible queries
CREATE INDEX idx_tbl_article_comments ON tbl_article USING gin(comments);
CREATE INDEX idx_tbl_article_custom_fields ON tbl_article USING gin(custom_fields);

-- Geographic index
CREATE INDEX idx_tbl_article_location ON tbl_article USING gist(location_coordinates);

-- Composite indexes for common queries
CREATE INDEX idx_tbl_article_status_published ON tbl_article(status, published_at DESC)
    WHERE status = 'published';
CREATE INDEX idx_tbl_article_category_status ON tbl_article(category_primary, status, published_at DESC);
CREATE INDEX idx_tbl_article_slug_pattern ON tbl_article (slug text_pattern_ops);



CREATE OR REPLACE FUNCTION calculate_trending_score(
    view_count BIGINT,
    like_count INTEGER,
    share_count INTEGER,
    comment_count INTEGER
)
    RETURNS NUMERIC
    IMMUTABLE
    LANGUAGE SQL
AS $$
SELECT view_count * 0.4 + like_count * 0.3 + share_count * 0.2 + comment_count * 0.1
$$;


-- Create trigger function for updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tbl_article_updated_at BEFORE UPDATE ON tbl_article
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to automatically calculate content metrics
CREATE OR REPLACE FUNCTION update_article_metrics()
    RETURNS TRIGGER AS $$
BEGIN
    -- Calculate word count from plain text content
    IF NEW.content_plain IS NOT NULL THEN
        NEW.word_count = array_length(
        string_to_array(
            regexp_replace(NEW.content_plain, '[^a-zA-Z0-9\s]', '', 'g'),
            ' '
            ),
            1
        );

        -- Estimate reading time (average 200 words per minute)
        NEW.reading_time_minutes = GREATEST(1, ROUND(NEW.word_count / 200.0));
    END IF;

    -- Update last_modified_at
    NEW.last_modified_at = NOW();

    -- Auto-generate slug if not provided
    IF NEW.slug IS NULL OR NEW.slug = '' THEN
        NEW.slug = lower(regexp_replace(NEW.title, '[^a-zA-Z0-9\s]', '', 'g'));
        NEW.slug = regexp_replace(NEW.slug, '\s+', '-', 'g');
        NEW.slug = substring(NEW.slug from 1 for 500);
    END IF;

    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_article_metrics_trigger BEFORE INSERT OR UPDATE ON tbl_article
    FOR EACH ROW EXECUTE FUNCTION update_article_metrics();


-- Helper function to add a comment (maintains JSONB array)
CREATE OR REPLACE FUNCTION add_comment(
    article_id UUID,
    author_name VARCHAR(100),
    author_email VARCHAR(255),
    comment_content TEXT,
    parent_comment_id UUID DEFAULT NULL
)
    RETURNS BOOLEAN AS $$
DECLARE
    new_comment JSONB;
BEGIN
    new_comment = jsonb_build_object(
            'id', gen_random_uuid(),
            'author', author_name,
            'email', author_email,
            'content', comment_content,
            'created_at', now(),
            'approved', false,
            'parent_id', parent_comment_id,
            'replies', '[]'::jsonb
                  );

    UPDATE tbl_article
    SET comments = comments || new_comment,
        comment_count = comment_count + 1
    WHERE id = article_id;

    RETURN FOUND;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION calculate_trending_score(
    view_count BIGINT,
    like_count INTEGER,
    share_count INTEGER,
    comment_count INTEGER
)
    RETURNS NUMERIC
    IMMUTABLE
AS $$
BEGIN
    RETURN view_count * 0.4 + like_count * 0.3 + share_count * 0.2 + comment_count * 0.1;
END;
$$ LANGUAGE plpgsql;


-- Helper function to add a version to history
CREATE OR REPLACE FUNCTION add_version(
    article_id UUID,
    version_title VARCHAR(500),
    version_content TEXT,
    changed_by VARCHAR(200),
    change_summary TEXT DEFAULT NULL
)
    RETURNS BOOLEAN AS $$
DECLARE
    new_version JSONB;
    version_num INTEGER;
BEGIN
    -- Get current version count
    SELECT jsonb_array_length(COALESCE(versions, '[]'::jsonb)) + 1
    INTO version_num
    FROM tbl_article
    WHERE id = article_id;

    new_version = jsonb_build_object(
    'version', version_num,
    'title', version_title,
    'content', version_content,
    'changed_by', changed_by,
    'changed_at', now(),
    'summary', change_summary
    );

    UPDATE tbl_article
    SET versions = COALESCE(versions, '[]'::jsonb) || new_version
    WHERE id = article_id;

    RETURN FOUND;
END;
$$ language 'plpgsql';

-- Views for common queries
CREATE VIEW published_articles AS
SELECT * FROM tbl_article
WHERE status = 'published'
  AND published_at <= NOW();

CREATE VIEW trending_articles AS
SELECT *,
       (view_count * 0.4 + like_count * 0.3 + share_count * 0.2 + comment_count * 0.1) as trend_score
FROM published_articles
WHERE published_at >= NOW() - INTERVAL '7 days'
ORDER BY trend_score DESC;

CREATE VIEW breaking_news AS
SELECT * FROM published_articles
WHERE is_breaking = TRUE
ORDER BY published_at DESC;

-- Sample insert with all the array and JSON features
/*
INSERT INTO tbl_article (
    title,
    slug,
    content,
    content_plain,
    author_name,
    author_email,
    category_primary,
    categories,
    tags,
    meta_keywords,
    status,
    featured_image_url,
    media_urls,
    media_types,
    media_captions,
    custom_fields
) VALUES (
    'Sample News Article Title',
    'sample-news-article-title',
    '<p>This is the <strong>HTML content</strong> of the article...</p>',
    'This is the plain text content of the article...',
    'John Doe',
    'john@example.com',
    'Technology',
    ARRAY['Technology', 'AI', 'Innovation'],
    ARRAY['artificial-intelligence', 'machine-learning', 'tech-news'],
    ARRAY['AI', 'technology', 'innovation', 'machine learning'],
    'published',
    'https://example.com/featured-image.jpg',
    ARRAY['https://example.com/image1.jpg', 'https://example.com/video1.mp4'],
    ARRAY['image', 'video'],
    ARRAY['Technology innovation', 'Demo video'],
    '{"sponsored_by": "TechCorp", "campaign_id": "tech2024"}'::jsonb
);
*/

INSERT INTO tbl_article (
    title,slug,subtitle,excerpt,content,content_plain,author_name,author_email,author_bio,author_avatar_url,author_social_links,
    author_is_verified,category_primary,category_secondary,categories,tags,source_name,source_domain,source_url,external_url,content_type,
    status,priority,language_code,country_code,featured_image_url,featured_image_alt,featured_image_caption,media_urls,
    media_types,media_captions,word_count,reading_time_minutes,published_at,meta_title,meta_description,meta_keywords,social_title,
    social_description,social_image_url,view_count,like_count,dislike_count,share_count,comment_count,is_featured,is_trending,allow_comments,location_city,
    location_state,location_country,location_coordinates,newsletter_sent,newsletter_sent_at,distribution_channels,bounce_rate,avg_time_on_page,search_keywords,
    referrer_domains,is_fact_checked,fact_checker_name,fact_check_date,quality_score,editorial_notes,syndication_rights,copyright_holder,license_type,custom_fields
)
VALUES (
   'Amazon Tests Drone Delivery in Kazakhstan to Cut Delivery Time by 60%',
   'amazon-drone-delivery-kazakhstan',
   'A major logistics shift is coming to Central Asia',
   'Amazon begins pilot drone delivery program in Almaty, aiming for faster rural delivery.',
   '<p>Amazon has started testing drone deliveries in Almaty, Kazakhstan, partnering with local logistics providers to speed up last-mile delivery. The pilot is expected to reduce delivery time in rural areas by over 60%.</p>',
   'Amazon has started testing drone deliveries in Almaty, Kazakhstan, partnering with local logistics providers to speed up last-mile delivery. The pilot is expected to reduce delivery time in rural areas by over 60%.',
   'Leyla Akhmetova',
   'leyla.akhmetova@logisticsnews.kz',
   'Senior logistics journalist with 12+ years in Central Asia transport reporting.',
   'https://example.com/images/leyla.jpg',
   ARRAY['https://twitter.com/leylaakh', 'https://linkedin.com/in/leylaakh'],
   TRUE,
   'Logistics',
   'Technology',
   ARRAY['Logistics', 'Technology', 'Kazakhstan'],
   ARRAY['amazon', 'drone-delivery', 'last-mile', 'central-asia', 'innovation'],
   'Reuters Logistics',
   'reuters.com',
   'https://www.reuters.com/amazon-drone-delivery-kazakhstan',
   'https://www.reuters.com/amazon-drone-delivery-kazakhstan',
   'breaking_news',
   'published',
   'high',
   'en',
   'KZ',
   'https://example.com/images/amazon-drone-kz.jpg',
   'Amazon Drone flying over Almaty suburbs',
   'An Amazon Prime Air drone conducts trial delivery in Kazakhstan',
   ARRAY[
       'https://example.com/images/warehouse.jpg',
       'https://example.com/videos/drone-footage.mp4'
       ],
   ARRAY['image', 'video'],
   ARRAY['Kazakhstan fulfillment warehouse', 'Drone in flight over mountains'],
   125,
   1,
   NOW() - INTERVAL '2 hours',
   'Amazon Launches Drone Delivery in Kazakhstan',
   'Amazon launches drone logistics test in Kazakhstan, targeting rural speed improvements.',
   ARRAY['drone', 'amazon', 'logistics', 'kazakhstan', 'delivery innovation'],
   'Amazon Drone Program Expands',
   'Amazon is now using drones for delivery in Central Asia',
   'https://example.com/social/amazon-drone.jpg',
   32890,
   1230,
   45,
   987,
   32,
   TRUE,
   TRUE,
   TRUE,
   'Almaty',
   'Almaty',
   'Kazakhstan',
   POINT(76.886, 43.238),
   TRUE,
   NOW() - INTERVAL '1 hour',
   ARRAY['website', 'newsletter', 'rss'],
   37.4,
   234,
   ARRAY['amazon drones', 'kazakhstan delivery', 'logistics tech'],
   ARRAY['google.com', 'linkedin.com'],
   TRUE,
   'Olzhas Isayev',
   NOW() - INTERVAL '30 minutes',
   8.7,
   'Verified by local logistics analyst and fact-checked with drone license records.',
   'shared',
   'Reuters Logistics',
   'creative_commons',
   '{"region":"Central Asia", "impact":"rural delivery"}'::jsonb
);
