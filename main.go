package main

import (
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	Token = "NzAzMDY2ODc0NjY0OTc2Mzk0.XqJMDQ.bIUc2QePzGLqcpVCo_RqxH8d8QE"
	//flag.StringVar(&Token, "t", "", "Bot Token")
	//flag.Parse()
}

var buffer = make([][]byte, 0)

var fullAccessRole = "daisy"

func main() {

	// Load the sound file.
	/* Lo dejo para cuando hagamos los audios
	err := loadSound()
	if err != nil {
		fmt.Println("Error loading sound: ", err)
		fmt.Println("Please copy $GOPATH/src/github.com/bwmarrin/examples/airhorn/airhorn.dca to this directory.")
		return
	}
	*/

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

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

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.

	if m.Author.ID == session.State.User.ID {
		return
	}

	// Find the channel that the message came from.
	channel, err := session.State.Channel(m.ChannelID)
	if err != nil {
		println("Could not find channel.")
		return
	}

	// Find the guild for that channel.
	guild, err := session.State.Guild(channel.GuildID)
	if err != nil {
		println("Could not find guild.")
		return
	}

	// Find the member that send the meesage.
	member, err := session.GuildMember(guild.ID,m.Author.ID)
	if err != nil {
		println("Could not find member.")
		return
	}

	fullAccess := getFullAccessStatus(session,guild,member)

	/*
	if strings.HasPrefix(strings.ToLower(m.Content), "airhorn"){

		// Look for the message sender in that guild'session current voice states.
		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playSound(session, guild.ID, vs.ChannelID)
				if err != nil {
					fmt.Println("Error playing sound:", err)
				}

				return
			}
		}
	}

	if strings.HasPrefix(strings.ToLower(m.Content),"hola"){
		users := parseUsers(m.Content)
		for _,v := range users{
			session.ChannelMessageSend(m.ChannelID,"Hola <@!"+ v + ">")
		}
		return
	}
    */



	//Para abajo las funciones que no pueden hacer todos xd
	//Juanma cdo refactorices esto me vas a putear mucho tkm
	if !fullAccess {
		session.ChannelMessage(channel.ID,"<@!"+m.Author.ID+"> ")
		return
	}


	noTengoGanasDePensarUnNombreDeVariableWip, err := regexp.MatchString("^muu+te$", m.Content)

	if err != nil{
		return
	}

	if noTengoGanasDePensarUnNombreDeVariableWip{
		muteChannel(session,guild,m.Author.ID,true)
		return
	}

	noTengoGanasDePensarUnNombreDeVariableWip, err = regexp.MatchString("^desmuu+te$", m.Content)

	if err != nil{
		return
	}

	if noTengoGanasDePensarUnNombreDeVariableWip{
		muteChannel(session,guild,m.Author.ID,false)
		return
	}

}

func muteChannel(s *discordgo.Session, guild *discordgo.Guild, userID string, mute bool){

	var channelToMuteAll string

	// Look for the message sender in that guild's current voice states.
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			channelToMuteAll = vs.ChannelID
			break
		}
	}

	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == channelToMuteAll {
			s.GuildMemberMute(vs.GuildID,vs.UserID,mute)
		}
	}

}

func getFullAccessStatus(session *discordgo.Session, guild *discordgo.Guild, member *discordgo.Member) bool{

	//Esto es porque me estaba olvidando del Jebus (?
	if guild.OwnerID == member.User.ID {
		return true
	}

	roles, err := session.GuildRoles(guild.ID)
	if err != nil {
		println("Could not find roles.")
		return false
	}

	adminRoleList := make([]string, 0, len(roles))

	for _, role := range roles{
		if role.Permissions & 8 != 0 || strings.ToLower(role.Name) == fullAccessRole{
			adminRoleList = adminRoleList[:len(adminRoleList)+1]
			adminRoleList[len(adminRoleList)-1] = role.ID
		}
	}

	for _,role := range member.Roles{
		if contains(adminRoleList, role) {
				return true
		}
	}

	return false
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// loadSound attempts to load an encoded sound file from disk.
// Lo dejo para cuando haya que usarlo xdd
func loadSound() error {

	file, err := os.Open("airhorn.dca")
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel.
// Lo dejo para cuando haya que usarlo xdd x2
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)
	
	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}