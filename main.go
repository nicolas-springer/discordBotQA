package main

import (
	"discordATLBot/bot"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var questions = []string{}
var answers = make(map[string][]string)

var usersscores = make(map[string]map[string]int)

var lastCommandTime = make(map[string]time.Time)

var serverIDs = []string{"1086067609289629797", "971518341158146139", "753316178348474369", "681712426755817504"} //pba-lumoy-todoc-p

type lastQuestion struct {
	correctIndex int
	msgQuestion  *discordgo.Message
}
type channeluserscore struct {
}

var lastQuestionServerID = make(map[string]lastQuestion)

func init() {
	godotenv.Load()
	for _, id := range serverIDs {
		usersscores[id] = make(map[string]int)
	}
}

func main() {

	questions = bot.LoadQuestions()
	answers = bot.LoadAnswers()

	discord, err := discordgo.New("Bot " + os.Getenv("BOT_KEY"))
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	discord.AddHandler(messageCreate)        // command !pregunta
	discord.AddHandler(onMessageReactionAdd) // reactions by users on msg from !pregunta

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}

	fmt.Println("Bot is running. Press CTRL-C to exit.")

	<-make(chan struct{})
	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Check if the message starts with the command prefix
	if !strings.HasPrefix(m.Content, "!pregunta") {
		return
	}

	found := false
	for _, serverID := range serverIDs {
		if m.GuildID == serverID {
			found = true

			/*lastTime, ok := lastCommandTime[serverID]
			if ok && time.Since(lastTime) < time.Minute { // Prevent command spam
				return
			}
			lastCommandTime[serverID] = time.Now()
			*/
			break
		}
	}
	if !found {
		return
	}

	// Choose a random question
	rand.Seed(time.Now().Unix())
	q := questions[rand.Intn(len(questions))]

	// Get the answers for the question
	a := make([]string, len(answers[q]))
	copy(a, answers[q])

	// Shuffle the answers
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})

	// Create the message with the question and answers
	message := fmt.Sprintf("**Pregunta de Certificación OCA: \n** %s\n", q)
	for i, answer := range a {
		message += fmt.Sprintf("%d. %s\n", i+1, answer)
	}

	correctIndex := 0
	for i, answer := range a {
		if answer == answers[q][0] {
			correctIndex = i
		}
	}
	fmt.Print(correctIndex)

	// Send the message
	msg, err := s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Println("Error sending message: ", err)
		return
	}

	err = s.MessageReactionAdd(m.ChannelID, msg.ID, "1️⃣")
	if err != nil {
		log.Println("Error al agregar la reacción 1️⃣:", err)
	}
	err = s.MessageReactionAdd(m.ChannelID, msg.ID, "2️⃣")
	if err != nil {
		log.Println("Error al agregar la reacción 2️⃣:", err)
	}
	err = s.MessageReactionAdd(m.ChannelID, msg.ID, "3️⃣")
	if err != nil {
		log.Println("Error al agregar la reacción 3️⃣:", err)
	}
	err = s.MessageReactionAdd(m.ChannelID, msg.ID, "4️⃣")
	if err != nil {
		log.Println("Error al agregar la reacción 4️⃣:", err)
	}

	// Handling reactions
	lastQuestionServerID[m.ChannelID] = lastQuestion{correctIndex, msg}

	// Send the correct answer
	time.Sleep(3 * time.Second)

	correctAnswer := fmt.Sprintf("**Respuesta: %v **%s \n", correctIndex+1, answers[q][0])
	_, errrs := s.ChannelMessageSend(m.ChannelID, correctAnswer)
	if errrs != nil {
		fmt.Println("Error sending message: ", errrs)
		return
	}
}

func onMessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	// Verificar que la reacción es a un mensaje que ha sido enviado por el bot
	if m.UserID == s.State.User.ID {
		return
	}

	// Verificar que la reacción es a la ultima pregunta de este servidor
	if m.MessageReaction.MessageID != lastQuestionServerID[m.ChannelID].msgQuestion.ID {
		return
	}

	// Verificar que el usuario no ha respondido ya, con un map de tiempo equivalente a la duracion de la aparicion de la respuesta
	/*for _, usuario := range usuariosCorrectos {
		if usuario.ID == m.UserID {
			return
		}
	}
	*/

	// Verificar si la respuesta es correct

	if lastQuestionServerID[m.ChannelID].correctIndex == 0 {
		switch m.Emoji.Name {
		case "1️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta correcta!")
			usersRigth, _ := s.MessageReactions(m.ChannelID, m.MessageID, "1️⃣", 100, "", "")
			for _, u := range usersRigth {
				println(u.Username, ":", u.ID)
			}
			for _, u := range usersRigth {

				if u.ID != s.State.User.ID {
					fmt.Println("en if: ", m.UserID, s.State.User.ID)
					if usersscores[m.ChannelID] == nil {
						usersscores[m.ChannelID] = make(map[string]int)
					}
					usersscores[m.ChannelID][u.ID] += 10
				}
			}
			data, err := json.MarshalIndent(usersscores, "", "  ")
			if err != nil {
				fmt.Println("Error al convertir a JSON:", err)
				return
			}
			fmt.Println(string(data))

		case "2️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta!")
		case "3️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta")
		case "4️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta")
		}
		return
	}
	if lastQuestionServerID[m.ChannelID].correctIndex == 1 {
		switch m.Emoji.Name {
		case "1️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta!")
		case "2️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta correcta!")
			s.ChannelMessageSend(m.ChannelID, "Respuesta correcta!")
			usersRigth, _ := s.MessageReactions(m.ChannelID, m.MessageID, "1️⃣", 100, "", "")

			for _, u := range usersRigth {
				println(u.Username, ":", u.ID)
			}

			for _, u := range usersRigth {

				if u.ID != s.State.User.ID {
					fmt.Println("en if: ", m.UserID, s.State.User.ID)
					if usersscores[m.ChannelID] == nil {
						usersscores[m.ChannelID] = make(map[string]int)
					}
					usersscores[m.ChannelID][u.ID] += 10
				}
			}

			data, err := json.MarshalIndent(usersscores, "", "  ")
			if err != nil {
				fmt.Println("Error al convertir a JSON:", err)
				return
			}
			fmt.Println(string(data))
		case "3️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta")
		case "4️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta")
		}

		return
	}
	if lastQuestionServerID[m.ChannelID].correctIndex == 2 {
		switch m.Emoji.Name {
		case "1️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta!")
		case "2️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta!")
		case "3️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta correcta")
			s.ChannelMessageSend(m.ChannelID, "Respuesta correcta!")
			usersRigth, _ := s.MessageReactions(m.ChannelID, m.MessageID, "1️⃣", 100, "", "")
			for _, u := range usersRigth {
				println(u.Username, ":", u.ID)
			}
			for _, u := range usersRigth {

				if u.ID != s.State.User.ID {
					fmt.Println("en if: ", m.UserID, s.State.User.ID)
					if usersscores[m.ChannelID] == nil {
						usersscores[m.ChannelID] = make(map[string]int)
					}
					usersscores[m.ChannelID][u.ID] += 10
				}
			}

			data, err := json.MarshalIndent(usersscores, "", "  ")
			if err != nil {
				fmt.Println("Error al convertir a JSON:", err)
				return
			}
			fmt.Println(string(data))
		case "4️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta")
		}
		return
	}
	if lastQuestionServerID[m.ChannelID].correctIndex == 3 {
		switch m.Emoji.Name {
		case "1️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta!")
		case "2️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta!")
		case "3️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta incorrecta")
		case "4️⃣":
			s.ChannelMessageSend(m.ChannelID, "Respuesta correcta")
			s.ChannelMessageSend(m.ChannelID, "Respuesta correcta!")
			usersRigth, _ := s.MessageReactions(m.ChannelID, m.MessageID, "1️⃣", 100, "", "")
			for _, u := range usersRigth {
				println(u.Username, ":", u.ID)
			}
			for _, u := range usersRigth {

				if u.ID != s.State.User.ID {
					fmt.Println("en if: ", m.UserID, s.State.User.ID)
					if usersscores[m.ChannelID] == nil {
						usersscores[m.ChannelID] = make(map[string]int)
					}
					usersscores[m.ChannelID][u.ID] += 10
				}
			}

			data, err := json.MarshalIndent(usersscores, "", "  ")
			if err != nil {
				fmt.Println("Error al convertir a JSON:", err)
				return
			}
			fmt.Println(string(data))
		}
		return
	}

}
