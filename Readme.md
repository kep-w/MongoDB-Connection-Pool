# Golang MongoDB connection pool

### Efficiently control the number of database connections and realize automatic connection management 

### Details:

Provide connection pool management through struct MongoPool. 

The connection pool parameters are set during initialization, including: 
- cached channel for storing connections, and the cache controls the maximum number of idle connections. 
- Database address, 
- timeout period, 
- maximum number of connections. 
for example:
`Golang
func init() {
	once.Do(func() {
		DBPool = &MongoPool{
			pool:        make(chan *mongo.Client, 10), 
			connections: 0,                            
			timeout:     10 * time.Second,
			uri:         "mongodb://user:password@127.0.0.1:27017/dbname",
			poolSize:    10, 
		}
	})
}
`

Note: You can use function once.Do to ensure that the connection pool is only created once.
