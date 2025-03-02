curl --location 'http://localhost:80/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(1 / 5 * 2) + (1 - 5) / 5 * ( 5 / 8)"
}