package types

//[CreateRoomCommand,["Players":["LumaLivy"],"Reservations":["0"],"NAT":[2],"Slots":[[1,3]],"RoomState":0,"HostMood":1,"PassedNoJoinPoint":0,"Location":[0x7f000001],"Language":1,"BuildVersion":289,"Search":""]]

type RoomUser struct {
	// The name of the user
	Name string
	// How many local players are associated with the user 0 for it's just the user on their own
	SubUsers int
	// NAT type for the user, 1 is direct internet access, 2 is behind a traditional NAT (i.e has a public IP),
	// 3 is CG NAT Hell
	NAT int
}
