# Delta-API

`docker build -t delta .`

`docker run --env-file .env  -p 8081:8080 -it delta`

Where 8080 is the port for the server, and  8081 is the port of the API


## TODO
- Add logging
- Catch all endpoint ? Like a 404