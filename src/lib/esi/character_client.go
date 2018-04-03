package esi

import (
	"lib/esi/response"
	"lib/esi/request"
	"fmt"
)

func (c *Client) CharacterSkills(cid uint, t string) (response.CharacterSkills, error) {

	res := response.CharacterSkills{}

	if err := c.rpc.Send(request.CharacterSkills{CharacterID: cid, Token: t}, &res); err != nil {

		fmt.Println(err.Error())

		return res, err

	}

	return res, nil

}