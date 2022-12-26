package telegram

const ( // for reply messages
	msgHello              = "Привіт! ☺️\nЯ бот для платного доступу на канал про харчування 🤍 Натискай на кнопку знизу, щоб оформити підписку на місяць 👇🏼"
	msgUnknownCommand     = "Невідома команда"
	msgRequestApproved    = "Запит схвалено 💌🗝"
	msgRequestDeclined    = "Запит на приєднання відхилено 😔\n Немає активної підписки 💔"
	msgUpdateSubscription = "Псс.. підписка на місяць закінчується через 3 дні, можливо варто продовжити? 😏"
	msgSubscriptionEnded  = "Підписка закінчилась. Можливо варто продовжити? 😏"
	msgAskQuestion        = "Напиши тут своє запитання і ми спробуємо його вирішити якомога швидше 😊"
	msgThnxForPayment     = "Дякую за оплату 🙌🏻 Тримай свій лінк на приєднання:\n"
	msgQuestionAccepted   = "Ми отримали твоє повідомлення і уже працюємо над вирішенням проблеми 😉"
	msgSendAnswer         = "Напиши відповідь одним повідомленням і я передам її користувачу."
)

const ( // for bot commands
	cmdStart       = "/start"
	cmdHelp        = "/help"
	cmdQuestion    = "Задати питання"
	cmdCurrentSubs = "Показати активні підписки"
)

const ( // for buttons
	btnOneMonth = "Купити підписку на 1 місяць 🪄"
)
