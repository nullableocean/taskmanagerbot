package user

import (
	"errors"
	"log"
	"strconv"
	"taskbot/domain"
	"taskbot/pkg/password"
	"taskbot/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	passLength = 12
)

type TelegramUserService struct {
	us    *UserService
	urepo repository.TelegramUserRepository
}

func NewTelegramUserService(us *UserService, repo repository.TelegramUserRepository) *TelegramUserService {
	return &TelegramUserService{
		us:    us,
		urepo: repo,
	}
}

func (s *TelegramUserService) CreateFromUpdate(update tgbotapi.Update) (domain.User, error) {
	if update.Message == nil {
		return domain.User{}, errors.New("message empty")
	}

	chat := update.Message.Chat
	tid := chat.ID
	username := update.ChatMember.From.UserName

	log.Printf("new user. chat_id: %d, uid: %d, username: %s",
		tid,
		update.ChatMember.From.ID,
		update.ChatMember.From.UserName,
	)

	if username == "" {
		username = s.getLabelUsername(tid)
	}

	newUser := domain.User{
		Username:   username,
		Name:       update.ChatMember.From.FirstName,
		Password:   s.getPass(),
		TelegramId: tid,
	}

	u, err := s.us.Save(newUser)

	if err != nil && errors.As(err, repository.ErrUsernameTaken) {
		newUser.Username = s.getLabelUsername(tid)
		return s.us.Save(newUser)
	}

	return u, err
}

func (s *TelegramUserService) FindByTelegramId(tid int64) (domain.User, error) {
	return s.urepo.GetByTelegramId(tid)
}

func (s *TelegramUserService) getPass() string {
	return password.Generate(passLength)
}

func (s *TelegramUserService) getLabelUsername(tid int64) string {
	tidStr := strconv.FormatInt(tid, 10)
	return "telegram_user_number_" + tidStr
}
