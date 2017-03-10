package Network


type MainData struct { // Dette som sendes mellom heisene og legges i backoup
    Source string
    Destination string
    Message_type int
    Data [][]int
}

