package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/AdrianOrlow/links-api/app/model"
	"github.com/AdrianOrlow/links-api/app/utils"
	"github.com/AdrianOrlow/links-api/config"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const oauthStateCookieName = "oauthstate"

var googleOauthConfig *oauth2.Config

func InitializeAuth(config *config.Config) {
	domain, exists := os.LookupEnv("GOOGLE_OAUTH_REDIRECT_DOMAIN")
	if !exists {
		log.Fatal("ENV VARIABLE DOESN'T EXISTS: GOOGLE_OAUTH_REDIRECT_DOMAIN")
	}

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  domain + "oauth/google/callback",
		ClientID:     config.GoogleOauthConfig.ClientID,
		ClientSecret: config.GoogleOauthConfig.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

func HandleGoogleLogin(_ *gorm.DB, w http.ResponseWriter, r *http.Request) {
	stateOauthCookie, err := generateStateOauthCookie(w)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	url := googleOauthConfig.AuthCodeURL(stateOauthCookie)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(_ *gorm.DB, w http.ResponseWriter, r *http.Request) {
	oauthState, err := r.Cookie(oauthStateCookieName)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	user, err := getUserInfo(r.FormValue("state"), r.FormValue("code"), oauthState)
	if err != nil {
		respondError(w, http.StatusTemporaryRedirect, err)
		return
	}

	token, err := utils.CreateLoginJWT(user.Email)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err)
		return
	}

	script := `
        <script>
        const receiveMessage = (event) => {
          const trustedOrigins = ["http://localhost:3000", "https://l.orlow.me"];
          if (!trustedOrigins.includes(event.origin)) {
            return;
          }
          event.source.postMessage(
            {
              token: "` + token.Token  + `",
              source: "api"
            },
            event.origin
          );
        }
		
        window.addEventListener("message", receiveMessage);
        </script>
	`

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(script))
}

func generateStateOauthCookie(w http.ResponseWriter) (string, error) {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	_, err := rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: oauthStateCookieName, Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state, err
}

func getUserInfo(state string, code string, oauthState *http.Cookie) (*model.User, error) {
	user := &model.User{}

	if state != oauthState.Value {
		return user, fmt.Errorf("InvalidOauthState")
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return user, err
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return user, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(contents, user)

	return user, nil
}