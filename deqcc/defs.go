package deqcc

import (
	"encoding/binary"
	"math"
	"unsafe"
)

type Int int32
type Float float32
type String int32
type Func int32
type Vec3 [3]float32
type Short int16
type UShort uint16

type GlobalVars struct {
	Pad               [28]Int
	Self              Int
	Other             Int
	World             Int
	Time              Float
	FrameTime         Float
	ForceRetouch      Float
	MapName           String
	Deathmatch        Float
	Coop              Float
	Teamplay          Float
	ServerFlags       Float
	TotalSecrets      Float
	TotalMonsters     Float
	FoundSecrets      Float
	KilledMonsters    Float
	Parm1             Float
	Parm2             Float
	Parm3             Float
	Parm4             Float
	Parm5             Float
	Parm6             Float
	Parm7             Float
	Parm8             Float
	Parm9             Float
	Parm10            Float
	Parm11            Float
	Parm12            Float
	Parm13            Float
	Parm14            Float
	Parm15            Float
	Parm16            Float
	VForward          Vec3
	VUp               Vec3
	VRight            Vec3
	TraceAllSolid     Float
	TraceStartSolid   Float
	TraceFraction     Float
	TraceEndpos       Vec3
	TracePlaneNormal  Vec3
	TracePlaneDist    Float
	TraceEnt          Int
	TraceInOpen       Float
	TraceInWater      Float
	MsgEntity         Int
	Main              Func
	StartFrame        Func
	PlayerPreThink    Func
	PlayerPostThink   Func
	ClientKill        Func
	ClientConnect     Func
	PutClientInServer Func
	ClientDisconnect  Func
	SetNewParms       Func
	SetChangeParms    Func
}

var GlobalVarsSz = int64(unsafe.Sizeof(GlobalVars{}))

type EntVars struct {
	modelindex    Float
	absmin        Vec3
	absmax        Vec3
	ltime         Float
	movetype      Float
	solid         Float
	origin        Vec3
	oldorigin     Vec3
	velocity      Vec3
	angles        Vec3
	avelocity     Vec3
	punchangle    Vec3
	classname     String
	model         String
	frame         Float
	skin          Float
	effects       Float
	mins          Vec3
	maxs          Vec3
	size          Vec3
	touch         Func
	use           Func
	think         Func
	blocked       Func
	nextthink     Float
	groundentity  Int
	health        Float
	frags         Float
	weapon        Float
	weaponmodel   String
	weaponframe   Float
	currentammo   Float
	ammo_shells   Float
	ammo_nails    Float
	ammo_rockets  Float
	ammo_cells    Float
	items         Float
	takedamage    Float
	chain         Int
	deadflag      Float
	view_ofs      Vec3
	button0       Float
	button1       Float
	button2       Float
	impulse       Float
	fixangle      Float
	v_angle       Vec3
	idealpitch    Float
	netname       String
	enemy         Int
	flags         Float
	colormap      Float
	team          Float
	max_health    Float
	teleport_time Float
	armortype     Float
	armorvalue    Float
	waterlevel    Float
	watertype     Float
	ideal_yaw     Float
	yaw_speed     Float
	aiment        Int
	goalentity    Int
	spawnflags    Float
	target        String
	targetname    String
	dmg_take      Float
	dmg_save      Float
	dmg_inflictor Int
	owner         Int
	movedir       Vec3
	message       String
	sounds        Float
	noise         String
	noise1        String
	noise2        String
	noise3        String
}

//go:generate stringer -type=Op
type Op UShort

const (
	OpDone Op = iota
	OpMulF
	OpMulV
	OpMulFV
	OpMulVF
	OpDivF
	OpAddF
	OpAddV
	OpSubF
	OpSubV
	OpEqF
	OpEqV
	OpEqS
	OpEqE
	OpEqFNC
	OpNeF
	OpNeV
	OpNeS
	OpNeE
	OpNeFNC
	OpLe
	OpGe
	OpLt
	OpGt
	OpLoadF
	OpLoadV
	OpLoadS
	OpLoadENT
	OpLoadFLD
	OpLoadFNC
	OpAddress
	OpStoreF
	OpStoreV
	OpStoreS
	OpStoreENT
	OpStoreFLD
	OpStoreFNC
	OpStorepF
	OpStorepV
	OpStorepS
	OpStorepENT
	OpStorepFLD
	OpStorepFNC
	OpReturn
	OpNotF
	OpNotV
	OpNotS
	OpNotENT
	OpNotFNC
	OpIf
	OpIfNot
	OpCall0
	OpCall1
	OpCall2
	OpCall3
	OpCall4
	OpCall5
	OpCall6
	OpCall7
	OpCall8
	OpState
	OpGoto
	OpAnd
	OpOr
	OpBitAnd
	OpBitOr
)

type Statement struct {
	Op      Op
	A, B, C Short
}

var StatementSz = int64(unsafe.Sizeof(Statement{}))

type Type UShort

const (
	TypeVoid Type = iota
	TypeString
	TypeFloat
	TypeVector
	TypeEntity
	TypeField
	TypeFunction
	TypePointer
)

func (t Type) String() string {
	var (
		raw   = t &^ DefSaveGlobal
		saved = t&DefSaveGlobal != 0
		base  string
	)

	switch raw {
	case TypeVoid:
		base = "void"
	case TypeString:
		base = "string"
	case TypeFloat:
		base = "float"
	case TypeVector:
		base = "vector"
	case TypeEntity:
		base = "entity"
	case TypeField:
		base = "field"
	case TypeFunction:
		base = "function"
	case TypePointer:
		base = "pointer"
	default:
		panic("unhandled")
	}
	if saved {
		return "+" + base
	}
	return base
}

func (t Type) Size() int {
	switch t &^ DefSaveGlobal {
	case 0, 1, 2, 4, 5, 6, 7:
		return 1
	case 3:
		return 3
	default:
		panic("unhandled")
	}
}

const DefSaveGlobal = 1 << 15

type Def struct {
	Typ Type // if DEF_SAVEGLOBGAL bit is set
	// the variable needs to be saved in savegames
	Offset UShort
	SName  Int
}

var DefSz = int64(unsafe.Sizeof(Def{}))

type EDef struct {
	Def
	Data interface{}
}

func NewEDef(d Def, p *ProgramsReader, e []Eval) (*EDef, error) {
	ofs := int(d.Offset)
	switch d.Typ &^ DefSaveGlobal {
	case TypeString:
		strOf := asInt32(e[ofs].Data[:])
		str, err := p.GetString(int(strOf))
		if err != nil {
			return nil, err
		}
		return &EDef{d, dataString(str)}, nil
	case TypeVector:
		var vec3 dataVec3
		for i := range vec3 {
			vec3[i] = asFloat32(e[ofs+i].Data[:])
		}
		return &EDef{d, vec3}, nil
	case TypeVoid:
		return &EDef{d, "(void)"}, nil
	case TypeFloat:
		data := asFloat32(e[ofs].Data[:])
		return &EDef{d, dataFloat(data)}, nil
	default:
		return nil, nil
	}
}

const MaxParms = 8

type Function struct {
	first_statement Int // negative numbers are builtins
	parm_start      Int
	locals          Int // total ints of parms + locals

	profile Int // runtime

	s_name Int
	s_file Int // source file defined in

	numparms  Int
	parm_size [MaxParms]byte
}

type Section struct{ Offset, Num Int }

type Programs struct {
	Version      Int
	CRC          Int     // check of header file
	Statements   Section // statement 0 is an error
	GlobalDefs   Section
	FieldDefs    Section
	Functions    Section // function 0 is an empty
	Strings      Section // first string is a null string
	Globals      Section
	EntityFields Int
}

var ProgramsSize = int64(unsafe.Sizeof(Programs{}))

type dataString string
type dataFloat float32
type dataVec3 [3]float32
type dataFunc struct{}
type dataVoid struct{}
type dataEntity struct{}
type dataField struct{}
type dataPointer struct{}

type Eval struct {
	Data [4]byte
	/*
		union
		string_t		string;
		float			_float;
		float			vector[3];
		func_t			function;
		int				_int;
		int				edict;
	*/
}

var EvalSz = int64(unsafe.Sizeof(Eval{}))

func asInt32(data []byte) int32     { return int32(binary.LittleEndian.Uint32(data)) }
func asFloat32(data []byte) float32 { return math.Float32frombits(uint32(asInt32(data))) }
