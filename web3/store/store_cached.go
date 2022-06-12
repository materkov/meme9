package store

type CachedItem struct {
	object Object
	err    error
}

type CachedStore struct {
	Store Store

	Needed   map[int]bool
	ObjCache map[int]CachedItem
}

func (c *CachedStore) ObjGet(id int) (Object, error) {
	item, ok := c.ObjCache[id]
	if ok {
		return item.object, item.err
	}

	c.Need(id)
	c.preload()

	item = c.ObjCache[id]
	return item.object, item.err
}

func (c *CachedStore) preload() {
	if len(c.Needed) == 0 {
		return
	}

	neededList := make([]int, len(c.Needed))

	idx := 0
	for objID := range c.Needed {
		neededList[idx] = objID
		idx++
	}

	c.Needed = map[int]bool{}

	objects, _ := c.Store.ObjGet(neededList)

	for _, id := range neededList {
		c.ObjCache[id] = CachedItem{object: objects[id]}
	}
}

func (c *CachedStore) Need(id int) {
	if _, ok := c.ObjCache[id]; !ok {
		c.Needed[id] = true
	}
}
