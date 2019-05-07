package preferences

//Get the setting group
type Preferences struct {
	//And the value
	Settings *SettingGroup `json:"settings"`

	//We can also old other groups
	Options *OptionGroup `json:"options"`
}
