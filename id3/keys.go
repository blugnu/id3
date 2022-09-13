package id3

type FrameKey int

const (
	UnknownKey FrameKey = iota

	// virtual keys
	DiscNo
	NumDiscs
	TrackNo
	NumTracks

	// v2.3.0 standard frames
	AENC
	APIC
	COMM
	COMR
	ENCR
	EQUA
	ETCO
	GEOB
	GRID
	IPLS
	LINK
	MCDI
	MLLT
	OWNE
	PCNT
	POPM
	POSS
	PRIV
	RBUF
	RVAD
	RVRB
	SYLT
	SYTC
	TALB
	TBPM
	TCMP
	TCOM
	TCON
	TCOP
	TDAT
	TDLY
	TENC
	TEXT
	TFLT
	TIME
	TIT1
	TIT2
	TIT3
	TKEY
	TLAN
	TLEN
	TMED
	TOAL
	TOFN
	TOLY
	TOPE
	TORY
	TOWN
	TPE1
	TPE2
	TPE3
	TPE4
	TPOS
	TPUB
	TRCK
	TRDA
	TRSN
	TRSO
	TSIZ
	TSO2
	TSOC
	TSRC
	TSSE
	TXXX
	TYER
	UFID
	USER
	USLT
	WCOM
	WCOP
	WOAF
	WOAR
	WOAS
	WORS
	WPAY
	WPUB
	WXXX

	// v2.4.0 standard frames
	ASPI
	EQU2
	TIPL
	RVA2
	SEEK
	SIGN
	TDEN
	TDOR
	TDRC
	TDRL
	TDTG
	TMCL
	TMOO
	TPRO
	TSOA
	TSOP
	TSOT
	TSST
)

type FrameKeySet []FrameKey

func (set FrameKeySet) Add(key FrameKey) FrameKeySet {
	if set.Contains(key) {
		return set
	}
	return append(set, key)
}

func (set FrameKeySet) Contains(key FrameKey) bool {
	for _, k := range set {
		if k == key {
			return true
		}
	}
	return false
}

func (set FrameKeySet) IsEmpty() bool {
	return len(set) == 0
}

func (set FrameKeySet) Remove(key FrameKey) FrameKeySet {
	if !set.Contains(key) {
		return set
	}
	newset := FrameKeySet{}
	for _, k := range set {
		if k == key {
			continue
		}
		newset = append(newset, k)
	}
	return newset
}

var v220keys = map[string]FrameKey{
	"CRA": AENC,
	"PIC": APIC,
	"COM": COMM,
	"IPL": IPLS,
	"LNK": LINK,
	"MCI": MCDI,
	"MLL": MLLT,
	"POP": POPM,
	"BUF": RBUF,
	"RVA": RVAD,
	"REV": RVRB,
	"SLT": SYLT,
	"STC": SYTC,
	"TAL": TALB,
	"TBP": TBPM,
	"TCP": TCMP,
	"TCM": TCOM,
	"TCO": TCON,
	"TCR": TCOP,
	"TDA": TDAT,
	"TDY": TDLY,
	"TEN": TENC,
	"TXT": TEXT,
	"TFT": TFLT,
	"TIM": TIME,
	"TT1": TIT1,
	"TT2": TIT2,
	"TT3": TIT3,
	"TKE": TKEY,
	"TLA": TLAN,
	"TLE": TLEN,
	"TMT": TMED,
	"TOT": TOAL,
	"TOF": TOFN,
	"TOL": TOLY,
	"TOA": TOPE,
	"TOR": TORY,
	"TP1": TPE1,
	"TP2": TPE2,
	"TP3": TPE3,
	"TP4": TPE4,
	"TPA": TPOS,
	"TPB": TPUB,
	"TRK": TRCK,
	"TRD": TRDA,
	"TSI": TSIZ,
	"TSA": TSOA,
	"TS2": TSO2,
	"TSP": TSOP,
	"TSC": TSOC,
	"TST": TSOT,
	"TRC": TSRC,
	"TSS": TSSE,
	"TXX": TXXX,
	"TYE": TYER,
	"UFI": UFID,
	"ULT": USLT,
}

var v230keys = map[string]FrameKey{
	"AENC": AENC,
	"APIC": APIC,
	"COMM": COMM,
	"COMR": COMR,
	"ENCR": ENCR,
	"EQUA": EQUA,
	"ETCO": ETCO,
	"GEOB": GEOB,
	"GRID": GRID,
	"IPLS": IPLS,
	"LINK": LINK,
	"MCDI": MCDI,
	"MLLT": MLLT,
	"OWNE": OWNE,
	"PCNT": PCNT,
	"POPM": POPM,
	"POSS": POSS,
	"PRIV": PRIV,
	"RBUF": RBUF,
	"RVAD": RVAD,
	"RVRB": RVRB,
	"SYLT": SYLT,
	"SYTC": SYTC,
	"TALB": TALB,
	"TBPM": TBPM,
	"TCMP": TCMP,
	"TCOM": TCOM,
	"TCON": TCON,
	"TCOP": TCOP,
	"TDAT": TDAT,
	"TDLY": TDLY,
	"TENC": TENC,
	"TEXT": TEXT,
	"TFLT": TFLT,
	"TIME": TIME,
	"TIT1": TIT1,
	"TIT2": TIT2,
	"TIT3": TIT3,
	"TKEY": TKEY,
	"TLAN": TLAN,
	"TLEN": TLEN,
	"TMED": TMED,
	"TOAL": TOAL,
	"TOFN": TOFN,
	"TOLY": TOLY,
	"TOPE": TOPE,
	"TORY": TORY,
	"TOWN": TOWN,
	"TPE1": TPE1,
	"TPE2": TPE2,
	"TPE3": TPE3,
	"TPE4": TPE4,
	"TPOS": TPOS,
	"TPUB": TPUB,
	"TRCK": TRCK,
	"TRDA": TRDA,
	"TRSN": TRSN,
	"TRSO": TRSO,
	"TSIZ": TSIZ,
	"TSOA": TSOA,
	"TSO2": TSO2,
	"TSOP": TSOP,
	"TSOC": TSOC,
	"TSOT": TSOT,
	"TSRC": TSRC,
	"TSSE": TSSE,
	"TXXX": TXXX,
	"TYER": TYER,
	"UFID": UFID,
	"USER": USER,
	"USLT": USLT,
	"WCOM": WCOM,
	"WCOP": WCOP,
	"WOAF": WOAF,
	"WOAR": WOAR,
	"WOAS": WOAS,
	"WORS": WORS,
	"WPAY": WPAY,
	"WPUB": WPUB,
	"WXXX": WXXX,
}

var v240keys = map[string]FrameKey{
	"ASPI": ASPI,
	"EQU2": EQU2,
	"RVA2": RVA2,
	"SEEK": SEEK,
	"SIGN": SIGN,
	"TDEN": TDEN,
	"TDOR": TDOR,
	"TDRC": TDRC,
	"TDRL": TDRL,
	"TDTG": TDTG,
	"TMCL": TMCL,
	"TMOO": TMOO,
	"TPRO": TPRO,
	"TSST": TSST,
}

func FrameKeyForID(id string) FrameKey {
	var key FrameKey
	var ok bool

	switch len(id) {
	case 3:
		key, ok = v220keys[id]

	default:
		key, ok = v230keys[id]
		if !ok {
			key, ok = v240keys[id]
		}
	}

	if ok {
		return key
	}

	return UnknownKey
}
