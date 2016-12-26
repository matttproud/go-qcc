package deqcc

import "unsafe"

type Int int32
type Float float32
type String int32
type Func int32
type Vec3 [3]float32
type Short int16
type UShort uint16

type GlobalVars struct {
	pad                [28]Int
	self               Int
	other              Int
	world              Int
	time               Float
	frametime          Float
	force_retouch      Float
	mapname            String
	deathmatch         Float
	coop               Float
	teamplay           Float
	serverflags        Float
	total_secrets      Float
	total_monsters     Float
	found_secrets      Float
	killed_monsters    Float
	parm1              Float
	parm2              Float
	parm3              Float
	parm4              Float
	parm5              Float
	parm6              Float
	parm7              Float
	parm8              Float
	parm9              Float
	parm10             Float
	parm11             Float
	parm12             Float
	parm13             Float
	parm14             Float
	parm15             Float
	parm16             Float
	v_forward          Vec3
	v_up               Vec3
	v_right            Vec3
	trace_allsolid     Float
	trace_startsolid   Float
	trace_fraction     Float
	trace_endpos       Vec3
	trace_plane_normal Vec3
	trace_plane_dist   Float
	trace_ent          Int
	trace_inopen       Float
	trace_inwater      Float
	msg_entity         Int
	main               Func
	StartFrame         Func
	PlayerPreThink     Func
	PlayerPostThink    Func
	ClientKill         Func
	ClientConnect      Func
	PutClientInServer  Func
	ClientDisconnect   Func
	SetNewParms        Func
	SetChangeParms     Func
}

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

func (t Type) String() string {
	var (
		raw   = t &^ DefSaveGlobal
		saved = t&DefSaveGlobal != 0
		base  string
	)

	switch raw {
	case 0:
		base = "void"
	case 1:
		base = "string"
	case 2:
		base = "float"
	case 3:
		base = "vector"
	case 4:
		base = "entity"
	case 5:
		base = "field"
	case 6:
		base = "function"
	case 7:
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
