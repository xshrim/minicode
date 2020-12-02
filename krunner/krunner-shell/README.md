### Krunner shell

return shell output
can run commands

---

### configure

all in bash file `~/.config/krunnershell.sh`


create actions
In krunner run `action:`

In bash file create one or two function for one action
 - match_action() {...}    return text list
 - run_action() {...}      optional, click on line
Optional for run() : in list match(), we can return action to run and title "action||title"

---

### install

After install run `kquitapp5 krunner` if not reboot

---

### Actions Examples

in file: `~/.config/krunnershell.sh`

 * list available actions with this krunner plugin `about:`
```bash
match_about(){
  cmds=($(set | awk -F'_| ' '/^match_.*\(\)/ {print $2}'));echo 'keys:' ${cmds[@]}
  uname -smrn
  echo "Memory used: $(/usr/bin/free -h | awk '/^Mem/ {print $3}')"
}
```

 * list browser tabs `tabs:`
```bash
match_tabs(){
  qdbus org.kde.plasma.browser_integration /TabsRunner org.kde.plasma.browser_integration.TabsRunner.GetTabs | sed 's|^title:|title:\|\||' | awk -F':' '/^(title|url)/ {print $(NF)}' | tac | awk 'ORS=NR%2?FS:RS' | sed 's|//|https://|'
}
run_tabs(){
  xdg-open "$1"
}
```

 * play web radio `radio:`
 ```bash
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
 ```

 * display man in browser `man:`
```bash
match_man(){
  while read -r name id sep txt; do
    echo -e "${name}||${name}:\t${txt}"
  done < <(man -k "$1")
}
run_man(){
  qdbus org.kde.plasmashell /org/kde/osdService showText 'help' "$1"
  man -Thtml "$1" >"/tmp/krunner-man.html" && xdg-open "/tmp/krunner-man.html"
}
```

 * display mounted partitions `df:` and open
 (action to run is different from title)
```bash
match_df() {
  while read -r line; do
    echo "${line##* }||${line}"
  done < <(df -x tmpfs| awk '/^\// {print $1"\t"$5" "$(NF)}')
}
run_df(){
  dolphin $1
}
```

 * systemd errors `err:`
```bash
match_err(){
  while read -r m j h m u mes; do
    echo "${j} ${h:0:5} ${u%%\[*} ${mes}"
  done < <(SYSTEMD_COLORS=0 journalctl -b0 -p3 -r -n12 --no-tail --no-pager)
}
```

 * display process `ps:`
 (action to run is different from title)
```bash
match_ps(){
  while read -r u p t c; do
    echo "$p||$u $p $c $t"
  done < <(ps -aux -t|grep $1|grep -v '"'|grep -v "grep"|awk '{print $1" "$2" "$8" "$(NF)}'|head -n10)
}
run_ps(){
  qdbus org.kde.klipper /klipper org.kde.klipper.klipper.setClipboardContents $1  #save PID
  qdbus org.kde.plasmashell /org/kde/osdService showText 'system'  "$1"
  ksysguard    # can paste PID in search bar
}
```

 * get env variables
```bash
match_env() {
  #show bash and zsh config !
  $SHELL -ic "/usr/bin/env"| grep "="| grep -i "$1" --color=never| grep -v "^_"
}
```

 * news from one forum `new:`
```bash
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
```

 * pacman infos `pac:`
```bash
match_pac(){
  pacman -Qq | grep $1
}
run_pac(){
  qdbus org.kde.klipper /klipper org.kde.klipper.klipper.setClipboardContents $1
  qdbus org.kde.plasmashell /org/kde/osdService showText 'manjaro' "$1"
  pamac-manager --details=$1
}
```

 * get/open git project `git:`
 (first call is long !!!)
 ```bash
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
 ```
