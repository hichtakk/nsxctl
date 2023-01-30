package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
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
		ValidArgsFunction: func (cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			Login()
			gw_names := []string{}
			switch tier {
			case 0:
				for _, gw := range nsxtclient.GetTier0Gateway("") {
					gw_names = append(gw_names, gw.Name)
				}
			case 1:
				for _, gw := range nsxtclient.GetTier1Gateway("") {
					gw_names = append(gw_names, gw.Name)
				}
			default:
				for _, gw := range nsxtclient.GetTier0Gateway("") {
					gw_names = append(gw_names, gw.Name)
				}
				for _, gw := range nsxtclient.GetTier1Gateway("") {
					gw_names = append(gw_names, gw.Name)
				}
			}
			return gw_names, cobra.ShellCompDirectiveNoFileComp
		},
		PreRunE: func(c *cobra.Command, args []string) error {
			if tier != -1 && (tier < 0 || tier > 1) {
				log.Fatalf("gateway tier must be specified by flag -t/--tier with value of 0 or 1.\n")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				gwName := args[0]
				switch tier {
				case 0:
					gw, err := nsxtclient.GetTier0GatewayFromName(gwName)
					if err != nil {
						log.Fatal(err)
					}
					locale_service_id := nsxtclient.GetLocaleService(gw.Id)
					interfaces := nsxtclient.GetInterface(gw.Id, locale_service_id[0])
					bgp := nsxtclient.GetBgpConfig(gw.Id, locale_service_id[0])
					gw.Print(interfaces, bgp)
				case 1:
					gw, err := nsxtclient.GetTier1GatewayFromName(gwName)
					if err != nil {
						log.Fatal(err)
					}
					gw.Print()
				default:
					gw0, err := nsxtclient.GetTier0GatewayFromName(gwName)
					if err != nil {
						gw1, err := nsxtclient.GetTier1GatewayFromName(gwName)
						if err != nil {
							log.Fatal(fmt.Sprintf("Error: Tier-0 or Tier-1 gateway '%s' is not found", gwName))
						}
						gw1.Print()
						return
					}
					locale_service_id := nsxtclient.GetLocaleService(gw0.Id)
					interfaces := nsxtclient.GetInterface(gw0.Id, locale_service_id[0])
					bgp := nsxtclient.GetBgpConfig(gw0.Id, locale_service_id[0])
					gw0.Print(interfaces, bgp)
				}
			} else {
				switch tier {
				case 0:
					gws := nsxtclient.GetTier0Gateway("")
					gws.Print(output)
				case 1:
					gws := nsxtclient.GetTier1Gateway("")
					gws.Print(output)
				default:
					w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
					w.Write([]byte(strings.Join([]string{"Tier", "ID", "Name", "HA Mode", "Failover Mode", "FW Enabled"}, "\t")+ "\n"))
					gws0 := nsxtclient.GetTier0Gateway("")
					sort.Slice(gws0, func(i, j int) bool { return gws0[i].Name > gws0[i].Name })
					for _, gw := range gws0 {
						w.Write([]byte(strings.Join([]string{"0", gw.Id, gw.Name, gw.HaMode, gw.FailoverMode, strconv.FormatBool(!gw.Firewall)}, "\t")+ "\n"))
					}
					gws1 := nsxtclient.GetTier1Gateway("")
					sort.Slice(gws1, func(i, j int) bool { return gws1[i].Name > gws1[i].Name })
					for _, gw := range gws1 {
						w.Write([]byte(strings.Join([]string{"1", gw.Id, gw.Name, "ACTIVE_STANDBY", gw.FailoverMode, strconv.FormatBool(!gw.Firewall)}, "\t")+ "\n"))
					}
					w.Flush()
				}
			}
		},
	}
	gatewayCmd.Flags().Int16VarP(&tier, "tier", "t", -1, "gateway tier type (0 or 1)")

	return gatewayCmd
}

func GetTier0GatewayNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	Login()
	t0gw_names := []string{}
	gws := nsxtclient.GetTier0Gateway("")
	for _, gw := range gws {
		t0gw_names = append(t0gw_names, gw.Name)
	}
	return t0gw_names, cobra.ShellCompDirectiveNoFileComp
}

func NewCmdShowRoutingTable() *cobra.Command {
	var ipv6 bool
	aliases := []string{"rt"}
	gatewayCmd := &cobra.Command{
		Use:               "routes ${TIER_0_GATEWAY_NAME}",
		Aliases:           aliases,
		Short:             fmt.Sprintf("show routing table of specified tier-0 gateways [%s]", strings.Join(aliases, ",")),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: GetTier0GatewayNames,
		Run: func(cmd *cobra.Command, args []string) {
			ecs := nsxtclient.GetEdgeCluster()
			ecss := structs.EdgeClusters(*ecs)
			gw, err := nsxtclient.GetTier0GatewayFromName(args[0])
			if err != nil {
				log.Fatal(err)
				return
			}
			routes := nsxtclient.GetRoutingTable(gw.Id)
			for _, er := range routes {
				ec := ecss.GetClusterById(er.GetEdgeClusterId())
				idx := er.GetEdgeClusterNodeIdx()
				node := nsxtclient.GetTransportNodeById(ec.Members[idx].Id)
				fmt.Printf("/edge-cluster/%v/node/%v\n", ec.Name, node.Name)
				if ipv6 == true {
					e := er.GetEntries(6)
					e.Print()
				} else {
					e := er.GetEntries(4)
					e.Print()
				}
			}
		},
	}
	gatewayCmd.Flags().BoolVarP(&ipv6, "ipv6", "6", false, "show IPv6 routes")

	return gatewayCmd
}

func NewCmdShowBgpAdvRoutes() *cobra.Command {
	aliases := []string{"adv"}
	gatewayCmd := &cobra.Command{
		Use:               "bgp-adv-routes ${TIER_0_GATEWAY_NAME}",
		Aliases:           aliases,
		Short:             fmt.Sprintf("show bgp advertise routes of specified tier-0 gateways [%s]", strings.Join(aliases, ",")),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: GetTier0GatewayNames,
		Run: func(cmd *cobra.Command, args []string) {
			locale := "default"
			gw, err := nsxtclient.GetTier0GatewayFromName(args[0])
			if err != nil {
				log.Fatal(err)
				return
			}
			neighbors := nsxtclient.GetBgpNeighbors(gw.Id, locale)
			for _, nb := range neighbors {
				fmt.Printf("BGP neighbor: %v, Remote ASN: %v\n", nb.Address, nb.Asn)
				nb_edges := nsxtclient.GetBgpNeighborsAdvRoutes(nb.Path)
				for _, e := range nb_edges {
					node := nsxtclient.GetTransportNodeById(e.EdgeId)
					fmt.Printf("Edge node: %v, Source IP: %v\n\n", node.Name, e.Source)
					entries := e.GetEntries()
					entries.Print()
					fmt.Println()
				}
			}
		},
	}

	return gatewayCmd
}

func GetGatewayNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	tier, err := cmd.Flags().GetInt16("tier")
	if err != nil {
		log.Fatal(err)
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	if tier == 0 {
		return GetTier0GatewayNames(cmd, args, toComplete)
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func NewCmdTopGateway() *cobra.Command {
	var tier int16
	var interval int
	aliases := []string{"gw"}
	gatewayCmd := &cobra.Command{
		Use:               "gateway",
		Aliases:           aliases,
		Short:             fmt.Sprintf("monitor logical gateways [%s]", strings.Join(aliases, ",")),
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: GetGatewayNames,
		Run: func(cmd *cobra.Command, args []string) {
			highlight_row = 0
			if tier == 0 {
				gw, err := nsxtclient.GetTier0GatewayFromName(args[0])
				if err != nil {
					log.Fatal(err)
					return
				}
				if interval < 3 {
					log.Fatalln("interval must be greater than 3 second")
				}
				runTop(gw, interval)
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
