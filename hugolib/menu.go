// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import "sort"

type MenuEntry struct {
	Url      string
	Name     string
	Menu     string
	PreName  string
	PostName string
	Weight   int
	Parent   string
	Children Menu
}

type Menu []*MenuEntry
type Menus map[string]*Menu
type PageMenus map[string]*MenuEntry

func (me *MenuEntry) AddChild(child *MenuEntry) {
	me.Children = append(me.Children, child)
	me.Children.Sort()
}

func (me *MenuEntry) HasChildren() bool {
	return me.Children != nil
}

//func (me *MenuEntry) RelUrl() string {
//link, err := p.permalink()
//if err != nil {
//return "", err
//}

//link.Scheme = ""
//link.Host = ""
//link.User = nil
//link.Opaque = ""
//return link.String(), nil
//}

func (m Menu) Add(me *MenuEntry) Menu {
	app := func(slice Menu, x ...*MenuEntry) Menu {
		n := len(slice) + len(x)
		if n > cap(slice) {
			size := cap(slice) * 2
			if size < n {
				size = n
			}
			new := make(Menu, size)
			copy(new, slice)
			slice = new
		}
		slice = slice[0:n]
		copy(slice[n-len(x):], x)
		return slice
	}

	m = app(m, me)
	m.Sort()
	return m
}

/*
 * Implementation of a custom sorter for Menu
 */

// A type to implement the sort interface for Menu
type MenuSorter struct {
	menu Menu
	by   MenuEntryBy
}

// Closure used in the Sort.Less method.
type MenuEntryBy func(m1, m2 *MenuEntry) bool

func (by MenuEntryBy) Sort(menu Menu) {
	ms := &MenuSorter{
		menu: menu,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ms)
}

var DefaultMenuEntrySort = func(m1, m2 *MenuEntry) bool {
	if m1.Weight == m2.Weight {
		return m1.Name < m2.Name
	} else {
		return m1.Weight < m2.Weight
	}
}

func (ms *MenuSorter) Len() int      { return len(ms.menu) }
func (ms *MenuSorter) Swap(i, j int) { ms.menu[i], ms.menu[j] = ms.menu[j], ms.menu[i] }

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (ms *MenuSorter) Less(i, j int) bool { return ms.by(ms.menu[i], ms.menu[j]) }

func (p Menu) Sort() {
	MenuEntryBy(DefaultMenuEntrySort).Sort(p)
}

func (p Menu) Limit(n int) Menu {
	if len(p) < n {
		return p[0:n]
	} else {
		return p
	}
}

func (p Menu) ByWeight() Menu {
	MenuEntryBy(DefaultMenuEntrySort).Sort(p)
	return p
}

func (p Menu) ByName() Menu {
	title := func(m1, m2 *MenuEntry) bool {
		return m1.Name < m2.Name
	}

	MenuEntryBy(title).Sort(p)
	return p
}

func (m Menu) Reverse() Menu {
	for i, j := 0, len(m)-1; i < j; i, j = i+1, j-1 {
		m[i], m[j] = m[j], m[i]
	}

	return m
}