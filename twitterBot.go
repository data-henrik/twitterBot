// Simple Twitter bot for testing
// IBM Cloud Code Engine subscriptions.
// The bot runs as http server (Code Engine app) and receives
// POST requests on the "tweet" route. The POST requests are
// sent be the ping subscription at the configured times.
//
// Written by Henrik Loeser, 2021

package main

import (
	"fmt"
	"html"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mmcdole/gofeed"

	"github.com/labstack/echo/v4"
)

// Variables, taken from .env or K8s secrets
var (
	TwitterAPIKey            string = os.Getenv("TWITTER_APIKEY")
	TwitterAPISecret         string = os.Getenv("TWITTER_APIKEY_SECRET")
	TwitterAccessToken       string = os.Getenv("TWITTER_ACCESS_TOKEN")
	TwitterAccessTokenSecret string = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	SecretKey                string = os.Getenv("SECRET_KEY")
)

// Allow to configure the blog feed and the tweet text. Set up defaults except for the secret.
type Tweet_Params struct {
	Secret       string `form:"SECRET_KEY" json:"secret_key"`
	RSSFeed      string `form:"FEED" json:"feed"`
	TweetString1 string `form:"TWEET_STRING1" json:"tweet_string1"`
	TweetString2 string `form:"TWEET_STRING2" json:"tweet_string2"`
	ItemRange    uint   `form:"ITEM_RANGE" json:"item_range"`
}

func (tp *Tweet_Params) fill() {
	if tp.RSSFeed == "" {
		tp.RSSFeed = "https://www.ibm.com/cloud/blog/rss"
	}
	if tp.TweetString1 == "" {
		tp.TweetString1 = "A recent #IBMCloud #blog is titled: %s. Read it at %s "
	}
	if tp.TweetString2 == "" {
		tp.TweetString2 = "Written in #GoLang and deployed on #CodeEngine. #IBM #news #cloud"
	}
	if tp.ItemRange == 0 {
		tp.ItemRange = 8
	}
}

// run the http server with two routes:
// 1) /:       Hello world
// 2) /tweet:  Send the Twitter status update (tweet)
func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/tweet", tweet)
	e.Logger.Fatal(e.Start(":8080"))

}

// compose a tweet based on the latest IBM Cloud blog feed
func getMessage(url string, msg1 string, msg2 string, itemRange uint) string {
	// the feed to use
	var tweet string

	// open the feed
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)

	// how many items are in the feed?
	numitem := len(feed.Items)
	if numitem > int(itemRange) {
		numitem = int(itemRange)
	}
	// fetch a random blog entry
	rand.Seed(time.Now().UnixNano())
	var rnum int = rand.Intn(numitem)

	// compose the actual tweet based on a snippet with title and link and a 2nd snippet
	tweet = fmt.Sprintf(msg1,
		html.UnescapeString(feed.Items[rnum].Title),
		html.UnescapeString(feed.Items[rnum].Link))
	tweet += fmt.Sprintf(msg2)
	tweet += fmt.Sprintf(" %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Println(tweet)
	return tweet
}

// Receive the request to tweet. First check for the passed secret.
// If it matches, proceed to set up the Twitter client and post the
// status update.
func tweet(c echo.Context) error {
	data := new(Tweet_Params)
	if bindErr := c.Bind(data); bindErr != nil {
		log.Println("Error binding")
		return bindErr
	}
	// fill with defaults if not passed in
	data.fill()
	if SecretKey == data.Secret {
		// we are good to go, get the message to tweet
		message := getMessage(data.RSSFeed, data.TweetString1, data.TweetString2, data.ItemRange)

		// set up the authorization for the Twitter client
		config := oauth1.NewConfig(TwitterAPIKey, TwitterAPISecret)
		token := oauth1.NewToken(TwitterAccessToken, TwitterAccessTokenSecret)
		httpClient := config.Client(oauth1.NoContext, token)

		// Twitter client
		client := twitter.NewClient(httpClient)
		tweet, resp, err := client.Statuses.Update(message, nil)
		if err != nil {
			log.Println(resp)
			log.Println(err)
			log.Fatal("Tweet failed")
			return err
		}
		// low level debugging for the logs... :)
		log.Printf("STATUSES SHOW:\n%+v\n", tweet)
		return c.String(http.StatusOK, "tweeted :)\n")
	} else {
		log.Fatal("No matching secret provided\n")
		return c.String(http.StatusUnauthorized, "No matching secret provided\n")
	}
}
