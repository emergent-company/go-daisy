package ui

import "strconv"

// documentDirectionJS returns the inline script body used by
// DocumentDirectionScript.
func documentDirectionJS(dir string) string {
	return "(function(){var target=" + strconv.Quote(dir) + ";if(!target)target='ltr';" +
		"if(document.documentElement.getAttribute('dir')!==target){" +
		"document.documentElement.setAttribute('dir',target);}})();"
}
