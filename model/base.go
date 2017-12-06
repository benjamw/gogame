package model

import (
	"context"

	"github.com/benjamw/golibs/db"
	"google.golang.org/appengine/datastore"
)

type Base struct {
	key   *datastore.Key `datastore:"-"`
	isNew bool           `datastore:"-"`
}

// EntityType returns the struct type used in the datastore
func (m *Base) EntityType() string {
	return "BASE"
}

// GetKey returns the datastore key for the struct
func (m *Base) GetKey() *datastore.Key {
	return m.key
}

// SetKey sets the datastore key on the struct
func (m *Base) SetKey(key *datastore.Key) error {
	m.key = key
	return nil
}

// IsNew returns true if this is a new (non-saved) struct
func (m *Base) IsNew() bool {
	if m.GetKey() == nil {
		return true
	}
	return m.isNew
}

// SetIsNew sets the isNew flag on the struct
func (m *Base) SetIsNew(isNew bool) {
	m.isNew = isNew
}

// PreSave sets some basic info before continuing on to Save
// PreSave is called in db.Save, and if an error is returned,
// will halt saving the struct
func (m *Base) PreSave(ctx context.Context) error {
	return nil
}

// PostSave performs any actions needed after the struct has
// been saved to the datastore
func (m *Base) PostSave(ctx context.Context) error {
	return nil
}

// PostLoad performs any actions after the struct has been
// loaded from the datastore
func (m *Base) PostLoad(ctx context.Context) error {
	fullKey := m.GetKey()
	if fullKey == nil {
		err := &db.MissingKeyError{}
		return err
	}

	return nil
}

// Transform will alter the datastore struct if needed
func (m *Base) Transform(ctx context.Context, pl datastore.PropertyList) error {
	return nil
}

// PreDelete performs any tasks need before the struct is
// deleted from the datastore.
// PreDelete gets called in db.Delete, and if an error is returned,
// will halt deletion of the struct
func (m *Base) PreDelete(ctx context.Context) error {
	return nil
}

// Collect gets a properly sized []db.Model ready for use in db.LoadMultiX
func (m *Base) Collect(num int) []db.Model {
	return nil

	/* concrete versions should look like the following:
	// replace 'Foo' with the concrete model name

	l := make([]*Foo, n)
	ml := make([]db.Model, n)
	for k := range l {
		v := new(Foo)
		ml[k] = db.Model(v)
	}

	return ml

	// also, see Scatter below

	*/
}

/*

A related function to go along with the b.Collect above
but not actually part of the model because it acts on FooList ( []Foo ):
// replace 'Foo' with the concrete model name

// Scatter splits the abstract []db.Model into concrete items in FooList
func (l *FooList) Scatter(ml []db.Model) {
	*l = make(FooList, 0)
	for k := range ml {
		v, ok := ml[k].(*Foo)
		if !ok {
			continue
		}

		*l = append(*l, *v)
	}

	return
}

*/
