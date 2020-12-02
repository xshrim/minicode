match_about(){
  cmds=($(set | awk -F'_| ' '/^match_.*\(\)/ {print $2}'));echo 'keys:' ${cmds[@]}
  uname -smrn
  echo "Memory used: $(/usr/bin/free -h | awk '/^Mem/ {print $3}')"
}

match_tabs(){
  qdbus org.kde.plasma.browser_integration /TabsRunner org.kde.plasma.browser_integration.TabsRunner.GetTabs | sed 's|^title:|title:\|\||' | awk -F':' '/^(title|url)/ {print $(NF)}' | tac | awk 'ORS=NR%2?FS:RS' | sed 's|//|https://|'
}
run_tabs(){
  xdg-open "$1"
}

match_fd() {
  while read -r line; do
    echo "$line"
  done < <(fd "$1" ~)
}
run_fd() {
  xdg-open $1
}

match_loc() {
  while read -r line; do
    echo "$line"
  done < <(locate "$1")
}
run_loc() {
  xdg-open $1
}

match_radio(){
 echo "http://ais-sa2-dal01-1.cdnstream.com:80/1989_64.aac||Rock mix||rock"
 echo "http://node1.mingusradio.com:7646/rock||Mingus radio||rock"
 echo "http://stream.punkrockers-radio.de:8000/mp3||PunkRockers||punk"
 echo "http://94.23.26.22:8090/live.mp3||Punk fm||punk"
}
run_radio(){
 if pacman -Qq clementine 2>/dev/null; then
   qdbus org.mpris.MediaPlayer2.clementine /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.OpenUri "$1" ||clementine "$1"
   return 0
 fi
 if pacman -Qq vlc 2>/dev/null; then
   qdbus org.mpris.MediaPlayer2.vlc /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.OpenUri "$1" || vlc "$1"
   return 0
 fi
}

match_man(){
  while read -r name id sep txt; do
    echo -e "${name}||${name}:\t${txt}"
  done < <(man -k "$1")
}
run_man(){
  qdbus org.kde.plasmashell /org/kde/osdService showText 'help' "$1"
  man -Thtml "$1" >"/tmp/krunner-man.html" && xdg-open "/tmp/krunner-man.html"
}


match_df() {
  while read -r line; do
    echo "${line##* }||${line}"
  done < <(df -x tmpfs| awk '/^\// {print $1"\t"$5" "$(NF)}')
}
run_df(){
  dolphin $1
}

match_err(){
  while read -r m j h m u mes; do
    echo "${j} ${h:0:5} ${u%%\[*} ${mes}"
  done < <(SYSTEMD_COLORS=0 journalctl -b0 -p3 -r -n12 --no-tail --no-pager)
}

match_env() {
  #show bash and zsh config !
  $SHELL -ic "/usr/bin/env"| grep "="| grep -i "$1" --color=never| grep -v "^_"
}

match_new(){
python - "$0"<<'EOF'
import json, urllib.request
with urllib.request.urlopen(f"https://forum.manjaro.org/c/announcements.json") as f_url:
    req = f_url.read()
topics = json.loads(req)['topic_list']['topics']
for topic in [t for t in topics if not t['title'].startswith('About') and not t['closed']]:
    print(f"https://forum.manjaro.org/t/{topic['slug']}/{topic['id']}/||{topic['fancy_title']}")
EOF
}
run_new(){
  xdg-open "$1"
}

match_pac(){
  pacman -Qq | grep $1
}
run_pac(){
  qdbus org.kde.klipper /klipper org.kde.klipper.klipper.setClipboardContents $1
  qdbus org.kde.plasmashell /org/kde/osdService showText 'manjaro' "$1"
  pamac-manager --details=$1
}

__getallgit(){
python - "$0"<<'EOF'
from pathlib import Path
wpath = "/home/Data/Patrick/workspace/"  # CHANGE PATH !
for p in Path(wpath).glob('**/.git/config'):
   with open(p, "r") as fp:
       for line in fp:
           if "url =" in line:
               ppath = str(p)[:-11]
               url = line.split("=", 2)[1].strip()
               print(f"{ppath}||{url.split('/')[-1]}\t{ppath[len(wpath):]}||{url}")
               break
EOF
}
match_git(){
 ## list git projects ##
 if ! [[ -s "/tmp/krunner-allgit" ]] && __getallgit > "/tmp/krunner-allgit"    # save first run of day
 grep "$1" "/tmp/krunner-allgit" --color=never
}
run_git(){
 qdbus org.kde.plasmashell /org/kde/osdService showText 'git' "$1"
 dolphin "$1"
 qdbus org.kde.klipper /klipper org.kde.klipper.klipper.setClipboardContents "$1"
 session_num=$(qdbus org.kde.konsole /Windows/1 org.kde.konsole.Window.newSession "git" "$1")
 qdbus org.kde.konsole  /Sessions/${session_num} org.kde.konsole.Session.runCommand "git status" &
}
