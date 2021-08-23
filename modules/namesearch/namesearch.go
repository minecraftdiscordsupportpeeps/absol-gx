package namesearch

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/lordralex/absol/api"
	"github.com/lordralex/absol/api/logger"
	"io"
	"net/http"
	"strings"
	"time"
)

type Module struct {
	api.Module
}

func (*Module) Load(ds *discordgo.Session) {
	api.RegisterCommand("ns", RunCommand)
	api.RegisterCommand("namesearch", RunCommand)

	api.RegisterIntentNeed(discordgo.IntentsGuildMessages, discordgo.IntentsDirectMessages)
}

func RunCommand(ds *discordgo.Session, mc *discordgo.MessageCreate, _ string, args []string) {
	if len(args) == 0 {
		return
	}

	response, err := http.Get("https://api.ashcon.app/mojang/v2/user/" + args[0])
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	type responseStruct struct {
		UUID            string `json:"uuid"`
		Username        string `json:"username"`
		UsernameHistory []struct {
			Username  string    `json:"username"`
			ChangedAt time.Time `json:"changed_at,omitempty"`
		} `json:"username_history"`
		Textures struct {
			Custom bool `json:"custom"`
			Slim   bool `json:"slim"`
			Skin   struct {
				URL  string `json:"url"`
				Data string `json:"data"`
			} `json:"skin"`
			Cape struct {
				URL  string `json:"url"`
				Data string `json:"data"`
			} `json:"cape"`
			Raw struct {
				Value     string `json:"value"`
				Signature string `json:"signature"`
			} `json:"raw"`
		} `json:"textures"`
		CreatedAt string `json:"created_at"`
	}

	var result responseStruct
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return
	}

	// set up the embed
	var fields []*discordgo.MessageEmbedField
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "UUID",
		Value:  result.UUID,
		Inline: false,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Current Username",
		Value:  result.Username,
		Inline: false,
	})

	// add image to embed
	embed := &discordgo.MessageEmbed{
		Title: "Result for `" + strings.Join(args, "") + "`",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: result.Textures.Skin.URL,
		},
		Fields: fields,
		Image: &discordgo.MessageEmbedImage{
			URL:    "https://mc-heads.net/body/" + result.UUID + "/800/left",
			Width:  800,
			Height: 800,
		},
	}

	send := &discordgo.MessageSend{
		Embed: embed,
	}

	_, err = ds.ChannelMessageSendComplex(mc.ChannelID, send)
	if err != nil {
		logger.Err().Printf("Failed to send message\n%s", err)
	}
}
