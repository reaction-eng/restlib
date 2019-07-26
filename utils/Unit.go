// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package utils

type Unit struct {
	//Store the name
	Name string

	//Store the symbol
	Symbol string

	//Store the conversation to si
	toSI float64
}

/**
Store the default name
*/
const UNKNOWN = "UNKNOWN"

//Store a list of constant units
var lengthUnits = map[string]Unit{
	"m": {
		"meter",
		"m",
		1.0,
	},
	"km": {
		"kilometer",
		"km",
		1000.0,
	},
	"cm": {
		"centimeter",
		"cm",
		0.01,
	},
	"mm": {
		"millimeter",
		"mm",
		0.001,
	},
	"ft": {
		"foot",
		"ft",
		0.3048,
	},
	"inch": {
		"inch",
		"in",
		0.0254,
	},
}

/**
Function to look up string from
*/
func GetLengthUnit(unitName string) Unit {

	//get the length based upon the units
	if unit, found := lengthUnits[unitName]; found {
		return unit
	} else {
		return Unit{
			UNKNOWN,
			"m",
			1.0,
		}
	}

}

/**
Convert to SI from the other unit
*/
func (unit *Unit) ToSI(input float64) float64 {
	return input * unit.toSI
}

/**
Convert to SI from the other unit
*/
func (unit *Unit) GetFactorSI() float64 {
	return unit.toSI
}
