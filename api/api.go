package api

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"fmt"

	"github.com/gorilla/mux"

	"my_api/cache"
)

//url for hatchways api
var url = "https://api.hatchways.io/assessment/blog/posts?tag="

//PingPacket is a struct representing response body from ping route
type PingPacket struct {
	Success bool `json:"success"`
}

// ErrorPacket is a struct representing error response body from posts route
type ErrorPacket struct {
	Error string `json:"error"`
}

// FetchPacket is a struct representing request body for posts route
type FetchPacket struct {
	Tags      string `json:"tags"`
	SortBy    string `json:"sortBy"`
	Direction string `json:"direction"`
}

//PostsPacket is a struct representing successful response body from posts route
type PostsPacket struct {
	Posts []Post `json:"posts"`
}

//Post is a struct representing a single post
type Post struct {
	Id         int      `json:"id"`
	Author     string   `json:"author"`
	AuthorId   int      `json:"authorId"`
	Likes      int      `json:"likes"`
	Popularity float32  `json:"popularity"`
	Reads      int      `json:"reads"`
	Tags       []string `json:"tags"`
}

//NewRouter sets up and returns a new gorilla mux router
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/ping", PingHandler).Methods("GET")
	//use cache middleware on the posts route
	r.HandleFunc("/api/posts", cache.Cached("1h", PostsGetHandler)).Methods("GET")
	return r
}

//PingHandler handles the ping route
func PingHandler(w http.ResponseWriter, r *http.Request) {
	//ping the hatchways api on port 80
	_, err := net.DialTimeout("tcp", "api.hatchways.io:80", 1*time.Second)
	succ := true
	if err != nil {
		succ = false
	}
	//encode a ping packet with status of the ping
	pack := PingPacket{
		Success: succ,
	}
	json.NewEncoder(w).Encode(pack)
}

//PostsGetHandler handles the posts route
func PostsGetHandler(w http.ResponseWriter, r *http.Request) {
	//decode the request body
	var pack FetchPacket
	json.NewDecoder(r.Body).Decode(&pack)

	//respond with 400 and error if tags are not present
	if pack.Tags == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorPacket{"Tags parameter is required"})
		return
	}

	//check SortBy paramater
	if pack.SortBy != "id" &&
		pack.SortBy != "reads" &&
		pack.SortBy != "likes" &&
		pack.SortBy != "popularity" {

		//if SortBy is not present, set SortBy to be 'id'
		if pack.SortBy == "" {
			pack.SortBy = "id"
		//respond with 400 and error of SortBy is invalid
		} else {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorPacket{"sortBy paramater is invalid"})
			return
		}
	}

	//check Direction paramater
	if pack.Direction != "asc" && pack.Direction != "desc" {
		//set direction to asc if not present in request
		if pack.Direction == "" {
			pack.Direction = "asc"
		//respond with 400 if direction is invalid
		} else {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorPacket{"direction paramater is invalid"})
			return
		}
	}

	//split tags string into string slice of tags
	tags := strings.Split(pack.Tags, ",")

	//concurrenlty make api calls.
	wg := sync.WaitGroup{}
	var m sync.Map
	for _, t := range tags {
		wg.Add(1)
		go func(tag string) {
			t0 := time.Now()
			fmt.Println(t0)
			resp, err := http.Get(url + tag)
			t1 := time.Now()
			fmt.Println(t1.Sub(t0))
			//fetch data from api. respond with 500 if error
			if err != nil {
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(ErrorPacket{err.Error()})
				return
			}
			//read response body. respond with 500 if error
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(ErrorPacket{err.Error()})
				return
			}
			//get response body as a PostsPacket (slice of Posts)
			var posts PostsPacket
			json.Unmarshal(data, &posts)
			//add all posts to a concurrent map to quickly remove duplicates
			for _, p := range posts.Posts {
				m.Store(p.Id, p)
			}
			wg.Done()
		}(t)
	}
	wg.Wait()

	//read map into a slice of Post
	var posts []Post
	m.Range(func(key, value interface{}) bool {
		posts = append(posts, value.(Post))
		return true
	})
	//sort the Post slice by pack.SortBy in the order of pack.Direction
	sort.Slice(posts, func(i, j int) bool {
		if pack.Direction == "desc" {
			var tmp = i
			i = j
			j = tmp
		}
		if pack.SortBy == "id" {
			return posts[i].Id < posts[j].Id
		} else if pack.SortBy == "reads" {
			return posts[i].Reads < posts[j].Reads
		} else if pack.SortBy == "likes" {
			return posts[i].Likes < posts[j].Likes
		} else if pack.SortBy == "popularity" {
			return posts[i].Popularity < posts[j].Popularity
		}
		return true
	})
	//encode posts
	json.NewEncoder(w).Encode(posts)
}
