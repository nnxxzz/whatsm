package consts

const (
	DbDialectDefault   = "sqlite3"
	DbAddressDefault   = "file:whatsmeow.db?_foreign_keys=on"
	DbDialectConfigKey = "whatsmeow.db.dialect"
	DbAddressConfigKey = "whatsmeow.db.address"

	AutoMarkMessageKey = "whatsmeow.autoMarkMessage"

	PlatformDefault       = "Linux"
	BusinessNameDefault   = "Chrome"
	PushNameDefault       = "Faraway"
	PlatformConfigKey     = "whatsmeow.client.platform"
	BusinessNameConfigKey = "whatsmeow.client.businessName"
	PushNameConfigKey     = "whatsmeow.client.pushName"

	ClientDisplayNameDefault   = "Chrome (Linux)"
	ClientDisplayNameConfigKey = "whatsmeow.clientDisplayName"

	MaxUserDefault   = 200
	MaxUserConfigKey = "whatsmeow.maxUser"

	UploadFileImage = 1
	UploadFileVideo = 2
	UploadFileAudio = 3

	NotifyHostDefault = ""
	NotifyPathDefault = ""
)
