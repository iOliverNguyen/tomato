package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Mode string

var (
	ModeWork       Mode = "work"
	ModeShortBreak Mode = "short-break"
	ModeLongBreak  Mode = "long-break"

	StateStopped = "[S]"
	StatePaused  = "[P]"
	StateRunning = "[R]"

	SepColon = ":"
	SepBreak = "ː"
	N        = 4

	DurationWork       time.Duration
	DurationShortBreak time.Duration
	DurationLongBreak  time.Duration
)

func (mode Mode) Duration() time.Duration {
	switch mode {
	case ModeWork:
		return DurationWork
	case ModeShortBreak:
		return DurationShortBreak
	case ModeLongBreak:
		return DurationLongBreak
	}
	panic("Unexpected")
}

func (mode Mode) Sep() string {
	switch mode {
	case ModeWork:
		return SepColon
	case ModeShortBreak, ModeLongBreak:
		return SepBreak
	}
	panic("Unexpected")
}

type Server struct {
	mode  Mode
	state string
	t     time.Time
	d     time.Duration // remaining duration
	count int
}

func NewServer() *Server {
	return &Server{
		mode:  ModeWork,
		state: StateStopped,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.Index)
	mux.HandleFunc("/status", s.Status)
	mux.HandleFunc("/time", s.Time)
	mux.HandleFunc("/action/start", s.ActionStart)
	mux.HandleFunc("/action/stop", s.ActionStop)

	return mux
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Tomato 1.0.0")
}

func (s *Server) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	s.updateStatus()
	if r.Header.Get("Accept") == "application/json" {
		w.Write(s.formatStatusJSON())
	} else {
		fmt.Fprint(w, s.formatStatus())
	}
}

func (s *Server) Time(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	s.updateStatus()
	fmt.Fprint(w, s.formatTimer())
}

// ActionStart starts or pauses the current interval.
func (s *Server) ActionStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	now := time.Now()
	switch s.state {
	case StateStopped:
		t := now.Add(s.mode.Duration())
		s.t = t
		s.state = StateRunning

	case StatePaused:
		t := now.Add(s.d)
		s.t = t
		s.state = StateRunning

	case StateRunning:
		s.updateStatus()
		if s.state == StateRunning {
			s.d = s.t.Sub(now)
			s.state = StatePaused
		}
	}

	str := s.formatStatus()
	log.Println(str)
	fmt.Fprint(w, str)
}

// ActionStop stops the current running interval or switch mode.
func (s *Server) ActionStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	switch s.state {
	case StateRunning, StatePaused:
		s.state = StateStopped
	case StateStopped:
		switch s.mode {
		case ModeWork:
			if s.count < N {
				s.mode = ModeShortBreak
			} else {
				s.mode = ModeLongBreak
			}

		case ModeShortBreak:
			s.mode = ModeWork

		case ModeLongBreak:
			s.mode = ModeShortBreak
		}
	}

	str := s.formatStatus()
	log.Println(str)
	fmt.Fprint(w, str)
}

func (s *Server) nextMode() {
	switch s.mode {
	case ModeShortBreak, ModeLongBreak:
		if s.mode == ModeLongBreak {
			s.count = 0
		}

		s.mode = ModeWork

	case ModeWork:
		s.count++
		if s.count < N {
			s.mode = ModeShortBreak
		} else {
			s.mode = ModeLongBreak
		}
	}
	panic("Unexpected")
}

func (s *Server) updateStatus() {
	switch s.state {
	case StateRunning:
		if time.Now().After(s.t) {
			s.state = StateStopped
			s.nextMode()
		}
	}
	log.Print(s.formatStatus())
}

func (s *Server) formatStatusJSON() []byte {
	data, _ := json.Marshal(map[string]interface{}{
		"mode":  s.mode,
		"state": s.state,
		"timer": s.formatTimer(),
		"i":     s.count,
		"n":     N,
	})
	return data
}

func (s *Server) formatStatus() string {
	return fmt.Sprintf("%v %v %d/%d %v", s.state, s.formatTimer(), s.count, N, s.mode)
}

func (s *Server) formatTimer() string {
	switch s.state {
	case StateStopped:
		return formatTimer(s.mode.Duration(), s.mode.Sep())
	case StatePaused:
		return formatTimer(s.d, s.mode.Sep())
	case StateRunning:
		return formatTimer(s.t.Sub(time.Now()), s.mode.Sep())
	}
	panic("Unexpected")
}

func formatTimer(d time.Duration, sep string) string {
	if d < 0 {
		d = 0
	}
	m := int(d / time.Minute)
	s := int((d % time.Minute) / time.Second)
	if m > 99 {
		m = 99
	}

	return fmt.Sprintf("%02d%s%02d", m, sep, s)
}

func parseDuration(s string) time.Duration {
	if s == "" {
		log.Fatalf("Invalid duration `%v`", s)
	}
	unit := time.Minute
	switch s[len(s)-1] {
	case 'm':
		s = s[:len(s)-1]
	case 's':
		unit = time.Second
		s = s[:len(s)-1]
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid duration `%v`", s)
	}

	if i <= 0 {
		log.Fatalf("Invalid duration `%v`", s)
	}
	return time.Duration(i) * unit
}

func main() {
	flListen := flag.String("listen", ":12321", "Address to listen on")

	flag.IntVar(&N, "n", N, "Number of intervals between long break")
	flag.StringVar(&SepColon, "colon", SepColon, "Custom separator")
	flag.StringVar(&SepBreak, "colon-alt", SepBreak, "Alternative separator for break modes")

	flDurationWork := flag.String("work", "25m", "Work interval")
	flDurationShortBreak := flag.String("short", "5m", "Short break interval")
	flDurationLongBreak := flag.String("long", "15m", "Long break interval")

	flag.Parse()

	if N <= 0 || N >= 10 {
		log.Fatalf("Invalid number of intervals (%v)", N)
	}
	DurationWork = parseDuration(*flDurationWork)
	DurationShortBreak = parseDuration(*flDurationShortBreak)
	DurationLongBreak = parseDuration(*flDurationLongBreak)
	log.Printf("Interval=%v ShortBreak=%v LongBreak=%v N=%v", DurationWork, DurationShortBreak, DurationLongBreak, N)

	s := NewServer()
	log.Printf("Server listen at %v", *flListen)
	err := http.ListenAndServe(*flListen, s.Handler())
	log.Fatal(err)
}
