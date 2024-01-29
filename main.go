package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alan-Daniels/job_share/storage"
	"github.com/robfig/cron/v3"
)

var people = []string{
	"Alan",
	"Jarrod",
	"Alex D",
	"Aidan",
	"Jordan",
	"Alex S",
	"Cooper",
}

var weeklyJobs = []string{
	"Make sure bins are done",
	"Deep clean kitchen",
	"Clean floors",
	"Clean floors",
	"Sweep outside",
	"Sweep outside",
	"Wash public cloths/laundry",
}

var dailyJobs = []string{
	"Put away dry dishes",
	"Wipe counters",
	"Clean stove",
	"Trash check throughout house",
	"Checkup outside spaces",
	"Feed fish",
	"turn lights & devices off",
}

type options struct {
	WeeklyDiscordUri string `json:"weekly_discord_uri"`
	DailyDiscordUri  string `json:"daily_discord_uri"`
}

func getOptions() (options, error) {
	var options options
	opt, err := os.ReadFile("/data/options.json")
	if err != nil {
		return options, err
	}
	err = json.Unmarshal(opt, &options)
	if err != nil {
		return options, err
	}

	return options, nil
}

func SelectRandom(people, jobs []string) []string {
	l := max(len(people), len(jobs))
	comb := make([]string, l)

	rand.Shuffle(l, func(i, j int) { jobs[i], jobs[j] = jobs[j], jobs[i] })
	for i := 0; i < l; i++ {
		comb[i] = jobs[i]
	}

	return comb
}

func SendJobsToDiscord(uri, speel string, actions []string) error {
	var msg = storage.DiscordMessage{
		Content: speel,
		Embeds: []storage.DiscordEmbed{
			{
				Fields: make([]storage.DiscordField, len(actions)),
			},
		},
	}
	for i := 0; i < len(actions); i++ {
		msg.Embeds[0].Fields[i].Name = people[i]
		msg.Embeds[0].Fields[i].Value = actions[i]
	}

	fmt.Printf("msg: %v\n", msg)

	bmsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(bmsg))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 200 {
		return fmt.Errorf("bad response, todo explain")
	}
	return nil
}

func doDailyJobs(discordUri string) error {
	err := SendJobsToDiscord(
		discordUri,
		"Here are daily (you may obviously trade jobs between yourselves).",
		SelectRandom(people, dailyJobs),
	)
	return err
}

func doWeeklyJobs(discordUri string) error {
	err := SendJobsToDiscord(
		discordUri,
		"Here are weekly (you may obviously trade jobs between yourselves).",
		SelectRandom(people, weeklyJobs),
	)
	return err
}

func main() {
	opt, err := getOptions()
	if err != nil {
		panic(err)
	}

	loc, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		panic(err)
	}

	cr := cron.New()

	// daily
	cr.AddFunc("@daily", func() {
		dt := time.Now().In(loc)
		fmt.Println("Time for @daily is: ", dt.String())
		err := doDailyJobs(opt.DailyDiscordUri)
		if err != nil {
			println(err)
		}
	})
	// weekly
	cr.AddFunc("@weekly", func() {
		dt := time.Now().In(loc)
		fmt.Println("Time for @weekly is: ", dt.String())
		err := doWeeklyJobs(opt.WeeklyDiscordUri)
		if err != nil {
			println(err)
		}
	})

	cr.Start()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Waiting for sigint or sigterm to exit...")
	<-done
	fmt.Println("Goodbye!")
}
