package preferences

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSettingGroup(t *testing.T) {
	// arrange
	// act
	settingGroup := newSettingGroup()

	// assert
	assert.NotNil(t, settingGroup)
	assert.NotNil(t, settingGroup.Settings)
	assert.NotNil(t, settingGroup.SubGroup)
	assert.Empty(t, settingGroup.Settings)
	assert.Empty(t, settingGroup.SubGroup)
}

func TestCheckSubStructureValid(t *testing.T) {
	// arrange
	settingGroup := &SettingGroup{}

	// act
	settingGroup.checkSubStructureValid()

	// assert
	assert.NotNil(t, settingGroup)
	assert.NotNil(t, settingGroup.Settings)
	assert.NotNil(t, settingGroup.SubGroup)
	assert.Empty(t, settingGroup.Settings)
	assert.Empty(t, settingGroup.SubGroup)
}

func TestGetSubGroup(t *testing.T) {
	testCases := []struct {
		currentSettingGroup *SettingGroup
		expectedSubGroup    *SettingGroup
	}{
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedSubGroup: &SettingGroup{
				Settings: map[string]string{
					"red": "yellow",
				},
				SubGroup: map[string]*SettingGroup{},
			},
		},
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{},
				SubGroup: map[string]*SettingGroup{},
			},
			expectedSubGroup: newSettingGroup(),
		},
	}
	for _, testCase := range testCases {
		// arrange
		// act
		result := testCase.currentSettingGroup.GetSubGroup("green")

		// assert
		assert.Equal(t, testCase.expectedSubGroup, result)
	}
}

func TestGetValueAsString(t *testing.T) {
	id := "blue"

	testCases := []struct {
		currentSettingGroup *SettingGroup
		expectedValue       string
		expectedError       error
	}{
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: errors.New("setting " + id + " not found"),
			expectedValue: "",
		},
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					id:      "123",
					"alpha": "true",
				},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: nil,
			expectedValue: "123",
		},
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					id:      "123",
					"alpha": "true",
				},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"id": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: nil,
			expectedValue: "123",
		},
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					"alpha": "true",
				},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"id": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: errors.New("setting " + id + " not found"),
			expectedValue: "",
		},
	}
	for _, testCase := range testCases {
		// arrange
		// act
		result, err := testCase.currentSettingGroup.GetValueAsString(id)

		// assert
		assert.Equal(t, testCase.expectedValue, result)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestGetValueAsBool(t *testing.T) {
	id := "blue"

	testCases := []struct {
		currentSettingGroup *SettingGroup
		expectedValue       bool
		expectedError       error
	}{
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: errors.New("setting " + id + " not found"),
			expectedValue: false,
		},
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					id:      "123",
					"alpha": "true",
				},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: nil,
			expectedValue: false,
		},
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					id:      "true",
					"alpha": "true",
				},
			},
			expectedError: nil,
			expectedValue: true,
		},
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					id:      "false",
					"alpha": "true",
				},
			},
			expectedError: nil,
			expectedValue: false,
		},
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					id:      "True",
					"alpha": "true",
				},
			},
			expectedError: nil,
			expectedValue: true,
		},
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					id:      "False",
					"alpha": "true",
				},
			},
			expectedError: nil,
			expectedValue: false,
		},
		{
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					"alpha": "true",
				},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"id": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: errors.New("setting " + id + " not found"),
			expectedValue: false,
		},
	}
	for _, testCase := range testCases {
		// arrange
		// act
		result, err := testCase.currentSettingGroup.GetValueAsBool(id)

		// assert
		assert.Equal(t, testCase.expectedValue, result)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestGetSettingAsString(t *testing.T) {
	testCases := []struct {
		tree                []string
		currentSettingGroup *SettingGroup
		expectedValue       string
		expectedError       error
	}{
		{
			tree: []string{"green", "blue"},
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: errors.New("setting blue not found"),
			expectedValue: "",
		},
		{
			tree: []string{"green", "blue"},
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red":  "yellow",
							"blue": "1232",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: nil,
			expectedValue: "1232",
		},
		{
			tree: []string{"blue"},
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					"blue": "1232",
				},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: nil,
			expectedValue: "1232",
		},
		{
			tree: []string{"green", "green", "blue"},
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					"blue": "1232",
				},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red": "yellow",
						},
						SubGroup: map[string]*SettingGroup{
							"green": &SettingGroup{
								Settings: map[string]string{
									"blue": "yellow",
								},
								SubGroup: map[string]*SettingGroup{},
							},
						},
					},
				},
			},
			expectedError: nil,
			expectedValue: "yellow",
		},
		{
			tree: []string{"blue", "green", "yellow"},
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					"blue": "1232",
				},
				SubGroup: map[string]*SettingGroup{},
			},
			expectedError: errors.New("setting yellow not found"),
			expectedValue: "",
		},
	}
	for _, testCase := range testCases {
		// arrange
		// act
		result, err := testCase.currentSettingGroup.GetSettingAsString(testCase.tree)

		// assert
		assert.Equal(t, testCase.expectedValue, result)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestGetSettingAsBool(t *testing.T) {
	testCases := []struct {
		tree                []string
		currentSettingGroup *SettingGroup
		expectedValue       bool
		expectedError       error
	}{
		{
			tree: []string{"green", "blue"},
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: errors.New("setting blue not found"),
			expectedValue: false,
		},
		{
			tree: []string{"green", "blue"},
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red":  "yellow",
							"blue": "true",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: nil,
			expectedValue: true,
		},
		{
			tree: []string{"blue"},
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					"blue": "true",
				},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red": "yellow",
						},
						SubGroup: map[string]*SettingGroup{},
					},
				},
			},
			expectedError: nil,
			expectedValue: true,
		},
		{
			tree: []string{"green", "green", "blue"},
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					"blue": "1232",
				},
				SubGroup: map[string]*SettingGroup{
					"green": &SettingGroup{
						Settings: map[string]string{
							"red": "yellow",
						},
						SubGroup: map[string]*SettingGroup{
							"green": &SettingGroup{
								Settings: map[string]string{
									"blue": "false",
								},
								SubGroup: map[string]*SettingGroup{},
							},
						},
					},
				},
			},
			expectedError: nil,
			expectedValue: false,
		},
		{
			tree: []string{"blue", "green", "yellow"},
			currentSettingGroup: &SettingGroup{
				Settings: map[string]string{
					"blue": "1232",
				},
				SubGroup: map[string]*SettingGroup{},
			},
			expectedError: errors.New("setting yellow not found"),
			expectedValue: false,
		},
	}
	for _, testCase := range testCases {
		// arrange
		// act
		result, err := testCase.currentSettingGroup.GetSettingAsBool(testCase.tree)

		// assert
		assert.Equal(t, testCase.expectedValue, result)
		assert.Equal(t, testCase.expectedError, err)
	}
}
