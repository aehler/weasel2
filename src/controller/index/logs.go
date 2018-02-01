package index

import (
	"app"
	"lib/logs"
	"strconv"
	"encoding/json"
	"fmt"
)

func Logs(c *app.Context) {

	var limit uint = 20

	service := c.Params.ByName("sname")
	offset := c.GetParam("offset")

	if offset == "" {
		offset = "0"
	}

	o, err := strconv.Atoi(offset)
	if err != nil {

		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	logEntries, err := logs.Logs.GetServiceLogs(service, limit, uint(o))
	if err != nil {

		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	c.RenderHTML("/log-list.html", map[string]interface {} {
		"service": service,
		"logs"   : logEntries,
		"limit"  : limit,
		"offset" : o,
	})

}

func LogDetails(c *app.Context) {

	id := c.Params.ByName("logId")

	logEntry, err := logs.Logs.GetServiceLogEntry(id)
	if err != nil {

		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return
	}

	var d []interface{}
	s := logEntry.Details

	if err := json.Unmarshal([]byte(logEntry.Details), &d); err == nil {

		s = ""

		for _, v := range d {

			s = fmt.Sprintf("%s\n\n%v", s, v)

		}

	}


	c.RenderHTML("/log-entry.html", map[string]interface {} {
		"entry"   : logEntry,
		"details" : s,
	})
}