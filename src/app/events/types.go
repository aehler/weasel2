package events

import (

)

const (
	BEEntityContract string = "contract"
	BEEntityGPZ string = "gpz"
	BEEntityProcedure string = "procedure"
)

var entities = map[string]string{
	BEEntityProcedure : BEEntityProcedure,
	BEEntityGPZ : BEEntityGPZ,
	BEEntityContract : BEEntityContract,
}
