package main

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	consumerKey := "XkiEGQ7M7uAt0Bu7VLvN3cp1F"
	consumerSecret := "RyOrPXUG5ebgp1BT4RGIJ9cTiqZBJlI7Y1XnRuLUYMT30PlWKN"
	accessToken := "179565456-8zA2I8T4slRSd7ehyzgwlkIeSqLf5yaIVzesxXke"
	accessSecret := "ReJvTyHONgCoP5Hj9iQIdalvEp2U0bCoMVY1VsiJh7iCx"

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		IncludeEmail: twitter.Bool(true),
	}
	user, _, _ := client.Accounts.VerifyCredentials(verifyParams)
	fmt.Printf("User's Name:%+v\n", user.ScreenName)

	searchParams := &twitter.SearchTweetParams{
		Query:      "#golang",
		Count:      5,
		ResultType: "recent",
		Lang:       "en",
	}
	searchResult, _, _ := client.Search.Tweets(searchParams)

	for _, tweet := range searchResult.Statuses {
		tweetid := tweet.ID
		a, b, c := client.Statuses.Retweet(tweetid, &twitter.StatusRetweetParams{})
		fmt.Println(a, b, c)

		fmt.Printf("RETWEETED: %+v\n", tweet.Text)
	}
}
