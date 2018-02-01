package registry

import (
	"sync"
	"log"
)

var defaultId string = "default-id::"
var mu = sync.Mutex{}

type RegChan map[string]map[string]chan interface{}

func NewRegChan() RegChan {

	defaultReg := make(map[string]map[string]chan interface{})

	defaultReg[defaultId] = make(map[string]chan interface{})

	return RegChan(defaultReg)
}

func (r RegChan) put(key, id string, v interface{}) {

	mu.Lock()
	defer mu.Unlock()

	if _, ok := r[id][key]; ok {
		r[id][key] <- v
	}
}

//Puts to a default collection
func (r RegChan) Put(key string, v interface{}) {

	r.put(key, defaultId, v)

}

//Puts to all existing collections
func (r RegChan) PutAll(key string, v interface{}){

	for _, id := range r.ListCollections() {
		r.put(key, id, v)
	}

}

//Puts to a specified collection
func (r RegChan) PutCollection(key, id string, v interface{}) {

	r.put(key, id, v)

}

func (r RegChan) register(key, id string) {

	mu.Lock()
	defer mu.Unlock()

	log.Println("registering", id, key)

	if _, ok := r[id][key]; !ok {
		r[id][key] = make(chan interface{})
	}

}

//Register dfault channel
func (r RegChan) Register(key string) {

	r.register(key, defaultId)

}

//Register a channel in collection. If a collection does not exist, it will not make a new collection, use RegisterCollection
func (r RegChan) RegisterInCollection(key, id string) {

	if id == defaultId {
		panic("RegChan.RegisterCollection: cannot use specified id")
	}

	r.register(key, id)

}

//Register a collection of named channels
func (r RegChan) RegisterCollection(id string) {

	mu.Lock()
	defer mu.Unlock()

	if _, ok := r[id]; !ok {
		r[id] = make(map[string]chan interface{})
	}

}

//List all registered collections
func (r RegChan) ListCollections() (res []string) {

	for k, _ := range r {

		res = append(res, k)

	}

	return

}

//List all keys in collection
func (r RegChan) ListKeys(id string) (res []string) {

	if l, ok := r[id]; ok {

		for k, _ := range l {

			res = append(res, k)

		}

	}

	return

}

//Get all channels from collection
func (r RegChan) Collection(id string) (map[string]chan interface{}) {

	mu.Lock()
	defer mu.Unlock()

	if l, ok := r[id]; ok {

		return l

	}

	return map[string]chan interface{}{}

}

//Drop collection
func (r RegChan) UnsetCollection(id string) {

	mu.Lock()
	defer mu.Unlock()

	delete(r, id)

}

//Drop a key in collection. If collection remains empty, it will drop it as well
func (r RegChan) Unset(key, id string) {

	mu.Lock()
	defer mu.Unlock()

	if _, ok := r[id]; ok {

		delete(r[id], key)

		if len(r[id]) == 0 {

			delete(r, id)

		}

	}

}

//Returns default channel
func (r RegChan) Get(key string) chan interface{} {

	mu.Lock()
	defer mu.Unlock()

	if _, ok := r[defaultId][key]; ok {
		return r[defaultId][key]
	}

	return nil
}