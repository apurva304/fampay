# Youtube Client
- We are using Youtube search api to fetch videos at a particular intervel.
- We have used custom round robin based loadbalance to use multiple api keys and make more number of call without qouta exceeding.
- Youtube client can be found in `youtube` directory. <br>

# Build and Run
- To build the project locally use `$ go build -o main.go`
- To run the project using the docker-compose.yml file use `$ docker-compose -f docker-compose.yml up`
- To hit the search api use `curl http://localhost:3000/videos?query=<query>`
- The above url can also take optional parameters like `pageNumber` and `pageItemCount`. By default pageNumber 1 and pageItemCount is set as 10
