package bot

import (
	"fmt"
	"log"
	"newscomuni/db"
	"newscomuni/rss"
	"strings"
	"sync"
	"time"

	"github.com/NicoNex/echotron/v3"
)

type Bot struct {
	chatID   int64
	database *db.DB
	echotron.API
}

type BotShim struct {
	token    string
	database *db.DB
}

func New(token string, database *db.DB) (b *BotShim) {
	return &BotShim{
		token,
		database,
	}
}

func (b *BotShim) Start(wg *sync.WaitGroup) {
	dsp := echotron.NewDispatcher(b.token, b.Dispatcher)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Bot started.")
		for {
			err := dsp.Poll()
			if err != nil {
				log.Printf(err.Error())
				// In case of issues wait 5 seconds before trying to reconnect.
				time.Sleep(5 * time.Second)
			}
		}
	}()
}

func (b *BotShim) Dispatcher(chatID int64) echotron.Bot {
	return &Bot{
		chatID:   chatID,
		database: b.database,
		API:      echotron.NewAPI(b.token),
	}
}

type Responder func()

func (b *Bot) showSubs() {
	towns, err := b.database.GetSubscriptionsByChatId(b.chatID)
	if err != nil {
		log.Print(err)
		return
	}

	if len(towns) == 0 {
		b.SendMessage(ConstantMessages["internal-error:no-towns"],
			b.chatID,
			nil,
		)
		return
	}

	keybuttons := [][]echotron.InlineKeyboardButton{}

	for _, v := range towns {
		callbackdata := fmt.Sprintf("unsub:%s", v.ID)

		button := []echotron.InlineKeyboardButton{{
			Text:         v.Name,
			CallbackData: callbackdata,
		}}

		keybuttons = append(keybuttons, button)
	}

	keyboard := echotron.InlineKeyboardMarkup{
		InlineKeyboard: keybuttons,
	}

	b.SendMessage("Sei iscritto a queste città. Toccane una per annullare l'iscrizione.", b.chatID,
		&echotron.MessageOptions{
			ReplyMarkup: keyboard,
		},
	)
}

func (b *Bot) showList() {
	towns, err := b.database.ListTowns()

	if err != nil {
		log.Print(err)
		return
	}

	keybuttons := [][]echotron.InlineKeyboardButton{}

	for _, v := range towns {
		callbackdata := fmt.Sprintf("sub:%s", v.ID)
		button := []echotron.InlineKeyboardButton{{
			Text:         v.Name,
			CallbackData: callbackdata,
		}}
		keybuttons = append(keybuttons, button)
	}

	keyboard := echotron.InlineKeyboardMarkup{
		InlineKeyboard: keybuttons,
	}

	b.SendMessage("Ecco la lista delle città disponibili. Toccane una per iscriverti. Se vuoi suggerire una città da aggiungere,invia un messaggio a @eliseomartelli.", b.chatID,
		&echotron.MessageOptions{
			ReplyMarkup: keyboard,
		},
	)
}

// Update implements echotron.Bot.
func (b *Bot) Update(update *echotron.Update) {
	var functionalMessages map[string]Responder = map[string]Responder{
		"/lista":      b.showList,
		"/iscrizioni": b.showSubs,
	}

	if update.Message != nil {
		log.Printf("Bot received an update. Update type: %s", "message")
		if val, ok := ConstantMessages[update.Message.Text]; ok {
			b.SendMessage(val, b.chatID, nil)
		} else if val, ok := functionalMessages[update.Message.Text]; ok {
			val()
		} else {
			b.SendMessage("Comando non riconosciuto. Riprova.", b.chatID, nil)
		}
	}

	if update.CallbackQuery != nil {
		log.Printf("Bot received an update. Update type: %s", "callback query")
		if strings.HasPrefix(update.CallbackQuery.Data, "sub:") {
			b.subscribe(update)
		}
		if strings.HasPrefix(update.CallbackQuery.Data, "unsub:") {
			b.usubscribe(update)
		}
		b.AnswerCallbackQuery(update.CallbackQuery.ID, nil)
	}
}

func (b *Bot) usubscribe(update *echotron.Update) {
	townId := strings.TrimPrefix(update.CallbackQuery.Data, "unsub:")
	town, err := b.database.GetTownByID(townId)
	if err != nil {
		log.Print(err)
		return
	}
	if err := b.database.RemoveSubscription(b.chatID, townId); err != nil {
		log.Print(err)
		return
	} else {
		msg := fmt.Sprintf("✅ Hai annullato l'iscrizione a: %s", town.Name)
		b.SendMessage(msg, b.chatID, nil)
		b.showSubs()
	}
}

func (b *Bot) subscribe(update *echotron.Update) {
	townId := strings.TrimPrefix(update.CallbackQuery.Data, "sub:")
	town, err := b.database.GetTownByID(townId)
	if err != nil {
		log.Print(err)
		return
	}
	if err := b.database.AddSubscription(b.chatID, town.ID); err != nil {
		log.Print(err)
	}
	msg := fmt.Sprintf("✅ Sei iscritto agli aggiornamenti di: %s", town.Name)
	b.SendMessage(msg, b.chatID, nil)
	b.showSubs()
}

func (b *BotShim) SendUpdates(parser *rss.RSS, town db.Town) {
	items, lastUpdated, err := parser.Parse(
		town.FeedUrl,
		town.LastUpdated,
	)
	if err != nil {
		log.Print(err)
		return
	}
	subscribers, err := b.database.GetSubscriberChatId(town.ID)
	if err != nil {
		log.Print(err)
		return
	}
	for _, subscriber := range subscribers {
		for _, item := range items {
			msg := fmt.Sprintf(
				"<b><i>%s</i></b>\n\n<b>%s</b>\n\n%s",
				town.Name,
				item.Title,
				item.Link,
			)
			echotron.NewAPI(b.token).SendMessage(
				msg,
				subscriber,
				&echotron.MessageOptions{
					ParseMode: "HTML",
				},
			)
		}
	}
	err = b.database.UpdateLastUpdated(town.ID, lastUpdated)
	if err != nil {
		log.Print(err)
		return
	}
}

func (b *BotShim) Scanner(wg *sync.WaitGroup, database *db.DB, rss *rss.RSS) {
	wg.Add(1)
	go func() {
		log.Printf("Scanner started.")
		defer wg.Done()
		for {
			towns, err := database.GetSubscribedTowns()
			if err != nil {
				log.Print(err)
				return
			}
			for _, town := range towns {
				wg.Add(1)
				go func(town db.Town) {
					defer wg.Done()
					b.SendUpdates(rss, town)
				}(town)
			}
			time.Sleep(10 * time.Minute)
		}
	}()
}
