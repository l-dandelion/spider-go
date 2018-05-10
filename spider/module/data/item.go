package data

/*
 * define the item type
 */
type Item map[string]interface{}

/*
 * check the item
 */
func (item Item) Valid() bool {
	return item != nil
}
