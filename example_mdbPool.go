// 数据库交互
package Pool

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"sync"
	"time"
)

var (
	DBPool *MongoPool
	once   sync.Once
)

func init() {
	once.Do(func() {
		DBPool = &MongoPool{
			pool:        make(chan *mongo.Client, 10), //最大闲置连接数
			connections: 0,                            //当前程序连接数
			timeout:     10 * time.Second,
			uri:         "mongodb://user:password@127.0.0.1:27017/dbname",
			poolSize:    10, //最大连接数 = 当前连接数+闲置连接数
		}
	})
}

type MongoData struct {
	Id    string `bson:"id"`
	Name  string `bson:"name"`
	Other string `bson:"other"`
}

func Find() (result MongoData, err error) {
	conn, err := DBPool.GetConnection()
	if err != nil {
		log.Printf("获取数据库连接失败，err=%v", err)
		return
	}
	defer DBPool.CloseConnection(conn)
	collection := GetCollection(conn, "testdb", "test")
	err = collection.FindOne(context.TODO(), bson.D{{"id", 1}}).Decode(&result)
	if err != nil {
		log.Printf("查询失败，err=%v", err)
		return
	}
	log.Printf("数据库查询成功，result=%#v", result)
	return
}
