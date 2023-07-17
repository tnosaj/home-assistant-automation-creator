package main

import (
	"flag"
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

type settings struct {
	interval  int
	bandwidth int
	max       int
}

func evaluateInputs() settings {
	var s settings

	flag.IntVar(&s.interval, "i", 5, "every N minutes to run automation")
	flag.IntVar(&s.bandwidth, "b", 50, "how wide is the W bandwidth")
	flag.IntVar(&s.max, "m", 1600, "max W")
	flag.Parse()
	return s
}

//template :=
// - id: '1689406697826'
//   alias: 1201-1400 Consumption
//   description: ''
//   trigger:
//   - platform: time_pattern
//     minutes: /5
//   condition:
//   - condition: numeric_state
//     entity_id: sensor.amis_power_buy_w
//     above: 1201
//     below: 1400
//   action:
//   - device_id: 0a98ea723382f808de0ffabb64e1f237
//     domain: number
//     entity_id: number.hms_1600_4t_limit_nonpersistent_absolute
//     type: set_value
//     value: 1300
//   mode: single
//```

type AutomationEntry struct {
	Id          string      `yaml:"id"`
	Alias       string      `yaml:"alias"`
	Description string      `yaml:"description"`
	Trigger     []Trigger   `yaml:"trigger"`
	Condition   []Condition `yaml:"condition"`
	Action      []Action    `yaml:"action"`
	Mode        string      `yaml:"mode"`
}
type Trigger struct {
	Platform string `yaml:"platform"`
	Minutes  string `yaml:"minutes"`
}

type Condition struct {
	Condition string `yaml:"condition"`
	Entity_id string `yaml:"entity_id"`
	Above     int    `yaml:"above,omitempty"`
	Below     int    `yaml:"below,omitempty"`
	State     string `yaml:"state,omitempty"`
}

type Action struct {
	Device_id string `yaml:"device_id"`
	Domain    string `yaml:"domain"`
	Entity_id string `yaml:"entity_id"`
	Type      string `yaml:"type"`
	Value     int    `yaml:"value"`
}

func main() {
	s := evaluateInputs()

	var AutomationEntries []AutomationEntry

	start := 0
	end := s.bandwidth
	watt := s.bandwidth - s.bandwidth/2
	for i := 0; i < s.max/s.bandwidth; i++ {
		var automationEntry AutomationEntry

		automationEntry.Id = fmt.Sprintf("from%04dto%04dconsumption", start, end)
		automationEntry.Alias = fmt.Sprintf("%04d - %04d Consumption", start, end)
		automationEntry.Description = ""
		automationEntry.Mode = "single"

		automationEntry.Trigger = append(
			automationEntry.Trigger,
			Trigger{
				Platform: "time_pattern",
				Minutes:  fmt.Sprintf("/%d", s.interval),
			},
		)
		automationEntry.Condition = append(
			automationEntry.Condition,
			Condition{
				Condition: "numeric_state",
				Entity_id: "sensor.amis_power_buy_w",
				Above:     start,
				Below:     end,
			},
			Condition{
				Condition: "state",
				Entity_id: "binary_sensor.hms_1600_4t_producing",
				State:     "on",
			},
		)
		automationEntry.Action = append(
			automationEntry.Action,
			Action{
				Device_id: "0a98ea723382f808de0ffabb64e1f237",
				Domain:    "number",
				Entity_id: "number.hms_1600_4t_limit_nonpersistent_absolute",
				Type:      "set_value",
				Value:     watt,
			},
		)

		fmt.Println(fmt.Sprintf("start: %d - end: %d - watt: %d", start, end, watt))
		AutomationEntries = append(AutomationEntries, automationEntry)
		start = start + s.bandwidth
		end = end + s.bandwidth
		watt = end - s.bandwidth/2
	}

	var automationEntry AutomationEntry
	automationEntry.Id = "maxconsumption"
	automationEntry.Alias = fmt.Sprintf("> %d Consumption", s.max)
	automationEntry.Description = ""
	automationEntry.Mode = "single"

	automationEntry.Trigger = append(
		automationEntry.Trigger,
		Trigger{
			Platform: "time_pattern",
			Minutes:  fmt.Sprintf("/%d", s.interval),
		},
	)
	automationEntry.Condition = append(
		automationEntry.Condition,
		Condition{
			Condition: "numeric_state",
			Entity_id: "sensor.amis_power_buy_w",
			Above:     1600,
		},
		Condition{
			Condition: "state",
			Entity_id: "binary_sensor.hms_1600_4t_producing",
			State:     "on",
		},
	)
	automationEntry.Action = append(
		automationEntry.Action,
		Action{
			Device_id: "0a98ea723382f808de0ffabb64e1f237",
			Domain:    "number",
			Entity_id: "number.hms_1600_4t_limit_nonpersistent_relative",
			Type:      "set_value",
			Value:     100,
		},
	)

	fmt.Println(fmt.Sprintf("start: %d - end: %d - watt: %d", start, end, watt))
	AutomationEntries = append(AutomationEntries, automationEntry)
	d, err := yaml.Marshal(&AutomationEntries)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%s", string(d))
	//template :=
	// - id: '1689406697826'
	//   alias: 1201-1400 Consumption
	//   description: ''
	//   trigger:
	//   - platform: time_pattern
	//     minutes: /5
	//   condition:
	//   - condition: numeric_state
	//     entity_id: sensor.amis_power_buy_w
	//     above: 1201
	//     below: 1400
	//   action:
	//   - device_id: 0a98ea723382f808de0ffabb64e1f237
	//     domain: number
	//     entity_id: number.hms_1600_4t_limit_nonpersistent_absolute
	//     type: set_value
	//     value: 1300
	//   mode: single
	//```
}
