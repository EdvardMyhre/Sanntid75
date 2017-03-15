package pending_manager

import (
	"fmt"
	"time"

	"../types"
)

type struct_pendingListFloor struct {
	upOrder                 uint8
	timestamp_upOrder       time.Time
	downOrder               uint8
	timestamp_downOrder     time.Time
	internalOrder           uint8
	timestamp_internalOrder time.Time
}

var pendingList [types.NUMBER_OF_FLOORS]struct_pendingListFloor

type struct_distributor_state struct {
	busy               uint8
	timestamp_busyTime time.Time
}

var distributor_state struct_distributor_state

func Pending_task_manager(channel_from_button_intermediary <-chan types.Task,
	channel_from_distributor <-chan types.Task, channel_to_distributor chan<- types.Task,
	channel_from_assigned_tasks_manager <-chan types.Task, channel_to_assigned_tasks_manager chan<- types.Task,
	channel_from_backup_manager <-chan types.MainData) {
	fmt.Println("Pending Task Manager Go Routine startup")

	for i := 0; i < len(pendingList); i++ {
		pendingList[i].upOrder = 0
		pendingList[i].downOrder = 0
		pendingList[i].internalOrder = 0
		pendingList[i].timestamp_upOrder = time.Time{}
		pendingList[i].timestamp_downOrder = time.Time{}
		pendingList[i].timestamp_internalOrder = time.Time{}
	}
	distributor_state.busy = 0
	distributor_state.timestamp_busyTime = time.Time{}

	assigned_order_pendingList_startingIndex := 0

	for {
		//Receive from button poller
		select {
		case message_buttonOrder := <-channel_from_button_intermediary:
			adjust_pendinglist(message_buttonOrder.Type, message_buttonOrder.Floor, 0, false)
		default:
		}

		//Receive from assigned task manager
		select {
		case message_buttonOrder := <-channel_from_assigned_tasks_manager:
			adjust_pendinglist(message_buttonOrder.Type, message_buttonOrder.Floor, message_buttonOrder.Assigned, true)
		default:
		}

		//Receive from distributor
		select {
		case message_distributorStatus := <-channel_from_distributor:
			if message_distributorStatus.Finished == 0 {
				distributor_state.busy = 255
				distributor_state.timestamp_busyTime = time.Now()
			} else if message_distributorStatus.Finished != 0 {
				distributor_state.busy = 0
				distributor_state.timestamp_busyTime = time.Time{}
			}
			if message_distributorStatus.Assigned != 0 {
				adjust_pendinglist(message_distributorStatus.Type, message_distributorStatus.Floor, 0, true)
			}
		default:
		}

		//Receive from backup manager
		select {
		case message_backup := <-channel_from_backup_manager:
			fmt.Println("Received message from backup manager. Must merge into pending list, and add timestamp to tasks", message_backup)
			for j := 0; j < len(message_backup.Data); j++ {
				if message_backup.Data[j][0] != types.BTN_TYPE_COMMAND {
					adjust_pendinglist(message_backup.Data[j][0], message_backup.Data[j][1], 255, true)
				}
			}
		default:
		}

		//Sending to assigned task manager
		assigned_newOrderBool := false
		assigned_newOrderType := 0
		for i := assigned_order_pendingList_startingIndex; i < len(pendingList); i++ {
			//Check for internal order
			if (pendingList[i].internalOrder != 0) && (time.Time.IsZero(pendingList[i].timestamp_internalOrder) || (time.Since(pendingList[i].timestamp_internalOrder) > types.TIMEOUT_PENDINGLIST_ORDER)) {

				assigned_newOrderBool = true
				assigned_newOrderType = types.BTN_TYPE_COMMAND
				assigned_order_pendingList_startingIndex = i
				break
			} else if (pendingList[i].upOrder != 0) && (time.Since(pendingList[i].timestamp_upOrder) > types.TIMEOUT_PENDINGLIST_ORDER) && !time.Time.IsZero(pendingList[i].timestamp_upOrder) {
				assigned_newOrderBool = true
				assigned_newOrderType = types.BTN_TYPE_UP
				assigned_order_pendingList_startingIndex = i
				break
			} else if (pendingList[i].downOrder != 0) && (time.Since(pendingList[i].timestamp_downOrder) > types.TIMEOUT_PENDINGLIST_ORDER) && !time.Time.IsZero(pendingList[i].timestamp_downOrder) {
				assigned_newOrderBool = true
				assigned_newOrderType = types.BTN_TYPE_DOWN
				assigned_order_pendingList_startingIndex = i
			}
		}
		if assigned_newOrderBool {
			var sendOrder types.Task
			sendOrder.Floor = assigned_order_pendingList_startingIndex
			sendOrder.Type = assigned_newOrderType
			select {
			case channel_to_assigned_tasks_manager <- sendOrder:
				adjust_pendinglist(sendOrder.Type, sendOrder.Floor, 0, true)
			case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
			}
			if assigned_order_pendingList_startingIndex >= (len(pendingList) - 1) {
				assigned_order_pendingList_startingIndex = 0
			} else {
				assigned_order_pendingList_startingIndex += 1
			}
		} else {
			assigned_order_pendingList_startingIndex = 0
		}

		//Sending to task distributor
		distribute_newOrderBool := false
		distribute_newOrderType := 0
		distribute_newOrderIndex := 0

		if (distributor_state.busy == 0) || (time.Since(distributor_state.timestamp_busyTime) > types.TIMEOUT_MODULE_DISTRIBUTOR) {

			for i := 0; i < len(pendingList); i++ {
				if (pendingList[i].upOrder != 0) && (time.Time.IsZero(pendingList[i].timestamp_upOrder)) {
					distribute_newOrderBool = true
					distribute_newOrderType = types.BTN_TYPE_UP
					distribute_newOrderIndex = i
					break
				} else if (pendingList[i].downOrder != 0) && (time.Time.IsZero(pendingList[i].timestamp_downOrder)) {
					distribute_newOrderBool = true
					distribute_newOrderType = types.BTN_TYPE_DOWN
					distribute_newOrderIndex = i
					break
				}
			}
		}

		if distribute_newOrderBool {
			var sendOrder types.Task
			sendOrder.Floor = distribute_newOrderIndex
			sendOrder.Type = distribute_newOrderType
			select {
			case channel_to_distributor <- sendOrder:
				adjust_pendinglist(sendOrder.Type, sendOrder.Floor, 0, true)
			case <-time.After(types.TIMEOUT_MESSAGE_SEND_WAITTIME):
				fmt.Println("Message to Distributor FAILED to send")
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

}

func adjust_pendinglist(adjust_type int, adjust_floor int, adjust_assigned int, adjust_timestamp bool) {
	if adjust_floor >= len(pendingList) || adjust_floor < 0 {
		//Illegal floor value

	} else if adjust_assigned == 0 {
		if adjust_type == types.BTN_TYPE_UP {
			pendingList[adjust_floor].upOrder = 255
			if adjust_timestamp {
				pendingList[adjust_floor].timestamp_upOrder = time.Now()
			}
		} else if adjust_type == types.BTN_TYPE_DOWN {
			pendingList[adjust_floor].downOrder = 255
			if adjust_timestamp {
				pendingList[adjust_floor].timestamp_downOrder = time.Now()
			}
		} else if adjust_type == types.BTN_TYPE_COMMAND {
			pendingList[adjust_floor].internalOrder = 255
			if adjust_timestamp {
				pendingList[adjust_floor].timestamp_internalOrder = time.Now()
			}
		}

	} else if adjust_assigned != 0 {
		//Order has been assigned. No longer pending. Remove from pending list. Remove timestamp if adjust_timestamp is true
		if adjust_type == types.BTN_TYPE_UP {
			pendingList[adjust_floor].upOrder = 0
			if adjust_timestamp {
				pendingList[adjust_floor].timestamp_upOrder = time.Time{}
			}
		} else if adjust_type == types.BTN_TYPE_DOWN {
			pendingList[adjust_floor].downOrder = 0
			if adjust_timestamp {
				pendingList[adjust_floor].timestamp_downOrder = time.Time{}
			}
		} else if adjust_type == types.BTN_TYPE_COMMAND {
			pendingList[adjust_floor].internalOrder = 0
			if adjust_timestamp {
				pendingList[adjust_floor].timestamp_internalOrder = time.Time{}
			}
		}
	}

}
