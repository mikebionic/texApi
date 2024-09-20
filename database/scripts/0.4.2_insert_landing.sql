-- Insert languages
INSERT INTO
    languages (code, name)
VALUES
    ('en', 'English'),
    ('ru', 'Russian');

-- Insert content types
INSERT INTO
    content_types (type_name, title, description)
VALUES
    (
        'about',
        'About Us',
        'Information about our company and who we are.'
    ),
    (
        'partner',
        'Our Partners',
        'Details about our trusted business partners.'
    ),
    (
        'achievement',
        'Our Achievements',
        'Notable awards and recognitions we have received.'
    ),
    (
        'faq',
        'Frequently Asked Questions',
        'Answers to common questions.'
    ),
    (
        'video',
        'Videos',
        'A collection of videos showcasing our operations and client feedback.'
    ),
    (
        'mission',
        'Our Mission',
        'Our core values and goals.'
    ),
    (
        'contact',
        'Contact Information',
        'Ways to get in touch with us.'
    ),
    (
        'how_we_work',
        'How We Work',
        'An overview of our process and methodology.'
    );

-- Insert content
INSERT INTO
    content (
        lang_id,
        content_type_id,
        title,
        subtitle,
        description,
        image_url,
        step
    )
VALUES
    -- About Us
    (
        1,
        1,
        'About Us',
        'Who We Are',
        'We are a leading logistics provider specializing in supply chain management, ensuring timely deliveries and exceptional service.',
        'about_en.jpg',
        1
    ),
    (
        2,
        1,
        'О нас',
        'Кто мы',
        'Мы ведущий поставщик логистических услуг, специализирующийся на управлении цепочками поставок и обеспечивающий своевременную доставку.',
        'about_ru.jpg',
        1
    ),
    -- Partners
    (
        1,
        2,
        'Partner A',
        'Our Trusted Partner',
        'A reliable logistics partner with years of experience in the industry.',
        'partner_a_en.jpg',
        2
    ),
    (
        2,
        2,
        'Партнёр A',
        'Наш надёжный партнёр',
        'Надёжный логистический партнёр с многолетним опытом работы в отрасли.',
        'partner_a_ru.jpg',
        2
    ),
    (
        1,
        2,
        'Partner B',
        'Leading Innovation',
        'Partner B is at the forefront of logistics technology, offering cutting-edge solutions.',
        'partner_b_en.jpg',
        3
    ),
    (
        2,
        2,
        'Партнёр B',
        'Лидеры инноваций',
        'Партнёр B находится на переднем крае логистических технологий, предлагая современные решения.',
        'partner_b_ru.jpg',
        3
    ),
    (
        1,
        2,
        'Partner C',
        'Global Reach',
        'Partner C operates across multiple continents, providing international logistics solutions.',
        'partner_c_en.jpg',
        4
    ),
    (
        2,
        2,
        'Партнёр C',
        'Глобальный охват',
        'Партнёр C работает на нескольких континентах, предоставляя международные логистические решения.',
        'partner_c_ru.jpg',
        4
    ),
    -- Achievements
    (
        1,
        3,
        'Best Logistics Provider 2023',
        'Our Achievements',
        'Awarded for excellence in logistics and customer service.',
        'achievement_2023_en.jpg',
        5
    ),
    (
        2,
        3,
        'Лучший поставщик логистических услуг 2023',
        'Наши достижения',
        'Награждён за выдающиеся результаты в логистике и обслуживании клиентов.',
        'achievement_2023_ru.jpg',
        5
    ),
    (
        1,
        3,
        'TEX Logistics Award',
        'Sustainability Efforts',
        'Recognized for our commitment to sustainable logistics practices.',
        'tex_award_en.jpg',
        6
    ),
    (
        2,
        3,
        'ТЕКС Логистика',
        'Устойчивое развитие',
        'Признаны за приверженность устойчивым логистическим практикам.',
        'tex_award_ru.jpg',
        6
    ),
    (
        1,
        3,
        'Fastest Growing Company 2022',
        'Growth Recognition',
        'Recognized as one of the fastest-growing logistics companies in the region.',
        'fastest_growth_en.jpg',
        7
    ),
    (
        2,
        3,
        'Самая быстрорастущая компания 2022',
        'Признание роста',
        'Признаны одной из самых быстрорастущих логистических компаний в регионе.',
        'fastest_growth_ru.jpg',
        7
    ),
    -- Video section
    (
        1,
        5,
        'How We Operate',
        'Our Process Overview',
        'A brief overview of our logistics process, highlighting our efficiency and dedication.',
        'video_overview_en.mp4',
        8
    ),
    (
        2,
        5,
        'Как мы работаем',
        'Обзор нашего процесса',
        'Краткий обзор нашего логистического процесса, подчеркивающий нашу эффективность и преданность делу.',
        'video_overview_ru.mp4',
        8
    ),
    (
        1,
        5,
        'Client Testimonials',
        'Hear From Our Clients',
        'Watch our clients share their experiences with our logistics solutions.',
        'client_testimonials_en.mp4',
        9
    ),
    (
        2,
        5,
        'Отзывы клиентов',
        'Услышьте от наших клиентов',
        'Смотрите, как наши клиенты делятся своим опытом работы с нашими логистическими решениями.',
        'client_testimonials_ru.mp4',
        9
    ),
    -- FAQs
    (
        1,
        4,
        'What services do you offer?',
        'Frequently Asked Questions',
        'We offer a range of logistics services including transportation, warehousing, and supply chain management.',
        '',
        10
    ),
    (
        2,
        4,
        'Какие услуги вы предлагаете?',
        'Часто задаваемые вопросы',
        'Мы предлагаем ряд логистических услуг, включая транспортировку, складирование и управление цепочками поставок.',
        '',
        10
    ),
    (
        1,
        4,
        'How can I track my shipment?',
        'Tracking Information',
        'You can track your shipment using our online portal or contact our support team.',
        '',
        11
    ),
    (
        2,
        4,
        'Как я могу отслеживать свою посылку?',
        'Информация для отслеживания',
        'Вы можете отслеживать свою посылку, используя наш онлайн-портал или связавшись с нашей службой поддержки.',
        '',
        11
    ),
    (
        1,
        4,
        'What is your return policy?',
        'Return Policy',
        'We have a flexible return policy. Contact us for more details.',
        '',
        12
    ),
    (
        2,
        4,
        'Какова ваша политика возврата?',
        'Политика возврата',
        'У нас гибкая политика возврата. Свяжитесь с нами для получения дополнительной информации.',
        '',
        12
    ),
    -- Our Mission
    (
        1,
        6,
        'Our Mission',
        'What Drives Us',
        'Our mission is to provide efficient and reliable logistics solutions to our clients.',
        '',
        13
    ),
    (
        2,
        6,
        'Наша миссия',
        'Что нас движет',
        'Наша миссия — предоставлять эффективные и надёжные логистические решения нашим клиентам.',
        '',
        13
    ),
    -- Contact Information
    (
        1,
        7,
        'Contact Us',
        'Get in Touch',
        'Feel free to reach out to us via phone or email for any inquiries.',
        '',
        14
    ),
    (
        2,
        7,
        'Связаться с нами',
        'Свяжитесь с нами',
        'Не стесняйтесь обращаться к нам по телефону или электронной почте с любыми вопросами.',
        '',
        14
    ),
    (
        1,
        7,
        'Office Locations',
        'Where to Find Us',
        'We have offices located in key cities around the globe.',
        '',
        15
    ),
    (
        2,
        7,
        'Офисы',
        'Где нас найти',
        'У нас есть офисы в ключевых городах по всему миру.',
        '',
        15
    );