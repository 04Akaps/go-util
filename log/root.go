package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Log struct {
	LogFile *os.File
}

func SetLog(p string) *Log {
	l := &Log{}
	if !strings.HasSuffix(p, ".txt") {
		panic(".txt is not suffixed at logName Env")
	} else {

		f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)

		if os.IsNotExist(err) {
			if f, err = os.Create(p); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}

		l.LogFile = f
		log.SetOutput(f)
		return l

	}
}

func (l *Log) InfoLog(w string) {
	msg := "[INFO] " + w
	fmt.Println(msg)
	log.Printf(msg)
}

func (l *Log) ErrLog(w string) {
	msg := "[ERR] " + w
	fmt.Println(msg)
	log.Printf(msg)
}

func (l *Log) CritLog(w string) {
	msg := "[Crit] " + w
	fmt.Println(msg)
	log.Printf(msg)
	panic(w)
}
