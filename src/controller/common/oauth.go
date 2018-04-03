package common

import (
	"app"
	"net/http"
	"bytes"
	"encoding/base64"
	"encoding/json"
	au "lib/auth"
	"errors"
	"strconv"
	"app/registry"
	"middleware/auth"
	"lib/esi"
)

var key = "83jfj39_aaHR#0&hh3SL"

func oauthCallback (c *app.Context){

	code := c.GetUrlParam("code")

	user := c.Get("user").(*au.User)

	//ssid, err := crypto.DecryptB64(c.GetUrlParam("state"), key)
	//if err != nil {
	//
	//	c.RenderHTML("/errors/500.html", map[string]interface {} {
	//		"Error" : err.Error(),
	//	})
	//
	//	c.Stop()
	//
	//	return
	//
	//}

	client := &http.Client{}

	reqBody := bytes.NewBuffer([]byte("grant_type=authorization_code&code="+code))

	req, err := http.NewRequest("POST", "https://login.eveonline.com/oauth/token", reqBody)

	if err != nil {

		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	auth_app := []byte("bbacb081bc184538b1c8aa360036bd82:8lZeUzm8eNmDc0OBzzLx0iXLEc2rpctP3s3BZvjt")


	enc := base64.URLEncoding.EncodeToString(auth_app)

	req.Header.Add("Authorization", "Basic "+enc)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "login.eveonline.com")

	if resp, err := client.Do(req); err == nil {

		result := struct{
			Access_token string
			Expires_in uint
			Token_type string
			Refresh_token string
		}{}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {

			c.RenderHTML("/errors/500.html", map[string]interface {} {
				"Error" : "Eve server responded with error code "+strconv.Itoa(int(resp.StatusCode))+" "+resp.Status+" while authorizing",
			})

			c.Stop()

			return

		}

		jd := json.NewDecoder(resp.Body)

		if errjd := jd.Decode(&result); errjd != nil {

			c.RenderHTML("/errors/500.html", map[string]interface {} {
				"Error" : errjd.Error(),
			})

			c.Stop()

			return

		}

		user.OAuthRToken = result.Refresh_token
		user.OAuthToken = result.Access_token

		//Obtaining Character ID
		if errc := charID(user); errc != nil {

			c.RenderHTML("/errors/500.html", map[string]interface {} {
				"Error" : errc.Error(),
			})

			c.Stop()

			return

		}

	} else {

		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	registry.Registry.Session.Replace(user.SessionID, auth.Auth{User: user, SSID: user.SessionID})

	skills, err := esi.MC.CharacterSkills(user.UserID, user.OAuthToken)
	if err != nil {

		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return

	}

	registry.Registry.Session.Upsert(user.SessionID+"_skills", skills)

	c.Redirect("/")

}

func charID(user *au.User) error {

	client := &http.Client{}

	reqBody := bytes.NewBuffer([]byte{})

	req, err := http.NewRequest("GET", "https://login.eveonline.com/oauth/verify", reqBody)

	if err != nil {

		return err
	}

	req.Header.Add("User-Agent", "bbacb081bc184538b1c8aa360036bd82")
	req.Header.Add("Authorization", "Bearer "+user.OAuthToken)
	req.Header.Add("Host", "login.eveonline.com")

	if resp, err := client.Do(req); err == nil {

		result := struct{
			CharacterID uint
			CharacterName string
			ExpiresOn string
			Scopes string
			TokenType string
			CharacterOwnerHash string
		}{}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {

			return errors.New("Eve server responded with error "+resp.Status+" while getting character data")

		}

		jd := json.NewDecoder(resp.Body)

		if errjd := jd.Decode(&result); errjd != nil {

			return errjd

		}

		user.UserID = result.CharacterID
		user.UserLastName = result.CharacterName

	} else {

		return err

	}

	return nil

}