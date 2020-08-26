# studio 
Webservice for handling the creation of classes for the studio, and bookings for the classes

## Running this webservice
- Ensure go is installed and that $GOROOT is exported as the appropriate directory
- In the terminal run `git clone git@github.com:FrancisLennon17/studio.git`
- Create a new file in `studio/config` called config.go and copy the contents of config.json.example into the file. Update that information with the appropriate configuration
- Ensure the DB has the required tables, see `DB scripts` heading for more
- from the newly created directory directory, run `go run main.go` (can optionally pass `-config={config.json location}`, but the file will be fetched from the config folder by default)

## DB scripts
Scripts for creating the required tables used in this webservice can be found in the  `/scripts` directory

## Endpoints
`POST /classes` to create classes.

Sample request body:
```
{
    "class_name": "Yoga",
    "start_date": "2020-10-01",
    "end_date": "2020-10-10",
    "capacity": 20
}
```

`POST /bookings` to create a booking for a date

Sample request body: 
```
{
    "name": "jimbo",
    "date": "2020-10-01"
}
```

## Postman Collection
https://www.getpostman.com/collections/63dc832c827574bf8ee7

## To Do
- Add swagger-parser comments for auto generation of swagger file
- Finish HTTP unit tests (Need to mock time.Now in sepearte interface also)
- Finish DB unit tests
- Add config package tests
- Add additional request body validation
- Implement consistant linting rules
- Add additional config for dbTimeout & max gym capacity
