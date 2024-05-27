-- +goose Up

-- insert into properties
INSERT INTO catalog.properties (property_name, property_value)
VALUES ('tos', '{"title":"Terms of Service","text":"By using our services, you agree to the following terms..."}');

-- +goose Down
DELETE
FROM catalog.properties
WHERE property_name = 'tos';
