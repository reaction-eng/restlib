// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package routing

import (
	"net/http"
)

type Route struct {
	Name           string
	Method         string
	Pattern        string
	HandlerFunc    http.HandlerFunc
	Public         bool
	ReqPermissions []string
}
