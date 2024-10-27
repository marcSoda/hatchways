package main

import (
	"bytes"
	"encoding/json"
	"my_api/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

//TestPing tests the ping route
func TestPing(t *testing.T) {
	//new request
	req, err := http.NewRequest("GET", "/api/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.PingHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "{\"success\":true}\n"

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//TestPostsGetHandlerLikesAsc tests post route on ordering by likes in ascending order
func TestPostsGetHandlerLikesAsc(t *testing.T) {
	pack := api.FetchPacket{
		Tags:      "science,tech",
		SortBy:    "likes",
		Direction: "asc",
	}
	postsTester(pack, 200, "NIL", t)
}

//TestPostsGetHandlerLikesDesc tests post route on ordering by likes in descending order
func TestPostsGetHandlerLikesDesc(t *testing.T) {
	pack := api.FetchPacket{
		Tags:      "science,tech",
		SortBy:    "likes",
		Direction: "desc",
	}
	postsTester(pack, 200, "NIL", t)
}

//TestPostsGetHandlerReadsAsc tests post route on ordering by reads in ascending order
func TestPostsGetHandlerReadsAsc(t *testing.T) {
	pack := api.FetchPacket{
		Tags:      "science,tech",
		SortBy:    "reads",
		Direction: "asc",
	}
	postsTester(pack, 200, "NIL", t)
}

//TestPostsGetHandlerReadsDesc tests post route on ordering by reads in descending order
func TestPostsGetHandlerReadsDesc(t *testing.T) {
	pack := api.FetchPacket{
		Tags:      "science,tech",
		SortBy:    "reads",
		Direction: "desc",
	}
	postsTester(pack, 200, "NIL", t)
}

//TestPostsGetHandlerPopularityAsc tests post route on ordering by popularity in ascending order
func TestPostsGetHandlerPopularityAsc(t *testing.T) {
	pack := api.FetchPacket{
		Tags:      "science,tech",
		SortBy:    "popularity",
		Direction: "asc",
	}
	postsTester(pack, 200, "NIL", t)
}

//TestPostsGetHandlerPopularityDesc tests post route on ordering by popularity in descenging order
func TestPostsGetHandlerPopularityDesc(t *testing.T) {
	pack := api.FetchPacket{
		Tags:      "science,tech",
		SortBy:    "popularity",
		Direction: "desc",
	}
	postsTester(pack, 200, "NIL", t)
}

//TestPostsGetHandlerIdAsc tests post route on ordering by id in ascending order
func TestPostsGetHandlerIdAsc(t *testing.T) {
	pack := api.FetchPacket{
		Tags:      "science,tech",
		SortBy:    "id",
		Direction: "asc",
	}
	postsTester(pack, 200, "NIL", t)
}

//TestPostsGetHandlerIdDesc tests post route on ordering by id in descending order
func TestPostsGetHandlerIdDesc(t *testing.T) {
	pack := api.FetchPacket{
		Tags:      "science,tech",
		SortBy:    "id",
		Direction: "desc",
	}
	postsTester(pack, 200, "NIL", t)
}

//TestMissingTags tests post route with missing tags
func TestMissingTags(t *testing.T) {
	pack := api.FetchPacket{
		SortBy:    "id",
		Direction: "desc",
	}
	postsTester(pack, 400, "{\"error\":\"Tags parameter is required\"}", t)
}

//TestInvalidSortBy tests post route with an invalid sortBy paramater
func TestInvalidSortBy(t *testing.T) {
	pack := api.FetchPacket{
		Tags:      "science",
		SortBy:    "d",
		Direction: "desc",
	}
	postsTester(pack, 400, "{\"error\":\"sortBy paramater is invalid\"}", t)
}

//TestInvalidDirection tests post route with an invalid direction paramater
func TestInvalidDirection(t *testing.T) {
	pack := api.FetchPacket{
		Tags:      "science",
		SortBy:    "id",
		Direction: "i",
	}
	postsTester(pack, 400, "{\"error\":\"direction paramater is invalid\"}", t)
}

//TestMissingSortBy tests post route with missing sortBy in request
func TestMissingSortBy(t *testing.T) {
	pack := api.FetchPacket{
		Tags:      "science",
		Direction: "desc",
	}
	postsTester(pack, 200, "NIL", t)
}

//TestMissingDirection tests post route with missing direction in request
func TestMissingDirection(t *testing.T) {
	pack := api.FetchPacket{
		Tags:   "science",
		SortBy: "id",
	}
	postsTester(pack, 200, "NIL", t)
}

//postsTester issues an api call and tests results
//Used to reduce code duplication for subsequent tests
//expectedCode is the expected http response code
//expectedMessage is the expected response body. Declare NIL if no comparison is necessary
func postsTester(pack api.FetchPacket, expectedCode int, expectedMessage string, t *testing.T) {
	//define struct representing request body
	packBytes, err := json.Marshal(pack)
	if err != nil {
		t.Fatal(err)
	}
	packBytesReader := bytes.NewReader(packBytes)

	//new request
	req, err := http.NewRequest("GET", "/api/posts", packBytesReader)
	if err != nil {
		t.Fatal(err)
	}
	tr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.PostsGetHandler)
	handler.ServeHTTP(tr, req)
	//ensure correct status
	if status := tr.Code; status != expectedCode {
		t.Errorf("wrong status code: got %v want %v. Body: %v", status, expectedCode, tr.Body)
	}

	//if there is an expected message, test it
	if expectedMessage != "NIL" {
		if tr.Body.String() != expectedMessage+"\n" {
			t.Errorf("unexpected body: got %v want %v",
				tr.Body.String(), expectedMessage)
		}
	}

	//get response body
	var posts []api.Post
	err = json.NewDecoder(tr.Body).Decode(&posts)

	//ensure correct order of posts
	for i := 1; i < len(posts)-1; i++ {
		var li = i - 1
		var ri = i
		if pack.Direction == "desc" {
			li = i
			ri = i - 1
		}
		if pack.SortBy == "reads" {
			if posts[li].Reads > posts[ri].Reads {
				t.Error("Order is wrong")
			}
		} else if pack.SortBy == "likes" {
			if posts[li].Likes > posts[ri].Likes {
				t.Error("Order is wrong")
			}
		} else if pack.SortBy == "popularity" {
			if posts[li].Popularity > posts[ri].Popularity {
				t.Error("Order is wrong")
			}
		} else {
			if posts[li].Id > posts[ri].Id {
				t.Error("Order is wrong")
			}
		}
	}
}
