// Code generated by ObjectBox; DO NOT EDIT.
// Learn more about defining entities and generating this file - visit https://golang.objectbox.io/entity-annotations

package storage

import (
	"errors"
	"github.com/google/flatbuffers/go"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/objectbox/fbutils"
)

type accountServerMapB_EntityInfo struct {
	objectbox.Entity
	Uid uint64
}

var AccountServerMapBBinding = accountServerMapB_EntityInfo{
	Entity: objectbox.Entity{
		Id: 1,
	},
	Uid: 1090041196471887670,
}

// AccountServerMapB_ contains type-based Property helpers to facilitate some common operations such as Queries.
var AccountServerMapB_ = struct {
	AccountServerMap_Servers *objectbox.PropertyString
	Id                       *objectbox.PropertyUint64
	Uid                      *objectbox.PropertyString
}{
	AccountServerMap_Servers: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     1,
			Entity: &AccountServerMapBBinding.Entity,
		},
	},
	Id: &objectbox.PropertyUint64{
		BaseProperty: &objectbox.BaseProperty{
			Id:     2,
			Entity: &AccountServerMapBBinding.Entity,
		},
	},
	Uid: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     3,
			Entity: &AccountServerMapBBinding.Entity,
		},
	},
}

// GeneratorVersion is called by ObjectBox to verify the compatibility of the generator used to generate this code
func (accountServerMapB_EntityInfo) GeneratorVersion() int {
	return 6
}

// AddToModel is called by ObjectBox during model build
func (accountServerMapB_EntityInfo) AddToModel(model *objectbox.Model) {
	model.Entity("AccountServerMapB", 1, 1090041196471887670)
	model.Property("AccountServerMap_Servers", 9, 1, 6802529252641519899)
	model.Property("Id", 6, 2, 4962380468247370576)
	model.PropertyFlags(1)
	model.Property("Uid", 9, 3, 1765808335051830653)
	model.PropertyFlags(2048)
	model.PropertyIndex(1, 976550826644671096)
	model.EntityLastPropertyId(3, 1765808335051830653)
}

// GetId is called by ObjectBox during Put operations to check for existing ID on an object
func (accountServerMapB_EntityInfo) GetId(object interface{}) (uint64, error) {
	return object.(*AccountServerMapB).Id, nil
}

// SetId is called by ObjectBox during Put to update an ID on an object that has just been inserted
func (accountServerMapB_EntityInfo) SetId(object interface{}, id uint64) error {
	object.(*AccountServerMapB).Id = id
	return nil
}

// PutRelated is called by ObjectBox to put related entities before the object itself is flattened and put
func (accountServerMapB_EntityInfo) PutRelated(ob *objectbox.ObjectBox, object interface{}, id uint64) error {
	return nil
}

// Flatten is called by ObjectBox to transform an object to a FlatBuffer
func (accountServerMapB_EntityInfo) Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) error {
	obj := object.(*AccountServerMapB)
	var offsetAccountServerMap_Servers = fbutils.CreateStringOffset(fbb, obj.AccountServerMap.Server)
	var offsetUid = fbutils.CreateStringOffset(fbb, obj.Uid)

	// build the FlatBuffers object
	fbb.StartObject(3)
	fbutils.SetUOffsetTSlot(fbb, 0, offsetAccountServerMap_Servers)
	fbutils.SetUint64Slot(fbb, 1, id)
	fbutils.SetUOffsetTSlot(fbb, 2, offsetUid)
	return nil
}

// Load is called by ObjectBox to load an object from a FlatBuffer
func (accountServerMapB_EntityInfo) Load(ob *objectbox.ObjectBox, bytes []byte) (interface{}, error) {
	if len(bytes) == 0 { // sanity check, should "never" happen
		return nil, errors.New("can't deserialize an object of type 'AccountServerMapB' - no data received")
	}

	var table = &flatbuffers.Table{
		Bytes: bytes,
		Pos:   flatbuffers.GetUOffsetT(bytes),
	}

	var propId = table.GetUint64Slot(6, 0)

	return &AccountServerMapB{
		AccountServerMap: AccountServerMap{
			Server: fbutils.GetStringSlot(table, 4),
		},
		Id:  propId,
		Uid: fbutils.GetStringSlot(table, 8),
	}, nil
}

// MakeSlice is called by ObjectBox to construct a new slice to hold the read objects
func (accountServerMapB_EntityInfo) MakeSlice(capacity int) interface{} {
	return make([]*AccountServerMapB, 0, capacity)
}

// AppendToSlice is called by ObjectBox to fill the slice of the read objects
func (accountServerMapB_EntityInfo) AppendToSlice(slice interface{}, object interface{}) interface{} {
	if object == nil {
		return append(slice.([]*AccountServerMapB), nil)
	}
	return append(slice.([]*AccountServerMapB), object.(*AccountServerMapB))
}

// Box provides CRUD access to AccountServerMapB objects
type AccountServerMapBBox struct {
	*objectbox.Box
}

// BoxForAccountServerMapB opens a box of AccountServerMapB objects
func BoxForAccountServerMapB(ob *objectbox.ObjectBox) *AccountServerMapBBox {
	return &AccountServerMapBBox{
		Box: ob.InternalBox(1),
	}
}

// Put synchronously inserts/updates a single object.
// In case the Id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the AccountServerMapB.Id property on the passed object will be assigned the new ID as well.
func (box *AccountServerMapBBox) Put(object *AccountServerMapB) (uint64, error) {
	return box.Box.Put(object)
}

// Insert synchronously inserts a single object. As opposed to Put, Insert will fail if given an ID that already exists.
// In case the Id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the AccountServerMapB.Id property on the passed object will be assigned the new ID as well.
func (box *AccountServerMapBBox) Insert(object *AccountServerMapB) (uint64, error) {
	return box.Box.Insert(object)
}

// Update synchronously updates a single object.
// As opposed to Put, Update will fail if an object with the same ID is not found in the database.
func (box *AccountServerMapBBox) Update(object *AccountServerMapB) error {
	return box.Box.Update(object)
}

// PutAsync asynchronously inserts/updates a single object.
// Deprecated: use box.Async().Put() instead
func (box *AccountServerMapBBox) PutAsync(object *AccountServerMapB) (uint64, error) {
	return box.Box.PutAsync(object)
}

// PutMany inserts multiple objects in single transaction.
// In case Ids are not set on the objects, they would be assigned automatically (auto-increment).
//
// Returns: IDs of the put objects (in the same order).
// When inserting, the AccountServerMapB.Id property on the objects in the slice will be assigned the new IDs as well.
//
// Note: In case an error occurs during the transaction, some of the objects may already have the AccountServerMapB.Id assigned
// even though the transaction has been rolled back and the objects are not stored under those IDs.
//
// Note: The slice may be empty or even nil; in both cases, an empty IDs slice and no error is returned.
func (box *AccountServerMapBBox) PutMany(objects []*AccountServerMapB) ([]uint64, error) {
	return box.Box.PutMany(objects)
}

// Get reads a single object.
//
// Returns nil (and no error) in case the object with the given ID doesn't exist.
func (box *AccountServerMapBBox) Get(id uint64) (*AccountServerMapB, error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		return nil, nil
	}
	return object.(*AccountServerMapB), nil
}

// GetMany reads multiple objects at once.
// If any of the objects doesn't exist, its position in the return slice is nil
func (box *AccountServerMapBBox) GetMany(ids ...uint64) ([]*AccountServerMapB, error) {
	objects, err := box.Box.GetMany(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*AccountServerMapB), nil
}

// GetManyExisting reads multiple objects at once, skipping those that do not exist.
func (box *AccountServerMapBBox) GetManyExisting(ids ...uint64) ([]*AccountServerMapB, error) {
	objects, err := box.Box.GetManyExisting(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*AccountServerMapB), nil
}

// GetAll reads all stored objects
func (box *AccountServerMapBBox) GetAll() ([]*AccountServerMapB, error) {
	objects, err := box.Box.GetAll()
	if err != nil {
		return nil, err
	}
	return objects.([]*AccountServerMapB), nil
}

// Remove deletes a single object
func (box *AccountServerMapBBox) Remove(object *AccountServerMapB) error {
	return box.Box.Remove(object)
}

// RemoveMany deletes multiple objects at once.
// Returns the number of deleted object or error on failure.
// Note that this method will not fail if an object is not found (e.g. already removed).
// In case you need to strictly check whether all of the objects exist before removing them,
// you can execute multiple box.Contains() and box.Remove() inside a single write transaction.
func (box *AccountServerMapBBox) RemoveMany(objects ...*AccountServerMapB) (uint64, error) {
	var ids = make([]uint64, len(objects))
	for k, object := range objects {
		ids[k] = object.Id
	}
	return box.Box.RemoveIds(ids...)
}

// Creates a query with the given conditions. Use the fields of the AccountServerMapB_ struct to create conditions.
// Keep the *AccountServerMapBQuery if you intend to execute the query multiple times.
// Note: this function panics if you try to create illegal queries; e.g. use properties of an alien type.
// This is typically a programming error. Use QueryOrError instead if you want the explicit error check.
func (box *AccountServerMapBBox) Query(conditions ...objectbox.Condition) *AccountServerMapBQuery {
	return &AccountServerMapBQuery{
		box.Box.Query(conditions...),
	}
}

// Creates a query with the given conditions. Use the fields of the AccountServerMapB_ struct to create conditions.
// Keep the *AccountServerMapBQuery if you intend to execute the query multiple times.
func (box *AccountServerMapBBox) QueryOrError(conditions ...objectbox.Condition) (*AccountServerMapBQuery, error) {
	if query, err := box.Box.QueryOrError(conditions...); err != nil {
		return nil, err
	} else {
		return &AccountServerMapBQuery{query}, nil
	}
}

// Async provides access to the default Async Box for asynchronous operations. See AccountServerMapBAsyncBox for more information.
func (box *AccountServerMapBBox) Async() *AccountServerMapBAsyncBox {
	return &AccountServerMapBAsyncBox{AsyncBox: box.Box.Async()}
}

// AccountServerMapBAsyncBox provides asynchronous operations on AccountServerMapB objects.
//
// Asynchronous operations are executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "execute & forget:" you gain faster put/remove operations as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
// In situations with (extremely) high async load, an async method may be throttled (~1ms) or delayed up to 1 second.
// In the unlikely event that the object could still not be enqueued (full queue), an error will be returned.
//
// Note that async methods do not give you hard durability guarantees like the synchronous Box provides.
// There is a small time window in which the data may not have been committed durably yet.
type AccountServerMapBAsyncBox struct {
	*objectbox.AsyncBox
}

// AsyncBoxForAccountServerMapB creates a new async box with the given operation timeout in case an async queue is full.
// The returned struct must be freed explicitly using the Close() method.
// It's usually preferable to use AccountServerMapBBox::Async() which takes care of resource management and doesn't require closing.
func AsyncBoxForAccountServerMapB(ob *objectbox.ObjectBox, timeoutMs uint64) *AccountServerMapBAsyncBox {
	var async, err = objectbox.NewAsyncBox(ob, 1, timeoutMs)
	if err != nil {
		panic("Could not create async box for entity ID 1: %s" + err.Error())
	}
	return &AccountServerMapBAsyncBox{AsyncBox: async}
}

// Put inserts/updates a single object asynchronously.
// When inserting a new object, the Id property on the passed object will be assigned the new ID the entity would hold
// if the insert is ultimately successful. The newly assigned ID may not become valid if the insert fails.
func (asyncBox *AccountServerMapBAsyncBox) Put(object *AccountServerMapB) (uint64, error) {
	return asyncBox.AsyncBox.Put(object)
}

// Insert a single object asynchronously.
// The Id property on the passed object will be assigned the new ID the entity would hold if the insert is ultimately
// successful. The newly assigned ID may not become valid if the insert fails.
// Fails silently if an object with the same ID already exists (this error is not returned).
func (asyncBox *AccountServerMapBAsyncBox) Insert(object *AccountServerMapB) (id uint64, err error) {
	return asyncBox.AsyncBox.Insert(object)
}

// Update a single object asynchronously.
// The object must already exists or the update fails silently (without an error returned).
func (asyncBox *AccountServerMapBAsyncBox) Update(object *AccountServerMapB) error {
	return asyncBox.AsyncBox.Update(object)
}

// Remove deletes a single object asynchronously.
func (asyncBox *AccountServerMapBAsyncBox) Remove(object *AccountServerMapB) error {
	return asyncBox.AsyncBox.Remove(object)
}

// Query provides a way to search stored objects
//
// For example, you can find all AccountServerMapB which Id is either 42 or 47:
//
//	box.Query(AccountServerMapB_.Id.In(42, 47)).Find()
type AccountServerMapBQuery struct {
	*objectbox.Query
}

// Find returns all objects matching the query
func (query *AccountServerMapBQuery) Find() ([]*AccountServerMapB, error) {
	objects, err := query.Query.Find()
	if err != nil {
		return nil, err
	}
	return objects.([]*AccountServerMapB), nil
}

// Offset defines the index of the first object to process (how many objects to skip)
func (query *AccountServerMapBQuery) Offset(offset uint64) *AccountServerMapBQuery {
	query.Query.Offset(offset)
	return query
}

// Limit sets the number of elements to process by the query
func (query *AccountServerMapBQuery) Limit(limit uint64) *AccountServerMapBQuery {
	query.Query.Limit(limit)
	return query
}
