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

var (
	TwitterAPIKey            string = os.Getenv("TWITTER_APIKEY")
	TwitterAPISecret         string = os.Getenv("TWITTER_APIKEY_SECRET")
	TwitterAccessToken       string = os.Getenv("TWITTER_ACCESS_TOKEN")
	TwitterAccessTokenSecret string = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
)

func main() {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/tweet", tweet)
	e.Logger.Fatal(e.Start(":8080"))

}

func tweet(c echo.Context) error {
	log.Printf("tweet - Key %s", TwitterAPIKey)
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
	fmt.Printf("STATUSES SHOW:\n%+v\n", tweet)
	return c.String(http.StatusOK, "tweeted :)")
}
