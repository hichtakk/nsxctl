package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/hichtakk/nsxctl/structs"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

func NewCmdShowGateway() *cobra.Command {
	var tier int16
	var output string
	gatewayCmd := &cobra.Command{
		Use:     "gateway",
		Aliases: []string{"gw"},
		Short:   "show logical gateways",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.Login(site.GetCredential())
			if tier < 0 || tier > 1 {
				log.Fatalf("error %d\n", tier)
			}
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var gws structs.Tier0Gateways
			if len(args) > 0 {
				gws = nsxtclient.GetGateway(tier, args[0])
			} else {
				gws = nsxtclient.GetGateway(tier, "")
			}
			gws.Print(output)
		},
	}
	gatewayCmd.Flags().Int16VarP(&tier, "tier", "t", -1, "gateway tier type (0 or 1)")
	gatewayCmd.MarkFlagRequired("tier")

	return gatewayCmd
}

func NewCmdTopGateway() *cobra.Command {
	var tier int16
	gatewayCmd := &cobra.Command{
		Use:     "gateway",
		Aliases: []string{"gw"},
		Short:   "monitor logical gateways",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.Login(site.GetCredential())
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if tier == 0 {
				gws := nsxtclient.GetGateway(tier, args[0])
				if len(gws) < 1 {
					log.Fatalln("Tier-0 gateway not found")
				} else if len(gws) > 1 {
					log.Fatalln("found multiple Tier-0 gateways")
				}
				runTop(gws[0])
			} else if tier == 1 {
				log.Fatalln("top tier-1 gateway is not implemented yet.")
			} else {
				os.Exit(-1)
			}
		},
	}
	gatewayCmd.Flags().Int16VarP(&tier, "tier", "t", -1, "gateway tier type (0 or 1)")
	gatewayCmd.MarkFlagRequired("tier")

	return gatewayCmd
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func displayHeader(s tcell.Screen, gw structs.Tier0Gateway, interval int) {
	w, _ := s.Size()
	int_msg := fmt.Sprintf("update interval: %ds", interval)
	emitStr(s, 0, 0, tcell.StyleDefault, "[Press ESC to exit]")
	emitStr(s, 0, 2, tcell.StyleDefault, fmt.Sprintf("ID: %s,", gw.Id))
	emitStr(s, 6+len(gw.Name), 2, tcell.StyleDefault, fmt.Sprintf("Name: %s", gw.Name))
	emitStr(s, 0, 3, tcell.StyleDefault, fmt.Sprintf("HA: %s,", gw.HaMode))
	emitStr(s, 6+len(gw.HaMode), 3, tcell.StyleDefault, fmt.Sprintf("Preempt: %s", gw.FailoverMode))
	emitStr(s, w-len(int_msg), 0, tcell.StyleDefault, int_msg)
	s.Show()
}

func update(s tcell.Screen, stats map[int]structs.RouterStats, last_stats map[int]structs.RouterStats) {
	w, _ := s.Size()
	max_ifname_len := 0
	for _, stat := range stats {
		port_id_slice := strings.Split(stat.PortId, "/")
		port_id := port_id_slice[len(port_id_slice)-1]
		if len(port_id) > max_ifname_len {
			max_ifname_len = len(port_id)
		}
	}
	x_ifname := 0
	x_time := max_ifname_len + 2
	x_tx := max_ifname_len + 9
	x_rx := max_ifname_len + 21
	emitStr(s, 0, 5, tcell.StyleDefault, "IfName")
	emitStr(s, x_time, 5, tcell.StyleDefault, "Time")
	emitStr(s, x_tx, 5, tcell.StyleDefault, "TX")
	emitStr(s, x_rx, 5, tcell.StyleDefault, "RX")
	for col := 0; col <= w; col++ {
		s.SetContent(col, 6, tcell.RuneHLine, nil, tcell.StyleDefault)
	}

	y := 7
	for i, stat := range stats {
		port_id_slice := strings.Split(stat.PortId, "/")
		port_id := port_id_slice[len(port_id_slice)-1]
		if last_stats == nil {
			emitStr(s, x_ifname, y+i, tcell.StyleDefault, port_id)
			emitStr(s, x_time, y+i, tcell.StyleDefault, "*")
			emitStr(s, x_tx, y+i, tcell.StyleDefault, "*")
			emitStr(s, x_rx, y+i, tcell.StyleDefault, "*")
		} else {
			timediff := stat.PerNodeStatistics[0].LastUpdate - last_stats[i].PerNodeStatistics[0].LastUpdate
			tx_bytes := stat.PerNodeStatistics[0].Tx.TotalBytes - last_stats[i].PerNodeStatistics[0].Tx.TotalBytes
			rx_bytes := stat.PerNodeStatistics[0].Rx.TotalBytes - last_stats[i].PerNodeStatistics[0].Rx.TotalBytes
			tx_bps := float64(tx_bytes<<3) / (float64(timediff) / 1000.0)
			rx_bps := float64(rx_bytes<<3) / (float64(timediff) / 1000.0)
			emitStr(s, x_ifname, y+i, tcell.StyleDefault, port_id)
			emitStr(s, x_time, y+i, tcell.StyleDefault, strconv.Itoa(int(timediff)))
			emitStr(s, x_tx, y+i, tcell.StyleDefault, strconv.FormatFloat(tx_bps, 'f', 2, 64))
			emitStr(s, x_rx, y+i, tcell.StyleDefault, strconv.FormatFloat(rx_bps, 'f', 2, 64))
		}
	}
	s.Show()
}

func runTop(gw structs.Tier0Gateway) {
	interval := 5
	encoding.Register()
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)

	// draw initial contents headers
	s.Clear()
	displayHeader(s, gw, interval)
	w, h := s.Size()
	drawBox(s, w/2-13, h/2-1, w/2+11, h/2+1, defStyle, " collecting statistics")
	s.Show()

	// initial interface list
	stats := nsxtclient.GetGatewayInterfaceStats(gw)
	update(s, stats, nil)

	// update stats loop
	go func(g structs.Tier0Gateway, stats map[int]structs.RouterStats, interval int) {
		last_stats := stats
		for range time.Tick(time.Duration(interval) * time.Second) {
			new_stats := nsxtclient.GetGatewayInterfaceStats(g)
			s.Clear()
			displayHeader(s, g, interval)
			update(s, new_stats, last_stats)
			last_stats = new_stats
		}
	}(gw, stats, interval)

	// event loop
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			displayHeader(s, gw, interval)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
		}
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}

/*
func runTop2() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	// Draw initial boxes
	drawBox(s, 1, 1, 42, 7, boxStyle, "Click and drag to draw a box")
	drawBox(s, 5, 9, 32, 14, boxStyle, "Press C to reset")

	// Event loop
	ox, oy := -1, -1
	quit := func() {
		s.Fini()
		os.Exit(0)
	}
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			button := ev.Buttons()
			// Only process button events, not wheel events
			button &= tcell.ButtonMask(0xff)

			if button != tcell.ButtonNone && ox < 0 {
				ox, oy = x, y
			}
			switch ev.Buttons() {
			case tcell.ButtonNone:
				if ox >= 0 {
					label := fmt.Sprintf("%d,%d to %d,%d", ox, oy, x, y)
					drawBox(s, ox, oy, x, y, boxStyle, label)
					ox, oy = -1, -1
				}
			}
		}
	}
}
*/
