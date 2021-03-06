package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func printTweet(tweets, name, id chan string) {
	for i := 0; i < 10; i++ {
		fmt.Println(<-tweets)
		fmt.Println(<-name)
		fmt.Println(<-id)
	}
}
func main() {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")

	tweetText := make(chan string, 10)
	tweetName := make(chan string, 10)
	tweetID := make(chan string, 10)
	flags.Parse(os.Args[1:])

	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		//fmt.Println(tweet.Text)
		tweetText <- tweet.Text
		tweetName <- tweet.User.ScreenName
		tweetID <- tweet.User.Name
		time.Sleep(500 * time.Millisecond)
	}

	go printTweet(tweetText, tweetName, tweetID)
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}

	fmt.Println("Starting Stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"bengaluru", "bangalore", "blr"},
		StallWarnings: twitter.Bool(true),
		Language:      []string{"en"},
	}

	searchParams := &twitter.SearchTweetParams{
		Query:      "#golang",
		Count:      5,
		ResultType: "recent",
		Lang:       "en",
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	searchResult, _, _ := client.Search.Tweets(searchParams)

	for _, tweet := range searchResult.Statuses {
		tweetID := tweet.ID
		a, b, c := client.Statuses.Retweet(tweetID, &twitter.StatusRetweetParams{})
		fmt.Println(a, b, c)

		fmt.Printf("RETWEETED: %+v\n", tweet.Text)
	}

	// USER (quick test: auth'd user likes a tweet -> event)
	// userParams := &twitter.StreamUserParams{
	// 	StallWarnings: twitter.Bool(true),
	// 	With:          "followings",
	// 	Language:      []string{"en"},
	// }
	// stream, err := client.Streams.User(userParams)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// SAMPLE
	// sampleParams := &twitter.StreamSampleParams{
	// 	StallWarnings: twitter.Bool(true),
	// }
	// stream, err := client.Streams.Sample(sampleParams)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Receive messages until stopped or stream quits

	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
