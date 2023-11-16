package models

type MessageType string

const (
	UnknownMessageType MessageType = "unknown"

	StartMessageType MessageType = "start_message"
	HelpMessageType  MessageType = "help_message"

	RadiusVariantMessageType MessageType = "radius_variant_message"
	RadiusChangeMessageType  MessageType = "radius_change_message"

	BuildingsMessageType               MessageType = "buildings_message"
	BuildingsWithPaginationMessageType MessageType = "buildings_with_pagination_message"

	NotFoundMessageType MessageType = "not_found"
	ErrorMessageType    MessageType = "error"
)
