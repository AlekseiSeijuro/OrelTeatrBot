package notifier

import (
	"OrelTeatrBot/internal/domain"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type store interface {
	MarkNotificationSent(ctx context.Context, chatID int64, sessionID int64, epoch int, messageID int) error
	MarkNotificationFailure(ctx context.Context, chatID int64, sessionID int64, epoch int) error
	CreateNotifications(ctx context.Context) (int64, error)
	GetNotificationData(ctx context.Context) ([]domain.NotificationData, error)
}
type tgSender interface {
	SendMessage(ctx context.Context, chatID int64, message string) (int, error)
}

type Notifier struct {
	sitePath       string
	notifyInterval time.Duration
	store          store
	tgSender       tgSender
}

func NewNotifier(sitePath string, notifyInterval time.Duration, s store, sender tgSender) *Notifier {
	return &Notifier{
		sitePath:       sitePath,
		notifyInterval: notifyInterval,
		store:          s,
		tgSender:       sender,
	}
}

func (tn *Notifier) Run(ctx context.Context) error {
	ticker := time.NewTicker(tn.notifyInterval)
	defer ticker.Stop()
	for {
		err := tn.doIteration(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			slog.Error("failed to do iteration", "err", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (tn *Notifier) doIteration(ctx context.Context) error {
	newNotificationsCount, err := tn.store.CreateNotifications(ctx)
	if err != nil {
		return fmt.Errorf("create notifications: %w", err)
	}
	slog.Info("New notifications created:", "new notifications", newNotificationsCount)

	notifications, err := tn.store.GetNotificationData(ctx)
	if err != nil {
		return fmt.Errorf("get notification data: %w", err)
	}
	groupedData := groupNotificationData(notifications)
	for _, notification := range groupedData {
		message := tn.buildMessage(notification)
		messageID, err := tn.tgSender.SendMessage(ctx, notification[0].ChatID, message)
		if err != nil {
			for _, n := range notification {
				err2 := tn.store.MarkNotificationFailure(ctx, notification[0].ChatID, n.SessionID, n.AvailEpoch)
				if err2 != nil {
					slog.Error("mark notification failure:", "err", err2)
				}
			}
			slog.Error("send message:", "err", err)
			continue
		}
		for _, n := range notification {
			err = tn.store.MarkNotificationSent(ctx, notification[0].ChatID, n.SessionID, n.AvailEpoch, messageID)
			if err != nil {
				slog.Error("mark notification sent:", "err", err)
			}
		}
	}
	return nil
}

func (ntf *Notifier) buildMessage(data []domain.NotificationData) string {
	var message strings.Builder
	message.WriteString(fmt.Sprintf("Есть доступные билеты на спектакль %s.\n Доступное время:\n", data[0].EventTitle))
	for _, nd := range data {
		message.WriteString(fmt.Sprintf("%s - %s/s%d\n", nd.DateText, ntf.sitePath, nd.SessionID))
	}
	return message.String()
}

// Строки уже отсортированы по ChatID и EventID
func groupNotificationData(notificationData []domain.NotificationData) [][]domain.NotificationData {
	var groups [][]domain.NotificationData
	for i, nd := range notificationData {
		if i > 0 && nd.ChatID == notificationData[i-1].ChatID && nd.EventID == notificationData[i-1].EventID {
			last := len(groups) - 1
			groups[last] = append(groups[last], nd)
			continue
		}
		groups = append(groups, []domain.NotificationData{nd})
	}
	return groups
}
