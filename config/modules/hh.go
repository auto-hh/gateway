package modules

import "gateway/pkg/config"

const (
	HHClientId     config.ConfigKey = "CLIENT_ID"
	HHClientSecret config.ConfigKey = "CLIENT_SECRET"
	HHAppName      config.ConfigKey = "APP_NAME"
	HHAppVersion   config.ConfigKey = "APP_VERSION"
	HHRedirectUri  config.ConfigKey = "REDIRECT_URI"
	HHDevContact   config.ConfigKey = "DEV_CONTACT"
	HHRawUrl       config.ConfigKey = "RAW_URL"

	DefaultAppName    string = "MyAPP"
	DefaultDevContact string = "dev@mail.ru"
	DefaultAppVersion string = "1.0.0"
)

type HHConfig struct {
	clientId     string
	clientSecret string
	appName      string
	appVersion   string
	redirectUri  string
	devContact   string
	rawUrl       string
}

func NewHHConfig() *HHConfig {
	clientId := HHClientId.MustGet()
	clientSecret := HHClientSecret.MustGet()
	redirectUri := HHRedirectUri.MustGet()
	devContact := HHDevContact.Get(DefaultDevContact)
	rawUrl := HHRawUrl.MustGet()
	appName := HHAppName.Get(DefaultAppName)
	appVersion := HHAppVersion.Get(DefaultAppVersion)

	return &HHConfig{
		clientId:     clientId,
		clientSecret: clientSecret,
		appName:      appName,
		appVersion:   appVersion,
		redirectUri:  redirectUri,
		devContact:   devContact,
		rawUrl:       rawUrl,
	}
}

func (hhConfig *HHConfig) GetAppName() string {
	return hhConfig.appName
}

func (hhConfig *HHConfig) GetAppVersion() string {
	return hhConfig.appVersion
}
func (hhConfig *HHConfig) GetRedirectUri() string {
	return hhConfig.redirectUri
}

func (hhConfig *HHConfig) GetDevContact() string {
	return hhConfig.devContact
}

func (hhConfig *HHConfig) GetRawUrl() string {
	return hhConfig.rawUrl
}

func (hhConfig *HHConfig) GetClientId() string {
	return hhConfig.clientId
}

func (hhConfig *HHConfig) GetClientSecret() string {
	return hhConfig.clientSecret
}
