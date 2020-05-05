# Mongo Collection

Extends mongodb driver for go --> [https://github.com/mongodb/mongo-go-driver](https://github.com/mongodb/mongo-go-driver)

Adds methods to read `SingleResult` and `Cursor` returned by find one/many

Abstracts the connection to mongo using ICollection so a test implementation can be used
