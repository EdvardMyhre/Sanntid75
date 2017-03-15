package distributor

import (
	"fmt"
	"time"

	"../types"
)

func Task_distributor(channel_from_task_manager <-chan types.Task, channel_to_task_manager chan<- types.Task,
	channel_from_network <-chan types.MainData, channel_to_network chan<- types.MainData) {

	const (
		waiting_for_newOrder        = 1
		confirm_order               = 2
		request_weights             = 4
		waiting_for_weightResponses = 8
		distribute_order            = 16
	)

	var task_distributor_state int = waiting_for_newOrder
	var currentOrder types.Task
	var weightResponses []types.MainData
	task_dist_timestamp := time.Now()
	var message_distributingOrder types.MainData

	var iterate_counter int
	iterate_counter = 0

	for {
		switch task_distributor_state {
		case waiting_for_newOrder:
			select {
			case inputOrder := <-channel_from_task_manager:
				currentOrder = inputOrder
				iterate_counter = 0
				task_distributor_state = confirm_order
			case <-channel_from_network:
				//Clear channel
			default:
			}
		case confirm_order:
			currentOrder.Finished = 0
			currentOrder.Assigned = 255
			select {
			case channel_to_task_manager <- currentOrder:
				message_distributingOrder.Destination = "broadcast"
				message_distributingOrder.Type = types.REQUEST_WEIGHT
				var networkOrder []int
				networkOrder = nil
				networkOrder = append(networkOrder, currentOrder.Type)
				networkOrder = append(networkOrder, currentOrder.Floor)
				networkOrder = append(networkOrder, 0)
				networkOrder = append(networkOrder, 0)

				message_distributingOrder.Data = [][]int{networkOrder}
				iterate_counter = 0
				task_distributor_state = request_weights
			case <-channel_from_network:
				//Clear channel
			case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
				if iterate_counter > 5 {
					task_distributor_state = waiting_for_newOrder
					fmt.Println("Distributor unable to send REQUEST WEIGHTS message to network")
				} else {
					iterate_counter += 1
				}
			}

		case request_weights:
			select {
			case channel_to_network <- message_distributingOrder:

				task_dist_timestamp = time.Now()
				weightResponses = nil

				task_distributor_state = waiting_for_weightResponses
			case <-channel_from_network:
			case <-time.After(1000 * 1000000 * time.Nanosecond):
				if iterate_counter > 5 {
					task_distributor_state = waiting_for_newOrder
					fmt.Println("Distributor unable to send REQUEST WEIGHTS message to network")
				} else {
					iterate_counter += 1
				}

			}

		case waiting_for_weightResponses:
			select {
			case message_distributingOrder := <-channel_from_network:
				if message_distributingOrder.Type == types.GIVE_WEIGHT {
					weightResponses = append(weightResponses, message_distributingOrder)
				}
			default:
				if time.Since(task_dist_timestamp) > types.TIMEOUT_NETWORK_MESSAGE_RESPONSE {
					task_distributor_state = distribute_order
					iterate_counter = 0
				}
			}
		case distribute_order:
			if len(weightResponses) <= 0 {
				fmt.Println("Did not receive any responses")
				task_distributor_state = waiting_for_newOrder
			} else {
				var responseChosen types.MainData = weightResponses[0]
				for i := 0; i < len(weightResponses); i++ {
					if weightResponses[i].Data[0][0] < responseChosen.Data[0][0] {
						responseChosen = weightResponses[i]
					}
				}

				message_distributingOrder.Destination = responseChosen.Source
				message_distributingOrder.Type = types.DISTRIBUTE_ORDER

				var networkOrder []int
				networkOrder = nil
				networkOrder = append(networkOrder, currentOrder.Type)
				networkOrder = append(networkOrder, currentOrder.Floor)
				networkOrder = append(networkOrder, 0)
				networkOrder = append(networkOrder, 0)

				message_distributingOrder.Data = [][]int{networkOrder}

				channel_to_network <- message_distributingOrder

				currentOrder.Finished = 1
				currentOrder.Assigned = 0
				select {
				case channel_to_task_manager <- currentOrder:
					task_distributor_state = waiting_for_newOrder
				case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
					if iterate_counter > 5 {
						fmt.Println("Failed to send DISTRIBUTE_ORDER")
						task_distributor_state = waiting_for_newOrder
					} else {
						iterate_counter += 1
					}
				}

			}

		default:
			task_distributor_state = waiting_for_newOrder

		}

		time.Sleep(time.Millisecond)
	}

}
