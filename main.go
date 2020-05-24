// conip prints a minimal-size string containing every IPv4 address.
//
// The particular sequence printed is a de Bruijn sequence B(256, 4) beginning
// with four zeros. With the default text output, the alphabet is the set
// {"0", "1", "2", ..., "255"}. A "." or newline character separates each
// sequence term. The output is around 14.2 GiB.
//
// With binary output, the alphabet is the set {0, 1, 2, ..., 255}, and each
// term is written as a single byte with no separating characters. The output
// is exactly 4 GiB plus three bytes.
//
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

// terms sends the successive terms of B(256, 4) to ch. It should be called in
// a separate goroutine.
//
// To find the terms of the de Bruijn sequence, we concatenate the symbols of
// each lexicographically succeeding Lyndon word of length 1, 2, or 4. A string
// is a Lyndon word if it is lexicographically the unique minimum of its
// rotations. Each single symbol is trivially a Lyndon word. A pair of symbols
// is a Lyndon word iff its first symbol is less than its second. So, the
// interesting case is a word of length 4, u = αβγδ:
//
//	1. If α > β or α > γ or α > δ, then u is not a Lyndon word.
//	2. If α = δ, then u is not a Lyndon word.
//	3. If α = γ, then u is a Lyndon word iff β < δ.
//	4. Otherwise, u is a Lyndon word.
//
// Duval provides an algorithm to produce the lexicographically succeeding
// Lyndon word of length at most n given a current Lyndon word other than the
// maximum one. It is straightforward to modify it to skip words of length 3.
func terms(ch chan<- byte) {
	ch <- 0
	u := [4]byte{}
	for u[0] != 0xff {
		if u[3] == 0xff {
			// If the last symbol is currently the maximal one, then Duval's
			// generation algorithm would lead us to send a 3-element word. We
			// don't want that, so just check whether we'd send a 2- or
			// 1-element word, and otherwise skip to the next 4-element one.
			if u[2] == 0xff {
				if u[1] == 0xff {
					// 1-element Lyndon word.
					log.Println("1-element", u[0])
					u[0]++
					u[1], u[2], u[3] = u[0], u[0], u[0]
					ch <- u[0]
					continue
				}
				// 2-element Lyndon word.
				u[1]++
				u[2], u[3] = u[0], u[1]
				ch <- u[0]
				ch <- u[1]
				continue
			}
			// Would-be 3-element.
			u[2]++
			u[3] = u[0]
		}
		// 4-element Lyndon word.
		u[3]++
		ch <- u[0]
		ch <- u[1]
		ch <- u[2]
		ch <- u[3]
	}
	// When the loop terminates, we repeat the first three terms of the
	// de Bruijn sequence to finish the cycle.
	ch <- 0
	ch <- 0
	ch <- 0
	close(ch)
}

func main() {
	bin := false
	nl := false
	buf := 0
	o := ""
	flag.BoolVar(&bin, "bin", false, "output binary if true, text if false")
	flag.BoolVar(&nl, "n", false, "in text mode, separate terms by lines instead of .")
	flag.IntVar(&buf, "buf", 4096, "output buffer size")
	flag.StringVar(&o, "o", "", "output file name; stdout if empty")
	flag.Parse()

	out := os.Stdout
	if o != "" {
		var err error
		out, err = os.Create(o)
		if err != nil {
			panic(err)
		}
	}
	w := bufio.NewWriterSize(out, buf)
	ch := make(chan byte, 4)
	go terms(ch)
	if bin {
		for term := range ch {
			if err := w.WriteByte(term); err != nil {
				panic(err)
			}
		}
	} else {
		encs := encd[:]
		if nl {
			encs = encn[:]
		}
		_ = encs[0xff] // ensure bounds checks are optimized out
		if _, err := w.WriteString(encs[<-ch][1:]); err != nil {
			panic(err)
		}
		for term := range ch {
			if _, err := w.WriteString(encs[term]); err != nil {
				panic(err)
			}
		}
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}
}

var encd = [256]string{
	".0", ".1", ".2", ".3", ".4", ".5", ".6", ".7", ".8", ".9", ".10", ".11", ".12", ".13", ".14", ".15",
	".16", ".17", ".18", ".19", ".20", ".21", ".22", ".23", ".24", ".25", ".26", ".27", ".28", ".29", ".30", ".31",
	".32", ".33", ".34", ".35", ".36", ".37", ".38", ".39", ".40", ".41", ".42", ".43", ".44", ".45", ".46", ".47",
	".48", ".49", ".50", ".51", ".52", ".53", ".54", ".55", ".56", ".57", ".58", ".59", ".60", ".61", ".62", ".63",
	".64", ".65", ".66", ".67", ".68", ".69", ".70", ".71", ".72", ".73", ".74", ".75", ".76", ".77", ".78", ".79",
	".80", ".81", ".82", ".83", ".84", ".85", ".86", ".87", ".88", ".89", ".90", ".91", ".92", ".93", ".94", ".95",
	".96", ".97", ".98", ".99", ".100", ".101", ".102", ".103", ".104", ".105", ".106", ".107", ".108", ".109", ".110", ".111",
	".112", ".113", ".114", ".115", ".116", ".117", ".118", ".119", ".120", ".121", ".122", ".123", ".124", ".125", ".126", ".127",
	".128", ".129", ".130", ".131", ".132", ".133", ".134", ".135", ".136", ".137", ".138", ".139", ".140", ".141", ".142", ".143",
	".144", ".145", ".146", ".147", ".148", ".149", ".150", ".151", ".152", ".153", ".154", ".155", ".156", ".157", ".158", ".159",
	".160", ".161", ".162", ".163", ".164", ".165", ".166", ".167", ".168", ".169", ".170", ".171", ".172", ".173", ".174", ".175",
	".176", ".177", ".178", ".179", ".180", ".181", ".182", ".183", ".184", ".185", ".186", ".187", ".188", ".189", ".190", ".191",
	".192", ".193", ".194", ".195", ".196", ".197", ".198", ".199", ".200", ".201", ".202", ".203", ".204", ".205", ".206", ".207",
	".208", ".209", ".210", ".211", ".212", ".213", ".214", ".215", ".216", ".217", ".218", ".219", ".220", ".221", ".222", ".223",
	".224", ".225", ".226", ".227", ".228", ".229", ".230", ".231", ".232", ".233", ".234", ".235", ".236", ".237", ".238", ".239",
	".240", ".241", ".242", ".243", ".244", ".245", ".246", ".247", ".248", ".249", ".250", ".251", ".252", ".253", ".254", ".255",
}

var encn = [256]string{
	"\n0", "\n1", "\n2", "\n3", "\n4", "\n5", "\n6", "\n7", "\n8", "\n9", "\n10", "\n11", "\n12", "\n13", "\n14", "\n15",
	"\n16", "\n17", "\n18", "\n19", "\n20", "\n21", "\n22", "\n23", "\n24", "\n25", "\n26", "\n27", "\n28", "\n29", "\n30", "\n31",
	"\n32", "\n33", "\n34", "\n35", "\n36", "\n37", "\n38", "\n39", "\n40", "\n41", "\n42", "\n43", "\n44", "\n45", "\n46", "\n47",
	"\n48", "\n49", "\n50", "\n51", "\n52", "\n53", "\n54", "\n55", "\n56", "\n57", "\n58", "\n59", "\n60", "\n61", "\n62", "\n63",
	"\n64", "\n65", "\n66", "\n67", "\n68", "\n69", "\n70", "\n71", "\n72", "\n73", "\n74", "\n75", "\n76", "\n77", "\n78", "\n79",
	"\n80", "\n81", "\n82", "\n83", "\n84", "\n85", "\n86", "\n87", "\n88", "\n89", "\n90", "\n91", "\n92", "\n93", "\n94", "\n95",
	"\n96", "\n97", "\n98", "\n99", "\n100", "\n101", "\n102", "\n103", "\n104", "\n105", "\n106", "\n107", "\n108", "\n109", "\n110", "\n111",
	"\n112", "\n113", "\n114", "\n115", "\n116", "\n117", "\n118", "\n119", "\n120", "\n121", "\n122", "\n123", "\n124", "\n125", "\n126", "\n127",
	"\n128", "\n129", "\n130", "\n131", "\n132", "\n133", "\n134", "\n135", "\n136", "\n137", "\n138", "\n139", "\n140", "\n141", "\n142", "\n143",
	"\n144", "\n145", "\n146", "\n147", "\n148", "\n149", "\n150", "\n151", "\n152", "\n153", "\n154", "\n155", "\n156", "\n157", "\n158", "\n159",
	"\n160", "\n161", "\n162", "\n163", "\n164", "\n165", "\n166", "\n167", "\n168", "\n169", "\n170", "\n171", "\n172", "\n173", "\n174", "\n175",
	"\n176", "\n177", "\n178", "\n179", "\n180", "\n181", "\n182", "\n183", "\n184", "\n185", "\n186", "\n187", "\n188", "\n189", "\n190", "\n191",
	"\n192", "\n193", "\n194", "\n195", "\n196", "\n197", "\n198", "\n199", "\n200", "\n201", "\n202", "\n203", "\n204", "\n205", "\n206", "\n207",
	"\n208", "\n209", "\n210", "\n211", "\n212", "\n213", "\n214", "\n215", "\n216", "\n217", "\n218", "\n219", "\n220", "\n221", "\n222", "\n223",
	"\n224", "\n225", "\n226", "\n227", "\n228", "\n229", "\n230", "\n231", "\n232", "\n233", "\n234", "\n235", "\n236", "\n237", "\n238", "\n239",
	"\n240", "\n241", "\n242", "\n243", "\n244", "\n245", "\n246", "\n247", "\n248", "\n249", "\n250", "\n251", "\n252", "\n253", "\n254", "\n255",
}
