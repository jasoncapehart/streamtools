package streamtools

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// generic block

type block struct {
	RuleChan chan simplejson.Json
	sigChan  chan os.Signal
	name     string
}

func (b *block) updateRule(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}
	rule, err := simplejson.NewJson(body)
	if err != nil {
		log.Fatalf(err.Error())
	}
	b.RuleChan <- *rule
	fmt.Fprintf(w, "thanks buddy")
}

func (b *block) listenForRules() {
	// listen for rule changes
	http.HandleFunc("/", b.updateRule)
	http.ListenAndServe(":8080", nil)
}

func newBlock(name string) *block {
	return &block{
		RuleChan: make(chan simplejson.Json, 1),
		sigChan:  make(chan os.Signal),
		name:     name,
	}
}

// inBlocks only have an input from stream tools

type inBlockRoutine func(inChan chan simplejson.Json, RuleChan chan simplejson.Json)

type inBlock struct {
	*block // embeds the block type, giving us updateRule and the sigChan for free
	inChan chan simplejson.Json
	f      inBlockRoutine
}

func (b *inBlock) Run(topic string) {
	// set block function going
	go b.f(b.inChan, b.RuleChan)
	// connect to NSQ
	go nsqReader(topic, b.name, b.inChan)
	// set the rule server going
	go b.listenForRules()
	// block until an os.signal
	<-b.sigChan
}

func NewInBlock(f inBlockRoutine, name string) *inBlock {
	b := newBlock(name)
	inChan := make(chan simplejson.Json)
	return &inBlock{b, inChan, f}
}

// outBlocks only have an output to streamtools

type outBlockRoutine func(outChan chan simplejson.Json, RuleChan chan simplejson.Json)

type outBlock struct {
	*block  // embeds the block type, giving us updateRule and the sigChan for free
	outChan chan simplejson.Json
	f       outBlockRoutine
}

func (b *outBlock) Run(topic string) {
	// set block function going
	go b.f(b.outChan, b.RuleChan)
	// connect to NSQ
	log.Println("starting topic:", topic, "with channel:", b.name)
	go nsqWriter(topic, b.name, b.outChan)
	// set the rule server going
	go b.listenForRules()
	// block until an os.signal
	<-b.sigChan
}

func NewOutBlock(f outBlockRoutine, name string) *outBlock {
	b := newBlock(name)
	outChan := make(chan simplejson.Json)
	return &outBlock{b, outChan, f}
}

// state blocks only have inbound data, but also have an API for data

type stateBlockRoutine func(inChan chan simplejson.Json, RuleChan chan simplejson.Json, queryChan chan stateQuery)

type stateBlock struct {
	*block
	inChan    chan simplejson.Json
	queryChan chan stateQuery // for requests to query the state
	f         stateBlockRoutine
}

func (b *stateBlock) Run(topic string) {
	go b.f(b.inChan, b.RuleChan, b.queryChan)
	go b.listenForRules()
	go b.listenForStateQuery()
	<-b.sigChan
}

type stateQuery struct {
	query        simplejson.Json
	responseChan chan simplejson.Json
}

func (b *stateBlock) queryState(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}
	query, err := simplejson.NewJson(body)
	if err != nil {
		log.Fatalf(err.Error())
	}
	q := stateQuery{
		query:        *query,
		responseChan: make(chan simplejson.Json),
	}
	b.queryChan <- q
	// block until the response
	response := <-q.responseChan
	msg, err := response.Encode()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Fprintf(w, string(msg))
}

func (b *stateBlock) listenForStateQuery() {
	http.HandleFunc("/state", b.queryState)
}

func NewStateBlock(f stateBlockRoutine, name string) *stateBlock {
	b := newBlock(name)
	inChan := make(chan simplejson.Json)
	queryChan := make(chan stateQuery)
	return &stateBlock{b, inChan, queryChan, f}
}

// transfer blocks have both inbound and outbound data

type transferBlockRoutine func(inChan chan simplejson.Json, outChan chan simplejson.Json, RuleChan chan *http.Request)