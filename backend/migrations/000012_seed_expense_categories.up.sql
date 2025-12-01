-- Seed default expense categories and subcategories
-- This migration creates a comprehensive list of expense categories for a professional services office

-- 1. Alquiler y Arrendamiento
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('a1111111-1111-1111-1111-111111111111', 'Alquiler y Arrendamiento', 'Gastos relacionados con alquiler de espacios y equipos', NULL, true, 1, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('Alquiler de local', NULL, 'a1111111-1111-1111-1111-111111111111', true, 1, NOW(), NOW()),
('Alquiler de equipo', NULL, 'a1111111-1111-1111-1111-111111111111', true, 2, NOW(), NOW()),
('Gastos de comunidad', NULL, 'a1111111-1111-1111-1111-111111111111', true, 3, NOW(), NOW());

-- 2. Suministros
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('a2222222-2222-2222-2222-222222222222', 'Suministros', 'Servicios básicos del local', NULL, true, 2, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('Electricidad', NULL, 'a2222222-2222-2222-2222-222222222222', true, 1, NOW(), NOW()),
('Agua', NULL, 'a2222222-2222-2222-2222-222222222222', true, 2, NOW(), NOW()),
('Gas', NULL, 'a2222222-2222-2222-2222-222222222222', true, 3, NOW(), NOW()),
('Internet y teléfono', NULL, 'a2222222-2222-2222-2222-222222222222', true, 4, NOW(), NOW());

-- 3. Material y Consumibles
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('a3333333-3333-3333-3333-333333333333', 'Material y Consumibles', 'Material necesario para la actividad', NULL, true, 3, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('Material clínico', NULL, 'a3333333-3333-3333-3333-333333333333', true, 1, NOW(), NOW()),
('Material de oficina', NULL, 'a3333333-3333-3333-3333-333333333333', true, 2, NOW(), NOW()),
('Productos de limpieza', NULL, 'a3333333-3333-3333-3333-333333333333', true, 3, NOW(), NOW()),
('Material desechable', NULL, 'a3333333-3333-3333-3333-333333333333', true, 4, NOW(), NOW());

-- 4. Personal
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('a4444444-4444-4444-4444-444444444444', 'Personal', 'Gastos relacionados con el personal', NULL, true, 4, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('Salarios y nóminas', NULL, 'a4444444-4444-4444-4444-444444444444', true, 1, NOW(), NOW()),
('Seguridad social', NULL, 'a4444444-4444-4444-4444-444444444444', true, 2, NOW(), NOW()),
('Formación', NULL, 'a4444444-4444-4444-4444-444444444444', true, 3, NOW(), NOW()),
('Dietas y gastos de viaje', NULL, 'a4444444-4444-4444-4444-444444444444', true, 4, NOW(), NOW());

-- 5. Servicios Profesionales
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('a5555555-5555-5555-5555-555555555555', 'Servicios Profesionales', 'Asesoramiento y servicios externos', NULL, true, 5, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('Asesoría fiscal', NULL, 'a5555555-5555-5555-5555-555555555555', true, 1, NOW(), NOW()),
('Asesoría laboral', NULL, 'a5555555-5555-5555-5555-555555555555', true, 2, NOW(), NOW()),
('Servicios jurídicos', NULL, 'a5555555-5555-5555-5555-555555555555', true, 3, NOW(), NOW()),
('Consultoría', NULL, 'a5555555-5555-5555-5555-555555555555', true, 4, NOW(), NOW());

-- 6. Seguros
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('a6666666-6666-6666-6666-666666666666', 'Seguros', 'Pólizas de seguro', NULL, true, 6, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('Seguro de responsabilidad civil', NULL, 'a6666666-6666-6666-6666-666666666666', true, 1, NOW(), NOW()),
('Seguro de local', NULL, 'a6666666-6666-6666-6666-666666666666', true, 2, NOW(), NOW()),
('Seguro de equipo', NULL, 'a6666666-6666-6666-6666-666666666666', true, 3, NOW(), NOW());

-- 7. Marketing y Publicidad
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('a7777777-7777-7777-7777-777777777777', 'Marketing y Publicidad', 'Promoción y comunicación', NULL, true, 7, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('Publicidad online', NULL, 'a7777777-7777-7777-7777-777777777777', true, 1, NOW(), NOW()),
('Publicidad offline', NULL, 'a7777777-7777-7777-7777-777777777777', true, 2, NOW(), NOW()),
('Redes sociales', NULL, 'a7777777-7777-7777-7777-777777777777', true, 3, NOW(), NOW()),
('Diseño gráfico', NULL, 'a7777777-7777-7777-7777-777777777777', true, 4, NOW(), NOW()),
('Web y hosting', NULL, 'a7777777-7777-7777-7777-777777777777', true, 5, NOW(), NOW());

-- 8. Mantenimiento y Reparaciones
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('a8888888-8888-8888-8888-888888888888', 'Mantenimiento y Reparaciones', 'Mantenimiento de equipos y local', NULL, true, 8, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('Mantenimiento de equipos', NULL, 'a8888888-8888-8888-8888-888888888888', true, 1, NOW(), NOW()),
('Reparaciones', NULL, 'a8888888-8888-8888-8888-888888888888', true, 2, NOW(), NOW()),
('Calibración de equipos', NULL, 'a8888888-8888-8888-8888-888888888888', true, 3, NOW(), NOW());

-- 9. Software y Tecnología
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('a9999999-9999-9999-9999-999999999999', 'Software y Tecnología', 'Herramientas digitales y tecnología', NULL, true, 9, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('Licencias de software', NULL, 'a9999999-9999-9999-9999-999999999999', true, 1, NOW(), NOW()),
('Aplicaciones', NULL, 'a9999999-9999-9999-9999-999999999999', true, 2, NOW(), NOW()),
('Servicios cloud', NULL, 'a9999999-9999-9999-9999-999999999999', true, 3, NOW(), NOW()),
('Hardware', NULL, 'a9999999-9999-9999-9999-999999999999', true, 4, NOW(), NOW());

-- 10. Impuestos y Tasas
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Impuestos y Tasas', 'Obligaciones fiscales', NULL, true, 10, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('IVA', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', true, 1, NOW(), NOW()),
('Tasas municipales', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', true, 2, NOW(), NOW()),
('Impuestos locales', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', true, 3, NOW(), NOW());

-- 11. Gastos Financieros
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('abbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Gastos Financieros', 'Comisiones y gastos bancarios', NULL, true, 11, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('Comisiones bancarias', NULL, 'abbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', true, 1, NOW(), NOW()),
('Intereses de préstamos', NULL, 'abbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', true, 2, NOW(), NOW()),
('Gastos de transferencias', NULL, 'abbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', true, 3, NOW(), NOW());

-- 12. Otros Gastos
INSERT INTO expense_categories (id, name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('accccccc-cccc-cccc-cccc-cccccccccccc', 'Otros Gastos', 'Gastos varios no categorizados', NULL, true, 12, NOW(), NOW());

INSERT INTO expense_categories (name, description, parent_id, is_active, sort_order, created_at, updated_at) VALUES
('Gastos varios', NULL, 'accccccc-cccc-cccc-cccc-cccccccccccc', true, 1, NOW(), NOW()),
('Imprevistos', NULL, 'accccccc-cccc-cccc-cccc-cccccccccccc', true, 2, NOW(), NOW());
