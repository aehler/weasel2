package common

import (
	"app"
	au "lib/auth"
	"app/form"
	"lib/common"
)

func register (c *app.Context){

	user := c.Get("user").(*au.User)

	newUser := au.RegisterForm{}

	regform := form.New("Регистрация в системе", "register", user.SessionID)

	if err := regform.MapStruct(newUser); err != nil {
		c.RenderError(err)
		return
	}

	if c.IsPost(){

		err := regform.ParseForm(&newUser, c.Request)

		if err != nil {
			c.RenderJSONError(err)
			return
		}

		if _, err := common.RegisterUser(&newUser); err != nil {
			c.RenderJSONError(err)
			return
		} else {

			c.RenderJSON(map[string]interface{}{
				"success" : true,
				"message" : "Вам было направлено письмо на адрес "+newUser.Email+". Для завершения регистрации пройдите по ссылке, указанной в письме.",
			})

		}

		return
	}

	c.RenderHTML("/register.html", map[string]interface{}{
		"form": regform.Context(),
	})

	return

}