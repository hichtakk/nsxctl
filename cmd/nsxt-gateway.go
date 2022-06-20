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

var unit string
var highlight_row int

func NewCmdShowGateway() *cobra.Command {
	var tier int16
	var output string
	aliases := []string{"gw"}
	gatewayCmd := &cobra.Command{
		Use:     "gateway -t/--tier [0|1] ${GATEWAY_NAME}",
		Aliases: aliases,
		Short:   fmt.Sprintf("show logical gateways [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.Login(site.GetCredential())
			if tier < 0 || tier > 1 {
				log.Fatalf("gateway tier must be specified by flag -t/--tier with value of 0 or 1.\n")
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
				gws.Print(output)
			} else {
				gws = nsxtclient.GetGateway(tier, "")
				gws.Print(output)
			}
		},
	}
	gatewayCmd.Flags().Int16VarP(&tier, "tier", "t", -1, "gateway tier type (0 or 1)")
	gatewayCmd.MarkFlagRequired("tier")

	return gatewayCmd
}

func NewCmdShowRoutingTable() *cobra.Command {
	var ipVer int
	aliases := []string{"rt"}
	gatewayCmd := &cobra.Command{
		Use:     "routes -g/--gateway ${TIER_0_GATEWAY_NAME}",
		Aliases: aliases,
		Short:   fmt.Sprintf("show routing table of specified tier-0 gateways [%s]", strings.Join(aliases, ",")),
		Args:    cobra.ExactArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.Login(site.GetCredential())
			/*
				if tier < 0 || tier > 1 {
					log.Fatalf("gateway tier must be specified by flag -t/--tier with value of 0 or 1.\n")
				}
			*/
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecs := nsxtclient.GetEdgeCluster()
			ecss := structs.EdgeClusters(*ecs)
			routes := nsxtclient.GetRoutingTable(args[0])
			for _, er := range routes {
				ec := ecss.GetClusterById(er.GetEdgeClusterId())
				idx := er.GetEdgeClusterNodeIdx()
				node := nsxtclient.GetTransportNodeById(ec.Members[idx].Id)
				fmt.Printf("/edge-cluster/%v/node/%v\n", ec.Name, node.Name)
				e := er.GetEntries(ipVer)
				e.Print()
			}
		},
	}
	gatewayCmd.Flags().IntVarP(&ipVer, "ip", "", 4, "ip address version (4 or 6)")
	//gatewayCmd.MarkFlagRequired("tier")

	return gatewayCmd
}

func NewCmdTopGateway() *cobra.Command {
	var tier int16
	var interval int
	aliases := []string{"gw"}
	gatewayCmd := &cobra.Command{
		Use:     "gateway",
		Aliases: aliases,
		Short:   fmt.Sprintf("monitor logical gateways [%s]", strings.Join(aliases, ",")),
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
			highlight_row = 0
			if tier == 0 {
				gws := nsxtclient.GetGateway(tier, args[0])
				if len(gws) < 1 {
					log.Fatalln("Tier-0 gateway not found")
				} else if len(gws) > 1 {
					log.Fatalln("found multiple Tier-0 gateways")
				}
				if interval < 3 {
					log.Fatalln("interval must be greater than 3 second")
				}
				runTop(gws[0], interval)
			} else if tier == 1 {
				log.Fatalln("top tier-1 gateway is not implemented yet.")
			} else {
				os.Exit(-1)
			}
		},
	}
	gatewayCmd.Flags().Int16VarP(&tier, "tier", "t", -1, "gateway tier type (0 or 1)")
	gatewayCmd.Flags().IntVarP(&interval, "interval", "i", 5, "update interval (minimum 3 sec)")
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
	emitStr(s, 0, 0, tcell.StyleDefault, "[Press ESC or Ctrl-C to exit]")
	emitStr(s, 0, 2, tcell.StyleDefault, fmt.Sprintf("ID: %s,", gw.Id))
	emitStr(s, 6+len(gw.Name), 2, tcell.StyleDefault, fmt.Sprintf("Name: %s", gw.Name))
	emitStr(s, 0, 3, tcell.StyleDefault, fmt.Sprintf("HA: %s,", gw.HaMode))
	emitStr(s, 6+len(gw.HaMode), 3, tcell.StyleDefault, fmt.Sprintf("Preempt: %s", gw.FailoverMode))
	emitStr(s, w-len(int_msg), 0, tcell.StyleDefault, int_msg)
	displayFooter(s)
	s.Show()
}

func displayFooter(s tcell.Screen) {
	w, h := s.Size()
	reverseStyle := tcell.StyleDefault.Reverse(true)
	for x := 0; x < w; x++ {
		s.SetContent(x, h-1, ' ', nil, reverseStyle)
	}
	footer_msg := "[display unit keys] b: bps, k: Kbps, m: Mbps, g: Gbps, 'space': toggle unit"
	emitStr(s, 0, h-1, reverseStyle, footer_msg)
	str_idx := 0
	selected_unit := ""
	switch unit {
	case "":
		str_idx = strings.Index(footer_msg, "b:")
		selected_unit = "b: bps"
	case "K":
		str_idx = strings.Index(footer_msg, "k:")
		selected_unit = "k: Kbps"
	case "M":
		str_idx = strings.Index(footer_msg, "m:")
		selected_unit = "m: Mbps"
	case "G":
		str_idx = strings.Index(footer_msg, "g:")
		selected_unit = "g: Gbps"
	}
	emitStr(s, str_idx, h-1, tcell.StyleDefault, selected_unit)

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
	x_tx := max_ifname_len + 2
	x_rx := max_ifname_len + 21
	emitStr(s, 0, 5, tcell.StyleDefault, "IfName")
	emitStr(s, x_tx, 5, tcell.StyleDefault, fmt.Sprintf("TX [%sbps]", unit))
	emitStr(s, x_tx+10, 5, tcell.StyleDefault, fmt.Sprintf("TX[%spps]", unit))
	emitStr(s, x_rx, 5, tcell.StyleDefault, fmt.Sprintf("RX [%sbps]", unit))
	emitStr(s, x_rx+10, 5, tcell.StyleDefault, fmt.Sprintf("RX[%spps]", unit))
	for col := 0; col <= w; col++ {
		s.SetContent(col, 6, tcell.RuneHLine, nil, tcell.StyleDefault)
	}

	if highlight_row > len(stats) {
		highlight_row = 0
	}
	if highlight_row < 0 {
		highlight_row = len(stats)
	}
	y := 7
	for i, stat := range stats {
		style := tcell.StyleDefault
		if i+1 == highlight_row {
			style = tcell.StyleDefault.Reverse(true)
			for col := 0; col <= w; col++ {
				s.SetContent(col, y+i, ' ', nil, style)
			}
		}
		port_id_slice := strings.Split(stat.PortId, "/")
		port_id := port_id_slice[len(port_id_slice)-1]
		if last_stats == nil {
			emitStr(s, x_ifname, y+i, style, port_id)
			emitStr(s, x_tx, y+i, style, "*")
			emitStr(s, x_tx+10, y+i, style, "*")
			emitStr(s, x_rx, y+i, style, "*")
			emitStr(s, x_rx+10, y+i, style, "*")
		} else {
			timediff := stat.PerNodeStatistics[0].LastUpdate - last_stats[i].PerNodeStatistics[0].LastUpdate
			tx_bytes := stat.PerNodeStatistics[0].Tx.TotalBytes - last_stats[i].PerNodeStatistics[0].Tx.TotalBytes
			rx_bytes := stat.PerNodeStatistics[0].Rx.TotalBytes - last_stats[i].PerNodeStatistics[0].Rx.TotalBytes
			tx_bps := float64(tx_bytes<<3) / (float64(timediff) / 1000.0)
			rx_bps := float64(rx_bytes<<3) / (float64(timediff) / 1000.0)
			tx_pckts := stat.PerNodeStatistics[0].Tx.TotalPackets - last_stats[i].PerNodeStatistics[0].Tx.TotalPackets
			rx_pckts := stat.PerNodeStatistics[0].Rx.TotalPackets - last_stats[i].PerNodeStatistics[0].Rx.TotalPackets
			u := 1.0
			if unit == "K" {
				u = 1000.0
			} else if unit == "M" {
				u = 1000000.0
			} else if unit == "G" {
				u = 1000000000.0
			}
			emitStr(s, x_ifname, y+i, style, port_id)
			emitStr(s, x_tx, y+i, style, strconv.FormatFloat(tx_bps/u, 'f', 2, 64))
			emitStr(s, x_tx+10, y+i, style, strconv.FormatFloat(float64(tx_pckts)/u, 'f', 2, 64))
			emitStr(s, x_rx, y+i, style, strconv.FormatFloat(rx_bps/u, 'f', 2, 64))
			emitStr(s, x_rx+10, y+i, style, strconv.FormatFloat(float64(rx_pckts)/u, 'f', 2, 64))
		}
	}
	s.Show()
}

func runTop(gw structs.Tier0Gateway, interval int) {
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

	// draw initial contents headers
	s.Clear()
	displayHeader(s, gw, interval)
	w, h := s.Size()
	reverseStyle := tcell.StyleDefault.Reverse(true)
	drawBox(s, w/2-13, h/2-1, w/2+11, h/2+1, reverseStyle, " collecting statistics")
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
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				os.Exit(0)
			} else if ev.Key() == tcell.KeyRight {
				updateBpsUnit(s, rune('R'))
			} else if ev.Key() == tcell.KeyLeft {
				updateBpsUnit(s, rune('L'))
			} else if ev.Key() == tcell.KeyUp {
				highlight_row--
			} else if ev.Key() == tcell.KeyDown {
				highlight_row++
			} else {
				switch ev.Rune() {
				case 'b', 'k', 'm', 'g', ' ':
					updateBpsUnit(s, ev.Rune())
				}
			}
		}
	}
}

func updateBpsUnit(s tcell.Screen, u rune) {
	w, h := s.Size()
	switch u {
	case 'b':
		unit = ""
	case 'k':
		unit = "K"
	case 'm':
		unit = "M"
	case 'g':
		unit = "G"
	case ' ':
		if unit == "K" {
			unit = "M"
		} else if unit == "M" {
			unit = "G"
		} else if unit == "G" {
			unit = ""
		} else {
			unit = "K"
		}
	case 'R':
		if unit == "K" {
			unit = "M"
		} else if unit == "M" {
			unit = "G"
		} else if unit == "G" {
			unit = ""
		} else {
			unit = "K"
		}
	case 'L':
		if unit == "K" {
			unit = ""
		} else if unit == "M" {
			unit = "K"
		} else if unit == "G" {
			unit = "M"
		} else {
			unit = "G"
		}
	}
	msg := fmt.Sprintf("unit has been changed to %sbps", unit)
	for x := 0; x < w; x++ {
		s.SetContent(x, h-2, ' ', nil, tcell.StyleDefault)
	}
	emitStr(s, 0, h-2, tcell.StyleDefault.Reverse(true), msg)
	displayFooter(s)
	s.Show()
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
