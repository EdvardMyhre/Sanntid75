
package Network

const(
	TIMEOUT_MESSAGE_RESPONSE = 2.5

)

const(
	DESTINATION_BROADCAST = "broadcast"
	DESTINATION_BACKUP = "backup"


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

		
		//Message types used by distributor
		MESSAGE_TYPE_REQUEST_WEIGHT = (ID_MODULE_ELEVATOR_CONTROLLER | 24)		//xxx 11000
		MESSAGE_TYPE_DISTRIBUTE_ORDER = (ID_MODULE_ELEVATOR_CONTROLLER | 20)		//xxx 10100
		MESSAGE_TYPE_GIVE_DISTRIBUTOR_STATUS = (ID_MODULE_TASK_MANAGER | 20)
		
		
		
		//Message types used by task manager
		MESSAGE_TYPE_DISTRIBUTE_NEWORDER = (ID_MODULE_DISTRIBUTOR | 24)		//xxx 11000
		MESSAGE_TYPE_REQUEST_DISTRIBUTOR_STATUS = (ID_MODULE_DISTRIBUTOR | 20)		//xxx 10100
		
		
		//Message types used by elevator controller
		MESSAGE_TYPE_GIVE_WEIGHT = (ID_MODULE_DISTRIBUTOR | 18)	//xxx 10010
		
		
		
		//Message types used by network module
		
	)
