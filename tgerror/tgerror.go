package tgerror

type TelegramError string

func (te TelegramError) Error() string {
	return string(te)
}
