all:V: blit.js test.js

blit.js: u.rj cpu.rj rom.rj mem.rj telnet.rj glrender.rj blit.rj
	ratjs $prereq > $target

test.js: u.rj cpu.rj test.rj test_cpu.rj test_timing.rj
	ratjs $prereq > $target
