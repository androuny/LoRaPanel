package backend

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Sensor struct {
	name string
	creationDate time.Time
	lastUpdate time.Time
}

type SensorData struct {
	Temperature float64 `json:"temperature"`
	Humidity float64 `json:"humidity"`
	SensorName string `json:"sensorName"`
	CreationDate time.Time `json:"creationDate"`
}

type MongoHandler struct {
	client *mongo.Client
	database string
}


func (s Sensor) asBson() (primitive.D) {
	return bson.D{{"name", s.name}, {"creationDate", s.creationDate}, {"lastUpdate", s.lastUpdate}}
}

func (s SensorData) asBson() (primitive.D) {
	return bson.D{{"temperature", s.Temperature},{"humidity", s.Humidity},{"sensorName", s.SensorName},{"creationDate", s.CreationDate}}
}

func NewMongoHandler(conn_uri string, database string) (MongoHandler, time.Duration, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn_uri))
	err = client.Ping(ctx, readpref.Primary())

	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	mhdl := MongoHandler{client, database}

	t := time.Now()
	elapsed := t.Sub(start)

	return mhdl, elapsed, err
}

func (mhdl MongoHandler) NewSensor(username string) (time.Duration, error) {
	start := time.Now()

	collection := mhdl.client.Database(mhdl.database).Collection("Sensors")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doesNameExist, _, _ := mhdl.DoesSensorNameExist(username)
	if (doesNameExist) { // name is taken
		err := errors.New("name for sensor already taken")
		t := time.Now()
		elapsed := t.Sub(start)
		return elapsed, err
	} else { // name isnt taken
		s := Sensor{username, time.Now(), time.Now()}
		data := s.asBson()
	
		_, err := collection.InsertOne(ctx, data)
		if err != nil {log.Fatalf("[DB Sensor] [INSERT] %v", err)}
	
		t := time.Now()
		elapsed := t.Sub(start)
		return elapsed, err
		}
}

func (mhdl MongoHandler) DoesSensorNameExist(username string) (bool, time.Duration, error) {
	start := time.Now()

	isNameFound := false
	s := Sensor{}

	collection := mhdl.client.Database(mhdl.database).Collection("Sensors")
	filter := bson.D{{"name", username}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&s)
	if err == mongo.ErrNoDocuments {
		// Do something when no record was found
		isNameFound = false
	} else { 
		// something is found
		isNameFound = true
	}

	t := time.Now()
	elapsed := t.Sub(start)
	return isNameFound, elapsed, err
}

func (mhdl MongoHandler) NewSensorData(username string, temp float64, humidity float64) (time.Duration, error) {
	start := time.Now()

	collection := mhdl.client.Database(mhdl.database).Collection("SensorData")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doesNameExist, _, _ := mhdl.DoesSensorNameExist(username)
	if (doesNameExist) { // name is taken
		s := SensorData{Temperature: temp, Humidity: humidity, SensorName: username, CreationDate: time.Now()}
		data := s.asBson()
	
		_, err := collection.InsertOne(ctx, data)
		if err != nil {log.Fatalf("[DB SensorData] [INSERT] %v", err)}
	
		t := time.Now()
		elapsed := t.Sub(start)
		return elapsed, err

	} else { // name isnt taken
		err := errors.New("sensor name doesnt exist")
		t := time.Now()
		elapsed := t.Sub(start)
		return elapsed, err
	}
}

func (mhdl MongoHandler) GetAllSensorData(sensorname string) ([]SensorData, time.Duration, error) {
	start := time.Now()

	collection := mhdl.client.Database(mhdl.database).Collection("SensorData")
	filter := bson.D{{"sensorName", sensorname}}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	var results []SensorData
	err = cur.All(ctx, &results)
	if err != nil {
		panic(err)
	}


	t := time.Now()
	elapsed := t.Sub(start)
	return results, elapsed, err
}

// start making database