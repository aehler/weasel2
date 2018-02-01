# srm

Запихнуть событие для personal_tasks.alerts
```
(business_events.Firer).Fire(business_events.Event{
		Object : p.Object,
		ObjectId : p.Id,
		BusinessEvent : "must_finish_proposals_in_24h",
		EventData : map[string]interface {}{
			"ao" : "Согласование тестовое",
		},
	})
```