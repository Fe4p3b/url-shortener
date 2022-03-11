// Package storage defines errors for storages.
package storage

import "errors"

var ErrorNoLinkFound = errors.New("link not found")
var ErrorDuplicateShortlink = errors.New("duplicate short link")

var ErrorMethodIsNotImplemented = errors.New("method is not implemented")
