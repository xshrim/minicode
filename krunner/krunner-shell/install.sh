#!/usr/bin/bash

echo "ps -faux  | grep -E 'krunner'"

sudo cp -v resources/org.manjaro.shell.krunner.service /usr/share/dbus-1/services/
#sudo cp resources/pamac-shell.service /usr/lib/systemd/user/
sudo cp -v  resources/plasma-runner-ushell.desktop /usr/share/kservices5/
sudo mkdir -p /usr/lib/qt/plugins/
sudo cp -v plugin/krunner_shell.py /usr/lib/qt/plugins/
sudo chmod +x /usr/lib/qt/plugins/krunner_shell.py

if [[ "$1" == "-l" ]]; then
    sudo rm /usr/share/dbus-1/services/org.manjaro.shell.krunner.service
    sudo rm /usr/lib/qt/plugins/krunner_shell.py
    #sudo rm /usr/lib/systemd/user/shell-krunner.service
    echo "python plugin/krunner_pamac.py # load dbus server in other terminal"
    #python plugin/krunner_shell.py &
fi

if [[ "$1" == "-r" ]]; then
    sudo rm -v /usr/share/dbus-1/services/org.manjaro.shell.krunner.service
    sudo rm -v /usr/share/kservices5/plasma-runner-ushell.desktop
    #sudo rm -v /usr/lib/systemd/user/shell-krunner.service
    sudo rm -v /usr/lib/qt/plugins/krunner_shell.py
    ps -aux | grep 'krunner'
fi

kquitapp5 krunner
#kstart5 krunner
qdbus org.kde.krunner /App org.kde.krunner.App.querySingleRunner usershell ls:/etc/pa
#qdbus --literal org.manjaro.shell.krunner /krunner org.kde.krunner1.Match "ls:/etc/pac"

