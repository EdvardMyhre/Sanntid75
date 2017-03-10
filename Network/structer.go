package Network

type InternalMessage struct { // Mellom moduler (ikke network)
	Message_type int
    Data []int
}

type MainData struct { // Dette som sendes mellom heisene og legges i backoup
    Source string
    Destination string
    Message_type int
    Data [][]int
}

