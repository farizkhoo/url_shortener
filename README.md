# url_shortener

Start api service:

`docker-compose run --rm --service-ports api`

Run migrations:

`make migrate`

Run webserver api

`make run/api`

Open a new terminal and test these endpoints with curl

`curl -X POST -H "Content-Type: application/json" --data “{\”url\”: \”https://www.google.com/\”}” http://localhost:3000/shorten_url/ -v`

Get key value from response json and enter into the following

`curl http://localhost:3000/shorten_url/${key} -v`
