package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/viper"
	"solace.dev/go/messaging"
	"solace.dev/go/messaging/pkg/solace/config"
	"solace.dev/go/messaging/pkg/solace/message"
	"solace.dev/go/messaging/pkg/solace/resource"
)

func main() {
	emoji.Println(":green_circle: Starting subscriber")

	// Read configuration
	viper.SetEnvPrefix("SOLACE")
	viper.AutomaticEnv()

	// Connect to event broker
	brokerConfig := config.ServicePropertyMap{
		config.TransportLayerPropertyHost:                viper.Get("HOST"),
		config.ServicePropertyVPNName:                    viper.Get("VPN"),
		config.AuthenticationPropertySchemeBasicUserName: viper.Get("USERNAME"),
		config.AuthenticationPropertySchemeBasicPassword: viper.Get("PASSWORD"),
	}

	messagingService, err := messaging.NewMessagingServiceBuilder().
		FromConfigurationProvider(brokerConfig).
		Build()

	if err != nil {
		errorHandler(err)
	}

	if err := messagingService.Connect(); err != nil {
		errorHandler(err)
	}

	if messagingService.IsConnected() {
		emoji.Println(":green_circle: Connected to event broker: ", viper.Get("HOST"))
	}

	// Create queue consumer
	queue := resource.QueueDurableExclusive(viper.GetString("QUEUE"))
	receiver, err := messagingService.CreatePersistentMessageReceiverBuilder().
		WithMessageAutoAcknowledgement().
		Build(queue)

	if err != nil {
		errorHandler(err)
	}

	defer func() {
		if err := recover(); err != nil {
			errorHandler(err)
		}
	}()

	if err := receiver.Start(); err != nil {
		errorHandler(err)
	}

	if receiver.IsRunning() {
		emoji.Println(":green_circle: Receiver is running")
	}

	if err := receiver.ReceiveAsync(messageHandler); err != nil {
		errorHandler(err)
	}

	emoji.Println(":green_circle: Bound to queue: ", viper.Get("QUEUE"))

	// Exit signals handler (terminate subscriber)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	fmt.Print("\r")
	emoji.Println(":orange_circle: Stopping subscriber")

	receiver.Terminate(0)

	if receiver.IsTerminated() {
		emoji.Println(":orange_circle: Receiver was terminated")
	}

	messagingService.Disconnect()

	if !messagingService.IsConnected() {
		emoji.Println(":orange_circle: Connection was closed")
	}
}

// Message handler (event processor)
func messageHandler(message message.InboundMessage) {
	emoji.Println(":purple_circle: Message received")

	if payload, ok := message.GetPayloadAsString(); ok {
		fmt.Println(payload)
	} else if payload, ok := message.GetPayloadAsBytes(); ok {
		fmt.Println(string(payload))
	}
}

// Error handler
func errorHandler(err interface{}) {
	emoji.Println(":red_circle: Error: ", err)
	os.Exit(1)
}
