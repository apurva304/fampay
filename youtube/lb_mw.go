package youtube

import (
	"sync"
	"time"
)

type loadBalancingMW struct {
	ytClient []Service
	counter  int
	mtx      sync.RWMutex
}

func NewLb(apiKeys []string) (r *loadBalancingMW, err error) {
	r = &loadBalancingMW{
		ytClient: make([]Service, len(apiKeys)),
		counter:  0,
	}

	for i, key := range apiKeys {
		client, err := NewService(key)
		if err != nil {
			return nil, err
		}
		r.ytClient[i] = client
	}
	return r, nil
}

func (svc *loadBalancingMW) Search(query string, publishedAfter time.Time) (err error) {
	svc.mtx.RLock()
	index := svc.counter % len(svc.ytClient)
	err = svc.ytClient[index].Search(query, publishedAfter)
	svc.mtx.RUnlock()

	switch err {
	case ErrQuotaExceeded:
		svc.remove(index)
		return
	default:
	}
	svc.incCounter()
	return
}

func (svc *loadBalancingMW) incCounter() {
	svc.mtx.Lock()
	defer svc.mtx.Unlock()

	svc.counter++
}

func (svc *loadBalancingMW) remove(index int) {
	svc.mtx.Lock()
	defer svc.mtx.Unlock()

	var newArr []Service
	for i, yt := range svc.ytClient {
		if index == i {
			continue
		}
		newArr = append(newArr, yt)
	}
	svc.ytClient = newArr
	svc.counter++
}

// func main() {
// 	lastSuccessFetchTime := time.Now().Add(-10 * time.Minute)

// 	tc := time.NewTicker(1 * time.Second)
// 	for {
// 		select {
// 		case <-tc.C:
// 			err = svc.Search("music", lastSuccessFetchTime)
// 			switch err {
// 			case nil:
// 				lastSuccessFetchTime = time.Now()
// 				fmt.Println("success", lastSuccessFetchTime)
// 			case youtube.ErrNotFound:
// 			case youtube.ErrQuotaExceeded:
// 				fmt.Println("limit")
// 			default:
// 				panic(err)
// 			}
// 		}
// 	}
// }
