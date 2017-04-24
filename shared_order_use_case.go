package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (c *Channel) exitChannel(number string) {
	c.Lock()
	defer c.Unlock()
	delete(c.Terminals, number)
}

func (c *Channel) openChannel(number string, channelUUID string) *Terminal {
	c.Lock()
	defer c.Unlock()
	t := &Terminal{
		Number: number,
		Sub:    orderClient.redisClient.Subscribe(channelUUID),
	}
	c.Terminals[number] = t

	return t
}

func (c *OrderClient) getChannel(uuid string, merchantUUID string, pin string) (*Channel, error) {
	c.Lock()
	defer c.Unlock()
	ch, ok := c.channels[uuid]
	if !ok {
		return nil, errors.New("channel_id not found")
	}
	if ch.MerchantID != merchantUUID {
		return nil, errors.New("merchant_id not authorized")
	}
	if ch.PIN != pin {
		return nil, errors.New("pin not authorized")
	}
	return ch, nil
}

func (c *OrderClient) createChannel(uuid string, merchantUUD string, number string) *Channel {
	c.Lock()
	defer c.Unlock()
	ch, ok := c.channels[uuid]
	if !ok {
		ch = &Channel{
			UUID:       uuid,
			MerchantID: merchantUUD,
			PIN:        number,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Terminals:  make(map[string]*Terminal),
		}
		c.channels[uuid] = ch
	}
	return ch
}

func shareOrder(w http.ResponseWriter, r *http.Request) {

	var orderUUID string
	number := r.FormValue("logic_number")
	merchantUUID := r.FormValue("merchant_id")

	if orderUUID = mux.Vars(r)["order_id"]; len(strings.TrimSpace(orderUUID)) == 0 {
		respondWithError(w, http.StatusBadRequest, "OrderId not found")
		return
	}

	// order em memoria no redis
	redisOrder, err := findOrderRedis(number, orderUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	log.Printf("<< redis: %s merchantUUID: %s", redisOrder.MerchantID, merchantUUID)

	// definir regra do que é share ou nao
	// err = validateOrder(redisOrder, merchantUUID)
	// if err != nil {
	// 	respondWithError(w, http.StatusConflict, err.Error())
	// 	return
	// }

	log.Printf("<<createChannel:  number: %s orderUUID: %s merchantUUID: %s", number, orderUUID, merchantUUID)

	ch := orderClient.createChannel(orderUUID, merchantUUID, number)

	respondWithJSON(w, http.StatusCreated, ch)
}
