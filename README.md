Create a new bot:
```
curl -v -X POST -d '{
	"min_price": 100,
	"max_price": 3000,
	"bed": 2,
	"bath": 1,
	"slack_token": "example_slack_token",
	"craigslist_site": "sfbay",
	"max_transit_distance": 2,
	"craigslist_housing_section": "apa",
	"neighborhoods": ["albany", "el cerrito", "berkeley north", "berkeley", "rockridge", "oakland lake merritt"],
	"areas": ["eby", "sby", "nby"],
	"transit_stations": {
		"oakland_19th_bart": [37.8118051, -122.2720873],
		"macarthur_bart": [37.8265657, -122.2686705],
		"rockridge_bart": [37.841286, -122.2566329],
		"downtown_berkeley_bart": [37.8629541, -122.276594],
		"north_berkeley_bart": [37.8713411, -122.2849758],
		"el_cerrito_plaza_bart": [37.902694, -122.298968]
	},
	"boxes": {
		"albany": [
			[37.898925, -122.373782],
			[37.866726, -122.281639]
		],
		"rockridge": [
			[37.83826, -122.24073],
			[37.84680, -122.25944]
		],
		"berkeley": [
			[37.86226, -122.25043],
			[37.86781, -122.26502]
		],
		"north_berkeley": [
			[37.86425, -122.26330],
			[37.87655, -122.28974]
		],
		"oakland": [
			[37.885438, -122.355881],
			[37.631714, -122.114793]
		],
		"richmond": [
			[37.77188, -122.47263],
			[37.78029, -122.51005]
		]
	}
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
	"min_price": 100,
	"max_price": 5000,
	"bed": 3,
	"bath": 1,
	"slack_token": "example_slack_token",
	"craigslist_site": "sfbay",
	"max_transit_distance": 2,
	"craigslist_housing_section": "apa",
	"neighborhoods": ["albany", "el cerrito", "berkeley north", "berkeley", "rockridge", "oakland lake merritt"],
	"areas": ["eby", "sby", "nby"],
	"transit_stations": {
		"oakland_19th_bart": [37.8118051, -122.2720873],
		"macarthur_bart": [37.8265657, -122.2686705],
		"rockridge_bart": [37.841286, -122.2566329],
		"downtown_berkeley_bart": [37.8629541, -122.276594],
		"north_berkeley_bart": [37.8713411, -122.2849758],
		"el_cerrito_plaza_bart": [37.902694, -122.298968]
	},
	"boxes": {
		"albany": [
			[37.898925, -122.373782],
			[37.866726, -122.281639]
		],
		"rockridge": [
			[37.83826, -122.24073],
			[37.84680, -122.25944]
		],
		"berkeley": [
			[37.86226, -122.25043],
			[37.86781, -122.26502]
		],
		"north_berkeley": [
			[37.86425, -122.26330],
			[37.87655, -122.28974]
		],
		"oakland": [
			[37.885438, -122.355881],
			[37.631714, -122.114793]
		],
		"richmond": [
			[37.77188, -122.47263],
			[37.78029, -122.51005]
		]
	}
}'
34.208.42.84:8080/update
```
Current implementation requires that you pass all data to update, this will eventually not be necessary (you will only need to pass data you're updating).