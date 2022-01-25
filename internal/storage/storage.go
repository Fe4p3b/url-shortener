package storage

import "errors"

var ErrorNoLinkFound = errors.New("not found")
var ErrorDuplicateShortlink = errors.New("no such link")

var ErrorMethodIsNotImplemented = errors.New("method has not been implemented")
