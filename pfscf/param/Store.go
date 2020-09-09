package param

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

// Store stores a list of parameter descriptions
type Store map[string]Entry

// NewStore creates a new store.
func NewStore() (s Store) {
	s = make(Store, 0)
	return s
}

// add adds an entry to the store and also sets the ID on the entry
func (s *Store) add(id string, e Entry) {
	utils.Assert(!utils.IsSet(e.ID()) || id == e.ID(), "ID must not be set here")
	if _, exists := (*s)[id]; exists {
		utils.Assert(false, "As we only call this from a map in yaml, duplicates should not occur")
	}

	if !utils.IsSet(e.ID()) {
		e.setID(id)
	}
	(*s)[id] = e
}

// Get returns the Entry matching the provided id.
func (s *Store) Get(id string) (e Entry, exists bool) {
	e, exists = (*s)[id]
	return
}

// UnmarshalYAML unmarshals a Parameter Store
func (s *Store) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type storeYAML map[string]entryYAML

	sy := make(storeYAML, 0)

	err = unmarshal(&sy)
	if err != nil {
		return err
	}

	*s = NewStore()
	for key, value := range sy {
		s.add(key, value.e)
	}

	return nil
}

// InheritFrom inherits entries from another param store. An error is returned in case
// an entry exists in both stores.
func (s *Store) InheritFrom(other *Store) (err error) {
	for otherID, otherEntry := range *other {
		if _, exists := (*s)[otherID]; exists {
			return fmt.Errorf("Duplicate parameter ID '%v' found while inheriting", otherID)
		}
		s.add(otherID, otherEntry.deepCopy())
	}

	return nil
}

// IsValid checks whether all entries are valid.
func (s *Store) IsValid() (err error) {
	for _, entry := range *s {
		if err = entry.isValid(); err != nil {
			return fmt.Errorf("Error while validating parameter definition '%v': %v", entry.ID(), err)
		}
	}
	return nil
}

// ValidateAndProcessArgs checks whether all
func (s *Store) ValidateAndProcessArgs(as *args.Store) (err error) {
	for _, key := range as.GetKeys() {
		paramEntry, pExists := s.Get(key)

		// check that all entries in the arg store have a corresponding parameter entry
		if !pExists {
			return fmt.Errorf("Error while validating argument '%v': No corresponding parameter registered for template", key)
		}

		// ask each type whether the provided argument is valid, and add entries to argStore if required
		if err = paramEntry.validateAndProcessArgs(as); err != nil {
			return fmt.Errorf("Error while validating argument '%v': %v", key, err)
		}
	}

	return nil
}

// GetExampleArguments returns an array containing all keys and example values for all parameters.
// The result can be passed to the ArgStore.
func (s *Store) GetExampleArguments() (result []string) {
	result = make([]string, 0)

	for _, entry := range *s {
		result = append(result, fmt.Sprintf("%v=%v", entry.ID(), entry.Example()))
	}

	return result
}

// GetSortedKeys returns the list of keys contained in this store as sorted list.
func (s *Store) GetSortedKeys() (result []string) {
	result = make([]string, 0)

	for key := range *s {
		result = append(result, key)
	}

	sort.Strings(result)

	return result
}

// Describe returns a short textual description of all parameters contained in this store.
// It returns the description as a multi-line string.
func (s *Store) Describe(verbose bool) (result string) {
	var sb strings.Builder

	for _, key := range s.GetSortedKeys() {
		entry, _ := s.Get(key)
		fmt.Fprintf(&sb, entry.describe(verbose))
	}

	return sb.String()
}
