package command_msg

import "fmt"

func CreateVariantRadiusMsg() string {
	return fmt.Sprint("Выберете радиус, по которому будет произведен поиск достопримечательностей:")
}

func CreateSuccessResponseChangeRadiusMsg() string {
	return fmt.Sprint("Радиус применен")
}
