package youtube

import (
	"fampay/domain"
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

func (svc *loadBalancingMW) Search(query string, publishedAfter time.Time) (vidoes []domain.Video, err error) {
	svc.mtx.RLock()
	if len(svc.ytClient) < 1 {
		//limit exceeded for all the provided keys
		err = ErrQuotaExceeded
		return
	}

	index := svc.counter % len(svc.ytClient)
	vidoes, err = svc.ytClient[index].Search(query, publishedAfter)
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
