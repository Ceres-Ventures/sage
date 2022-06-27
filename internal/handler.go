package internal

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var (
	match = regexp.MustCompile("^!(help|status)(.*)")
)

func (s *Sage) handle(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	log.Debug().Msg("calling Sage.handle")
	if msg.Author.ID == sess.State.User.ID {
		return
	}

	matches := match.FindAllStringSubmatch(msg.Content, -1)
	if len(matches) > 0 {
		// we got some messages
		for _, m := range matches {
			cmd, args := strings.TrimSpace(m[1]), strings.Split(strings.TrimSpace(m[2]), " ")
			log.Debug().Interface("args", args).Str("cmd", cmd).Msg("processing command")
			_ = sendReaction(s.discordSession, msg, "⚙️")
			switch cmd {
			case "status":
				sendStatus(s, sess, msg)
			case "help":
				sendHelp(sess, msg)
			default:
				_ = sendReaction(s.discordSession, msg, "❓")
				replyToMessage(sess, msg, "No such command")
				sendHelp(sess, msg)
			}
		}
	}
}

// sendReaction will send a given reaction to a given message
// if the author of the message is a bot, no reaction will be added
func sendReaction(s *discordgo.Session, m *discordgo.MessageCreate, reaction string) error {
	log.Debug().Str("reaction", reaction).Msg("Calling sendReaction")
	// do not send reactions to other bots
	if m.Author.Bot {
		return nil
	}

	return s.MessageReactionAdd(m.ChannelID, m.ID, reaction)
}

func removeReaction(s *discordgo.Session, m *discordgo.MessageCreate, reaction string) {
	log.Debug().Str("reaction", reaction).Msg("Calling removeReaction")
	e := s.MessageReactionRemove(m.ChannelID, m.ID, reaction, "@me")
	if e != nil {
		log.Error().Err(e).Str("reaaction", reaction).Msg("Failed to remove reaction")
	}
}

func replyToMessage(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	log.Debug().Msg("Calling replyToMessage")
	_, err := s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
	if err != nil {
		log.Error().Err(err).Msg("Failed to send channel message")
	}
}

func sendError(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	_ = sendReaction(s, m, "❌")

	replyToMessage(s, m, fmt.Sprintf("---> **Sage has experienced a mental fart.** <---\nHere are the details:\n\n\t`%s`\n\n", msg))
}

func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate, msg string) error {
	log.Debug().Msg("Calling sendMessage")
	if m.Author.Bot {
		return nil
	}

	directMessageChannel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(directMessageChannel.ID, msg)
	if err != nil {
		return err
	}
	return nil
}

//go:embed status.md
var status string

//go:embed help.md
var help string

func sendHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Debug().Msg("Calling sendHelp")
	replyToMessage(s, m, help)
}

func sendStatus(sage *Sage, s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Debug().Msg("Calling sendStatus")

	t, err := template.New("status").Parse(status)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse template")
		sendError(s, m, err.Error())
		return
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, sage.blockChainManager.GetChains())
	if err != nil {
		log.Error().Err(err).Msg("failed to execute template")
		sendError(s, m, err.Error())
		return
	}
	replyToMessage(s, m, tpl.String())
}
