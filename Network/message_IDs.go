
package Network


const(
	DESTINATION_BROADCAST = "BROADCAST"
	DESTINATION_BACKUP = "BACKUP"


)



//Module IDs
const (
	ID_MODULE_NETWORK = 0						//000 xxxxx
	ID_MODULE_DISTRIBUTOR = 96					//011 xxxxx
	ID_MODULE_TASK_MANAGER = 160				//101 xxxxx
	ID_MODULE_ELEVATOR_CONTROLLER = 192			//110 xxxxx
	
)
 
//Message Constants
const (

		//DISTRIBUTOR: Input Messages
		//From Task Manager
		ID_MSG_TYPE_DISTRIBUTOR_NEW_COMMAND = (ID_MODULE_DISTRIBUTOR | 24)		//xxx 11000
		ID_MSG_TYPE_DISTRIBUTOR_REQUEST_STATUS = (ID_MODULE_DISTRIBUTOR | 20)		//xxx 10100
		//From Network Module
		ID_MSG_TYPE_DISTRIBUTOR_RESPONSE_WEIGHTS = (ID_MODULE_DISTRIBUTOR | 18)	//xxx 10010
		//DISTRUBUTOR Messages DONE
		
		
		
		
		
		//ELEVATOR CONTROLLER: Input Messages
		//From Distributor (Via Network)
		ID_MSG_TYPE_ELEVATOR_CONTROLLER_REQUEST_WEIGHTS = (ID_MODULE_ELEVATOR_CONTROLLER | 24)		//xxx 11000
		ID_MSG_TYPE_ELEVATOR_HANDLER_DISTRIBUTE_ORDER = (ID_MODULE_ELEVATOR_CONTROLLER | 20)		//xxx 10100
		//From Network Module
		
		//From Task Manager
		
		//ELEVATOR CONTROLLER Messages DONE
		
		
		
		
		
		
				
		//TASK MANAGER: Input Messages
		//From Distributor
		ID_MSG_TYPE_TASK_MANAGER_GIVE_STATUS = (ID_MODULE_TASK_MANAGER | 20)
		//From Elevator Controller
		
		//From Network Module
		
		
		
	)
