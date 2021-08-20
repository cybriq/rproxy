package main

import (
	_l "github.com/cybriq/log"

	"github.com/cybriq/opts/version"
)

var log = _l.Get(_l.Add(version.PathBase))
