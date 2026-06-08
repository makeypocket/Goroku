package inline

import (
	"fmt"
	"regexp"
	"testing"
)

func TestBotIDPatternOriginal(t *testing.T) {
	// Let's define the original botIDPattern
	originalBotIDPattern := `<a class="tm-row tm-row-link" href="/botfather/bot/(\d+)">` +
		`<img class="tm-row-pic tm-row-pic-user" src="[^"]+">` +
		`<div> <div class="tm-row-value">[^<]*</div>` +
		`<div class="tm-row-description">@%s</div> </div></a>`

	// This is the HTML with newlines/tabs/spaces that we might expect
	htmlWithNewlines := `<a class="tm-row tm-row-link" href="/botfather/bot/123456">
    <img class="tm-row-pic tm-row-pic-user" src="https://telegram.org/img/logo.png">
    <div>
        <div class="tm-row-value">Goroku Bot</div>
        <div class="tm-row-description">@Goroku_123456_bot</div>
    </div>
</a>`

	customBot := "Goroku_123456_bot"
	reOriginal := regexp.MustCompile(fmt.Sprintf(originalBotIDPattern, regexp.QuoteMeta(customBot)))
	if reOriginal.MatchString(htmlWithNewlines) {
		t.Log("Original pattern matched html with newlines (unexpected if it's strict on spaces/newlines)")
	} else {
		t.Log("Original pattern did not match html with newlines (as expected)")
	}
}

func TestBotIDPatternNew(t *testing.T) {
	newBotIDPattern := `(?s)<a[^>]*href="/botfather/bot/(\d+)"[^>]*>(?:[^<]|<[^/]|</[^a]|</[aA][^>])*@%s.*?</a>`

	htmlWithNewlines := `<a class="tm-row tm-row-link" href="/botfather/bot/123456">
    <img class="tm-row-pic tm-row-pic-user" src="https://telegram.org/img/logo.png">
    <div>
        <div class="tm-row-value">Goroku Bot</div>
        <div class="tm-row-description">@Goroku_123456_bot</div>
    </div>
</a>`

	customBot := "Goroku_123456_bot"
	reNew := regexp.MustCompile(fmt.Sprintf(newBotIDPattern, regexp.QuoteMeta(customBot)))
	matches := reNew.FindStringSubmatch(htmlWithNewlines)
	if len(matches) > 1 && matches[1] == "123456" {
		t.Logf("New pattern matched successfully! Bot ID: %s", matches[1])
	} else {
		t.Errorf("New pattern failed to match or group 1 was not correct: %v", matches)
	}

	// Let's test it with multiple bots to ensure it doesn't cross anchors
	htmlMultipleBots := `<a class="tm-row tm-row-link" href="/botfather/bot/111111">
    <img class="tm-row-pic" src="1.png">
    <div>
        <div class="tm-row-description">@FirstBot_bot</div>
    </div>
</a>
<a class="tm-row tm-row-link" href="/botfather/bot/222222">
    <img class="tm-row-pic" src="2.png">
    <div>
        <div class="tm-row-description">@SecondBot_bot</div>
    </div>
</a>`

	reSecond := regexp.MustCompile(fmt.Sprintf(newBotIDPattern, regexp.QuoteMeta("SecondBot_bot")))
	matchesSecond := reSecond.FindStringSubmatch(htmlMultipleBots)
	if len(matchesSecond) > 1 && matchesSecond[1] == "222222" {
		t.Logf("New pattern matched second bot correctly! Bot ID: %s", matchesSecond[1])
	} else {
		t.Errorf("New pattern failed to match second bot correctly: %v", matchesSecond)
	}
}
