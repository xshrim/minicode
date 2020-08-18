#!/usr/bin/sh

# from python:2.7.4-slim

jupyter kernelspec list

# update source
cat /dev/null > /etc/apt/sources.list
echo 'deb https://mirrors.ustc.edu.cn/debian/ buster main contrib non-free' >> /etc/apt/sources.list
echo 'deb https://mirrors.ustc.edu.cn/debian/ buster-updates main contrib non-free' >> /etc/apt/sources.list
echo 'deb https://mirrors.ustc.edu.cn/debian-security/ buster/updates main contrib non-free' >> /etc/apt/sources.list
apt update
apt install openssh-server openssh-client iputils-ping iproute2 netcat nmap wget curl xz-utils zip vim zsh -y
/etc/ssh/sshd_config 
PermitRootLogin  yes

# install jupyterlab
pip install jupyterlab
pip install --pre jupyter-lsp
pip install python-language-server[all]
jupyter kernelspec list

jupyter lab --allow-root --no-browser --NotebookApp.token='' --NotebookApp.notebook_dir="/workspace"  --ip=0.0.0.0 --port=8888
# jupyter lab --generate-config
# "set the c.NotebookApp.token parameter to an empty string"


# install
# use local server

# install octave

# install git gcc g++ lua php
apt install git gcc g++ lua5.3 php php-zmq -y

# install zeromq4
apt install libtool pkg-config build-essential autoconf automake uuid-dev libsodium-dev
wget https://github.com/zeromq/libzmq/releases/download/v4.3.2/zeromq-4.3.2.tar.gz
tar xvf zeromq-4.3.2.tar.gz
cd zeromq-4.3.2
# ./configure --without-libsodium
./configure
make install
ldconfig
ldconfig -p|grep zmq
cd ..
rm -rf zeromq-4.3.2 zeromq-4.3.2.tar.gz

# install nodejs
# apt install nodejs -y
wget https://nodejs.org/dist/v12.13.1/node-v12.13.1-linux-x64.tar.xz
tar xvf node-v12.13.1-linux-x64.tar.xz
mv node-v12.13.1-linux-x64 /usr/local/share/node
echo 'export PATH=/usr/local/share/node/bin:$PATH' >> /etc/profile
rm -rf node-v12.13.1-linux-x64.tar.xz

# install golang
wget https://studygolang.com/dl/golang/go1.13.4.linux-amd64.tar.gz
tar xvf go1.13.4.linux-amd64.tar.gz
mv go /usr/local/share/go
mkdir -p /opt/gopath
echo 'export GOPATH=/opt/gopath' >> /etc/profile
echo 'export PATH=$GOPATH/bin:$PATH' >> /etc/profile
echo 'export PATH=/usr/local/share/go/bin:$PATH' >> /etc/profile
rm -rf go1.13.4.linux-amd64.tar.gz

# install java
wget --no-check-certificate -c --no-cookies --header "Cookie: oraclelicense=accept-securebackup-cookie" https://download.oracle.com/otn/java/jdk/11.0.5+10/e51269e04165492b90fa15af5b4eb1a5/jdk-11.0.5_linux-x64_bin.tar.gz
tar xvf jdk-11.0.5_linux-x64_bin.tar.gz
mv jdk-11.0.5 /usr/local/share/jdk
echo 'export JAVA_HOME=/usr/local/share/jdk' >> /etc/profile
echo 'export PATH=$JAVA_HOME/bin:$PATH' >> /etc/profile
echo 'export CLASSPATH=.:$JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar' >> /etc/profile
# echo 'export PATH=/usr/local/share/jdk/bin:$PATH' >> /etc/profile
rm -rf jdk-11.0.5_linux-x64_bin.tar.gz


# vim kernel
git clone https://github.com/mattn/vim_kernel.git
cd vim_kernel
python setup.py install
python -m vim_kernel.install
cd ..

# bash kernel
pip install bash_kernel
python -m bash_kernel.install --prefix /usr/local/

# zsh kernel
pip install notebook zsh_jupyter_kernel
python3 -m zsh_jupyter_kernel.install --prefix /usr/local/

# octave kernel
pip install octave_kernel

# c kernel
pip install jupyter-c-kernel
install_c_kernel --prefix /usr/local/

# lua kernel
pip install ilua

# javascript kernel
npm install -g --unsafe-perm ijavascript
ijsinstall --install=global

# typescript kernel
npm install -g --unsafe-perm itypescript
its --install=global

# php kernel
curl -sS https://getcomposer.org/installer | php 
mv composer.phar /usr/bin/composer
wget https://litipk.github.io/Jupyter-PHP-Installer/dist/jupyter-php-installer.phar
chmod a+x jupyter-php-installer.phar
./jupyter-php-installer.phar install

# perl6
# apt install perl6
# zef install Jupyter::Kernel

# go kernel
go get -u github.com/gopherdata/gophernotes
mkdir /usr/local/share/jupyter/kernels/gophernotes
cp $GOPATH/src/github.com/gopherdata/gophernotes/kernel/* /usr/local/share/jupyter/kernels/gophernotes

# java kernel
mkdir ijava
cd ijava
wget https://github.com/SpencerPark/IJava/releases/download/v1.3.0/ijava-1.3.0.zip 
unzip ijava-1.3.0.zip
python install.py --prefix /usr/local/
cd ..


# jupyter kernelspec list
bash           /usr/local/share/jupyter/kernels/bash
c              /usr/local/share/jupyter/kernels/c
gophernotes    /usr/local/share/jupyter/kernels/gophernotes
java           /usr/local/share/jupyter/kernels/java
javascript     /usr/local/share/jupyter/kernels/javascript
jupyter-php    /usr/local/share/jupyter/kernels/jupyter-php
lua            /usr/local/share/jupyter/kernels/lua
octave         /usr/local/share/jupyter/kernels/octave
python3        /usr/local/share/jupyter/kernels/python3
typescript     /usr/local/share/jupyter/kernels/typescript
vim_kernel     /usr/local/share/jupyter/kernels/vim_kernel
zsh            /usr/local/share/jupyter/kernels/zsh


# jupyter labextension list
@hadim/jupyter-archive v0.5.4  enabled  OK
@jupyter-widgets/jupyterlab-manager v1.1.0  enabled  OK
@jupyter-widgets/midicontrols v0.1.2  enabled  OK
@jupyterlab/celltags v0.2.0  enabled  OK
@jupyterlab/commenting-extension v0.2.1  enabled  OK
@jupyterlab/dataregistry-extension v3.0.0  enabled  OK
@jupyterlab/geojson-extension v1.0.0  enabled  OK
@jupyterlab/git v0.8.2  enabled  OK
@jupyterlab/github v1.0.1  enabled  OK
@jupyterlab/google-drive v1.0.0  enabled  OK
@jupyterlab/jupyterlab-telemetry v0.2.0  enabled  OK
@jupyterlab/latex v1.0.0  enabled  OK
@jupyterlab/metadata-extension v2.0.0  enabled  OK
@jupyterlab/mp4-extension v0.1.0  enabled  OK
@jupyterlab/pullrequests v0.2.0  enabled  OK
@jupyterlab/shortcutui v0.4.0  enabled  OK
@jupyterlab/toc v2.0.0-rc.0  enabled  OK
@krassowski/jupyterlab-lsp v0.6.1  enabled  OK
@krassowski/jupyterlab_go_to_definition v0.7.1  enabled  OK
@telamonian/theme-darcula v1.1.3  enabled  OK
@voziq/custom_flbrowser v1.1.3  enabled  OK
jupyterlab-drawio v0.6.0  enabled  OK
jupyterlab-favorites v1.0.0  enabled  OK
jupyterlab-gitlab v0.3.0  enabled  OK
jupyterlab-jupytext v1.0.2  enabled  OK
jupyterlab-kernelspy v1.1.0  enabled  OK
jupyterlab-python-file v0.3.0  enabled  OK
jupyterlab-recents v1.0.1  enabled  OK
jupyterlab_templates v0.2.0  enabled  OK
jupyterlab_toastify v2.3.2  enabled  OK
jupyterlab_vim v0.11.0  enabled  OK
sophon-notebook-htmlviewer-extension v1.0.2  enabled  OK
verdant-log v1.0.0  enabled  OK





