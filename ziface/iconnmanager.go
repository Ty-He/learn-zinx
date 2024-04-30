package ziface

type IConnManager interface {
    // insert an elemment
    Insert(IConnection)

    // erase an element
    Remove(uint32)

    // get an element
    Get(uint32) (IConnection, error) 

    // get the total size of conn
    Len() int

    // clear map
    Clear()
}
