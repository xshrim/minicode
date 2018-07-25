#                                    LS_COLORS
# Maintainers: Magnus Woldrich <m@japh.se>,
#              Ryan Delaney <ryan.delaney@gmail.com> OpenGPG: 0D98863B4E1D07B6
#         URL: https://github.com/trapd00r/LS_COLORS
#     Version: 0.254
#     Updated: Tue Mar 29 21:25:30 AEST 2016
#
#   This is a collection of extension:color mappings, suitable to use as your
#   LS_COLORS environment variable. Most of them use the extended color map,
#   described in the ECMA-48 document; in other words, you'll need a terminal
#   with capabilities of displaying 256 colors.
#
#   As of this writing, around 300 different filetypes/extensions is supported.
#   That's indeed a lot of extensions, but there's a lot more! Therefore I need
#   your help.
#
#   Fork this project on github, add the extensions you are missing, and send me
#   a pull request.
#
#   For files that usually ends up next to each other, like html, css and js,
#   try to pick colors that fit nicely together. Filetypes with multiple
#   possible extensions, like htm and html, should have the same color.

# This program is distributed in the hope that it will be useful, but WITHOUT ANY
# WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
# PARTICULAR PURPOSE.  See the Perl Artistic License for more details.
#
# This program is free software: you can redistribute it and/or modify it under
# the terms of the Perl Artistic License as published by the Perl Foundation,
# either version 1.0 of the License, or (at your option) any later version.
#
# You should have received a copy of the Perl Artistic License along
# with this program.  If not, see <http://www.perlfoundation.org/artistic_license_1_0>.

#introduce
# A color init string consists of one or more of the following numeric codes:
# Attribute codes:
# 00=none 01=bold 04=underscore 05=blink 07=reverse 08=concealed
# Text color codes:
# 30=black 31=red 32=green 33=yellow 34=blue 35=magenta 36=cyan 37=white
# Background color codes:
# 40=black 41=red 42=green 43=yellow 44=blue 45=magenta 46=cyan 47=white

# The format for 256 color escape codes is 38;5;colorN for foreground colors and 48;5;colorN for background colors.
# examples:
# .mp3  38;5;160                   # Set fg color to color 160
# .flac 48;5;240                   # Set bg color to color 240
# .ogg  38;5;160;48;5;240          # Set fg color 160 *and* bg color 240.
# .wav  01;04;05;38;5;160;48;5;240 # Pure madness: make bold (01), underlined (04), blink (05), fg color 160, and bg color 240!

# core {{{1
BLK                   38;5;68
CAPABILITY            38;5;17
CHR                   38;5;5;1
DIR                   38;5;4
DOOR                  38;5;127
EXEC                  38;5;128;1
FIFO                  38;5;13
FILE                  38;5;10
LINK                  target
MULTIHARDLINK         38;5;222;1
# "NORMAL don't reset the bold attribute -
# https://github.com/trapd00r/LS_COLORS/issues/11
#NORMAL                38;5;254
NORMAL                38;5;10
ORPHAN                48;5;1;38;5;232;1
OTHER_WRITABLE        38;5;220;1
SETGID                48;5;3;38;5;0
SETUID                38;5;220;1;3;100;1
SOCK                  38;5;197
STICKY                38;5;86;48;5;234
STICKY_OTHER_WRITABLE 48;5;235;38;5;139;3

*LS_COLORS 48;5;89;38;5;197;1;3;4;7 # :-)
# }}}
# documents {{{1
*README               38;5;181;1
*README.rst           38;5;181;1
*LICENSE              38;5;181;1
*COPYING              38;5;181;1
*INSTALL              38;5;181;1
*COPYRIGHT            38;5;181;1
*AUTHORS              38;5;181;1
*HISTORY              38;5;181;1
*CONTRIBUTORS         38;5;181;1
*PATENTS              38;5;181;1
*VERSION              38;5;181;1
*NOTICE               38;5;181;1
*CHANGES              38;5;181;1
.log                  38;5;180
# plain-text {{{2
.txt                  38;5;183
# markup {{{2
.etx                  38;5;123
.info                 38;5;123
.markdown             38;5;123
.md                   38;5;123
.mkd                  38;5;123
.nfo                  38;5;123
.pod                  38;5;123
.rst                  38;5;123
.tex                  38;5;123
.textile              38;5;123
# key-value, non-relational data {{{2
.json                 38;5;135
.msg                  38;5;135
.pgn                  38;5;135
.rss                  38;5;135
.xml                  38;5;135
.yaml                 38;5;135
.yml                  38;5;135
.rdata                38;5;135
# }}}
# binary {{{2
.cbr                  38;5;141
.cbz                  38;5;141
.chm                  38;5;141
.djvu                 38;5;141
.pdf                  38;5;141
# wps
.wps                  38;5;117
# words {{{3
.docm                 38;5;111;4
.doc                  38;5;111
.docx                 38;5;111
.eps                  38;5;111
.ps                   38;5;111
.odb                  38;5;111
.odt                  38;5;111
.rtf                  38;5;111
# presentation {{{3
.odp                  38;5;122
.pps                  38;5;122
.ppt                  38;5;122
.pptx                 38;5;122
#   Powerpoint show
.ppts                 38;5;122
#   Powerpoint with enabled macros
.pptxm                38;5;122;4
#   Powerpoint show with enabled macros
.pptsm                38;5;122;4
# spreadsheet {{{3
.csv                  38;5;139
#   Open document spreadsheet
.ods                  38;5;132
.xla                  38;5;133
#   Excel spreadsheet
.xls                  38;5;178
.xlsx                 38;5;178
#   Excel spreadsheet with macros
.xlsxm                38;5;178;4
#   Excel module
.xltm                 38;5;151;4
.xltx                 38;5;151
# }}}
# }}}
# configs {{{2
*cfg                  38;5;190
*conf                 38;5;190
*rc                   38;5;190
.cf                   38;5;190
.cfg                  38;5;190
.conf                 38;5;190
.rc                   38;5;190
.ini                  38;5;190
.plist                38;5;190
#   vim
.viminfo              38;5;166
.vimrc                38;5;166
.zshrc                38;5;166
.bashrc               38;5;166
.xonshrc              38;5;166
#   cisco VPN client configuration
.pcf                  38;5;172
#   adobe photoshop proof settings file
.psf                  38;5;175
# }}}
# }}}
# code {{{1
# version control {{{2
.git                  38;5;213
.gitignore            38;5;225
.gitattributes        38;5;225
.gitmodules           38;5;225

# shell {{{2
.awk                  38;5;49
.bash                 38;5;49
.sed                  38;5;49
.sh                   38;5;49
.csh                  38;5;49
.ksh                  38;5;49
.fish                 38;5;49
.zsh                  38;5;49
.vim                  38;5;49

# interpreted {{{2
.ahk                  38;5;195
# python
.py                   38;5;35
# perl
.pl                   38;5;93
.t                    38;5;93
# sql
.msql                 38;5;97
.mysql                38;5;97
.pgsql                38;5;97
.sql                  38;5;97
#   Tool Command Language
.tcl                  38;5;95;1
# R language
.r                    38;5;63
# GrADS script
.gs                   38;5;89
# repo
.repo                 38;5;98
# compiled {{{2
#
#   assembly language
.asm                  38;5;100
#   LISP
.cl                   38;5;96
.lisp                 38;5;96
#   lua
.lua                  38;5;103
#   Moonscript
.moon                 38;5;116
#   C
.c                    38;5;142
.h                    38;5;136
.obj                  38;5;144
.tcc                  38;5;136
#   C++
.c++                  38;5;10
.h++                  38;5;137
.hpp                  38;5;137
.hxx                  38;5;137
.ii                   38;5;138
#   method file for Objective C
.m                    38;5;201
#   Csharp
.cc                   38;5;184
.cs                   38;5;184
.cp                   38;5;214
.cpp                  38;5;214
.cxx                  38;5;184
.asp                  38;5;211
.aspx                 38;5;212
#   Crystal
.cr                   38;5;159
#   Google golang
.go                   38;5;112
#   scala
.scala                38;5;9
#   Solidity
.sol                  38;5;90
#   fortran
.f                    38;5;58
.for                  38;5;58
.ftn                  38;5;58
#   pascal
.s                    38;5;30
#   Rust
.rs                   38;5;79
#   D
.d                    38;5;23
#   Swift
.swift                38;5;33
#   ?
.sx                   38;5;81
#   interface file in GHC - https://github.com/trapd00r/LS_COLORS/pull/9
.hi                   38;5;110
#   haskell
.hs                   38;5;69
.lhs                  38;5;69

# binaries {{{2
# compiled apps for interpreted languages
.pyc                  38;5;14
# }}}
# html {{{2
.css                  38;5;216;1
.less                 38;5;226;1
.sass                 38;5;227;1
.scss                 38;5;228;1
.htm                  38;5;200;1
.html                 38;5;201;1
.jhtm                 38;5;165;1
.mht                  38;5;199;1
.eml                  38;5;162;1
.mustache             38;5;161;1
# }}}
# java {{{2
.coffee               38;5;205;1
.java                 38;5;99;1
.js                   38;5;3;1
.mjs                  38;5;203;1
.jsm                  38;5;204;1
.jsm                  38;5;205;1
.jsp                  38;5;219;1
# }}}
# php {{{2
.php                  38;5;25
#   CakePHP view scripts and helpers
.ctp                  38;5;24
#   Twig template engine
.twig                 38;5;67
# }}}
# vb/a {{{2
.vb                   38;5;146
.vba                  38;5;146
.vbs                  38;5;146
# 2}}}
# Build stuff {{{2
*Dockerfile           38;5;88
*Makefile             38;5;88
*MANIFEST             38;5;88
*pm_to_blib           38;5;88
.dockerfile           38;5;88
.dockerignore         38;5;88
.makefile             38;5;88
.manifest             38;5;88
.pm_to_blib           38;5;88
# automake
.am                   38;5;56
.in                   38;5;56
.hin                  38;5;56
.scan                 38;5;56
.m4                   38;5;56
.old                  38;5;56
.out                  38;5;56
.skip                 38;5;224
# }}}
# patch files {{{2
.diff                 48;5;23;38;5;24
.patch                48;5;29;38;5;30;1
#}}}
# graphics {{{1
.bmp                  38;5;64
.tiff                 38;5;64
.tif                  38;5;64
.tga                  38;5;64
.cdr                  38;5;64
.gif                  38;5;64
.ico                  38;5;64
.jpeg                 38;5;64
.jpg                  38;5;64
.fli                  38;5;64
.nth                  38;5;64
.png                  38;5;64
.psd                  38;5;64
.pcx                  38;5;64
.mng                  38;5;64
.pbm                  38;5;64
.pgm                  38;5;64
.xpm                  38;5;64
.ppm                  38;5;64
# }}}
# vector {{{1
.ai                   38;5;51
.eps                  38;5;51
.epsf                 38;5;51
.drw                  38;5;51
.ps                   38;5;51
.svg                  38;5;51
# }}}
# video {{{1
.avi                  38;5;75
.divx                 38;5;75
.ifo                  38;5;75
.m2v                  38;5;75
.m4v                  38;5;75
.mkv                  38;5;75
.mov                  38;5;75
.mp4                  38;5;75
.mpeg                 38;5;75
.mpg                  38;5;75
.ogm                  38;5;75
.rmvb                 38;5;75
.sample               38;5;75
.wmv                  38;5;75
  # mobile/streaming {{{2
.3g2                  38;5;73
.3gp                  38;5;73
.gp3                  38;5;73
.webm                 38;5;73
.gp4                  38;5;73
.asf                  38;5;73
.flv                  38;5;73
.ts                   38;5;73
.ogv                  38;5;73
.f4v                  38;5;73
  # }}}
  # lossless {{{2
.vob                  38;5;85;1
# }}}
# audio {{{1
.3ga                  38;5;72;1
.s3m                  38;5;72;1
.aac                  38;5;72;1
.au                   38;5;72;1
.dat                  38;5;72;1
.dts                  38;5;72;1
.fcm                  38;5;72;1
.m4a                  38;5;72;1
.mid                  38;5;72;1
.midi                 38;5;72;1
.mod                  38;5;72;1
.mp3                  38;5;72;1
.mp4a                 38;5;72;1
.oga                  38;5;72;1
.ogg                  38;5;72;1
.opus                 38;5;72;1
.s3m                  38;5;72;1
.sid                  38;5;72;1
.wma                  38;5;72;1
# lossless
.ape                  38;5;80;1
.aiff                 38;5;80;1
.cda                  38;5;80;1
.flac                 38;5;80;1
.alac                 38;5;80;1
.midi                 38;5;80;1
.pcm                  38;5;80;1
.wav                  38;5;80;1
.wv                   38;5;80;1
.wvc                  38;5;80;1

# }}}
# fonts {{{1
.afm                  38;5;157
.fon                  38;5;157
.font                 38;5;157
.fnt                  38;5;157
.pfb                  38;5;157
.pfm                  38;5;157
.ttf                  38;5;157
.ttc                  38;5;157
.otf                  38;5;157
.eot                  38;5;157
.aiff                 38;5;157
.woff                 38;5;157
.woff2                38;5;157
#   postscript fonts
.pfa                  38;5;156
# }}}
# archives {{{1
.7z                   38;5;130
.a                    38;5;130
.arj                  38;5;130
.bz2                  38;5;130
.cpio                 38;5;130
.gz                   38;5;130
.bz2                  38;5;130
.tb2                  38;5;130
.tz2                  38;5;130
.tbz2                 38;5;130
.lrz                  38;5;130
.lz                   38;5;130
.lzh                  38;5;130
.lzma                 38;5;130
.lzo                  38;5;130
.rar                  38;5;130
.s7z                  38;5;130
.sz                   38;5;130
.tar                  38;5;130
.taz                  38;5;130
.tgz                  38;5;130
.xz                   38;5;130
.z                    38;5;130
.zip                  38;5;130
.zipx                 38;5;130
.zoo                  38;5;130
.zpaq                 38;5;130
.zz                   38;5;130
  # packaged apps {{{2
.apk                  38;5;160
.deb                  38;5;160
.rpm                  38;5;160
.jad                  38;5;160
.jar                  38;5;160
.cab                  38;5;160
.mdf                  38;5;160
.pak                  38;5;160
.pk3                  38;5;160
.vdf                  38;5;160
.vpk                  38;5;160
.bsp                  38;5;160
.dmg                  38;5;160
# }}}
# windows
.exe                  38;5;52
.msi                  38;5;52
.rsp                  38;5;52
.btm                  38;5;52
.dll                  38;5;52
.osx                  38;5;52
.ocx                  38;5;52
.cmd                  38;5;52
.bat                  38;5;52
.reg                  38;5;52
# segments from 0 to three digits after first extension letter {{{2
.r[0-9]{0,2}          38;5;145
.zx[0-9]{0,2}         38;5;145
.z[0-9]{0,2}          38;5;145
# partial files
.part                 38;5;24
  # }}}
# partition images {{{2
.iso                  38;5;94
.bin                  38;5;94
.nrg                  38;5;94
.qcow                 38;5;94
.sparseimage          38;5;94
.toast                38;5;94
.vcd                  38;5;94
.vmdk                 38;5;94
.vdi                  38;5;94
.vdisk                38;5;94
.box                  38;5;94
.img                  38;5;94
.vbox                 38;5;94
# }}}
# databases {{{2
.accdb                38;5;60
.accde                38;5;60
.accdr                38;5;60
.accdt                38;5;60
.db                   38;5;60
.fmp12                38;5;60
.fp7                  38;5;60
.localstorage         38;5;60
.mdb                  38;5;60
.mde                  38;5;60
.sqlite               38;5;60
.typelib              38;5;60
# NetCDF database
.nc                   38;5;61
# }}}
# tempfiles {{{1
# undo files
.pacnew               38;5;230
.un~                  38;5;230
.orig                 38;5;230
# backups
.bup                  38;5;231
.test                 38;5;231
.bak                  38;5;231
.res                  38;5;231
.hd                   38;5;231
.new                  38;5;231
.old                  38;5;231
.o                    38;5;231 #   *nix Object file (shared libraries, core dumps etc)
.rlib                 38;5;231 #   Static rust library
# temporary files
.swp                  38;5;229
.swo                  38;5;229
.tmp                  38;5;229
.temp                  38;5;229
.sassc                38;5;229
# state files
.et                   38;5;193
.pid                  38;5;193
.state                38;5;193
*lockfile             38;5;193
# error logs
.err                  38;5;187;1
.error                38;5;187;1
.stderr               38;5;187;1
# state dumps
.dump                 38;5;188
.stackdump            38;5;188
.zcompdump            38;5;188
.zwc                  38;5;188
# tcpdump, network traffic capture
.pcap                 38;5;102
.cap                  38;5;102
.dmp                  38;5;102
# macOS
.ds_store             38;5;174
.localized            38;5;174
.cfusertextencoding   38;5;174
# }}}
# hosts {{{1
# /etc/hosts.{deny,allow}
.allow                38;5;153
.deny                 38;5;152
# }}}
# systemd {{{1
# http://www.freedesktop.org/software/systemd/man/systemd.unit.html
.service              38;5;43
*@.service            38;5;43
.socket               38;5;43
.sock                 38;5;43
.swap                 38;5;43
.device               38;5;43
.mount                38;5;43
.automount            38;5;43
.target               38;5;43
.path                 38;5;43
.timer                38;5;43
.snapshot             38;5;43
# }}}
# transaction {{{1
.tx                   38;5;60
.block                38;5;60
.idx                  38;5;60
.wiz                  38;5;60
# }}}
# metadata {{{1
.application          38;5;155
.cue                  38;5;155
.description          38;5;155
.directory            38;5;155
.m3u                  38;5;155
.m3u8                 38;5;155
.md5                  38;5;155
.properties           38;5;155
.sfv                  38;5;155
.srt                  38;5;155
.theme                38;5;155
.torrent              38;5;155
.urlview              38;5;155
# }}}
# encrypted data {{{1
.asc                  38;5;192;3
.bfe                  38;5;192;3
.enc                  38;5;192;3
.gpg                  38;5;192;3
.signature            38;5;192;3
.sig                  38;5;192;3
.p12                  38;5;192;3
.pem                  38;5;192;3
.crt                  38;5;192;3
.pgp                  38;5;192;3
.asc                  38;5;192;3
.enc                  38;5;192;3
.sig                  38;5;192;3
# 1}}}
# emulators {{{1
.32x                  38;5;107
.cdi                  38;5;107
.fm2                  38;5;107
.rom                  38;5;107
.sav                  38;5;107
.st                   38;5;107
  # atari
.a00                  38;5;108
.a52                  38;5;108
.a64                  38;5;108
.a78                  38;5;108
.adf                  38;5;108
.atr                  38;5;108
  # nintendo
.gb                   38;5;109
.gba                  38;5;109
.gbc                  38;5;109
.gel                  38;5;109
.gg                   38;5;109
.ggl                  38;5;109
.ipk                  38;5;109 # Nintendo (DS Packed Images)
.j64                  38;5;109
.nds                  38;5;109
.nes                  38;5;109
  # Sega
.sms                  38;5;194
# }}}
# unsorted {{{1
#
#   Portable Object Translation for GNU Gettext
.pot                  38;5;101
#   CAD files for printed circuit boards
.pcb                  38;5;101
#   groff (rendering app for texinfo)
.mm                   38;5;101
#   perldoc
.pod                  38;5;101
#   GIMP brush
.gbr                  38;5;101
# printer spool file
.spl                  38;5;101
#   GIMP project file
.scm                  38;5;101
# RStudio project file
.rproj                38;5;11
#   Nokia Symbian OS files
.sis                  38;5;65
.1p                   38;5;65
.3p                   38;5;65
.cnc                  38;5;65
.def                  38;5;65
.ex                   38;5;65
.example              38;5;65
.feature              38;5;65
.ger                  38;5;65
.map                  38;5;65
.mf                   38;5;65
.mfasl                38;5;65
.mi                   38;5;65
.mtx                  38;5;65
.pc                   38;5;65
.pi                   38;5;65
.plt                  38;5;65
.pm                   38;5;65
.rb                   38;5;65
.rdf                  38;5;65
.rst                  38;5;65
.ru                   38;5;65
.sch                  38;5;65
.sty                  38;5;65
.sug                  38;5;65
.t                    38;5;65
.tdy                  38;5;65
.tfm                  38;5;65
.tfnt                 38;5;65
.tg                   38;5;65
.vcard                38;5;65
.vcf                  38;5;65
.xln                  38;5;65
#   AppCode files
.iml                  38;5;12
#   Xcode files
.xcconfig             38;5;158
.entitlements         38;5;158
.strings              38;5;158
.storyboard           38;5;158
.xcsettings           38;5;158
.xib                  38;5;158
# }}}
# termcap {{{1
TERM ansi
TERM color-xterm
TERM con132x25
TERM con132x30
TERM con132x43
TERM con132x60
TERM con80x25
TERM con80x28
TERM con80x30
TERM con80x43
TERM con80x50
TERM con80x60
TERM cons25
TERM console
TERM cygwin
TERM dtterm
TERM Eterm
TERM eterm-color
TERM gnome
TERM gnome-256color
TERM jfbterm
TERM konsole
TERM kterm
TERM linux
TERM linux-c
TERM mach-color
TERM mlterm
TERM putty
TERM rxvt
TERM rxvt-256color
TERM rxvt-cygwin
TERM rxvt-cygwin-native
TERM rxvt-unicode
TERM rxvt-unicode-256color
TERM rxvt-unicode256
TERM screen
TERM screen-256color
TERM screen-256color-bce
TERM screen-bce
TERM screen-w
TERM screen.linux
TERM screen.rxvt
TERM terminator
TERM vt100
TERM xterm
TERM xterm-16color
TERM xterm-256color
TERM xterm-88color
TERM xterm-color
TERM xterm-debian
# }}}


# vim: ft=dircolors:fdm=marker:et:sw=2:
