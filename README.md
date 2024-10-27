# Golang Backend Assessment - Blog Posts

## Intro
I very much enjoyed this project. I find the project format to be much more indicative of merit than short coding interviews. I have no recommendations for improvement. I was able to finish all aspects of the project, incuding the bonus.

## Timeline
- Base API with concurrency: ~ 2.5 hours
- Testing:                   ~ 1 hour
- Bonus:                     ~ 1 hours
- Total:                     ~ 4.5 hours
This summer, I hold a research position, a teaching assistant position, and am taking classes at Lehigh University. I was unable to finish this project in one sitting due to these responsibilities.

## Structure
My repository is organized modularly with future development in mind. I separated the API and Middleware into separate packages main for this reason. A future developer could easily add packages for static files, a database, sessions, etc.

The code is well documented and properly formatted.

## API
The routes are fairly straightforward.
The most notable aspect of the API is the concurrency built in to the posts route.
I checked for race conditions and was able to produce none.

## Testing
Rather than having a thousand line file with huge amounts of unnecessary code duplication, I created one function that makes all of the API calls for the posts route. It is called by 13 functions that test the posts route. This approach saved over 1000 lines of code and tons of unnecessary duplication.

## Bonus
I found a great resource for server-side-caching at https://github.com/goenning/go-cache-demo
The implementation is quite simple.
There are many libraries for server-side-caching in go, but none of the worked exactly as I wanted them to. I was able to adapt the code found in the aforementioned repository to suit my needs perfectly.
Caching improves speeds significantly. Without caching, my calls to the posts route averaged about 150ms. After caching, the response time averaged around 0.50ms.

## Running the API
- Go version: `go1.18.3 linux/amd64`
- Operating system: `Arch Linux 5.18.9-arch1-1`
- Dependencies are listed in go.mod. The only dependency is gorilla mux.
- After cloning, the application should be able to be run by running `go run main.go`. The API is running on port 1701. "Running" is printed to the console when it starts up.
- To run tests, simply run `go test -v`

## Conclusion
Thank you very much for your time. I found this assessment extremely enjoyable. I am excited to hear back!
