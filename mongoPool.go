package Pool

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)


type MongoPool struct {
	pool        chan *mongo.Client //存放连接的管道, 缓存控制最大闲置连接数
	timeout     time.Duration      //超时
	uri         string             //地址
	connections int                //当前系统连接数
	poolSize    int                //最大连接数
}


func (mp *MongoPool) getContextTimeOut() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), mp.timeout)
	return ctx
}

func (mp *MongoPool) createToChan() {
	var conn *mongo.Client
	conn, e := mongo.NewClient(options.Client().ApplyURI(mp.uri))
	if e != nil {
		log.Fatalf("Create the Pool failed，err=%v", e)
	}
	e = conn.Connect(mp.getContextTimeOut())
	if e != nil {
		log.Fatalf("Create the Pool failed，err=%v", e)
	}
	mp.pool <- conn
	mp.connections++
}

func (mp *MongoPool) CloseConnection(conn *mongo.Client) error {
	select {
	case mp.pool <- conn:
		return nil
	default:
		if err := conn.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Close the Pool failed，err=%v", err)
			return err
		}
		mp.connections--
		return nil
	}
}

func (mp *MongoPool) GetConnection() (*mongo.Client, error) {
	for {
		select {
		case conn := <-mp.pool:
			err := conn.Ping(mp.getContextTimeOut(), readpref.Primary())
			if err != nil {
				log.Fatalf("获取连接池连接失败，err=%v", err)
				return nil, err
			}
			return conn, nil
		default:
			if mp.connections < mp.poolSize {
				mp.createToChan()
			}
		}
	}
}

func GetCollection(conn *mongo.Client, dbname, collection string) *mongo.Collection {
	return conn.Database(dbname).Collection(collection)
}
