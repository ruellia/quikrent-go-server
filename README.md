Create a new bot:
```
curl -v -X POST -d '{
    "min_price":300,
    "max_price":500,
    "bed":3,
    "bath":1,
    "slack_token":"xoxp-139832760016-139832760080-141230786695-779c19d54d6410ed09f73f795d649cf2",
    "neighborhoods":
        ["berkeley north", "berkeley", "rockridge"]
}'
ec2-52-25-39-194.us-west-2.compute.amazonaws.com:8080/create
```

Delete a bot:
```
curl -v -X POST -d '{
"slack_token":"test"
}'
ec2-52-25-39-194.us-west-2.compute.amazonaws.com:8080/delete
```