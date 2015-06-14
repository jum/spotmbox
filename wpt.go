package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bytbox/go-mail"
)

type WayPoint struct {
	XMLName   string    `json:"-" xml:"wpt"`
	Name      string    `json:"name" xml:"name"`
	Latitude  float64   `json:"lat" xml:"lat,attr"`
	Longitude float64   `json:"lon" xml:"lon,attr"`
	Time      time.Time `json:"time" xml:"time"`
	Comment   string    `json:"cmt" xml:"cmt"`
	Type      string    `json:"type" xml:"-"`
	Messenger string    `json:"messenger" xml:"-"`
}

func ParseSpotMessage(hdr []mail.Header, text string) (w WayPoint, err error) {
	//debug("h %#v\n", hdr)
	//debug("text %s\n", text)
	for _, h := range hdr {
		//debug("%#v\n", h)
		switch strings.ToLower(h.Key) {
		case "x-spot-latitude":
			w.Latitude, err = strconv.ParseFloat(h.Value, 64)
			if err != nil {
				return
			}
		case "x-spot-longitude":
			w.Longitude, err = strconv.ParseFloat(h.Value, 64)
			if err != nil {
				return
			}
		case "x-spot-time":
			var secs int64
			secs, err = strconv.ParseInt(h.Value, 10, 64)
			if err != nil {
				return
			}
			w.Time = time.Unix(secs, 0)
		case "x-spot-messenger":
			w.Messenger = h.Value
		case "x-spot-type":
			w.Type = h.Value
		}
	}
	eqRegex := regexp.MustCompile("(=[0-9a-fA-F][0-9a-fA-F])")
	if err != nil {
		return
	}
	lines := strings.Split(text, "\r\n")
	for i := len(lines) - 2; i >= 0; i-- {
		if len(lines[i]) >= 1 {
			if lines[i][len(lines[i])-1] == '=' {
				lines[i] = lines[i][0:len(lines[i])-1] + lines[i+1]
				lines = append(lines[:i+1], lines[i+2:]...)
			}
		}
	}
	for _, l := range lines {
		//debug("l %s\n", l)
		f := strings.Split(l, ":")
		debug("f %#v\n", f)
		if len(f) == 2 {
			switch f[0] {
			case " Message":
				w.Comment = eqRegex.ReplaceAllStringFunc(f[1], func(s string) string {
					val, err := strconv.ParseUint(s[1:], 16, 8)
					if err != nil {
						panic(err.Error())
					}
					bin := []byte{byte(val)}
					return string(bin)
				})
				if strings.HasPrefix(w.Comment, "Message/ Nachricht") {
					w.Comment = w.Comment[18:]
				}
			}
		}
	}
	return
}
