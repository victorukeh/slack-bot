package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/shomali11/slacker"
	"github.com/spf13/viper"
)

func test() {

}

// Here's what the function does:
//  1. It takes in a channel of slacker.CommandEvent values called "analyticsChannel".
//  2. It continuously listens on the channel for incoming command events.
//  3. When it receives a command event, it prints out some information about the event to the console,
//     including the command text, user ID, and channel ID.
//  4. The function continues listening on the channel for more events until it is closed.
func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
	}
}

func viperEnvVariable(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func main() {
	slackBottToken := viperEnvVariable("SLACK_BOT_TOKEN")
	slackAppToken := viperEnvVariable("SLACK_APP_TOKEN")
	bot := slacker.NewClient(slackBottToken, slackAppToken)
	go printCommandEvents(bot.CommandEvents())

	bot.Command("My year of birth is <year>", &slacker.CommandDefinition{
		Description: "yob calculator",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				println("error")
			}
			age := 2023 - yob
			r := fmt.Sprintf("age is %d", age)
			response.Reply(r)
		},
	})
	// Context is often used in Go to control the timeout, cancellation, and passing values between function calls.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println("Ctx: ", ctx)
	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
