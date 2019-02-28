package manager

func init(){
	queues = make(map[string]*Queue)
	go mgc()
}