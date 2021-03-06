// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package arm64

import (
	"cmd/compile/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/arm64"
)

const (
	LeftRdwr  uint32 = gc.LeftRead | gc.LeftWrite
	RightRdwr uint32 = gc.RightRead | gc.RightWrite
)

// This table gives the basic information about instruction
// generated by the compiler and processed in the optimizer.
// See opt.h for bit definitions.
//
// Instructions not generated need not be listed.
// As an exception to that rule, we typically write down all the
// size variants of an operation even if we just use a subset.
//
// The table is formatted for 8-space tabs.
var progtable = [arm64.ALAST & obj.AMask]obj.ProgInfo{
	obj.ATYPE:     {Flags: gc.Pseudo | gc.Skip},
	obj.ATEXT:     {Flags: gc.Pseudo},
	obj.AFUNCDATA: {Flags: gc.Pseudo},
	obj.APCDATA:   {Flags: gc.Pseudo},
	obj.AUNDEF:    {Flags: gc.Break},
	obj.AUSEFIELD: {Flags: gc.OK},
	obj.ACHECKNIL: {Flags: gc.LeftRead},
	obj.AVARDEF:   {Flags: gc.Pseudo | gc.RightWrite},
	obj.AVARKILL:  {Flags: gc.Pseudo | gc.RightWrite},
	obj.AVARLIVE:  {Flags: gc.Pseudo | gc.LeftRead},

	// NOP is an internal no-op that also stands
	// for USED and SET annotations, not the Power opcode.
	obj.ANOP:                {Flags: gc.LeftRead | gc.RightWrite},
	arm64.AHINT & obj.AMask: {Flags: gc.OK},

	// Integer
	arm64.AADD & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.ASUB & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.ANEG & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AAND & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AORR & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AEOR & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AMUL & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.ASMULL & obj.AMask: {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AUMULL & obj.AMask: {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.ASMULH & obj.AMask: {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AUMULH & obj.AMask: {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.ASDIV & obj.AMask:  {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AUDIV & obj.AMask:  {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.ALSL & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.ALSR & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AASR & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.ACMP & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead},
	arm64.AADC & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite | gc.UseCarry},
	arm64.AROR & obj.AMask:   {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AADDS & obj.AMask:  {Flags: gc.SizeQ | gc.LeftRead | gc.RegRead | gc.RightWrite | gc.SetCarry},

	// Floating point.
	arm64.AFADDD & obj.AMask:  {Flags: gc.SizeD | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AFADDS & obj.AMask:  {Flags: gc.SizeF | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AFSUBD & obj.AMask:  {Flags: gc.SizeD | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AFSUBS & obj.AMask:  {Flags: gc.SizeF | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AFNEGD & obj.AMask:  {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite},
	arm64.AFNEGS & obj.AMask:  {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite},
	arm64.AFSQRTD & obj.AMask: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite},
	arm64.AFMULD & obj.AMask:  {Flags: gc.SizeD | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AFMULS & obj.AMask:  {Flags: gc.SizeF | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AFDIVD & obj.AMask:  {Flags: gc.SizeD | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AFDIVS & obj.AMask:  {Flags: gc.SizeF | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm64.AFCMPD & obj.AMask:  {Flags: gc.SizeD | gc.LeftRead | gc.RegRead},
	arm64.AFCMPS & obj.AMask:  {Flags: gc.SizeF | gc.LeftRead | gc.RegRead},

	// float -> integer
	arm64.AFCVTZSD & obj.AMask:  {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AFCVTZSS & obj.AMask:  {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AFCVTZSDW & obj.AMask: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AFCVTZSSW & obj.AMask: {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AFCVTZUD & obj.AMask:  {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AFCVTZUS & obj.AMask:  {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AFCVTZUDW & obj.AMask: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AFCVTZUSW & obj.AMask: {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},

	// float -> float
	arm64.AFCVTSD & obj.AMask: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AFCVTDS & obj.AMask: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},

	// integer -> float
	arm64.ASCVTFD & obj.AMask:  {Flags: gc.SizeQ | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.ASCVTFS & obj.AMask:  {Flags: gc.SizeQ | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.ASCVTFWD & obj.AMask: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.ASCVTFWS & obj.AMask: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AUCVTFD & obj.AMask:  {Flags: gc.SizeQ | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AUCVTFS & obj.AMask:  {Flags: gc.SizeQ | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AUCVTFWD & obj.AMask: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm64.AUCVTFWS & obj.AMask: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},

	// Moves
	arm64.AMOVB & obj.AMask:  {Flags: gc.SizeB | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	arm64.AMOVBU & obj.AMask: {Flags: gc.SizeB | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	arm64.AMOVH & obj.AMask:  {Flags: gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	arm64.AMOVHU & obj.AMask: {Flags: gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	arm64.AMOVW & obj.AMask:  {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	arm64.AMOVWU & obj.AMask: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	arm64.AMOVD & obj.AMask:  {Flags: gc.SizeQ | gc.LeftRead | gc.RightWrite | gc.Move},
	arm64.AFMOVS & obj.AMask: {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Move | gc.Conv},
	arm64.AFMOVD & obj.AMask: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Move},

	// Jumps
	arm64.AB & obj.AMask:    {Flags: gc.Jump | gc.Break},
	arm64.ABL & obj.AMask:   {Flags: gc.Call},
	arm64.ABEQ & obj.AMask:  {Flags: gc.Cjmp},
	arm64.ABNE & obj.AMask:  {Flags: gc.Cjmp},
	arm64.ABGE & obj.AMask:  {Flags: gc.Cjmp},
	arm64.ABLT & obj.AMask:  {Flags: gc.Cjmp},
	arm64.ABGT & obj.AMask:  {Flags: gc.Cjmp},
	arm64.ABLE & obj.AMask:  {Flags: gc.Cjmp},
	arm64.ABLO & obj.AMask:  {Flags: gc.Cjmp},
	arm64.ABLS & obj.AMask:  {Flags: gc.Cjmp},
	arm64.ABHI & obj.AMask:  {Flags: gc.Cjmp},
	arm64.ABHS & obj.AMask:  {Flags: gc.Cjmp},
	arm64.ACBZ & obj.AMask:  {Flags: gc.Cjmp},
	arm64.ACBNZ & obj.AMask: {Flags: gc.Cjmp},
	obj.ARET:                {Flags: gc.Break},
	obj.ADUFFZERO:           {Flags: gc.Call},
	obj.ADUFFCOPY:           {Flags: gc.Call},
}

func proginfo(p *obj.Prog) {
	info := &p.Info
	*info = progtable[p.As&obj.AMask]
	if info.Flags == 0 {
		gc.Fatalf("proginfo: unknown instruction %v", p)
	}

	if (info.Flags&gc.RegRead != 0) && p.Reg == 0 {
		info.Flags &^= gc.RegRead
		info.Flags |= gc.RightRead /*CanRegRead |*/
	}

	if (p.From.Type == obj.TYPE_MEM || p.From.Type == obj.TYPE_ADDR) && p.From.Reg != 0 {
		info.Regindex |= RtoB(int(p.From.Reg))
		if p.Scond != 0 {
			info.Regset |= RtoB(int(p.From.Reg))
		}
	}

	if (p.To.Type == obj.TYPE_MEM || p.To.Type == obj.TYPE_ADDR) && p.To.Reg != 0 {
		info.Regindex |= RtoB(int(p.To.Reg))
		if p.Scond != 0 {
			info.Regset |= RtoB(int(p.To.Reg))
		}
	}

	if p.From.Type == obj.TYPE_ADDR && p.From.Sym != nil && (info.Flags&gc.LeftRead != 0) {
		info.Flags &^= gc.LeftRead
		info.Flags |= gc.LeftAddr
	}

	if p.As == obj.ADUFFZERO {
		info.Reguse |= RtoB(arm64.REGRT1)
		info.Regset |= RtoB(arm64.REGRT1)
	}

	if p.As == obj.ADUFFCOPY {
		// TODO(austin) Revisit when duffcopy is implemented
		info.Reguse |= RtoB(arm64.REGRT1) | RtoB(arm64.REGRT2) | RtoB(arm64.REG_R5)

		info.Regset |= RtoB(arm64.REGRT1) | RtoB(arm64.REGRT2)
	}
}
