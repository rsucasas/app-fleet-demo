package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

var (
	RabbitMq_Host string = "localhost"
	RabbitMq_Port string = "5672"
	RabbitMq_User string = "guest"
	RabbitMq_Pass string = "guest"
)

func main() {
	flag.StringVar(&RabbitMq_Host, "rabbitmq-host", LookupEnvOrString("RABBITMQ_HOST", RabbitMq_Host), "RabbitMQ Host")
	flag.StringVar(&RabbitMq_Port, "rabbitmq-port", LookupEnvOrString("RABBITMQ_PORT", RabbitMq_Port), "RabbitMQ port")
	flag.StringVar(&RabbitMq_User, "rabbitmq-user", LookupEnvOrString("RABBITMQ_USER", RabbitMq_User), "RabbitMQ User")
	flag.StringVar(&RabbitMq_Pass, "rabbitmq-pass", LookupEnvOrString("RABBITMQ_PASS", RabbitMq_Pass), "RabbitMQ Password")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/publisher", homeHandler)
	mux.HandleFunc("/publisher/publish", publishHandler)
	fileServer := http.FileServer(http.Dir("./assets/"))
	mux.Handle("/publisher/assets/", http.StripPrefix("/publisher/assets", fileServer))

	// start web server
	log.Println("Starting RabbitMQ Demo app - Publisher on: 4001")
	log.Printf("RabbitMQ Instance: %s:%s\n", RabbitMq_Host, RabbitMq_Port)
	log.Printf("RabbitMQ User: %s\n", RabbitMq_User)

	err := http.ListenAndServe("0.0.0.0:4001", mux)
	log.Fatal(err)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./home.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		queueValue := r.PostForm.Get("queue")
		contentValue := r.PostForm.Get("content")
		publish(queueValue, contentValue)
		http.Redirect(w, r, "/publisher", http.StatusSeeOther)
	} else {
		http.Error(w, "Internal Server Error", 500)
	}
}

func publish(queueValue, contentValue string) {
	log.Printf("Publish to queue [%s] message [%s]\n", queueValue, contentValue)
	conn, err := amqp.Dial(RabbitMQInstanceConnectionPath())
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer conn.Close()

	log.Println("Successfully Connected to our RabbitMQ Instance")

	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queueValue, // queue name
		false,      // durable
		false,      // auto delete
		false,      // exclusive
		false,      // no wait
		nil,        // arguments
	)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println(queue)
	err = channel.Publish(
		"",
		queueValue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(contentValue),
		},
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	log.Printf("Successfully Published Message to Queue [%s]\n", queueValue)
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func RabbitMQInstanceConnectionPath() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", RabbitMq_User, RabbitMq_Pass, RabbitMq_Host, RabbitMq_Port)
}
