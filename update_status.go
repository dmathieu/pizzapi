package pizzapi

import (
	"log"
	"time"
)

func UpdateStatus(interval time.Duration) chan bool {
	ticker := time.NewTicker(interval)
	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, o := range loadedOrders {
					updateStatus(o)
				}

				log.Printf("count#update_status.run")
			case <-quit:
				ticker.Stop()
				close(quit)
				return
			}
		}
	}()

	return quit
}

func updateStatus(order *Order) {
	now := time.Now()
	var currentType *orderType

	for _, t := range orderTypes {
		if t.Name == order.Status {
			currentType = t
			continue
		}

		if currentType != nil && order.CreatedAt.Add(t.SetAfter).Before(now) {
			order.Status = t.Name
			log.Printf("count#update_status.process order_id=%d status=%s", order.Id, t.Name)
			return
		}
	}
}
