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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/labstack/echo/v4"

	_ "github.com/joho/godotenv/autoload"
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

// Receive the request to tweet. First check for the passed secret.
// If it matches, proceed to set up the Twitter client and post the
// status update.
func tweet(c echo.Context) error {
	PassedSecret := c.FormValue("SECRET_KEY")
	if SecretKey == PassedSecret {
		config := oauth1.NewConfig(TwitterAPIKey, TwitterAPISecret)
		token := oauth1.NewToken(TwitterAccessToken, TwitterAccessTokenSecret)
		// http.Client will automatically authorize Requests
		httpClient := config.Client(oauth1.NoContext, token)

		// twitter client
		client := twitter.NewClient(httpClient)
		message := fmt.Sprintf("I am running on #IBMCloud #CodeEngine using #Golang and testing #event subscriptions. The time is %s", time.Now())
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
