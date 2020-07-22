package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	Token = "NzAzMDY2ODc0NjY0OTc2Mzk0.Xxhnwg.lyDNuqp-MGhQNHrQaciEmkLJ2a8"
	//flag.StringVar(&Token, "t", "", "Bot Token")
	//flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(move)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func move(s *discordgo.Session, m *discordgo.VoiceStateUpdate){
	guild, err := s.GuildChannels(m.GuildID)
	var ida, idb string
	if err == nil {
		for _, v := range guild {
			if v.Type == discordgo.ChannelTypeGuildVoice {
				if v.Name == "A" {
					ida = v.ID
				}
				if v.Name == "B" {
					idb = v.ID
				}
			}
		}
		if ida != "" && idb != ""{
			if m.ChannelID == ida && m.SelfMute {
				s.GuildMemberMove(m.GuildID, m.UserID, idb)
			}
		}else{
			println("No se encuentran los canales")
		}
	}
}


// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(strings.ToLower(m.Content),"hola"){
		users := parseUsers(m.Content)
		for _,v := range users{
			s.ChannelMessageSend(m.ChannelID,"Hola <@!"+ v + ">")
		}
	}
}

func parseUsers(s string) (ret []string){
	for strings.Index(s,"<@!") != -1{
		s = s[strings.Index(s,"<@!"):]
		if strings.Index(s,">") != -1{
			_, err := strconv.Atoi(s[3:strings.Index(s,">")])
			if err == nil {
				ret = append(ret,s[3:strings.Index(s,">")])
			}
		}
		s = s[1:]
	}
	return
}