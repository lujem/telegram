package telegram_test

import (
	"fmt"
	. "github.com/nasermirzaei89/telegram"
	"net/http"
	"net/http/httptest"
	"testing"
)

func except(t *testing.T, actual, excepted interface{}) {
	if fmt.Sprintf("%+v", excepted) != fmt.Sprintf("%+v", actual) {
		t.Errorf("\nexcepted: %+v\nactual:   %+v", excepted, actual)
	}
}

func notExcept(t *testing.T, actual, notExcepted interface{}) {
	if fmt.Sprintf("%+v", notExcepted) == fmt.Sprintf("%+v", actual) {
		t.Errorf("\nnot excepted: %+v\nactual:       %+v", notExcepted, actual)
	}
}

func TestGetMe(t *testing.T) {
	token := "someToken"
	response := []byte(`{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"Test Bot","username":"TestBot"}}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := fmt.Sprintf("/bot%s/getMe", token)
		if r.URL.String() != u {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"ok":false,"error_code":401,"description":"Unauthorized"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	}))

	defer server.Close()

	BaseURL = server.URL

	// success
	bot := New(token)

	res, err := bot.GetMe()
	except(t, err, nil)

	except(t, res.IsOK(), true)
	except(t, res.Error(), nil)
	notExcept(t, res.GetUser(), nil)
	except(t, res.GetUser().ID, 1)
	except(t, res.GetUser().IsBot, true)
	except(t, res.GetUser().FirstName, "Test Bot")
	except(t, res.GetUser().LastName, nil)
	notExcept(t, res.GetUser().Username, nil)
	except(t, *res.GetUser().Username, "TestBot")
	except(t, res.GetUser().LanguageCode, nil)

	// fail
	bot = New("invalidToken")

	res, err = bot.GetMe()
	except(t, err, nil)

	except(t, res.IsOK(), false)
	except(t, res.GetUser(), nil)
	notExcept(t, res.Error(), nil)
	except(t, res.Error().GetErrorCode(), http.StatusUnauthorized)
	except(t, res.Error().GetDescription(), "Unauthorized")
	except(t, res.Error().GetParameters(), nil)
}

func TestSendMessage(t *testing.T) {
	token := "someToken"
	response := []byte(`{"ok":true,"result":[{"update_id":109399605,
  "message":{"message_id":1234,"from":{"id":123456789,"is_bot":false,"first_name":"John","last_name":"Doe","username":"johndoe","language_code":"en"},"chat":{"id":123456789,"first_name":"John","last_name":"Doe","username":"johndoe","type":"private"},"date":1234567890,"text":"Test"}}]}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := fmt.Sprintf("/bot%s/getUpdates", token)
		if r.URL.String() != u {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"ok":false,"error_code":401,"description":"Unauthorized"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	}))

	defer server.Close()

	BaseURL = server.URL

	// success
	bot := New(token)

	res, err := bot.GetUpdates()
	except(t, err, nil)

	except(t, res.IsOK(), true)
	except(t, res.Error(), nil)
	notExcept(t, res.GetUpdates(), nil)
	for _, u := range res.GetUpdates() {
		notExcept(t, u.UpdateID, 0)
	}
	//except(t, res.GetUser().ID, 1)
	//except(t, res.GetUser().IsBot, true)
	//except(t, res.GetUser().FirstName, "Test Bot")
	//except(t, res.GetUser().LastName, nil)
	//notExcept(t, res.GetUser().Username, nil)
	//except(t, *res.GetUser().Username, "TestBot")
	//except(t, res.GetUser().LanguageCode, nil)

	// fail
	bot = New("invalidToken")

	res, err = bot.GetUpdates()
	except(t, err, nil)

	except(t, res.IsOK(), false)
	except(t, len(res.GetUpdates()), 0)
	notExcept(t, res.Error(), nil)
	except(t, res.Error().GetErrorCode(), http.StatusUnauthorized)
	except(t, res.Error().GetDescription(), "Unauthorized")
	except(t, res.Error().GetParameters(), nil)
}