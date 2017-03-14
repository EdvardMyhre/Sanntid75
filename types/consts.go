package types

import "time"

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
	//First number is gives wait time in MILLISECONDS
	TIMEOUT_BACKUP_RESPONSE   = int64(500 * 1000000)
	RETRY_BACKUP_RESPONSE     = 5 * 1000000 * time.Nanosecond
	TIMEOUT_AMANAGER_WAITTIME = 500 * 1000000 * time.Nanosecond
	PAUSE_AMAGER              = 5 * 1000000 * time.Nanosecond
	PAUSE_ELEVATOR            = 5 * 1000000 * time.Nanosecond
	PAUSE_NET_LISTNER         = 2 * 1000000 * time.Nanosecond

	TIMEOUT_NETWORK_MESSAGE_RESPONSE = 2000 * 1000000 * time.Nanosecond
	TIMEOUT_MESSAGE_SEND_WAITTIME    = 50 * 1000000 * time.Nanosecond
	TIMEOUT_MODULE_DISTRIBUTOR       = 2000 * 1000000 * time.Nanosecond
	TIMEOUT_PENDINGLIST_ORDER        = 5000 * 1000000 * time.Nanosecond
)

const (
	DESTINATION_BROADCAST = "broadcast"
	DESTINATION_BACKUP    = "backup"
)

//Module IDs
const (
	ID_MODULE_NETWORK        = 0   //000 xxxxx
	ID_MODULE_DISTRIBUTOR    = 96  //011 xxxxx
	ID_MODULE_BACKUP_MANAGER = 160 //101 xxxxx
	ID_MODULE_AMANAGER       = 192 //110 xxxxx

)

//Message Constants
const (
	//Message types used by distributor
	REQUEST_WEIGHT   = (ID_MODULE_AMANAGER | 24) //xxx 11000
	DISTRIBUTE_ORDER = (ID_MODULE_AMANAGER | 20) //xxx 10100

	//Message types used by task manager
	GIVE_BACKUP = (ID_MODULE_AMANAGER | 30) //xxx 11110 //Dette er ønsket svar på REQUEST_BACKUP.
	ACK_BACKUP  = (ID_MODULE_AMANAGER | 29) //xxx 11101 //Dette er ønsket svar på PUSH_BACKUP.

	//Message types used by assigned tasks manager
	GIVE_WEIGHT    = (ID_MODULE_DISTRIBUTOR | 18)    //xxx 10010
	REQUEST_BACKUP = (ID_MODULE_BACKUP_MANAGER | 20) //xxx 10100
	PUSH_BACKUP    = (ID_MODULE_BACKUP_MANAGER | 18) //xxx 10010
	SET_LIGHT      = (ID_MODULE_AMANAGER | 18)       //xxx 10010
	TASK_ASSIGNED  = (ID_MODULE_AMANAGER | 17)       //xxx 10001

	//Message types used by network module
	BACKUP_LOST             = (ID_MODULE_BACKUP_MANAGER | 24) //xxx 11000
	IS_MY_BACKUP_ALIVE      = (ID_MODULE_NETWORK | 24)        //xxx 11000
	IS_MY_BACKUP_ALIVE_TRUE = (ID_MODULE_NETWORK | 20)        //xxx 10100
	YOU_ARE_MY_BACKUP       = (ID_MODULE_NETWORK | 18)        //xxx 10010
)
