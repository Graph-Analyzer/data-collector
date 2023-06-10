package input

import (
	"fmt"
	"graph-analyzer/data-collector/input/jalapeno"
	"graph-analyzer/data-collector/repository"
	"strconv"
	"sync"

	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	log "github.com/sirupsen/logrus"
)

type Jalapeno struct {
	Host             string
	RequestPort      int64
	SubscriptionPort int64
}

func (s *Jalapeno) read(repo repository.GraphRepository, networkName string) {
	log.Infoln("Using Jalapeno")

	requestURL := fmt.Sprintf("%s:%s", s.Host, strconv.FormatInt(s.RequestPort, 10))
	subscriptionURL := fmt.Sprintf("%s:%s", s.Host, strconv.FormatInt(s.SubscriptionPort, 10))

	wg := new(sync.WaitGroup)
	wg.Add(2)

	nodeChannel := make(chan *jagw.LsNodeEvent)
	edgeChannel := make(chan jalapeno.EdgeEvent)

	go jalapeno.SubscribeService(wg, nodeChannel, edgeChannel, subscriptionURL, requestURL)

	jalapeno.RequestService(repo, requestURL, networkName)
	jalapeno.ProcessData(repo, nodeChannel, edgeChannel, networkName)

	wg.Wait()
}
