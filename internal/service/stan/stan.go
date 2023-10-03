package stan

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"l0/internal/config"
	"log"
)

type OrderService interface {
	Create(orderUID, data string)
	Get(orderUID string) (string, error)
}

type Service struct {
	orderSvc OrderService

	NatsConnect *nats.Conn
	StanConnect stan.Conn
	Sub         stan.Subscription

	url       string
	clusterID string
	clientID  string
	subject   string
}

func New(cfg *config.Stan, orderSvc OrderService) *Service {
	return &Service{
		orderSvc:  orderSvc,
		url:       cfg.URL,
		clusterID: cfg.ClusterID,
		clientID:  cfg.ClientID,
		subject:   cfg.Subject,
	}
}

func (s *Service) Start() {

	nc, err := nats.Connect(s.url, nats.Name("Orders reader"))
	if err != nil {
		log.Fatal(err)
	}

	sc, err := stan.Connect(
		s.clusterID, s.clientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("connection lost, reason: %v", reason)
		}),
	)
	if err != nil {
		log.Fatalf("can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, s.url)
	}

	log.Printf("connected to %s clusterID: [%s] clientID: [%s]\n", s.url, s.clusterID, s.clientID)

	sub, err := sc.QueueSubscribe(
		s.subject,
		"",
		s.handleMessage,
		stan.StartAt(pb.StartPosition_NewOnly),
		stan.DurableName(""),
	)
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("listening on [%s], clientID=[%s]\n", s.subject, s.clientID)

	s.Sub = sub
	s.NatsConnect = nc
	s.StanConnect = sc
}

func (s *Service) Stop() {
	// s.Sub.Unsubscribe()
	s.StanConnect.Close()
	s.NatsConnect.Close()
}
