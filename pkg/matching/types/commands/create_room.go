package commands

type CreateRoom struct {
	Players           []string `json:"Players"`
	Reservations      []string `json:"Reservations"`
	NAT               []int    `json:"NAT"`
	Slots             [][]int  `json:"Slots"`
	RoomState         int      `json:"RoomState"`
	HostMood          int      `json:"HostMood"`
	PassedNoJoinPoint int      `json:"PassedNoJoinPoint"`
	Location          []uint32 `json:"Location"`
	Language          int      `json:"Language"`
	BuildVersion      int      `json:"BuildVersion"`
	Search            string   `json:"Search"`
}
