package adifparser

import (
	"fmt"
)

const (
	ADIFBoolean = iota
	ADIFNumber
	ADIFString
	ADIFDate
	ADIFTime
	ADIFLocation
)

type fieldMetadata struct {
	name     string
	datatype int
}

var typeCodeMap = map[byte]int{
	'A': ADIFString,
	'B': ADIFBoolean,
	'N': ADIFNumber,
	'S': ADIFString,
	'D': ADIFDate,
	'T': ADIFTime,
	'M': ADIFString,
	'L': ADIFLocation,
}

var ADIFfieldOrder []string
var ADIFfieldInfo map[string]fieldMetadata

func addField(name string, datatype int) {
	for _, n := range ADIFfieldOrder {
		if name == n {
			panic(fmt.Sprintf("Duplicate field name %s.", name))
		}
	}
	ADIFfieldOrder = append(ADIFfieldOrder, name)
	ADIFfieldInfo[name] = fieldMetadata{name, datatype}
}

func isStandardADIFField(name string) bool {
	for _, n := range ADIFfieldOrder {
		if name == n {
			return true
		}
	}
	return false
}

func init() {
	ADIFfieldInfo = make(map[string]fieldMetadata)

	// Common fields first
	addField("call", ADIFString)
	addField("station_callsign", ADIFString)
	addField("band", ADIFString)
	addField("freq", ADIFNumber)
	addField("mode", ADIFString)
	addField("qso_date", ADIFDate)
	addField("qso_date_off", ADIFDate)
	addField("time_on", ADIFTime)
	addField("time_off", ADIFTime)
	// Other fields alphabetically
	addField("address", ADIFString)
	addField("age", ADIFNumber)
	addField("a_index", ADIFNumber)
	addField("ant_az", ADIFNumber)
	addField("ant_el", ADIFNumber)
	addField("ant_path", ADIFString)
	addField("arrl_sect", ADIFString)
	addField("band_rx", ADIFString)
	addField("check", ADIFString)
	addField("class", ADIFString)
	addField("cnty", ADIFString)
	addField("comment", ADIFString)
	addField("cont", ADIFString)
	addField("contacted_op", ADIFString)
	addField("contest_id", ADIFString)
	addField("country", ADIFString)
	addField("cqz", ADIFString)
	addField("credit_submitted", ADIFString)
	addField("credit_granted", ADIFString)
	addField("distance", ADIFNumber)
	addField("dxcc", ADIFString)
	addField("email", ADIFString)
	addField("eq_call", ADIFString)
	addField("eqsl_qslrdate", ADIFString)
	addField("eqsl_qslsdate", ADIFString)
	addField("eqsl_qsl_rcvd", ADIFString)
	addField("eqsl_qsl_sent", ADIFString)
	addField("force_init", ADIFBoolean)
	addField("freq_rx", ADIFNumber)
	addField("gridsquare", ADIFString)
	addField("guest_op", ADIFString)
	addField("iota", ADIFString)
	addField("iota_island_id", ADIFString)
	addField("ituz", ADIFNumber)
	addField("k_index", ADIFNumber)
	addField("lat", ADIFLocation)
	addField("lon", ADIFLocation)
	addField("lotw_qslrdate", ADIFDate)
	addField("lotw_qslsdate", ADIFDate)
	addField("lotw_qsl_rcvd", ADIFString)
	addField("lotw_qsl_sent", ADIFString)
	addField("max_bursts", ADIFNumber)
	addField("ms_shower", ADIFString)
	addField("my_city", ADIFString)
	addField("my_cnty", ADIFString)
	addField("my_country", ADIFString)
	addField("my_cq_zone", ADIFNumber)
	addField("my_gridsquare", ADIFString)
	addField("my_iota", ADIFString)
	addField("my_iota_island_id", ADIFString)
	addField("my_itu_zone", ADIFNumber)
	addField("my_lat", ADIFLocation)
	addField("my_lon", ADIFLocation)
	addField("my_name", ADIFString)
	addField("my_postal_code", ADIFString)
	addField("my_rig", ADIFString)
	addField("my_sig", ADIFString)
	addField("my_sig_info", ADIFString)
	addField("my_state", ADIFString)
	addField("my_street", ADIFString)
	addField("name", ADIFString)
	addField("notes", ADIFString)
	addField("nr_bursts", ADIFNumber)
	addField("nr_pings", ADIFNumber)
	addField("operator", ADIFString)
	addField("owner_callsign", ADIFString)
	addField("pfx", ADIFString)
	addField("precedence", ADIFString)
	addField("prop_mode", ADIFString)
	addField("public_key", ADIFString)
	addField("qslmsg", ADIFString)
	addField("qslrdate", ADIFDate)
	addField("qslsdate", ADIFDate)
	addField("qsl_rcvd", ADIFString)
	addField("qsl_rcvd_via", ADIFString)
	addField("qsl_sent", ADIFString)
	addField("qsl_sent_via", ADIFString)
	addField("qsl_via", ADIFString)
	addField("qso_complete", ADIFString)
	addField("qso_random", ADIFBoolean)
	addField("qth", ADIFString)
	addField("rig", ADIFString)
	addField("rst_rcvd", ADIFString)
	addField("rst_sent", ADIFString)
	addField("rx_pwr", ADIFNumber)
	addField("sat_mode", ADIFString)
	addField("sat_name", ADIFString)
	addField("sfi", ADIFNumber)
	addField("sig", ADIFString)
	addField("sig_info", ADIFString)
	addField("srx", ADIFNumber)
	addField("srx_string", ADIFString)
	addField("state", ADIFString)
	addField("stx", ADIFNumber)
	addField("stx_string", ADIFString)
	addField("swl", ADIFBoolean)
	addField("ten_ten", ADIFNumber)
	addField("tx_pwr", ADIFNumber)
	addField("ve_prov", ADIFString)
	addField("web", ADIFString)
}
