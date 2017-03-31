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
ec2-52-25-39-194.us-west-2.compute.amazonaws.com:8080/create
```

Delete a bot:
```
curl -v -X POST -d '{
    "slack_token":"example_slack_token"
}'
ec2-52-25-39-194.us-west-2.compute.amazonaws.com:8080/delete
```
