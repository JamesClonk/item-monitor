package main

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/JamesClonk/item-monitor/config"
	"github.com/JamesClonk/item-monitor/log"
	"github.com/Masterminds/sprig"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// watch items, call webhooks with payload
	itemWatcher()
}

func itemWatcher() {
	for {
		for _, monitor := range config.Get().Monitors {
			if len(monitor.Name) > 0 {
				log.Debugf("going through [%s] items ...", monitor.Name)
			}

			for _, item := range monitor.Items {
				log.Debugf("getting body for item [%v], querying [%v] ...", item.Name, item.URL)

				resp, err := http.Get(item.URL)
				if err != nil {
					log.Errorf("could not get body for item [%v]: %v", item.Name, err)
					continue
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Errorf("could not read body of item [%v]: %v", item.Name, err)
					continue
				}
				now := time.Now()

				// run regex to extract value
				rx, err := regexp.Compile(item.Regex)
				if err != nil {
					log.Errorf("could not compile regex [%v]: %v", item.Regex, err)
					continue
				}

				result := rx.FindStringSubmatch(string(body))
				currentValue, err := strconv.ParseFloat(result[1], 64)
				if err != nil {
					log.Errorf("could not parse value [%v]: %v", result, err)
					continue
				}

				if currentValue < item.Value {
					log.Infof("notify about item [%s], because value is [%.2f]", item.Name, currentValue)
					for _, hook := range monitor.Webhooks {
						// parse webhook template, fill it with data
						var data bytes.Buffer
						tmpl := template.Must(template.New("webhook").Funcs(sprig.FuncMap()).Parse(hook.Template))
						if err := tmpl.Execute(
							&data,
							struct {
								Name         string
								URL          string
								Regex        string
								Value        float64
								CurrentValue float64
							}{
								Name:         item.Name,
								URL:          item.URL,
								Regex:        item.Regex,
								Value:        item.Value,
								CurrentValue: currentValue,
							},
						); err != nil {
							log.Errorf("could not prepare webhook template for [%s]: %v", hook.URL, err)
							continue
						}

						// post parsed template to webhook URL
						func() {
							resp, err := http.Post(hook.URL, "application/json", &data)
							if err != nil {
								log.Errorf("could not post monitoring of [%s] to [%s]: %v", item.Name, hook.URL, err)
								return
							}
							defer resp.Body.Close()

							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								log.Errorf("could not read post body for monitoring of [%s] to [%s]: %v", item.Name, hook.URL, err)
								return
							}

							if resp.StatusCode > 299 {
								log.Errorf("monitoring failed: %s", string(body))
							} else {
								log.Info("monitoring success!")
							}
						}()
					}
				}
				item.LastCheck = now
			}
		}

		sleepDuration := time.Duration(config.Get().Interval.Minutes()+float64(rand.Intn(5)+5)) * time.Minute
		log.Debugf("sleep for [%v]", sleepDuration)
		time.Sleep(sleepDuration)
	}
}
