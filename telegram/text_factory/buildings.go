package command_msg

import (
	"arch-telegram-service/models"
	"fmt"
)

func CreateBuildingsMessage(buildings []models.Building) string {
	var msg string
	for _, building := range buildings {
		msg += CreateBuildingMsg(building)
	}

	return msg
}

func CreateBuildingMsg(building models.Building) string {
	text := fmt.Sprintf("📍 *%s*", building.Name) + "\n" +
		fmt.Sprintf("Расстояние: _%.1f_ м.", building.Distance) + "\n" +
		fmt.Sprintf("Адрес: [%s](%s)", building.Address, building.LinkMapAddress) + "\n" +
		fmt.Sprintf("*Описание: *_%s_", building.Description) + "\n" +
		fmt.Sprintf("[Более подробно по ссылке](%s)", building.Link) +
		"\n" + "\n"

	return text
}

func NotFoundBuildingsMsg() string {
	text := "Ничего не найдено поблизости 😔"

	return text
}
