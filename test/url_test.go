package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	const subCommand = cmd.CommandURL
	testCases := []struct {
		input    []string
		expected string
	}{
		{
			[]string{runCommand, mainGo, subCommand, "-e", "https://reurl.cc/MNN9Gv"},
			"https://www.setn.com/News.aspx?NewsID=1161776&utm_campaign=viewallnews\n",
		},
		{
			[]string{runCommand, mainGo, subCommand, "-e", "https://bit.ly/3ogGuB1"},
			"https://theinitium.com/article/20220721-mainland-covid-prolonged-grief-disorder/?utm_source=Telegram&utm_medium=Telegram&utm_campaign=Telegram\n",
		},
		{
			[]string{runCommand, mainGo, subCommand, "-e", "https://youtu.be/uLGSEoN5KwI"},
			"https://www.youtube.com/watch?v=uLGSEoN5KwI&feature=youtu.be\n",
		},
		{
			[]string{runCommand, mainGo, subCommand, "-e", "https://utm.to/48vy8a"},
			"https://www.storm.mg/lifestyle/4366824?utm_source=telegram&utm_medium=post\n",
		},
		{
			[]string{runCommand, mainGo, subCommand, "-e", "https://lihi1.cc/itv4p"},
			"https://www.businessweekly.com.tw/focus/indep/6007870?utm_source=Line&utm_medium=social&utm_content=bw&utm_campaign=content\n",
		},
		{
			[]string{runCommand, mainGo, subCommand, "-e", "https://linshibi.pros.is/4c7llt"},
			"https://open.firstory.me/story/cl5k47dii00eo01zxamrp684m/platforms\n",
		},
		{
			[]string{runCommand, mainGo, subCommand, "-e", "https://spoti.fi/3O6QgAb"},
			"https://open.spotify.com/episode/0XOruTQxsN295v0ePD2YAk\n",
		},
		{
			[]string{runCommand, mainGo, subCommand, "-e", "https://redd.it/p2xbpj"},
			"https://www.reddit.com/comments/p2xbpj\n",
		},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			got, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(testCases[i].input, err)
			}
			assert.Equal(t, testCases[i].expected, string(got))
		})
	}
}

func TestURLBinary(t *testing.T) {
	const subCommand = cmd.CommandURL
	args := [][]string{
		{subCommand, "-e", "https://goo.gl/maps/b37Aq3Anc7taXQDd9"},
		{subCommand, "-e", "https://reurl.cc/7peeZl"},
		{subCommand, "-e", "https://bit.ly/3gk7w5x"},
	}
	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
