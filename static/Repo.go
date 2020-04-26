// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package static

type Repo interface {
	/**
	Get the public static
	*/
	GetStaticPublicDocument(path string) (string, error)

	/**
	Get the public static
	*/
	GetStaticPrivateDocument(path string) (string, error)
}
