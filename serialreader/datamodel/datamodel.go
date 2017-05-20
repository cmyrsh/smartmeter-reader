package datamodel

import (
	"fmt"
	//"log"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
The code consists of (up to) 6 group sub-identifiers marked by letters A to F.
All these may or may not be present in the identifier (e.g. groups A and B are often omitted).
In order to decide to which group the sub-identifier belongs, the groups are separated by unique separators:

A-B:C.D.E*F

- The A group defines the medium (0=abstract objects, 1=electricity, 6=heat, 7=gas, 8=water, ...)
- The B group defines the channel. Each device with multiple channels generating measurement results, can separate the results into the channels.
- The C group defines the physical value (current, voltage, energy, level, temperature, ...)
- The D group defines the quantity computation output of specific algorythm
- The E group specifies the measurement type defined by groups A to D into individual measurements (e.g. switching ranges)
- The F group separates the results partly defined by groups A to E. The typical usage is the specification of individual time ranges.

*/

var (
	p_version = regexp.MustCompile("1-3:0\\.2\\.8\\(|\\)")

	p_timestamp  = regexp.MustCompile("0-0:1\\.0\\.0\\(|\\)")
	p_equip_id   = regexp.MustCompile("0-0:96\\.1\\.1\\(|\\)")
	p_eto_trf1   = regexp.MustCompile("1-0:1\\.8\\.1\\(|\\)")
	p_efrom_trf1 = regexp.MustCompile("1-0:2\\.8\\.1\\(|\\)")

	p_eto_trf2   = regexp.MustCompile("1-0:1\\.8\\.2\\(|\\)")
	p_efrom_trf2 = regexp.MustCompile("1-0:2\\.8\\.2\\(|\\)")

	p_trf_id = regexp.MustCompile("0-0:96\\.14\\.0\\(|\\)")

	p_actpwr_to   = regexp.MustCompile("1-0:1\\.7\\.0\\(|\\)")
	p_actpwr_from = regexp.MustCompile("1-0:2\\.7\\.0\\(|\\)")

	p_count_pwrflr = regexp.MustCompile("0-0:96\\.7\\.21\\(|\\)")
	/* value with multiple paranthesis*/
	p_pwrflr_log = regexp.MustCompile("0-0:96\\.7\\.9\\(|\\)|\\(")

	p_count_vltsags   = regexp.MustCompile("1-0:32\\.32\\.0\\(|\\)")
	p_count_vltswells = regexp.MustCompile("1-0:32\\.36\\.0\\(|\\)")

	p_txt_msgcoder = regexp.MustCompile("0-0:96\\.13\\.1\\(|\\)")
	p_txt_msgmax   = regexp.MustCompile("0-0:96\\.13\\.0\\(|\\)")

	p_inst_curr = regexp.MustCompile("1-0:31\\.7\\.0\\(|\\)")

	p_inst_activepwr_in  = regexp.MustCompile("1-0:21\\.17\\.0\\(|\\)")
	p_inst_activepwr_out = regexp.MustCompile("1-0:22\\.17\\.0\\(|\\)")

	p_g_devicetype = regexp.MustCompile("0-1:24\\.1\\.0\\(|\\)")
	p_g_eqip_id    = regexp.MustCompile("0-1:96\\.1\\.0\\(|\\)")

	/* value with multiple paranthesis*/
	p_g_time_value = regexp.MustCompile("0-1:24\\.2\\.1\\(|\\)|\\(")

	m3remove = regexp.MustCompile("\\*m3")

	scremove = regexp.MustCompile("\\*s")

	kWremove = regexp.MustCompile("\\*kW")

	kWhremove = regexp.MustCompile("\\*kWh")

	ampremove = regexp.MustCompile("\\*A")

	gasreader_date_layout = "060102150405S"

	gasreader_unit_length = 4

	powerfailure_unit_length = 8

	powerfailure_date_layout = "060102150405W"
)

const (
	const_p_version = "1-3:0.2.8"

	const_p_timestamp  = "0-0:1.0.0"
	const_p_equip_id   = "0-0:96.1.1"
	const_p_eto_trf1   = "1-0:1.8.1"
	const_p_efrom_trf1 = "1-0:2.8.1"

	const_p_eto_trf2   = "1-0:1.8.2"
	const_p_efrom_trf2 = "1-0:2.8.2"

	const_p_trf_id = "0-0:96.14.0"

	const_p_actpwr_to   = "1-0:1.7.0"
	const_p_actpwr_from = "1-0:2.7.0"

	const_p_count_pwrflr = "0-0:96.7.21"
	/* value with multiple paranthesis*/
	const_p_pwrflr_log = "0-0:96.7.9"

	const_p_count_vltsags   = "1-0:32.32.0"
	const_p_count_vltswells = "1-0:32.36.0"

	const_p_txt_msgcoder = "0-0:96.13.1"
	const_p_txt_msgmax   = "0-0:96.13.0"

	const_p_inst_curr = "1-0:31.7.0"

	const_p_inst_activepwr_in  = "1-0:21.17.0"
	const_p_inst_activepwr_out = "1-0:22.17.0"

	const_p_g_devicetype = "0-1:24.1.0"
	const_p_g_eqip_id    = "0-1:96.1.0"

	/* value with multiple paranthesis*/
	const_p_g_time_value = "0-1:24.2.1"
)

type ValueFilter struct {
	regex_list []regexp.Regexp
}
type PowerFailure struct {
	Failure_time     time.Time
	Failure_duration int64
}

type GasReading struct {
	Capture_time time.Time
	Value_m3     float64
}

type P1Telegram struct {
	Manufacture_spec                      string
	Version_info                          string
	Date_timestamp                        time.Time
	Equipment_id                          string
	Electricity_incoming_tarfif1_kwh      float64
	Electricity_outgoing_tarfif1_kwh      float64
	Electricity_incoming_tarfif2_kwh      float64
	Electricity_outgoing_tarfif2_kwh      float64
	Tarrif_id                             string
	Actualpower_incoming_kW               float64
	Actualpower_outgoing_kW               float64
	Power_failure_count                   int64
	Power_failure_history                 []PowerFailure
	L1_voltage_sag_count                  int64
	L1_voltage_swell_count                int64
	Text_msg                              string
	Text_msg_maxchars                     int64
	Instantaneous_current_amp             float64
	Instantaneous_activepower_incoming_kW float64
	Instantaneous_activepower_outgoing_kW float64
	Gas_device_type                       string
	Gas_device_id                         string
	Gas_readings                          []GasReading
	Crc_string                            string

	/*
		/XMX5LGBBFFB231164713 -- manufacturing info

		1-3:0.2.8(42) -- version info for p1 output
		0-0:1.0.0(170326100029S) -- date timestamp -
		0-0:96.1.1(4530303034303031353734323134303134) - equipment identifier
		1-0:1.8.1(001878.531*kWh) - electricity to client - tarrif 1
		1-0:2.8.1(000000.000*kWh) - electricity by client - tarrif 1
		1-0:1.8.2(002292.468*kWh) - electricity to client - tarrif 2
		1-0:2.8.2(000000.000*kWh) - electricity by client - tarrif 2
		0-0:96.14.0(0001) - tarrif indicator
		1-0:1.7.0(00.235*kW) - actual power delivered - 1 watt resolution
		1-0:2.7.0(00.000*kW) - actual power received - 1 watt resolution
		0-0:96.7.21(00003) - number of power failures in any phase
		0-0:96.7.9(00001) - number of long power failures in any phase
		1-0:99.97.0(1)(0-0:96.7.19)(150327111101W)(0000005576*s) - power failure event log
		1-0:32.32.0(00000) - number of voltage sags in phase L1
		1-0:32.36.0(00000) - number of voltage swells in phase L1
		0-0:96.13.1() - text message coders
		0-0:96.13.0() - text message max
		1-0:31.7.0(001*A) - instantaneous current
		1-0:21.7.0(00.235*kW) - instataneous active power P+
		1-0:22.7.0(00.000*kW) - instataneous active power P-
		0-1:24.1.0(003) - Gas meter Device Type
		0-1:96.1.0(4730303233353631323231303339333134) - gas meter equipment identifier
		0-1:24.2.1(170326090000S)(03360.854*m3) - gas meter capture time and value
		!C1C8

	*/
}

func (target *P1Telegram) PopulateFromLine(data string) error {
	str_e := strings.TrimSpace(data)
	//log.Println(str_e)
	if strings.HasPrefix(str_e, "/") {
		target.Manufacture_spec = strings.TrimPrefix(str_e, "/")
		return nil
	}

	if strings.HasPrefix(str_e, const_p_version) {
		target.Version_info = p_version.Split(str_e, -1)[1]
		return nil
	}
	if strings.HasPrefix(str_e, const_p_timestamp) {
		var ts_string = p_timestamp.Split(str_e, -1)[1]
		t, err := time.Parse(gasreader_date_layout, ts_string)
		if err != nil {
			log.Fatal(err)
		}
		target.Date_timestamp = t
		return nil
	}
	if strings.HasPrefix(str_e, const_p_equip_id) {
		//log.Println("equipid : ", str_e)
		//target.Equipment_id = p_equip_id.Split(str_e, -1)[1]
		return nil
	}
	if strings.HasPrefix(str_e, const_p_eto_trf1) {
		//log.Println("const_p_eto_trf1 : ", str_e)
		target.Electricity_incoming_tarfif1_kwh, _ = strconv.ParseFloat(kWhremove.ReplaceAllString(p_eto_trf1.Split(str_e, -1)[1], ""), 0)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_efrom_trf1) {
		target.Electricity_outgoing_tarfif1_kwh, _ = strconv.ParseFloat(kWhremove.ReplaceAllString(p_efrom_trf1.Split(str_e, -1)[1], ""), 0)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_eto_trf2) {
		//log.Println("const_p_eto_trf2 : ", str_e)
		target.Electricity_incoming_tarfif2_kwh, _ = strconv.ParseFloat(kWhremove.ReplaceAllString(p_eto_trf2.Split(str_e, -1)[1], ""), 0)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_efrom_trf2) {
		target.Electricity_outgoing_tarfif2_kwh, _ = strconv.ParseFloat(kWhremove.ReplaceAllString(p_efrom_trf2.Split(str_e, -1)[1], ""), 0)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_trf_id) {
		target.Tarrif_id = p_trf_id.Split(str_e, -1)[0]
		return nil
	}
	if strings.HasPrefix(str_e, const_p_actpwr_to) {
		target.Actualpower_incoming_kW, _ = strconv.ParseFloat(kWremove.ReplaceAllString(p_actpwr_to.Split(str_e, -1)[1], ""), 0)
		return nil

	}
	if strings.HasPrefix(str_e, const_p_actpwr_from) {
		target.Actualpower_outgoing_kW, _ = strconv.ParseFloat(kWremove.ReplaceAllString(p_actpwr_from.Split(str_e, -1)[1], ""), 0)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_count_pwrflr) {
		target.Power_failure_count, _ = strconv.ParseInt(p_count_pwrflr.Split(str_e, -1)[1], 10, 0)
		return nil
	}
	/*---
		1-0:99.97.0(1)(0-0:96.7.19)(150327111101W)(0000005576*s) - power failure event log
	----*/
	if strings.HasPrefix(str_e, const_p_pwrflr_log) {
		target.Power_failure_history = extractPowerFailureLog(str_e)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_count_vltsags) {
		target.L1_voltage_sag_count, _ = strconv.ParseInt(p_count_vltsags.Split(str_e, -1)[1], 10, 0)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_count_vltswells) {
		target.L1_voltage_swell_count, _ = strconv.ParseInt(p_count_vltswells.Split(str_e, -1)[1], 10, 0)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_txt_msgcoder) {
		target.Text_msg = p_txt_msgcoder.Split(str_e, -1)[0]
		return nil
	}
	if strings.HasPrefix(str_e, const_p_txt_msgmax) {
		target.Text_msg_maxchars, _ = strconv.ParseInt(p_txt_msgmax.Split(str_e, -1)[1], 10, 0)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_inst_curr) {
		target.Instantaneous_current_amp, _ = strconv.ParseFloat(ampremove.ReplaceAllString(p_inst_curr.Split(str_e, -1)[1], ""), 0)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_inst_activepwr_in) {
		target.Instantaneous_activepower_incoming_kW, _ = strconv.ParseFloat(p_inst_activepwr_in.Split(str_e, -1)[1], 0)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_inst_activepwr_out) {
		target.Instantaneous_activepower_outgoing_kW, _ = strconv.ParseFloat(kWremove.ReplaceAllString(p_inst_activepwr_out.Split(str_e, -1)[1], ""), 0)
		return nil
	}
	if strings.HasPrefix(str_e, const_p_g_devicetype) {
		target.Gas_device_type = p_g_devicetype.Split(str_e, -1)[1]
		return nil
	}
	if strings.HasPrefix(str_e, const_p_g_eqip_id) {
		target.Gas_device_id = p_g_eqip_id.Split(str_e, -1)[1]
		return nil
	}
	/*-- create gas readings object--*/

	if strings.HasPrefix(str_e, const_p_g_time_value) {
		target.Gas_readings, _ = extractGasReadings(str_e)
		return nil
	}
	if strings.HasPrefix(str_e, "!") {
		target.Crc_string = strings.TrimPrefix(str_e, "!")
		return nil
	}

	return nil
}


func extractGasReadings(data_raw string) ([]GasReading, error) {

	var dataline = strings.TrimSpace(data_raw)
	var str_arr = p_g_time_value.Split(dataline, -1)
	var gas_reading = make([]GasReading, 0) //[math.Ceil(len(str_arr) / unit_length)]datamodel.GasReading{}
	var append_id = 0
	for str_g_idx, strg_val := range str_arr {

		//println("extractGasReadings : ", strg_val)

		if len(strg_val) > 0 {
			switch id := str_g_idx % gasreader_unit_length; id {
			case 1:
				t, err := time.Parse(gasreader_date_layout, strg_val)
				if err != nil {
					//return err
				}
				gr := GasReading{}
				gr.Capture_time = t
				gas_reading = append(gas_reading, gr)
				append_id = len(gas_reading) - 1
				continue
			case 3:
				val := m3remove.ReplaceAllString(strg_val, "")
				gas_reading[append_id].Value_m3, _ = strconv.ParseFloat(val, 0)
				continue
			default:
				// nothing
			}
		}

	}

	return gas_reading, nil
}

/*---
	1-0:99.97.0(1)(0-0:96.7.19)(150327111101W)(0000005576*s) - power failure event log
----*/
func extractPowerFailureLog(data_raw string) []PowerFailure {
	var dataline = strings.TrimSpace(data_raw)
	var str_arr = p_pwrflr_log.Split(dataline, -1)
	var power_failure = []PowerFailure{}
	var append_id = 0
	for str_g_idx, strg_val := range str_arr {
		//println("extractPowerFailureLog : ", strg_val)

		if len(strg_val) > 0 {
			switch id := str_g_idx % powerfailure_unit_length; id {
			case 2:
				//power_failure[str_g_idx].Failure_time = strg_val
				t, _ := time.Parse(powerfailure_date_layout, strg_val)
				pf := PowerFailure{}
				pf.Failure_time = t
				append_id = len(power_failure) - 1
			case 3:
				//power_failure[str_g_idx].Failure_duration = strg_val
				val := scremove.ReplaceAllString(strg_val, "")
				power_failure[append_id].Failure_duration, _ = strconv.ParseInt(val, 10, 0)
			default:
				continue
			}
		}

	}

	return power_failure
}
func (telegram P1Telegram) String() string {
	return strings.TrimSpace(fmt.Sprintf("%#v", telegram))
}

func (filter *ValueFilter) find(data string) int {
	for regex_idx, regex_value := range filter.regex_list {
		if regex_value.MatchString(data) {
			return regex_idx
		}

	}
	return -1

}
