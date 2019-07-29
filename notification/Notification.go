// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package notification

import (
	"time"
)

type Notification struct {
	Priority   int
	Send       time.Time
	Expiration time.Time
	Message    string
	UserID     int
}
