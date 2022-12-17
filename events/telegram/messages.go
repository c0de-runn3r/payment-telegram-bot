package telegram

const ( // for reply messages
	msgHello              = "Привіт 👋\nЯ бот для платного доступу до каналу. В мене можна оформити підписку.\n" + msgHelp
	msgHelp               = "Для оформлення підписки нажми на кнопку знизу."
	msgUnknownCommand     = "Невідома команда"
	msgRequestApproved    = "Запит на приєднання схвалено!"
	msgRequestDeclined    = "Запит на приєднання відхилено. Немає активної підписки!"
	msgUpdateSubscription = "Підписка закінчується через 3 дні. Можливо час продовжити? 😏"
	msgSubscriptionEnded  = "Підписка закінчилась. Можливо варто продовжити? 😏"
)

const ( // for bot commands
	cmdStart = "/start"
	cmdHelp  = "/help"

	cmdCurrentSubs = "Показати підписників"
)

const ( // for buttons
	btnBuySubscription = "Купити підписку 💵"

	btnOneMonth    = "1 місяць 🦭"
	btnThreeMonths = "3 місяці 🦓"
	btnSixMonths   = "6 місяців 🐳"
)
