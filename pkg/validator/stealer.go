package validator

import "regexp"

func ValStealerLog(logdata string) bool {
	re := regexp.MustCompile(`([a-zA-Z0-9+.-]+://[^:/\s]+(?:/[^:\s]*)?):([^:]+):([^\n]*)`)
	return re.MatchString(logdata)
}
