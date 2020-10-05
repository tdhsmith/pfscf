package content

import (
	"fmt"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

type entryYAML struct {
	e Entry
}

func (ey *entryYAML) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	// determine type of entry
	type entryTypeYAML struct{ Type string }
	var ety entryTypeYAML
	err = unmarshal(&ety)
	if err != nil {
		return err
	}
	if !utils.IsSet(ety.Type) {
		return fmt.Errorf("Missing or empty 'type' field")
	}

	// read concrete object based on type information
	switch ety.Type {
	case typeText:
		ey.e = newText()
	case typeRectangle:
		ey.e = newRectangle()
	case typeTrigger:
		ey.e = newTrigger()
	case typeChoice:
		ey.e = newChoice()
	case typeMultiline:
		ey.e = newMultiline()
	default:
		return fmt.Errorf("Unknown content type: '%v'", ety.Type)
	}

	if err = unmarshal(ey.e); err != nil {
		return err
	}

	return nil
}
