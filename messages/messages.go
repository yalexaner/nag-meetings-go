package messages

const (
	Start                    = "Я буду отправлять уведомления на встречи STB Daily Meetings. Если ты ещё не авторизован, напиши @yalexaner.\n\nЧтобы подписаться на уведомления, нужно отправить /subscribe. Чтобы отписаться, отправь /unsubscribe"
	Subscribed               = "Подписал на получение уведомлений"
	ErrorSubscribing         = "Не удалось подписаться на уведомления"
	Unsubscribed             = "Больше не буду присылать уведомления"
	ErrorUnsubscribing       = "Не получилось удалить из базы на получение уведомлений"
	NotAuthorized            = "Ты не авторизован. Напиши @yalexaner"
	UnknownCommand           = "Такую команду не знаю. Доступны только /start, /subscribe и /unsubscribe"
	UnknownError             = "Произошла какая-то ошибка"
	AllUsersAuthorized       = "Все пользователи авторизованы"
	GetUnathorizedUsersError = "Не удалось получить список неавторизованных пользователей"
	AuthorizeUserQuestion    = "Авторизовать пользователя?"
	Yes                      = "Да"
	No                       = "Нет"
	AuthorizeUserError       = "Не удалось авторизовать пользователя"
	RemoveUserError          = "Не удалось удалить пользователя из базы"
)
