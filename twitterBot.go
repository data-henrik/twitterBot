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
func getMessage() string {
	// the feed to use
	var url = "https://www.ibm.com/cloud/blog/rss"
	var tweet string
	// we also fetch a random other blog entry
	rand.Seed(time.Now().UnixNano())
	var rnum int = rand.Intn(9) + 1

	// open the feed
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)

	// compose the actual message based on two entries and the time for uniqueness
	tweet = fmt.Sprintf("The latest #IBMCloud #blog is titled: %s. Read it at %s ",
		html.UnescapeString(feed.Items[0].Title),
		html.UnescapeString(feed.Items[0].Link))
	tweet += fmt.Sprintf("An older blog is '%s' available at %s #news #cloud",
		html.UnescapeString(feed.Items[rnum].Title),
		html.UnescapeString(feed.Items[rnum].Link))
	tweet += fmt.Sprintf(" %s", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(tweet)
	return tweet
}

// Receive the request to tweet. First check for the passed secret.
// If it matches, proceed to set up the Twitter client and post the
// status update.
func tweet(c echo.Context) error {
	PassedSecret := c.FormValue("SECRET_KEY")
	if SecretKey == PassedSecret {
		// we are good to go, get the message to tweet
		message := getMessage()

		// set up the authorization for the Twitter client
		config := oauth1.NewConfig(TwitterAPIKey, TwitterAPISecret)
		token := oauth1.NewToken(TwitterAccessToken, TwitterAccessTokenSecret)
		httpClient := config.Client(oauth1.NoContext, token)

		// Twitter client
		client := twitter.NewClient(httpClient)
		tweet, _, err := client.Statuses.Update(message, nil)
		if err != nil {
			log.Fatal(TwitterAPIKey)
			return err
		}
		// low level debugging for the logs... :)
		fmt.Printf("STATUSES SHOW:\n%+v\n", tweet)
		return c.String(http.StatusOK, "tweeted :)")
	} else {
		return c.String(http.StatusUnauthorized, "No matching secret provided")
	}
}
