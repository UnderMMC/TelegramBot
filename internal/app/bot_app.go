package app

import (
	"TelegrammBot/internal/domain/entity"
	"TelegrammBot/internal/domain/repository"
	"TelegrammBot/internal/domain/service"
	"context"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

type BotService interface {
	ShowAllRepairment(ctx context.Context) (string, error)
}

type BotApp struct {
	service BotService
}

func NewBotApp() *BotApp {
	return &BotApp{}
}

var filters entity.SearchFilters

func (a *BotApp) Run() {
	var err error
	connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Ошибка поделючения к БД: %v", err)
	}
	defer db.Close()

	repo := repository.NewBotRepository(db)
	serv := service.NewBotService(repo)
	a.service = serv

	// Получение токена из переменной окружения (с .env файла)
	err = godotenv.Load()
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN environment variable not set")
	}

	// Создание бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Error creating bot ", err)
	}

	// Обновление бота
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.Message != nil && update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую, "+update.Message.From.FirstName+"!")
			bot.Send(msg)

			// Кнопки для выбора города
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Челябинск", "Chelyabinsk"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Москва", "Moscow"),
				),
			)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, выберите нужный город:")
			msg.ReplyMarkup = keyboard
			bot.Send(msg)

		} else if update.CallbackQuery != nil {
			callBack := update.CallbackQuery

			if filters.City == "" {
				switch callBack.Data {
				case "Chelyabinsk":
					filters.City = callBack.Data
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы выбрали город: "+filters.City+".\nТеперь выберите услугу:")

					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Ремонт механики", "mechanics"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Ремонт пластика", "plastic"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Перебортовка", "tire"),
						),
					)
					msg.ReplyMarkup = keyboard
					bot.Send(msg)
				default:
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Пожалуйста, выберите один из предложенных городов.")
					bot.Send(msg)
				}

			} else {
				// Обработка выбора услуги
				switch callBack.Data {
				case "mechanics", "plastic", "tire":
					filters.RepairType = callBack.Data

					// Создаем контекст с фильтрами
					ctx := context.WithValue(context.Background(), "searchFilters", filters)

					// Выводим результат
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вот кого удалось найти по вашему запросу:")
					bot.Send(msg)

					// Получение и вывод данных через сервис
					result, err := a.service.ShowAllRepairment(ctx)
					if err != nil {
						log.Println("Не удалось получить данные:", err)
					} else {
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, result)
						bot.Send(msg)
					}

					// Сброс фильтров после завершения
					filters = entity.SearchFilters{}
				}
			}

			// Подтверждаем callback запрос
			// bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Выбор принят"))
		}
	}
}

func commands(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Команды:\n"+
		"/start - вывод списков поставщиков услуг;\n"+
		"/city - выбор города;\n")
	bot.Send(msg)
}
