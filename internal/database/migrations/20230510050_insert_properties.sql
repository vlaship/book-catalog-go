-- +goose Up

-- insert into properties
INSERT INTO catalog.properties (property_name, property_value)
VALUES ('assumptions', '{
  "full_time_government_obligations": 15,
  "health_insurance": 10,
  "contractor_multiplier": 2,
  "contractor_working_weeks": 48
}');

-- +goose Down
DELETE
FROM catalog.properties
WHERE property_name = 'assumptions';
