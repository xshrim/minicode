#!/bin/env python3

"""
https://invent.kde.org/frameworks/krunner/-/blob/master/src/data/org.kde.krunner1.xml
https://api.kde.org/frameworks/krunner/html/runnercontext_8cpp_source.html#l00374
https://techbase.kde.org/Development/Tutorials/Plasma4/AbstractRunner

"""
'''
qdbus org.kde.krunner /App org.kde.krunner.App.querySingleRunner kshell ls:/etc/pa
qdbus --literal org.manjaro.shell.krunner /krunner org.kde.krunner1.Match "ls:/etc/pac"
'''
'''
lire ~/.config/kshell.sh

# creer 2 funtions par plugin
### ls: ###
run_ls (){
  dolphin $1
}
match_ls(){
  # retourne une liste
}

main(){
  kshell.sh ls match        (match_ls)
  kshell.sh ls run /etc/    (run_ls /etc/)
}
'''

from pathlib import Path
import subprocess
from dataclasses import dataclass
import dbus.service
from dbus.mainloop.glib import DBusGMainLoop
import gi
from gi.repository import GLib


DBusGMainLoop(set_as_default=True)

objpath = "/krunner"
iface = "org.kde.krunner1"


class Runner(dbus.service.Object):
    def __init__(self):
        """ set dbus and load alpm """
        dbus.service.Object.__init__(self, dbus.service.BusName("org.manjaro.shell.krunner", dbus.SessionBus()), objpath)
        # read config match_*** functions """
        self.config = f"{Path.home()}/.config/krunnershell.sh"
        self.actions = ()
        self.prefix = ""
        self.size = 0
        self.loadConfig()

    def loadConfig(self):
        self.actions = ()
        self.prefix = ""
        self.size = 0
        if not Path(self.config).exists():
            with open(self.config, "w") as file_out:
                cmd=r"cmds=($(set | awk -F'_| ' '/^match_.*\(\)/ {print $2}'));echo 'keys:' ${cmds[@]}"
                file_out.write("\nmatch_about(){\n\t"+ cmd +"\n\tuname -smrn\n}\n")
        try:
            with open(self.config, "r") as file_in:
                for line in file_in:
                    if line.startswith('match_'):
                        name = line[6:].rsplit('(', 1)[0]
                        if name:
                            self.actions += (name,)
            self.size = Path(self.config).stat().st_size
            print(self.actions)
        except FileNotFoundError:
            pass

    @dbus.service.method(iface, in_signature='s', out_signature='a(sssida{sv})')
    def Match(self, query: str):
        """ get the matches and returns a packages list """
        self.prefix = ""
        query = query.strip().lower()

        try:
            if self.size != Path(self.config).stat().st_size:
                self.loadConfig()
        except FileNotFoundError:
            self.loadConfig()

        #print(f"match: {query}...")
        if not self.actions or not ":" in query:
            return []

        #TEST plugin without datas
        #return [tuple([n, n, "", 32, 0.5, {"subtext": ""}]) for n in self.actions]

        self.prefix = query.split(':', 1)[0].replace(" ", "_")
        if self.prefix not in self.actions:
            return []
        query = query[len(self.prefix)+1:].strip()
        #return [tuple([query, query, "", 32, 0.5, {"subtext": ""}])]

        out, _ = subprocess.Popen(
            f"bash -c 'source {self.config} && match_{self.prefix} \"{query}\"'",
            shell=True, text=True,
            stdout=subprocess.PIPE
        ).communicate()
        print(f"stdout cmd: match_{self.prefix} \"{query}\" : {out}")

        ret = []
        rel = 1
        for node in out.splitlines():
            if "||" in node:    #if title != to run, split line and set [0] in data(link) and [1] in titre
                data = node.split("||", 2)
                for _ in range(len(data), 3):
                    data.append("")
            else:
                data = [node, node, ""]
            ret.append(tuple([data[0], data[1], "", 32, rel, {"subtext": data[2]}]))
            rel = rel - 0.02

        return ret


    @dbus.service.method(iface, in_signature='ss')
    def Run(self, data: str, _action_id: str):
        """ click on item, run command """
        pid = subprocess.Popen(f"bash -c 'source {self.config} && run_{self.prefix} {data}'", shell=True).pid
        #pid = subprocess.Popen(f"xdg-open \"{data}\"", shell=True).pid



runner = Runner()
loop = GLib.MainLoop()
loop.run()
