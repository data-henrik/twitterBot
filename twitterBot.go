// Simple Twitter bot for testing
// IBM Cloud Code Engine subscriptions.
// This version of the bot runs as job. That is, the program is called,
// processes the tweet request and exits.
//
// All necessary parameters are passed as environment variables.
//
// Written by Henrik Loeser, 2021-2022

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mmcdole/gofeed"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Variables, taken from .env or K8s secrets
// Allow to configure the blog feed and the tweet text..
var (
	TwitterAPIKey            string = os.Getenv("TWITTER_APIKEY")
	TwitterAPISecret         string = os.Getenv("TWITTER_APIKEY_SECRET")
	TwitterAccessToken       string = os.Getenv("TWITTER_ACCESS_TOKEN")
	TwitterAccessTokenSecret string = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	RSSFeed      string = getEnv("TWITTER_RSSFeed", "https://www.ibm.com/cloud/blog/rss")
	TweetString1 string = getEnv("TWITTER_TweetString1", "A recent #IBMCloud #blog is titled: %s. Read it at %s ")
	TweetString2 string = getEnv("TWITTER_TweetString2", "Written in #GoLang and deployed on #CodeEngine. #IBM #news #cloud")
	ItemRange    string = getEnv("TWITTER_ItemRange", "8")

	CE_DATA string = os.Getenv("CE_DATA")
)

type Tweet_Params struct {
	RSSFeed      string `json:"feed"`
	TweetString1 string `json:"tweet_string1"`
	TweetString2 string `json:"tweet_string2"`
	ItemRange    uint   `json:"item_range,omitempty"`
}

// Fill in defaults if not provided with the current request
func fill(tp *Tweet_Params) {
	if tp.RSSFeed != "" {
		RSSFeed = tp.RSSFeed
	}
	if tp.TweetString1 != "" {
		TweetString1 = tp.TweetString1
	}
	if tp.TweetString2 != "" {
		TweetString2 = tp.TweetString2
	}
	if tp.ItemRange > 0 {
		ItemRange = strconv.FormatUint(uint64(tp.ItemRange), 10)
	}
}

// compose a tweet based on the configured feed
func getMessage(url string, msg1 string, msg2 string, itemRange int64) string {
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

// Main - get ready to tweet
// If it matches, proceed to set up the Twitter client and post the
// status update.
func main() {
	// log for --debug parameter
	debugPtr := flag.Bool("debug", false, "debug mode: only print the composed message")
	flag.Parse()

	// read in possible data coming from event producer
	if CE_DATA != "" {
		data := new(Tweet_Params)
		err := json.Unmarshal([]byte(CE_DATA), &data)
		if err != nil {
			log.Println(err)
			log.Fatal(("error reading JSON input"))
		}
		// fill with defaults if not passed in
		fill(data)
	}
	// compose the tweet
	IRange, _ := strconv.ParseInt(ItemRange, 10, 32)
	message := getMessage(RSSFeed, TweetString1, TweetString2, IRange)

	if !*debugPtr {
		// set up the authorization for the Twitter client
		config := oauth1.NewConfig(TwitterAPIKey, TwitterAPISecret)
		token := oauth1.NewToken(TwitterAccessToken, TwitterAccessTokenSecret)
		httpClient := config.Client(oauth1.NoContext, token)

		// Twitter client
		client := twitter.NewClient(httpClient)
		// now tweet by posting a status update
		tweet, resp, err := client.Statuses.Update(message, nil)
		if err != nil {
			log.Println(resp)
			log.Println(err)
			log.Fatal("Tweet failed")
		}
		// low level debugging for the logs... :)
		log.Printf("STATUSES SHOW:\n%+v\n", tweet)
	}
	// return success indicator as response
}
