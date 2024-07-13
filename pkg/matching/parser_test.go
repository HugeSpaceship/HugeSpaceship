package matching

import (
	"github.com/HugeSpaceship/HugeSpaceship/pkg/matching/types/commands"
	"github.com/stretchr/testify/assert"
	"net/netip"
	"testing"
)

const testMatch = `[CreateRoomCommand,["Players":["LumaLivy"],"Reservations":["0"],"NAT":[2],"Slots":[[1,3]],"RoomState":0,"HostMood":1,"PassedNoJoinPoint":0,"Location":[0x7f000001],"Language":1,"BuildVersion":289,"Search":""]]`

func TestGetCommand(t *testing.T) {
	command := GetCommand([]byte(testMatch))
	assert.EqualValues(t, "CreateRoomCommand", command)
}

func TestUnmarshal(t *testing.T) {

	createRoom, err := Unmarshal[commands.CreateRoom]([]byte(testMatch))
	assert.Nil(t, err)
	//assert.Equal(t, command, "CreateRoomCommand")
	assert.Equal(t, "LumaLivy", createRoom.Players[0])
	assert.Equal(t, 2, createRoom.NAT[0])

	ip, ok := IPFromLocation(createRoom.Location[0])
	assert.True(t, ok)
	assert.Equal(t, netip.MustParseAddr("127.0.0.1"), ip)
}
