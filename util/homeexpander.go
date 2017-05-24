package util

import homedir "github.com/mitchellh/go-homedir"

// HomeExpander is a helper to optionally expand home.
//
// This func is a specialized helper, just intended for use in regloader or
// autoload funcs.
//
// DontExpandHome causes this func to return the original Path. Override is an
// optional configuration value sometimes provided in the config to override
// the dontExpandHome value.
//
// If any error is encountered, home is not expanded.
func HomeExpander(path string, dontExpandHome bool, dontExpandOverride *bool) string {
	if dontExpandOverride != nil {
		if *dontExpandOverride {
			return path
		}
	} else if dontExpandHome {
		return path
	}

	p, err := homedir.Expand(path)
	if err != nil {
		return path
	}

	return p
}
