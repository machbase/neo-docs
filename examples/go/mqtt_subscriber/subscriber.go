package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/tidwall/gjson"
)

func main() {
	wg := sync.WaitGroup{}

	// paho mqtt client options
	opts := paho.NewClientOptions()
	opts.SetCleanSession(true)
	opts.SetConnectRetry(false)
	opts.SetAutoReconnect(false)
	opts.SetProtocolVersion(4)
	opts.SetClientID("dash-consumer")
	opts.AddBroker("127.0.0.1:4057")
	opts.SetKeepAlive(3 * time.Second)
	opts.SetUsername("user")
	opts.SetPassword("pass")

	// connect to server with paho mqtt client
	client := paho.NewClient(opts)
	connectToken := client.Connect()
	connectToken.WaitTimeout(1 * time.Second)
	if connectToken.Error() != nil {
		panic(connectToken.Error())
	}

	alive := true
	const interval = 1 * time.Second

	// context
	ctx, cancel := context.WithCancel(context.Background())

	// make terminal interface
	term, err := tcell.New()
	if err != nil {
		panic(err)
	}
	defer term.Close()

	// line chart
	lchart, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
	)
	if err != nil {
		panic(err)
	}

	// terminal container
	cont, err := container.New(
		term,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		container.PlaceWidget(lchart),
	)
	if err != nil {
		panic(err)
	}

	// subscribe to receive data
	client.Subscribe("db/reply", 1, func(_ paho.Client, msg paho.Message) {
		buff := msg.Payload()
		str := string(buff)
		vSuccess := gjson.Get(str, "success")
		if vSuccess.Bool() {
			//nrow := gjson.Get(str, "data.rows.#").Int()
			timesRs := gjson.Get(str, "data.rows.#.1")
			times := make([]string, 0)
			timesRs.ForEach(func(k gjson.Result, v gjson.Result) bool {
				tick := time.Unix(0, v.Int())
				times = append(times, tick.Format("15:04:05"))
				return true
			})
			xlabels := make(map[int]string)
			for i, s := range times {
				xlabels[i] = s
			}
			valuesRs := gjson.Get(str, "data.rows.#.2")
			values := make([]float64, 0)
			valuesRs.ForEach(func(k gjson.Result, v gjson.Result) bool {
				values = append(values, v.Float())
				return true
			})

			if err := lchart.Series("first", values,
				linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(33))),
				linechart.SeriesXLabels(xlabels),
			); err != nil {
				panic(err)
			}
		} else {
			fmt.Println("RECV:", str)
			return
		}
	})

	// start consumer
	wg.Add(1)
	go func() {
		for alive {
			now := time.Now().UTC()
			from := now.Add(-10 * time.Minute)

			queryStr := fmt.Sprintf(`select NAME, TIME, VALUE from TAGDATA where TIME between %d and %d AND name = '%s'`,
				from.UnixNano(), now.UnixNano(), "series-1")
			//fmt.Println("SQL", queryStr)

			jsonStr := fmt.Sprintf(`{ "q": "%s" }`, queryStr)
			rt := client.Publish("db/query", 1, false, []byte(jsonStr))

			// if publish was not successful
			if !rt.WaitTimeout(1 * time.Second) {
				fmt.Println("no reponse from server")
				cancel()
			} else if err := rt.Error(); err != nil {
				fmt.Println("ERR", err.Error())
				cancel()
			}

			time.Sleep(interval)
		}
		wg.Done()
	}()

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
			// wait generator to finish
			alive = false
			wg.Wait()
		}
	}

	termOpts := []termdash.Option{
		termdash.KeyboardSubscriber(quitter),
		termdash.RedrawInterval(interval),
	}
	if err := termdash.Run(ctx, term, cont, termOpts...); err != nil {
		panic(err)
	}

	// disconnect mqtt connection
	client.Disconnect(100)
}