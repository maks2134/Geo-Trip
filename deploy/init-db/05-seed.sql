-- 05-seed.sql: Initial demo dataset for full walkthrough

-- 1. Demo User: traveler@atlasflow.io (password: nomad12345)
INSERT INTO users (id, email, password_hash, display_name, created_at, updated_at)
VALUES (
    'usr-demo-1',
    'traveler@atlasflow.io',
    '$2a$10$w85J5sK.GjRfgJqL9bZc..F1jR6hO7gP7kM9.h8V9m1Q9N2v5p2G2',
    'Alex Nomad',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- 2. Visited Countries for World Scratch Map
INSERT INTO visited_countries (user_id, country_iso2, visited_at) VALUES
('usr-demo-1', 'FR', NOW() - INTERVAL '60 days'),
('usr-demo-1', 'JP', NOW() - INTERVAL '30 days'),
('usr-demo-1', 'IT', NOW() - INTERVAL '120 days'),
('usr-demo-1', 'ES', NOW() - INTERVAL '90 days')
ON CONFLICT DO NOTHING;

-- 3. Geo Places
INSERT INTO places (id, owner_id, name, category, status, country_iso2, city, address, geom, created_at, updated_at) VALUES
('00000000-0000-0000-0000-000000000001', 'usr-demo-1', 'Senso-ji Temple', 'sight', 'visited', 'JP', 'Tokyo', '2 Chome-3-1 Asakusa, Taito City', ST_SetSRID(ST_MakePoint(139.7967, 35.7148), 4326), NOW(), NOW()),
('00000000-0000-0000-0000-000000000002', 'usr-demo-1', 'Shibuya Sky', 'sight', 'visited', 'JP', 'Tokyo', '2 Chome-24-12 Shibuya', ST_SetSRID(ST_MakePoint(139.7013, 35.6595), 4326), NOW(), NOW()),
('00000000-0000-0000-0000-000000000003', 'usr-demo-1', 'Tsukiji Outer Market', 'food', 'visited', 'JP', 'Tokyo', '4 Chome-16-2 Tsukiji, Chuo City', ST_SetSRID(ST_MakePoint(139.7708, 35.6655), 4326), NOW(), NOW()),
('00000000-0000-0000-0000-000000000004', 'usr-demo-1', 'Shinjuku Gyoen National Garden', 'nature', 'planned', 'JP', 'Tokyo', '11 Naitomachi, Shinjuku City', ST_SetSRID(ST_MakePoint(139.7101, 35.6852), 4326), NOW(), NOW()),
('00000000-0000-0000-0000-000000000005', 'usr-demo-1', 'Park Hyatt Tokyo', 'stay', 'planned', 'JP', 'Tokyo', '3 Chome-7-1-2 Nishishinjuku', ST_SetSRID(ST_MakePoint(139.6914, 35.6856), 4326), NOW(), NOW()),
('00000000-0000-0000-0000-000000000006', 'usr-demo-1', 'Eiffel Tower', 'sight', 'visited', 'FR', 'Paris', 'Champ de Mars, 5 Av. Anatole France', ST_SetSRID(ST_MakePoint(2.2945, 48.8584), 4326), NOW(), NOW()),
('00000000-0000-0000-0000-000000000007', 'usr-demo-1', 'Louvre Museum', 'sight', 'visited', 'FR', 'Paris', '75001 Paris', ST_SetSRID(ST_MakePoint(2.3376, 48.8606), 4326), NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 4. Trip Core: Tokyo Trip
INSERT INTO trips (id, owner_id, title, description, start_date, end_date, status, cover_url, currency, budget_limit, share_slug, is_public, created_at, updated_at)
VALUES (
    'trip-tokyo-2026',
    'usr-demo-1',
    'Tokyo Sakura Adventure',
    'A 7-day culinary, cultural, and photography expedition during the cherry blossom season.',
    '2026-04-01 00:00:00+00',
    '2026-04-07 23:59:59+00',
    'active',
    'https://images.unsplash.com/photo-1503899036084-c55cdd92da26?auto=format&fit=crop&w=1200&q=80',
    'USD',
    3500.00,
    'tokyo-2026-demo',
    TRUE,
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- Public Share Record
INSERT INTO trip_shares (id, trip_id, slug, hide_budget, hide_private_docs, created_at)
VALUES ('share-tokyo-demo', 'trip-tokyo-2026', 'tokyo-2026-demo', TRUE, TRUE, NOW())
ON CONFLICT (slug) DO NOTHING;

-- Days for Tokyo Trip
INSERT INTO trip_days (id, trip_id, day_index, date, title, notes, created_at) VALUES
('day-tokyo-1', 'trip-tokyo-2026', 1, '2026-04-01 00:00:00+00', 'Historic Asakusa & Night in Shibuya', 'Arrive at Haneda, train to hotel, evening street walk.', NOW()),
('day-tokyo-2', 'trip-tokyo-2026', 2, '2026-04-02 00:00:00+00', 'Seafood Market & Green Parks', 'Early morning sushi, relax under blooming cherry trees in Shinjuku.', NOW())
ON CONFLICT (id) DO NOTHING;

-- Itinerary items with Lexorank string positions
INSERT INTO itinerary_items (id, trip_id, day_id, place_id, place_name, latitude, longitude, category, address, position, duration_minutes, cost, notes, created_at, updated_at) VALUES
('it-1', 'trip-tokyo-2026', 'day-tokyo-1', '00000000-0000-0000-0000-000000000001', 'Senso-ji Temple', 35.7148, 139.7967, 3, '2 Chome-3-1 Asakusa', 'a0', 90, 0, 'Morning photo walk before the crowds arrive.', NOW(), NOW()),
('it-2', 'trip-tokyo-2026', 'day-tokyo-1', '00000000-0000-0000-0000-000000000002', 'Shibuya Sky', 35.6595, 139.7013, 3, '2 Chome-24-12 Shibuya', 'a5', 60, 25.00, 'Sunset reservation tickets booked online.', NOW(), NOW()),
('it-3', 'trip-tokyo-2026', 'day-tokyo-2', '00000000-0000-0000-0000-000000000003', 'Tsukiji Outer Market', 35.6655, 139.7708, 1, '4 Chome-16-2 Tsukiji', 'b0', 75, 45.00, 'Fresh tamagoyaki and sashimi bowls.', NOW(), NOW()),
('it-4', 'trip-tokyo-2026', 'day-tokyo-2', '00000000-0000-0000-0000-000000000004', 'Shinjuku Gyoen Garden', 35.6852, 139.7101, 4, '11 Naitomachi', 'b5', 120, 5.00, 'Hanami picnic with matcha lattes.', NOW(), NOW()),
('it-5', 'trip-tokyo-2026', NULL, '00000000-0000-0000-0000-000000000005', 'Park Hyatt Tokyo', 35.6856, 139.6914, 2, '3 Chome-7-1-2 Nishishinjuku', 'c0', 60, 420.00, 'Backlog wishlist hotel.', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Expenses
INSERT INTO expenses (id, trip_id, day_id, place_item_id, category, amount, currency, title, notes, created_at) VALUES
('exp-1', 'trip-tokyo-2026', 'day-tokyo-1', 'it-2', 4, 50.00, 'USD', 'Shibuya Sky Observation Deck', '2 tickets at sunset', NOW()),
('exp-2', 'trip-tokyo-2026', 'day-tokyo-2', 'it-3', 2, 90.00, 'USD', 'Tsukiji Seafood Breakfast', 'Sashimi bowl + grilled scallops', NOW()),
('exp-3', 'trip-tokyo-2026', NULL, NULL, 1, 1260.00, 'USD', 'Hotel Shinjuku 3 Nights', 'Prepaid booking', NOW()),
('exp-4', 'trip-tokyo-2026', NULL, NULL, 3, 140.00, 'USD', 'Suica & Narita Express Train', 'Transportation cards', NOW())
ON CONFLICT (id) DO NOTHING;

-- Checklists
INSERT INTO checklist_items (id, trip_id, section, title, is_done, position, created_at, updated_at) VALUES
('chk-1', 'trip-tokyo-2026', 1, 'Check passport validity (>6 months)', TRUE, '10', NOW(), NOW()),
('chk-2', 'trip-tokyo-2026', 1, 'Complete Visit Japan Web QR code', TRUE, '20', NOW(), NOW()),
('chk-3', 'trip-tokyo-2026', 1, 'Notify bank about travel in Japan', FALSE, '30', NOW(), NOW()),
('chk-4', 'trip-tokyo-2026', 2, 'Pack universal power adapter (Type A)', TRUE, '10', NOW(), NOW()),
('chk-5', 'trip-tokyo-2026', 2, 'Comfortable walking sneakers', TRUE, '20', NOW(), NOW()),
('chk-6', 'trip-tokyo-2026', 2, 'Camera + extra SD cards', FALSE, '30', NOW(), NOW()),
('chk-7', 'trip-tokyo-2026', 3, 'Purchase e-SIM at Haneda airport', FALSE, '10', NOW(), NOW()),
('chk-8', 'trip-tokyo-2026', 3, 'Withdraw yen cash at 7-Eleven ATM', FALSE, '20', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
