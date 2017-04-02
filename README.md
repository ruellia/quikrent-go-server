Create a new bot:
```
curl -v -X POST -d '{
    "min_price":1000,
    "max_price":3000,
    "bed":3,
    "bath":1,
    "slack_token":"example_slack_token",
    "neighborhoods":
        ["berkeley north", "berkeley", "rockridge"]
}'
34.208.42.84:8080/create
```

Delete an existing bot:
```
curl -v -X POST -d '{
    "slack_token":"example_slack_token"
}'
34.208.42.84:8080/delete
```

Update an existing bot:
```
curl -v -X POST -d '{
    "min_price":1000,
    "max_price":3000,
    "bed":3,
    "bath":1,
    "slack_token":"example_slack_token",
    "neighborhoods":
        ["berkeley north", "berkeley", "rockridge"]
}'
34.208.42.84:8080/create
```