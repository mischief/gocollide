package main

import (
	"flag"
	"fmt"
	flow "github.com/trustmaster/goflow"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
)

var (
	NPROC = 1
)

func init() {
	NPROC = runtime.NumCPU()
	runtime.GOMAXPROCS(NPROC)

	finish = make(chan bool)
}

var (
	profilefile = flag.String("profile", "", "profiling output file")
	dictfile    = flag.String("dict", "", "dictionary file")
	hashtype    = flag.String("hash", "MD5", "type of hash. one of MD5, SHA1, SHA256, SHA512")
	target      = flag.String("target", "beef", "substring to search for")

	finish chan bool
)

type CollisionApp struct {
	flow.Graph
}

func NewCollisionApp(dict, hashtype, word string) *CollisionApp {
	n := new(CollisionApp)

	n.InitGraphState()

	n.Add(NewHasher(hashtype), "hasher")
	n.Add(NewComparator(word), "comparator")
	n.Add(new(Printer), "printer")

	n.Connect("hasher", "Result", "comparator", "Word", make(chan HashResult))
	n.Connect("comparator", "Result", "printer", "Word", make(chan HashResult))

	n.MapInPort("In", "hasher", "Word")

	return n
}

func (ca *CollisionApp) Finish() {
	finish <- true
}

func main() {
	flag.Parse()

	// capture Ctrl-c on this channel
	killsig := make(chan os.Signal, 1)
	signal.Notify(killsig, os.Interrupt)

	if *profilefile != "" {

		f, err := os.Create(*profilefile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *dictfile == "" {
		log.Fatal("missing argument: -dict")
	}

	log.Printf("dict: %s hashtype: %s target: %s procs: %d", *dictfile, *hashtype, *target, NPROC)

	coll := NewCollisionApp(*dictfile, *hashtype, *target)
	wordchan := make(chan string, 10)
	coll.SetInPort("In", wordchan)

	flow.RunNet(coll)

  log.Printf("Loading dictionary")

	dict, err := LoadDict(*dictfile)
	if err != nil {
		log.Fatal(err)
	}

  log.Printf("Producing work..")

	tries := 0
	tries2 := 0
	for _, line := range dict {
		for _, other_line := range dict {
			select {
			case wordchan <- fmt.Sprintf("%s %s", line, other_line):
			case _ = <-killsig:
				close(wordchan)
				goto done
			default:
			}

			tries++
			if tries%1000000 == 0 {
				tries = 0
				tries2++
				log.Printf("Hashed %d million strings", tries2)
			}
		}
	}

done:
	<-finish
}
