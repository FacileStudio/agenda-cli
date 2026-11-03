package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/FacileStudio/porte/devicelogin"
	"github.com/FacileStudio/porte/loopback"
)

// obtainToken runs the device login, falling back to the loopback listener
// when useListener asks for it or the server has no device exchange.
func obtainToken(api *client.Client, serverURL string, noBrowser, useListener bool) (string, error) {
	if !useListener {
		session, err := devicelogin.Login(context.Background(), devicelogin.Options{
			AppOrigin:   serverURL,
			OpenBrowser: func(link string) bool { return !noBrowser && loopback.OpenBrowser(link) },
		})
		fallback, err := deviceOutcome(err, serverURL)
		if !fallback {
			return session.Token, err
		}
		ui.PrintWarn("This server has no device login yet, using the browser listener.")
	}
	authCode, err := obtainAuthCode(serverURL, noBrowser)
	if err != nil {
		return "", err
	}
	token, err := exchangeToken(api, authCode)
	if err != nil {
		return "", err
	}
	return *token, nil
}

// deviceOutcome reads the device login's error: whether to fall back to the
// listener, and the error to report otherwise.
func deviceOutcome(err error, serverURL string) (bool, error) {
	switch {
	case errors.Is(err, devicelogin.ErrExchangeNotMounted):
		return true, nil
	case errors.Is(err, devicelogin.ErrAccountNotLinked):
		site := strings.TrimSuffix(serverURL, "/api")
		return false, fmt.Errorf("sign in once at %s in a browser, then run login again", site)
	}
	return false, err
}
