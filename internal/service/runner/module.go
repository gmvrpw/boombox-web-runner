package runner

import (
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"gmvr.pw/boombox-web-runner/pkg/model"
)

func authorize(browser *rod.Browser, module *model.Module) error {
	err := browser.SetCookies(
		*&[]*proto.NetworkCookieParam{{Name: module.Auth.Name, Value: module.Auth.Value, Domain: module.Auth.Domain}},
	)
	if err != nil {
		return err
	}

	return nil
}

func play(timecode uint64, page *rod.Page, module *model.Module) error {
	if module.PlayAction.Selector == "" {
		return nil
	}

	play, err := page.Element(module.PlayAction.Selector)
	if err != nil {
		return err
	}

	play.MustClick()
	return nil
}

func now(page *rod.Page, module *model.Module) (uint64, error) {
	timecode, err := extractTime(module.PlaybackTime.Selector, page)
	if err != nil {
		return 0, err
	}

	return uint64(timecode), err
}

func finished(page *rod.Page, module *model.Module) (bool, error) {
	if module.RemainingTime.Selector != "" {
		remaining, err := extractTime(module.RemainingTime.Selector, page)
		if err == nil {
			return remaining < 1, nil
		}
	}

	playback, err := extractTime(module.PlaybackTime.Selector, page)
	if err != nil {
		return false, err
	}

	duration, err := extractTime(module.DurationTime.Selector, page)
	if err != nil {
		return false, err
	}

	return duration-playback < 1, nil
}

func extractTime(selector model.Selector, page *rod.Page) (int64, error) {
	el, err := page.Element(selector)
	if err != nil {
		return 0, err
	}

	text, err := el.Text()
	if err != nil {
		return 0, err
	}

	t, err := time.Parse(time.TimeOnly, text)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}
