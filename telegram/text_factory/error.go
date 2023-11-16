package command_msg

import "fmt"

func CreateErrorMsg() string {
	return fmt.Sprint("Упс, у нас технические неполадки. Попробуйте повторить запрос позднее :(")
}
