package adifparser

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
	ADIFfieldOrder = append(ADIFfieldOrder, name)
	ADIFfieldInfo[name] = fieldMetadata{name, datatype}
}

func init() {
	addField("address", ADIFString)
}
