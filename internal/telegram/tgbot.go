package telegram

import (
	"OrelTeatrBot/internal/domain"
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Store interface {
	UpdateUserActivity(ctx context.Context, user domain.User) (err error)
}

type TgBot struct {
	store Store
	bot   *bot.Bot
}

func NewBot(token string, store Store) (*TgBot, error) {

	tgBot := TgBot{store: store}
	opts := []bot.Option{
		bot.WithDefaultHandler(tgBot.handler),
	}
	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating telegram bot: %w", err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, tgBot.startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/stop", bot.MatchTypeExact, tgBot.stopHandler)
	tgBot.bot = b
	return &tgBot, nil
}

func (t *TgBot) Start(ctx context.Context) {
	t.bot.Start(ctx)
}

func (t *TgBot) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		_, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Я присылаю уведомления о билетах. Команды: /start, /stop"),
		})
		if err != nil {
			slog.Error("error sending message", "err", err)
		}
	}
}

func (t *TgBot) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var username, firstName string
	if update.Message.From != nil {
		username = update.Message.From.Username
		firstName = update.Message.From.FirstName
	}
	user := domain.User{ChatID: update.Message.Chat.ID, Username: username, IsActive: true}
	err := t.store.UpdateUserActivity(ctx, user)
	if err != nil {
		slog.Error("error updating user activity:", "err", err)
		_, err = t.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Возникла техническая ошибка. Просьба повторить позже."),
		})
		if err != nil {
			slog.Error("error sending message", "err", err)
		}
		return
	}
	_, err = t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf("Привет, %s! Вы подписались на получение уведомлений о билетах на сайте "+
			"театра \"Свободное пространство\". Вы будете получать такие уведомления как только появятся билеты на какое-либо представление. "+
			"Желаем успехов с покупкой!",
			firstName),
	})
	if err != nil {
		slog.Error("error sending message", "err", err)
	}
}

func (t *TgBot) stopHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var username string
	if update.Message.From != nil {
		username = update.Message.From.Username
	}
	user := domain.User{ChatID: update.Message.Chat.ID, Username: username, IsActive: false}
	err := t.store.UpdateUserActivity(ctx, user)
	if err != nil {
		slog.Error("error updating user activity:", "err", err)
		_, err = t.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Возникла техническая ошибка. Просьба повторить позже."),
		})
		if err != nil {
			slog.Error("error sending message", "err", err)
		}
		return
	}
	_, err = t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Вы отписались от получения уведомлений( Надеемся на скорое возвращение."),
	})
	if err != nil {
		slog.Error("error sending message", "err", err)
	}
}

func (t *TgBot) SendMessage(ctx context.Context, chatID int64, message string) (int, error) {
	m, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      message,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return 0, fmt.Errorf("error sending message: %w", err)
	}
	return m.ID, nil
}
