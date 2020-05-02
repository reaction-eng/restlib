package preferences_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/reaction-eng/restlib/preferences"
	"github.com/stretchr/testify/assert"
)

func TestLoadOptionsGroup(t *testing.T) {
	// arrange
	jsonString := `{
		  "id": "root",
		  "name": "Den Preferences",
		  "description": "Use this section to set and update CAWS Den Preferences",
		  "options": [
			{
			  "id":"showIntroOnLaunch",
			  "hidden": true,
			  "type": "bool",
			  "defaultValue": "true"
			},
			{
			  "id":"siteAgreement",
			  "hidden": true,
			  "type": "bool",
			  "defaultValue": "false"
			}
		  ],
		
		  "subgroups": [
			{
			  "id": "communication",
			  "name":"Communication",
			  "Description": "Control how the CAWS Den communicates with you.",
			  "Options":[
				{
				  "id": "emailNews",
				  "name":"Email News & Updates",
				  "description": "Should the CAWS Den email you News & Updates when posted?",
				  "type": "bool",
				  "defaultValue": "true"
				},
				{
				  "id": "inneedEmail",
				  "name":"In-Need of Foster Email Frequency",
				  "description": "How often should the CAWS Den email you about animals in-need of foster?",
				  "type": "string",
				  "defaultValue": "weekly",
				  "selection": ["daily", "weekly", "monthly", "never"]
				},
				{
				  "id": "inneedSpecies",
				  "name":"In-Need of Foster Email Species",
				  "description": "Would you like to receive in-need emails about Cats, Dogs, or both?",
				  "type": "string",
				  "defaultValue": "Cat,Dog",
				  "selection": ["Cat,Dog", "Dog", "Cat"]
				}
			  ]
		  }
		  ]
		}`

	file, err := ioutil.TempFile(os.TempDir(), "optionTest-")
	assert.Nil(t, err)
	defer os.Remove(file.Name())

	file.WriteString(jsonString)
	file.Close()

	// act
	optionsGroup, err := preferences.LoadOptionsGroup(file.Name())

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, optionsGroup)
	assert.Equal(t, "root", optionsGroup.Id)
	assert.Equal(t, "Use this section to set and update CAWS Den Preferences", optionsGroup.Description)
	assert.Equal(t, "Den Preferences", optionsGroup.Name)

	assert.Equal(t, 2, len(optionsGroup.Options))
	assert.Equal(t, "", optionsGroup.Options[0].Name)
	assert.Equal(t, "", optionsGroup.Options[0].Description)
	assert.Equal(t, "showIntroOnLaunch", optionsGroup.Options[0].Id)
	assert.Equal(t, "true", optionsGroup.Options[0].DefaultValue)
	assert.Equal(t, true, optionsGroup.Options[0].Hidden)
	assert.Equal(t, float64(0), optionsGroup.Options[0].MaxValue)
	assert.Equal(t, float64(0), optionsGroup.Options[0].MinValue)
	assert.Nil(t, optionsGroup.Options[0].Selection)
	assert.Equal(t, preferences.Bool, optionsGroup.Options[0].Type)

	assert.Equal(t, "", optionsGroup.Options[1].Name)
	assert.Equal(t, "", optionsGroup.Options[1].Description)
	assert.Equal(t, "siteAgreement", optionsGroup.Options[1].Id)
	assert.Equal(t, "false", optionsGroup.Options[1].DefaultValue)
	assert.Equal(t, true, optionsGroup.Options[1].Hidden)
	assert.Equal(t, float64(0), optionsGroup.Options[1].MaxValue)
	assert.Equal(t, float64(0), optionsGroup.Options[1].MinValue)
	assert.Nil(t, optionsGroup.Options[1].Selection)
	assert.Equal(t, preferences.Bool, optionsGroup.Options[1].Type)

	assert.Equal(t, 1, len(optionsGroup.SubGroups))
	assert.Equal(t, "communication", optionsGroup.SubGroups[0].Id)
	assert.Equal(t, "Communication", optionsGroup.SubGroups[0].Name)
	assert.Equal(t, "Control how the CAWS Den communicates with you.", optionsGroup.SubGroups[0].Description)
	assert.Nil(t, optionsGroup.SubGroups[0].SubGroups)

	assert.Equal(t, 3, len(optionsGroup.SubGroups[0].Options))
	assert.Equal(t, "Email News & Updates", optionsGroup.SubGroups[0].Options[0].Name)
	assert.Equal(t, "Should the CAWS Den email you News & Updates when posted?", optionsGroup.SubGroups[0].Options[0].Description)
	assert.Equal(t, "emailNews", optionsGroup.SubGroups[0].Options[0].Id)
	assert.Equal(t, "true", optionsGroup.SubGroups[0].Options[0].DefaultValue)
	assert.Equal(t, false, optionsGroup.SubGroups[0].Options[0].Hidden)
	assert.Equal(t, float64(0), optionsGroup.SubGroups[0].Options[0].MaxValue)
	assert.Equal(t, float64(0), optionsGroup.SubGroups[0].Options[0].MinValue)
	assert.Nil(t, optionsGroup.SubGroups[0].Options[0].Selection)
	assert.Equal(t, preferences.Bool, optionsGroup.SubGroups[0].Options[0].Type)

	assert.Equal(t, "In-Need of Foster Email Frequency", optionsGroup.SubGroups[0].Options[1].Name)
	assert.Equal(t, "How often should the CAWS Den email you about animals in-need of foster?", optionsGroup.SubGroups[0].Options[1].Description)
	assert.Equal(t, "inneedEmail", optionsGroup.SubGroups[0].Options[1].Id)
	assert.Equal(t, "weekly", optionsGroup.SubGroups[0].Options[1].DefaultValue)
	assert.Equal(t, false, optionsGroup.SubGroups[0].Options[1].Hidden)
	assert.Equal(t, float64(0), optionsGroup.SubGroups[0].Options[1].MaxValue)
	assert.Equal(t, float64(0), optionsGroup.SubGroups[0].Options[1].MinValue)
	assert.Equal(t, []string{"daily", "weekly", "monthly", "never"}, optionsGroup.SubGroups[0].Options[1].Selection)
	assert.Equal(t, preferences.String, optionsGroup.SubGroups[0].Options[1].Type)

	assert.Equal(t, "In-Need of Foster Email Species", optionsGroup.SubGroups[0].Options[2].Name)
	assert.Equal(t, "Would you like to receive in-need emails about Cats, Dogs, or both?", optionsGroup.SubGroups[0].Options[2].Description)
	assert.Equal(t, "inneedSpecies", optionsGroup.SubGroups[0].Options[2].Id)
	assert.Equal(t, "Cat,Dog", optionsGroup.SubGroups[0].Options[2].DefaultValue)
	assert.Equal(t, false, optionsGroup.SubGroups[0].Options[2].Hidden)
	assert.Equal(t, float64(0), optionsGroup.SubGroups[0].Options[2].MaxValue)
	assert.Equal(t, float64(0), optionsGroup.SubGroups[0].Options[2].MinValue)
	assert.Equal(t, []string{"Cat,Dog", "Dog", "Cat"}, optionsGroup.SubGroups[0].Options[2].Selection)
	assert.Equal(t, preferences.String, optionsGroup.SubGroups[0].Options[2].Type)
}

func TestLoadOptionsGroup_ThrowsErrorForMissingFile(t *testing.T) {
	// arrange
	file := "fake/file.json"
	// act
	optionsGroup, err := preferences.LoadOptionsGroup(file)

	// assert
	assert.NotNil(t, err)
	assert.Nil(t, optionsGroup)
}

func TestLoadOptionsGroup_ThrowsErrorForBadJson(t *testing.T) {
	// arrange
	jsonString := `{
		  "id": "root",
		  "name": "Den Preferences",
		  "description": "Use this section to set and update CAWS Den Preferences",
		  "options": [
			}`
	file, err := ioutil.TempFile(os.TempDir(), "optionTest-")
	assert.Nil(t, err)
	defer os.Remove(file.Name())

	file.WriteString(jsonString)
	file.Close()

	// act
	optionsGroup, err := preferences.LoadOptionsGroup(file.Name())

	// assert
	assert.NotNil(t, err)
	assert.Nil(t, optionsGroup)
}
