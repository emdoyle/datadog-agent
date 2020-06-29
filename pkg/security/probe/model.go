//go:generate go run github.com/DataDog/datadog-agent/pkg/security/generators/accessors -output model_accessors.go

package probe

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os/user"
	"strconv"
	"time"

	"github.com/DataDog/datadog-agent/pkg/security/secl/eval"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var NotEnoughData = errors.New("not enough data")

type Model struct {
	event *Event
}

func (m *Model) SetEvent(event interface{}) {
	m.event = event.(*Event)
}

func (m *Model) GetEvent() eval.Event {
	return m.event
}

type ChmodEvent struct {
	Mode        int32  `field:"mode"`
	Dev         uint32 `field:"-"`
	Inode       uint64 `field:"inode"`
	PathnameStr string `field:"filename" handler:"ResolveInode,string"`
}

func (e *ChmodEvent) marshalJSON(resolvers *Resolvers) ([]byte, error) {
	if e.Inode == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	buf.WriteRune('{')
	fmt.Fprintf(&buf, `"filename":"%s",`, e.ResolveInode(resolvers))
	fmt.Fprintf(&buf, `"inode":%d,`, e.Inode)
	fmt.Fprintf(&buf, `"mode":%d`, e.Mode)
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *ChmodEvent) UnmarshalBinary(data []byte) (int, error) {
	if len(data) < 16 {
		return 0, NotEnoughData
	}
	e.Mode = int32(byteOrder.Uint32(data[0:4]))
	e.Dev = byteOrder.Uint32(data[4:8])
	e.Inode = byteOrder.Uint64(data[8:16])
	return 16, nil
}

func (e *ChmodEvent) ResolveInode(resolvers *Resolvers) string {
	if len(e.PathnameStr) == 0 {
		e.PathnameStr = resolvers.DentryResolver.Resolve(e.Dev, e.Inode)
	}
	return e.PathnameStr
}

type ChownEvent struct {
	UID         int32  `field:"uid"`
	GID         int32  `field:"gid"`
	Dev         uint32 `field:"-"`
	Inode       uint64 `field:"inode"`
	PathnameStr string `field:"filename" handler:"ResolveInode,string"`
}

func (e *ChownEvent) marshalJSON(resolvers *Resolvers) ([]byte, error) {
	if e.Inode == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	buf.WriteRune('{')
	fmt.Fprintf(&buf, `"filename":"%s",`, e.ResolveInode(resolvers))
	fmt.Fprintf(&buf, `"inode":%d,`, e.Inode)
	fmt.Fprintf(&buf, `"uid":%d,`, e.UID)
	fmt.Fprintf(&buf, `"gid":%d`, e.GID)
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *ChownEvent) UnmarshalBinary(data []byte) (int, error) {
	if len(data) < 24 {
		return 0, NotEnoughData
	}
	e.UID = int32(byteOrder.Uint32(data[0:4]))
	e.GID = int32(byteOrder.Uint32(data[4:8]))
	e.Dev = byteOrder.Uint32(data[12:16])
	e.Inode = byteOrder.Uint64(data[16:24])
	return 24, nil
}

func (e *ChownEvent) ResolveInode(resolvers *Resolvers) string {
	if len(e.PathnameStr) == 0 {
		e.PathnameStr = resolvers.DentryResolver.Resolve(e.Dev, e.Inode)
	}
	return e.PathnameStr
}

type OpenEvent struct {
	Flags       uint32 `yaml:"flags" field:"flags"`
	Mode        uint32 `yaml:"mode" field:"mode"`
	Inode       uint64 `field:"inode"`
	Dev         uint32 `field:"-"`
	PathnameStr string `field:"filename" handler:"ResolveInode,string"`
	BasenameStr string `field:"basename" handler:"ResolveBasename,string"`
}

func (e *OpenEvent) marshalJSON(resolvers *Resolvers) ([]byte, error) {
	if e.Inode == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	buf.WriteRune('{')
	fmt.Fprintf(&buf, `"filename":"%s",`, e.ResolveInode(resolvers))
	fmt.Fprintf(&buf, `"inode":%d,`, e.Inode)
	fmt.Fprintf(&buf, `"mode":%d,`, e.Mode)
	fmt.Fprintf(&buf, `"flags":%d`, e.Flags)
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *OpenEvent) ResolveInode(resolvers *Resolvers) string {
	if len(e.PathnameStr) == 0 {
		e.PathnameStr = resolvers.DentryResolver.Resolve(e.Dev, e.Inode)
	}
	return e.PathnameStr
}

func (e *OpenEvent) ResolveBasename(resolvers *Resolvers) string {
	if len(e.BasenameStr) == 0 {
		e.BasenameStr = resolvers.DentryResolver.GetName(e.Dev, e.Inode)
	}
	return e.BasenameStr
}

func (e *OpenEvent) UnmarshalBinary(data []byte) (int, error) {
	if len(data) < 20 {
		return 0, NotEnoughData
	}
	e.Flags = byteOrder.Uint32(data[0:4])
	e.Mode = byteOrder.Uint32(data[4:8])
	e.Inode = byteOrder.Uint64(data[8:16])
	e.Dev = byteOrder.Uint32(data[16:20])
	return 20, nil
}

type MkdirEvent struct {
	Mode        int32  `field:"mode"`
	Dev         uint32 `field:"-"`
	Inode       uint64 `field:"inode"`
	PathnameStr string `field:"filename" handler:"ResolveInode,string"`
}

func (e *MkdirEvent) marshalJSON(resolvers *Resolvers) ([]byte, error) {
	if e.Inode == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	buf.WriteRune('{')
	fmt.Fprintf(&buf, `"filename":"%s",`, e.ResolveInode(resolvers))
	fmt.Fprintf(&buf, `"inode":%d,`, e.Inode)
	fmt.Fprintf(&buf, `"mode":%d`, e.Mode)
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *MkdirEvent) UnmarshalBinary(data []byte) (int, error) {
	if len(data) < 16 {
		return 0, NotEnoughData
	}
	e.Mode = int32(byteOrder.Uint32(data[0:4]))
	e.Dev = byteOrder.Uint32(data[4:8])
	e.Inode = byteOrder.Uint64(data[8:16])
	return 16, nil
}

func (e *MkdirEvent) ResolveInode(resolvers *Resolvers) string {
	if len(e.PathnameStr) == 0 {
		e.PathnameStr = resolvers.DentryResolver.Resolve(e.Dev, e.Inode)
	}
	return e.PathnameStr
}

type RmdirEvent struct {
	Dev         uint32 `field:"-"`
	Inode       uint64 `field:"inode"`
	PathnameStr string `field:"filename" handler:"ResolveInode,string"`
}

func (e *RmdirEvent) marshalJSON(resolvers *Resolvers) ([]byte, error) {
	if e.Inode == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	buf.WriteRune('{')
	fmt.Fprintf(&buf, `"filename":"%s",`, e.ResolveInode(resolvers))
	fmt.Fprintf(&buf, `"inode":%d`, e.Inode)
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *RmdirEvent) ResolveInode(resolvers *Resolvers) string {
	if len(e.PathnameStr) == 0 {
		e.PathnameStr = resolvers.DentryResolver.Resolve(e.Dev, e.Inode)
	}
	return e.PathnameStr
}

func (e *RmdirEvent) UnmarshalBinary(data []byte) (int, error) {
	if len(data) < 12 {
		return 0, NotEnoughData
	}
	e.Inode = byteOrder.Uint64(data[0:8])
	e.Dev = byteOrder.Uint32(data[8:12])
	return 12, nil
}

type UnlinkEvent struct {
	Inode       uint64 `field:"inode"`
	Dev         uint32 `field:"-"`
	PathnameStr string `field:"filename" handler:"ResolveInode,string"`
}

func (e *UnlinkEvent) marshalJSON(resolvers *Resolvers) ([]byte, error) {
	if e.Inode == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	buf.WriteRune('{')
	fmt.Fprintf(&buf, `"filename":"%s",`, e.ResolveInode(resolvers))
	fmt.Fprintf(&buf, `"inode":%d`, e.Inode)
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *UnlinkEvent) UnmarshalBinary(data []byte) (int, error) {
	if len(data) < 12 {
		return 0, NotEnoughData
	}
	e.Inode = byteOrder.Uint64(data[0:8])
	e.Dev = byteOrder.Uint32(data[8:12])
	return 12, nil
}

func (e *UnlinkEvent) ResolveInode(resolvers *Resolvers) string {
	if len(e.PathnameStr) == 0 {
		e.PathnameStr = resolvers.DentryResolver.Resolve(e.Dev, e.Inode)
	}
	return e.PathnameStr
}

type RenameEvent struct {
	Dev               uint32 `field:"-"`
	SrcInode          uint64 `field:"old_inode"`
	SrcPathnameStr    string `field:"old_filename" handler:"ResolveSrcInode,string"`
	TargetInode       uint64 `field:"new_inode"`
	TargetPathnameStr string `field:"new_filename" handler:"ResolveTargetInode,string"`
}

func (e *RenameEvent) marshalJSON(resolvers *Resolvers) ([]byte, error) {
	if e.SrcInode == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	buf.WriteRune('{')
	fmt.Fprintf(&buf, `"old_inode":%d,`, e.SrcInode)
	fmt.Fprintf(&buf, `"old_filename":"%s",`, e.ResolveSrcInode(resolvers))
	fmt.Fprintf(&buf, `"new_inode":%d,`, e.TargetInode)
	fmt.Fprintf(&buf, `"new_filename":"%s"`, e.ResolveTargetInode(resolvers))
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *RenameEvent) UnmarshalBinary(data []byte) (int, error) {
	if len(data) < 24 {
		return 0, NotEnoughData
	}
	e.Dev = byteOrder.Uint32(data[0:4])
	// padding
	e.SrcInode = byteOrder.Uint64(data[8:16])
	e.TargetInode = byteOrder.Uint64(data[16:24])
	return 24, nil
}

func (e *RenameEvent) ResolveSrcInode(resolvers *Resolvers) string {
	if len(e.SrcPathnameStr) == 0 {
		e.SrcPathnameStr = resolvers.DentryResolver.Resolve(0xffffffff, e.SrcInode)
	}
	return e.SrcPathnameStr
}

func (e *RenameEvent) ResolveTargetInode(resolvers *Resolvers) string {
	if len(e.TargetPathnameStr) == 0 {
		e.TargetPathnameStr = resolvers.DentryResolver.Resolve(e.Dev, e.TargetInode)
	}
	return e.TargetPathnameStr
}

type UtimesEvent struct {
	Atime       time.Time
	Mtime       time.Time
	Inode       uint64 `field:"inode"`
	Dev         uint32 `field:"-"`
	PathnameStr string `field:"filename" handler:"ResolveInode,string"`
}

func (e *UtimesEvent) marshalJSON(resolvers *Resolvers) ([]byte, error) {
	if e.Inode == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	buf.WriteRune('{')
	fmt.Fprintf(&buf, `"filename":"%s",`, e.ResolveInode(resolvers))
	fmt.Fprintf(&buf, `"inode":%d`, e.Inode)
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *UtimesEvent) UnmarshalBinary(data []byte) (int, error) {
	if len(data) < 48 {
		return 0, NotEnoughData
	}

	timeSec := byteOrder.Uint64(data[0:8])
	timeNsec := byteOrder.Uint64(data[8:16])
	e.Atime = time.Unix(int64(timeSec), int64(timeNsec))

	timeSec = byteOrder.Uint64(data[16:24])
	timeNsec = byteOrder.Uint64(data[24:32])
	e.Mtime = time.Unix(int64(timeSec), int64(timeNsec))

	e.Dev = byteOrder.Uint32(data[36:40])
	e.Inode = byteOrder.Uint64(data[40:48])

	return 48, nil
}

func (e *UtimesEvent) ResolveInode(resolvers *Resolvers) string {
	if len(e.PathnameStr) == 0 {
		e.PathnameStr = resolvers.DentryResolver.Resolve(e.Dev, e.Inode)
	}
	return e.PathnameStr
}

type ContainerEvent struct {
	ID     string   `yaml:"id" field:"id" event:"container"`
	Labels []string `yaml:"labels" field:"labels" event:"container"`
}

type KernelEvent struct {
	Type      uint64 `field:"type" handler:"ResolveType,string"`
	Timestamp uint64 `field:"-"`
	Retval    int64  `field:"retval"`
}

func (k *KernelEvent) ResolveType(resolvers *Resolvers) string {
	return ProbeEventType(k.Type).String()
}

func (k *KernelEvent) marshalJSON(resolvers *Resolvers) ([]byte, error) {
	if k.Type == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	buf.WriteRune('{')
	fmt.Fprintf(&buf, `"type":%d,`, k.Type) // TODO(sbaubeau): use resolved type
	fmt.Fprintf(&buf, `"timestamp":%d,`, k.Timestamp)
	fmt.Fprintf(&buf, `"retval":%d`, k.Retval)
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (k *KernelEvent) UnmarshalBinary(data []byte) (int, error) {
	if len(data) < 24 {
		return 0, NotEnoughData
	}
	k.Type = byteOrder.Uint64(data[0:8])
	k.Timestamp = byteOrder.Uint64(data[8:16])
	k.Retval = int64(byteOrder.Uint64(data[16:24]))
	return 24, nil
}

type ProcessEvent struct {
	Pidns       uint64 `field:"pidns"`
	Comm        string `field:"name" handler:"ResolveComm,string"`
	TTYName     string `field:"tty_name" handler:"ResolveTTY,string"`
	Pid         uint32 `field:"pid"`
	Tid         uint32 `field:"tid"`
	UID         uint32 `field:"uid"`
	GID         uint32 `field:"gid"`
	User        string `field:"user" handler:"ResolveUser,string"`
	Group       string `field:"group" handler:"ResolveGroup,string"`
	PathnameStr string `field:"filename" handler:"ResolveInode,string"`

	CommRaw    [16]byte `field:"-"`
	TTYNameRaw [64]byte `field:"-"`
	Inode      uint64   `field:"-"`
}

func (p *ProcessEvent) ResolveInode(resolvers *Resolvers) string {
	if len(p.PathnameStr) == 0 {
		p.PathnameStr = resolvers.DentryResolver.Resolve(0xffffffff, p.Inode)
	}
	return p.PathnameStr
}

func (p *ProcessEvent) marshalJSON(resolvers *Resolvers) ([]byte, error) {
	if p.Pid == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	buf.WriteRune('{')
	fmt.Fprintf(&buf, `"pidns":%d,`, p.Pidns)
	fmt.Fprintf(&buf, `"name":"%s",`, p.GetComm())
	if tty := p.GetTTY(); tty != "" {
		fmt.Fprintf(&buf, `"tty_name":"%s",`, tty)
	}
	fmt.Fprintf(&buf, `"pid":%d,`, p.Pid)
	fmt.Fprintf(&buf, `"tid":%d,`, p.Tid)
	fmt.Fprintf(&buf, `"uid":%d,`, p.UID)
	fmt.Fprintf(&buf, `"gid":%d`, p.GID)
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (p *ProcessEvent) ResolveTTY(resolvers *Resolvers) string {
	return p.GetTTY()
}

func (p *ProcessEvent) GetTTY() string {
	if len(p.TTYName) == 0 {
		p.TTYName = string(bytes.Trim(p.TTYNameRaw[:], "\x00"))
	}
	return p.TTYName
}

func (p *ProcessEvent) ResolveComm(resolvers *Resolvers) string {
	return p.GetComm()
}

func (p *ProcessEvent) GetComm() string {
	if len(p.Comm) == 0 {
		p.Comm = string(bytes.Trim(p.CommRaw[:], "\x00"))
	}
	return p.Comm
}

func (p *ProcessEvent) ResolveUser(resolvers *Resolvers) string {
	u, err := user.LookupId(strconv.Itoa(int(p.UID)))
	if err == nil {
		p.User = u.Username
	}
	return p.User
}

func (p *ProcessEvent) ResolveGroup(resolvers *Resolvers) string {
	g, err := user.LookupGroupId(strconv.Itoa(int(p.GID)))
	if err == nil {
		p.Group = g.Name
	}
	return p.Group
}

func (p *ProcessEvent) UnmarshalBinary(data []byte) (int, error) {
	if len(data) < 104 {
		return 0, NotEnoughData
	}
	p.Pidns = byteOrder.Uint64(data[0:8])
	if err := binary.Read(bytes.NewBuffer(data[8:24]), byteOrder, &p.CommRaw); err != nil {
		return 8, err
	}
	if err := binary.Read(bytes.NewBuffer(data[24:88]), byteOrder, &p.TTYNameRaw); err != nil {
		return 8 + len(p.CommRaw), err
	}
	p.Pid = byteOrder.Uint32(data[88:92])
	p.Tid = byteOrder.Uint32(data[92:96])
	p.UID = byteOrder.Uint32(data[96:100])
	p.GID = byteOrder.Uint32(data[100:104])
	return 104, nil
}

// genaccessors
type Event struct {
	ID        string         `yaml:"id" field:"-"`
	Event     KernelEvent    `yaml:"event" field:"event"`
	Process   ProcessEvent   `yaml:"process" field:"process" event:"*"`
	Container ContainerEvent `yaml:"container" field:"container"`
	Chmod     ChmodEvent     `yaml:"chmod" field:"chmod" event:"chmod"`
	Chown     ChownEvent     `yaml:"chown" field:"chown" event:"chown"`
	Open      OpenEvent      `yaml:"open" field:"open" event:"open"`
	Mkdir     MkdirEvent     `yaml:"mkdir" field:"mkdir" event:"mkdir"`
	Rmdir     RmdirEvent     `yaml:"rmdir" field:"rmdir" event:"rmdir"`
	Rename    RenameEvent    `yaml:"rename" field:"rename" event:"rename"`
	Unlink    UnlinkEvent    `yaml:"unlink" field:"unlink" event:"unlink"`
	Utimes    UtimesEvent    `yaml:"utimes" field:"utimes" event:"utimes"`

	resolvers *Resolvers `field:"-"`
}

func (e *Event) String() string {
	d, _ := json.Marshal(e)
	return string(d)
}

func (e *Event) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteRune('{')
	fmt.Fprintf(&buf, `"id":"%s",`, e.ID)
	fmt.Fprintf(&buf, `"type":"%s",`, e.GetType())

	entries := []struct {
		field      string
		marshalFnc func(resolvers *Resolvers) ([]byte, error)
	}{
		{
			field:      "event",
			marshalFnc: e.Event.marshalJSON,
		},
		{
			field:      "process",
			marshalFnc: e.Process.marshalJSON,
		},
		{
			field:      "file",
			marshalFnc: e.Chmod.marshalJSON,
		},
		{
			field:      "file",
			marshalFnc: e.Chown.marshalJSON,
		},
		{
			field:      "file",
			marshalFnc: e.Open.marshalJSON,
		},
		{
			field:      "file",
			marshalFnc: e.Mkdir.marshalJSON,
		},
		{
			field:      "file",
			marshalFnc: e.Rmdir.marshalJSON,
		},
		{
			field:      "file",
			marshalFnc: e.Unlink.marshalJSON,
		},
		{
			field:      "file",
			marshalFnc: e.Rename.marshalJSON,
		},
		{
			field:      "utimes",
			marshalFnc: e.Utimes.marshalJSON,
		},
	}

	var prev bool
	for _, entry := range entries {
		d, err := entry.marshalFnc(e.resolvers)
		if err != nil {
			return nil, err
		}
		if d != nil {
			if prev {
				buf.WriteRune(',')
			}
			buf.WriteString(`"` + entry.field + `":`)
			buf.Write(d)
			prev = true
		}
	}
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *Event) GetType() string {
	return ProbeEventType(e.Event.Type).String()
}

func (e *Event) GetID() string {
	return e.ID
}

func (e *Event) UnmarshalBinary(data []byte) (int, error) {
	offset, err := e.Process.UnmarshalBinary(data)
	if err != nil {
		return offset, err
	}

	return offset, nil
}

func NewEvent(resolvers *Resolvers) *Event {
	id, _ := uuid.NewRandom()
	return &Event{
		ID:        id.String(),
		resolvers: resolvers,
	}
}