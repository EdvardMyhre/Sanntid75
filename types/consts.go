package types

const (
	ElevatorType = 0 //1 is simulation, 0 is comedi
)

const (
	NUMBER_OF_FLOORS    = 4
	NUMBER_OF_BTN_TYPES = 3
)

const (
	ET_COMEDI = 0
	ET_SIM    = 1
)

const (
	MOTOR_DIR_DOWN = -1
	MOTOR_DIR_STOP = 0
	MOTOR_DIR_UP   = 1
)

const (
	BTN_TYPE_UP      = 0
	BTN_TYPE_DOWN    = 1
	BTN_TYPE_COMMAND = 2
)

const (
	LAMP_OFF = 0
	LAMP_ON  = 1
)

const (
	TIMEOUT_MESSAGE_RESPONSE = 2.5
)

const (
	DESTINATION_BROADCAST = "broadcast"
	DESTINATION_BACKUP    = "backup"
)

//Module IDs
const (
	ID_MODULE_NETWORK      = 0   //000 xxxxx
	ID_MODULE_DISTRIBUTOR  = 96  //011 xxxxx
	ID_MODULE_TASK_MANAGER = 160 //101 xxxxx
	ID_MODULE_AMANAGER     = 192 //110 xxxxx

)

//Message Constants
const (
	//Message types used by distributor
	MESSAGE_TYPE_REQUEST_WEIGHT          = (ID_MODULE_AMANAGER | 24) //xxx 11000
	MESSAGE_TYPE_DISTRIBUTE_ORDER        = (ID_MODULE_AMANAGER | 20) //xxx 10100
	MESSAGE_TYPE_GIVE_DISTRIBUTOR_STATUS = (ID_MODULE_TASK_MANAGER | 24)

	//Message types used by task manager
	MESSAGE_TYPE_DISTRIBUTE_NEWORDER        = (ID_MODULE_DISTRIBUTOR | 24) //xxx 11000
	MESSAGE_TYPE_REQUEST_DISTRIBUTOR_STATUS = (ID_MODULE_DISTRIBUTOR | 20) //xxx 10100

	//Message types used by elevator controller
	MESSAGE_TYPE_GIVE_WEIGHT = (ID_MODULE_DISTRIBUTOR | 18) //xxx 10010

	//Message types used by assigned tasks manager
	MESSAGE_TYPE_REQUEST_BACKUP = (ID_MODULE_TASK_MANAGER | 20) //xxx 10100
	MESSAGE_TYPE_PUSH_BACKUP    = (ID_MODULE_TASK_MANAGER | 18) //xxx 10010
	MESSAGE_TYPE_SET_LIGHT      = (ID_MODULE_AMANAGER | 24)     //xxx 11000

	//Message types used by network module
)
