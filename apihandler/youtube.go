package apihandler

import (
	pb "PageService/proto"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

var youtubeApiKeys = [...]string{
	"AIzaSyCKUyMUlRTHMG9LFSXPYEDQYn7BCfjFQyI",
	"AIzaSyCNGkNspHPreQQPdT-q8KfQznq4S2YqjgU",
	"AIzaSyABJehNy0EEzzKl-I7hXkvYeRwIupl2RYA"}

// Youtube PlaylistApi
type YtPlayist struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind    string `json:"kind"`
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
				Standard struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"standard"`
				Maxres struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"maxres"`
			} `json:"thumbnails"`
			ChannelTitle string `json:"channelTitle"`
			PlaylistID   string `json:"playlistId"`
			Position     int    `json:"position"`
			ResourceID   struct {
				Kind    string `json:"kind"`
				VideoID string `json:"videoId"`
			} `json:"resourceId"`
		} `json:"snippet"`
		ContentDetails struct {
			VideoID          string    `json:"videoId"`
			VideoPublishedAt time.Time `json:"videoPublishedAt"`
		} `json:"contentDetails"`
	} `json:"items"`
}

// youtube channhels
type YTChannel struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
			} `json:"thumbnails"`
			ChannelTitle         string `json:"channelTitle"`
			LiveBroadcastContent string `json:"liveBroadcastContent"`
		} `json:"snippet"`
	} `json:"items"`
}

// youtube search
type YTSearch struct {
	Kind          string `json:"kind"`
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		Kind string `json:"kind"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

//TODO need to add YOUTUBE KEY Swapper.

func youtubePlaylist(playlistId string, contentChan chan *pb.Content, errChan chan error) {

	defer close(contentChan)

	if req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/playlistItems", nil); err != nil {
		errChan <- err
	} else {
		q := req.URL.Query()
		q.Add("key", youtubeApiKeys[0])
		q.Add("playlistId", playlistId)
		q.Add("part", "snippet,contentDetails")
		q.Add("maxResults", "50")
		req.URL.RawQuery = q.Encode()
		client := &http.Client{}
		if resp, err := client.Do(req); err != nil {
			errChan <- err
		} else if resp.StatusCode != 200 {
			errChan <- status.Error(codes.NotFound, "Playlist Data not found")
		} else {

			var playlistResp YtPlayist
			if err = json.NewDecoder(resp.Body).Decode(&playlistResp); err != nil {
				errChan <- err
			}
			if err = resp.Body.Close(); err != nil {
				errChan <- err
			}

			for _, item := range playlistResp.Items {
				var contentTile pb.Content
				contentTile.Title = item.Snippet.Title
				contentTile.Poster = []string{item.Snippet.Thumbnails.Medium.URL}
				contentTile.Portriat = []string{item.Snippet.Thumbnails.Medium.URL}
				contentTile.IsDetailPage = false
				contentTile.Type = pb.TileType_ImageTile

				var play pb.Play
				play.Package = "com.google.android.youtube"
				if item.Snippet.ResourceID.Kind == "youtube#video" {
					play.Target = item.Snippet.ResourceID.VideoID
					play.Source = "Youtube"
					play.Type = "CWYT_VIDEO"
					contentTile.Play = []*pb.Play{&play}
					contentChan <- &contentTile
				}
			}
		}
	}
}

func youtubeChannel(channelId string, contentChan chan *pb.Content, errChan chan error) {

	defer close(contentChan)

	if req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/search", nil); err != nil {
		errChan <- err
	} else {
		q := req.URL.Query()
		q.Add("key", youtubeApiKeys[0])
		q.Add("channelId", channelId)
		q.Add("part", "snippet")
		q.Add("maxResults", "50")

		req.URL.RawQuery = q.Encode()

		client := &http.Client{}
		if resp, err := client.Do(req); err != nil {
			errChan <- err
		} else if resp.StatusCode != 200 {
			errChan <- status.Error(codes.NotFound, "Channel Data not found")
		} else {

			var playlistResp YTChannel
			if err = json.NewDecoder(resp.Body).Decode(&playlistResp); err != nil {
				errChan <- err
			}
			if err = resp.Body.Close(); err != nil {
				errChan <- err
			}

			for _, item := range playlistResp.Items {
				var contentTile pb.Content
				contentTile.Title = item.Snippet.Title
				contentTile.Poster = []string{item.Snippet.Thumbnails.Medium.URL}
				contentTile.Portriat = []string{item.Snippet.Thumbnails.Medium.URL}
				contentTile.IsDetailPage = false
				contentTile.Type = pb.TileType_ImageTile

				var play pb.Play
				play.Package = "com.google.android.youtube"
				if item.ID.Kind == "youtube#video" {
					play.Target = item.ID.VideoID
					play.Source = "Youtube"
					play.Type = "CWYT_VIDEO"
					contentTile.Play = []*pb.Play{&play}
					contentChan <- &contentTile
				}
			}
		}
	}
}

func youtubeSearch(searchQuery string, contentChan chan *pb.Content, errChan chan error) {

	defer close(contentChan)

	if req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/search", nil); err != nil {
		errChan <- err
	} else {

		q := req.URL.Query()
		q.Add("key", youtubeApiKeys[0])
		q.Add("maxResults", "50")
		q.Add("q", searchQuery)
		q.Add("part", "snippet")

		req.URL.RawQuery = q.Encode()
		client := &http.Client{}
		if resp, err := client.Do(req); err != nil {
			errChan <- err
		} else if resp.StatusCode != 200 {
			errChan <- status.Error(codes.NotFound, "Search Data not found")
		} else {

			var searchResp YTSearch
			if err = json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
				errChan <- err
			}

			if err = resp.Body.Close(); err != nil {
				errChan <- err
			}

			for _, item := range searchResp.Items {
				var contentTile pb.Content
				contentTile.Title = item.Snippet.Title
				contentTile.Poster = []string{item.Snippet.Thumbnails.Medium.URL}
				contentTile.Portriat = []string{item.Snippet.Thumbnails.Medium.URL}
				contentTile.IsDetailPage = false
				contentTile.Type = pb.TileType_ImageTile

				var play pb.Play
				play.Package = "com.google.android.youtube"
				if item.ID.Kind == "youtube#video" {
					play.Target = item.ID.VideoID
					play.Source = "Youtube"
					play.Type = "CWYT_VIDEO"
					contentTile.Play = []*pb.Play{&play}
					contentChan <- &contentTile
				}
			}
		}
	}
}
