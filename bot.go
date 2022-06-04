package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var Token = flag.String("token", "", "Bot access token")
var GuildID = flag.String("guild", "", "Guild ID")
var s *discordgo.Session

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "sweep-nickname",
		Description: "特定のロールが付与されたユーザー以外の全ユーザーのニックネームを消し去る",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "exclude-role",
				Description: "ロール付きを除外するか否か",
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Required:    false,
			},
		},
	},
}
var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"sweep-nickname": sweep_nickname,
}

func init() {
	flag.Parse()
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + *Token)
	if err != nil {
		fmt.Println("Cannot create a new instance: ", err)
		return
	}
}

func init() {
	err := s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}

		registeredCommands[i] = cmd
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stop

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	s.Close()
}

func sweep_nickname(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	var after string
	var excludeRole bool
	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Name == "exclude-role" {
			excludeRole = opt.BoolValue()
		}
	}
	sweepCount := 0
	//ギルドの全ユーザーに対する処理が終わるまで繰り返す
	for {
		members, err := discord.GuildMembers(i.GuildID, after, 1000)
		if err != nil {
			fmt.Println("Cannot get guild data: ", err)
			return
		}

		for _, member := range members {
			if member.Nick == "" {
				continue
			}

			if excludeRole && len(member.Roles) == 0 {
				continue
			}

			fmt.Printf("[NICK]%s#%s: %s\n", member.User.Username, member.User.Discriminator, member.Nick)

			err := discord.GuildMemberNickname(i.GuildID, member.User.ID, "")
			if err != nil {
				fmt.Println("[ERR]: ", err)
				continue
			}
			sweepCount++
		}

		//1000以下であれば、それ以上処理するユーザーはいないので終了
		if len(members) <= 1000 {
			break
		}
		//次の処理の開始位置を設定
		after = members[999].User.ID
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%d人のニックネームを除去しました", sweepCount),
		},
	})
}
