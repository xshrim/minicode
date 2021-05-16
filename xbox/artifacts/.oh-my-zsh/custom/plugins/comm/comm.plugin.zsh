# PWD
local cdir=$(dirname $0)
# zsh文档: https://github.com/goreliu/zshguide
########################################
# Git
########################################

# Default values for the appearance of the prompt. Configure at will.
ZSH_THEME_GIT_PROMPT_PREFIX="%{$fg_bold[blue]%}("
ZSH_THEME_GIT_PROMPT_SUFFIX="%{$fg_bold[blue]%})%{$reset_color%}"
#ZSH_THEME_GIT_PROMPT_SEPARATOR="ǀ"
#ZSH_THEME_GIT_PROMPT_BRANCH="%{$fg_bold[red]%}%{$reset_color%}"
ZSH_THEME_GIT_PROMPT_STAGED="%{$fg[yellow]%}ǀ%{$fg[cyan]%}⚑"
ZSH_THEME_GIT_PROMPT_CONFLICTS="%{$fg[yellow]%}ǀ%{$fg[red]%}⚡"
ZSH_THEME_GIT_PROMPT_ADDED="%{$fg[yellow]%}ǀ%{$fg[magenta]%}✚"
ZSH_THEME_GIT_PROMPT_DELETED="%{$fg[yellow]%}ǀ%{$fg[red]%}✖"
ZSH_THEME_GIT_PROMPT_CHANGED="%{$fg[yellow]%}ǀ%{$fg[yellow]%}✱"
ZSH_THEME_GIT_PROMPT_MODIFIED="%{$fg[yellow]%}ǀ%{$fg[yellow]%}✤"
ZSH_THEME_GIT_PROMPT_RENAMED="%{$fg[yellow]%}ǀ%{$fg[yellow]%}⬙"
ZSH_THEME_GIT_PROMPT_BEHIND="↓"
ZSH_THEME_GIT_PROMPT_AHEAD="↑"
ZSH_THEME_GIT_PROMPT_CLEAN="%{$fg[yellow]%}ǀ%{$fg_bold[green]%}✔"
ZSH_THEME_GIT_PROMPT_DIRTY="%{$fg[yellow]%}ǀ%{$fg_bold[red]%}✘"
ZSH_THEME_GIT_PROMPT_UNTRACKED="%{$fg[yellow]%}ǀ%{$fg[white]%}●"
ZSH_THEME_GIT_PROMPT_UNMERGED="%{$fg[yellow]%}ǀ%{$fg[red]%}✦"

########################################


########################################
# Powerline
#
# vim:ft=zsh ts=2 sw=2 sts=2
#
# ZYSzys's Theme - https://github.com/ZYSzys/zys-zsh-theme
# A Powerline-inspired theme for ZSH
#
# # README
#
# In order for this theme to render correctly, you will need a
# [Powerline-patched font](https://github.com/Lokaltog/powerline-fonts).
# Make sure you have a recent version: the code points that Powerline
# uses changed in 2012, and older versions will display incorrectly,
# in confusing ways.
#
# In addition, I recommend the
# [Solarized theme](https://github.com/altercation/solarized/) and, if you're
# using it on Mac OS X, [iTerm 2](http://www.iterm2.com/) over Terminal.app -
# it has significantly better color fidelity.
#
# # Goals
#
# The aim of this theme is to only show you *relevant* information. Like most
# prompts, it will only show git information when in a git working directory.
# However, it goes a step further: everything from the current user and
# hostname to whether the last call exited with an error to whether background
# jobs are running in this shell will all be displayed automatically when
# appropriate.

### Segment drawing
# A few utility functions to make it easy and re-usable to draw segmented prompts
###########################################

CURRENT_BG='NONE'
RCURRENT_BG='NONE'

# Special Powerline characters

() {
  local LC_ALL="" LC_CTYPE="en_US.UTF-8"
  # NOTE: This segment separator character is correct.  In 2012, Powerline changed
  # the code points they use for their special characters. This is the new code point.
  # If this is not working for you, you probably have an old version of the
  # Powerline-patched fonts installed. Download and install the new version.
  # Do not submit PRs to change this unless you have reviewed the Powerline code point
  # history and have new information.
  # This is defined using a Unicode escape sequence so it is unambiguously readable, regardless of
  # what font the user is viewing this source code in. Do not replace the
  # escape sequence with a single literal character.
  # Do not change this! Do not make it '\u2b80'; that is the old, wrong code point.
  SEGMENT_SEPARATOR=$'\ue0b0'
  RSEGMENT_SEPARATOR=$'\ue0b2'
}

# Begin a segment
# Takes two arguments, background and foreground. Both can be omitted,
# rendering default background/foreground.
prompt_segment() {
  local bg fg
  [[ -n $1 ]] && bg="%K{$1}" || bg="%k"
  [[ -n $2 ]] && fg="%F{$2}" || fg="%f"
  if [[ $CURRENT_BG != 'NONE' && $1 != $CURRENT_BG ]]; then
    echo -n " %{$bg%F{$CURRENT_BG}%}$SEGMENT_SEPARATOR%{$fg%} "
  else
    echo -n "%{$bg%}%{$fg%} "
  fi
  CURRENT_BG=$1
  [[ -n $3 ]] && echo -n $3
}

rprompt_segment() {
  local bg fg
  [[ -n $1 ]] && bg="%K{$1}" || bg="%k"
  [[ -n $2 ]] && fg="%F{$2}" || fg="%f"

  if [[ $RCURRENT_BG != 'NONE' && $1 != $RCURRENT_BG ]]; then
    echo -n " %{%K{$RCURRENT_BG}%F{$1}%}$RSEGMENT_SEPARATOR%{$fg%}"
    echo -n "%{$bg%}%{$fg%} "
  else
    echo -n "%{$bg%}%{$fg%} "
  fi
  RCURRENT_BG=$1
  [[ -n $3 ]] && echo -n $3
}

prompt_start() {
  if [[ $RCURRENT_BG == 'NONE' || $1 != $RCURRENT_BG ]]; then
    RCURRENT_BG=$1
    echo -n " %{%k%F{$RCURRENT_BG}%}$RSEGMENT_SEPARATOR"
  else
    echo -n "%{%k%}"
  fi
  #echo -n "%{%f%}"
}

# End the prompt, closing any open segments
prompt_end() {
  if [[ -n $CURRENT_BG ]]; then
    echo -n " %{%k%F{$CURRENT_BG}%}$SEGMENT_SEPARATOR"
  else
    echo -n "%{%k%}"
  fi
  echo -n "%{%f%}"
  CURRENT_BG=''
}

### Prompt components
# Each component will draw itself, and hide itself if no information needs to be shown

# Context: user@hostname (who am I and where am I)
prompt_context() {
  # if [[ "$USER" != "$DEFAULT_USER" || -n "$SSH_CLIENT" ]]; then
  #   prompt_segment black default "%(!.%{%F{yellow}%}.)$USER@%m"
  # fi
}

# Git: branch/detached head, dirty status
prompt_git() {
  (( $+commands[git] )) || return
  local PL_BRANCH_CHAR
  () {
    local LC_ALL="" LC_CTYPE="en_US.UTF-8"
    PL_BRANCH_CHAR=$'\ue0a0'         # 
  }
  local ref dirty mode repo_path
  repo_path=$(git rev-parse --git-dir 2>/dev/null)

  if $(git rev-parse --is-inside-work-tree >/dev/null 2>&1); then
    dirty=$(parse_git_dirty)
    ref=$(git symbolic-ref HEAD 2> /dev/null) || ref="➦ $(git rev-parse --short HEAD 2> /dev/null)"
    if [[ -n $dirty ]]; then
      prompt_segment yellow black
    else
      prompt_segment green black
    fi

    if [[ -e "${repo_path}/BISECT_LOG" ]]; then
      mode=" <B>"
    elif [[ -e "${repo_path}/MERGE_HEAD" ]]; then
      mode=" >M<"
    elif [[ -e "${repo_path}/rebase" || -e "${repo_path}/rebase-apply" || -e "${repo_path}/rebase-merge" || -e "${repo_path}/../.dotest" ]]; then
      mode=" >R>"
    fi

    setopt promptsubst
    autoload -Uz vcs_info

    zstyle ':vcs_info:*' enable git
    zstyle ':vcs_info:*' get-revision true
    zstyle ':vcs_info:*' check-for-changes true
    zstyle ':vcs_info:*' stagedstr '✚'
    zstyle ':vcs_info:*' unstagedstr '●'
    zstyle ':vcs_info:*' formats ' %u%c'
    zstyle ':vcs_info:*' actionformats ' %u%c'
    vcs_info
    echo -n "${ref/refs\/heads\//$PL_BRANCH_CHAR }${vcs_info_msg_0_%% }${mode}"
  fi
}

prompt_bzr() {
    (( $+commands[bzr] )) || return
    if (bzr status >/dev/null 2>&1); then
        status_mod=`bzr status | head -n1 | grep "modified" | wc -m`
        status_all=`bzr status | head -n1 | wc -m`
        revision=`bzr log | head -n2 | tail -n1 | sed 's/^revno: //'`
        if [[ $status_mod -gt 0 ]] ; then
            prompt_segment yellow black
            echo -n "bzr@"$revision "✚ "
        else
            if [[ $status_all -gt 0 ]] ; then
                prompt_segment yellow black
                echo -n "bzr@"$revision

            else
                prompt_segment green black
                echo -n "bzr@"$revision
            fi
        fi
    fi
}

prompt_hg() {
  (( $+commands[hg] )) || return
  local rev status
  if $(hg id >/dev/null 2>&1); then
    if $(hg prompt >/dev/null 2>&1); then
      if [[ $(hg prompt "{status|unknown}") = "?" ]]; then
        # if files are not added
        prompt_segment red white
        st='±'
      elif [[ -n $(hg prompt "{status|modified}") ]]; then
        # if any modification
        prompt_segment yellow black
        st='±'
      else
        # if working copy is clean
        prompt_segment green black
      fi
      echo -n $(hg prompt "☿ {rev}@{branch}") $st
    else
      st=""
      rev=$(hg id -n 2>/dev/null | sed 's/[^-0-9]//g')
      branch=$(hg id -b 2>/dev/null)
      if `hg st | grep -q "^\?"`; then
        prompt_segment red black
        st='±'
      elif `hg st | grep -q "^[MA]"`; then
        prompt_segment yellow black
        st='±'
      else
        prompt_segment green black
      fi
      echo -n "☿ $rev@$branch" $st
    fi
  fi
}

# Dir: current working directory
prompt_dir() {
  prompt_segment blue black '%~'
}

# Virtualenv: current working virtualenv
prompt_virtualenv() {
  local virtualenv_path="$VIRTUAL_ENV"
  if [[ -n $virtualenv_path && -n $VIRTUAL_ENV_DISABLE_PROMPT ]]; then
    prompt_segment blue black "(`basename $virtualenv_path`)"
  fi
}

# Status:
# - was there an error
# - am I root
# - are there background jobs?
prompt_status() {
  local symbols
  symbols=()
  [[ $RETVAL -ne 0 ]] && symbols+="%{%F{red}%}✘"
  [[ $UID -eq 0 ]] && symbols+="%{%F{yellow}%}⚡"
  [[ $(jobs -l | wc -l) -gt 0 ]] && symbols+="%{%F{cyan}%}⚙"

  [[ -n "$symbols" ]] && prompt_segment black default "$symbols"
}

# Date:
# Time:
rprompt_dtime() {
  rprompt_segment blue black '%D %T'
}

# status
rprompt_status() {
rprompt_segment yellow black '%?'
}

## Main prompt
build_prompt() {
  RETVAL=$?
  prompt_status
  prompt_virtualenv
  prompt_context
  prompt_dir
  prompt_git
  prompt_bzr
  prompt_hg
  prompt_end
}

build_rprompt() {
    prompt_start blue
    rprompt_dtime
    rprompt_status
}

########################################
# Settings
########################################

#color{{{
autoload colors zsh/terminfo
if [[ "$terminfo[colors]" -ge 8 ]]; then
colors
fi
for color in RED GREEN YELLOW BLUE MAGENTA CYAN WHITE; do
eval _$color='%{$terminfo[bold]$fg[${(L)color}]%}'
eval $color='%{$fg[${(L)color}]%}'
(( count = $count + 1 ))
done
FINISH="%{$terminfo[sgr0]%}"
#}}}

# symbols to choose from:
# ☀ ✹ ☄ ♆ ♀ ♁ ♐ ♇ ♈ ♉ ✹ ♚ ♛ ♜ ♝ ♞ ♟ ♠ ♣ ▾⚢ ⚲ ⚳ ⚴ ⚥ ☿ ⚤ ⚦ ⚒ ⚑ ⚐ ♺ ♻ ♼ ⓐ ⑃ ᐅ ☰ ☱ ☲ ☳ ☴ ☵ ☶ ☷ 
# ✡ ✔ ✘ ✗ ✖ ✚ ✱ ✤ ✦ ● ❤ ➜ ➟ ➼ ➤ ☁ ◀▶▲✂ ✎ ✐ 𝝙 ♒ ⑊ # ⑄ # ⓧ ⑂⨍ ⨎ ⨏ ⨷ ◯ ⩚ ⩛ ⩡ ⩱ ⩲ ± ⩵  ⩶ ⨠ ❮ » ‹ ⏎
# ⬅ ⬆ ⬇ ⬈ ⬉ ⬊ ⬋ ⬒ ⬓ ⬔ ⬕ ⬖ ⬗ ⬘ ⬙ ⬟  ⬤ ☀ ☂ ✭ ⚡ * Ⓞ ⓣ Ⓓ ⓜ ⓤ ⓡ ⑁ 〒 ǀ ǁ ǂ ĭ Ť Ŧ ↵ ⌚ ➦  ⚙  ✖ ⨀ ⨁ ⨂ ☛ ☚ ☎
# ★ ☯ ☮ ☣ ☢ ❥ ✧ ✪ ♬ ♫ ♪ ♩ ✉ ✯ ✮ ✿ ❀ ❁ ❂ ❉ ❄ ❅ ❆ ✝ ✓ ☑ ☐ ღ ⁂ ☽☾ϟ✢✠✝✛✜☩☨☥✁☊✒✑✏✎✐⌛✇✈

# ✄⏎⇧⇪⌂⌘⌫♈♉♋♌☺☹⚢✶⌑➢➣➥➦➪➩➨↰↱↲↳↴↶↷⇕⇖⇗⇘⇙⇚⇛⇧⇨⇦⇩⇪➔➘➙➚➛➜➝➞⏎➟➠➡☇☊☋⇍⇎⇏⇕➳➷➸➹➺➻
# ⤴⤵↵↔◆◇◢◣◤◥◊❖⎔⊞⊿◁◃◂◄▲△▴▼▾▿◮◭◐◑◒⦿◕◔⊖⊘『』
# ﹃﹄【】︼︻︽︾︿﹀﹃﹄╽╾╼╿╏╍╌┅┄┆┇┈┉⋓⋒⋑⋐┊┋┥┤┴┵┶┧┷┸┨┲┱┰┠┟┞┝├┬┫├┍┎┏┐┑┒┓└└┕┖┗┘┙┚┛

# http://zsh.sourceforge.net/Doc/Release/Prompt-Expansion.html
#命令提示符 {{{
precmd () {
#local gitprompt=$(git_super_status)
local gitprompt=$(git_prompt_info)
local zero='%([BSUbfksu]|([FK]|){*})'
local gitpromptsize=${#${(S%%)gitprompt//$~zero/}}

local smile="%(?,$GREEN☀%{$reset_color%},$RED☂%{$reset_color%})"
local count_db_wth_char=${#${${(%):-%/}//[[:ascii:]]/}}
HBAR=" -"

local theme="default"

if [[ $theme == "complex" ]]; then
    local leftsize=${#${(%):-♆<%M %/}}+$gitpromptsize+$count_db_wth_char
    local rightsize=${#${(%):-%D %T>♆}}

    FILLBAR="\${(l.(($COLUMNS - ($leftsize + $rightsize +2)))..${HBAR}.)}"

    #RPROMPT=$(echo "%(?..$RED%?$FINISH)")
    PROMPT=$(echo "$_BLUE♆<$_CYAN%M $_GREEN%/${gitprompt} $_YELLOW${(e)FILLBAR} $_MAGENTA%D %T$_BLUE>♆$FINISH
$fg_bold[yellow][ $MAGENTA%n $BLUE%h ${smile}$fg_bold[yellow] ] $_RED%(!.➤.⨠) $FINISH")
elif [[ $theme == "simple" ]]; then
    local leftsize=${#${(%):-♆<%M✦%n %/}}+$gitpromptsize+$count_db_wth_char
    local rightsize=${#${(%):-%D %T %h>♆}}+2

    FILLBAR="\${(l.(($COLUMNS - ($leftsize + $rightsize +2)))..${HBAR}.)}"

    #RPROMPT=$(echo "%(?..$RED%?$FINISH)")
    PROMPT=$(echo "$_BLUE♆<$_CYAN%M$YELLOW✦$MAGENTA%n $_GREEN%/${gitprompt} $_YELLOW${(e)FILLBAR} ${smile} $_MAGENTA%D %T %h$_BLUE>♆$FINISH
$_RED%(!.➤.⨠) $FINISH")
elif [[ $theme == "classic" ]]; then
    PROMPT="%{${fg_bold[blue]}%}✿ %{${fg_bold[red]}%}%m%(!.☢.❖)%n %{${fg_bold[magenta]}%}:: %{${fg_bold[yellow]}%}%~%{${fg_bold[cyan]}%}$(git_prompt_info) %(?,%{${fg_bold[green]}%}☀%{$reset_color%},%{${fg_bold[red]}%}☂%{$reset_color%}) %{${fg_bold[blue]}%}»%{${reset_color}%} "
elif [[ $theme == "power" ]]; then
    PROMPT='%{%f%b%k%}$(build_prompt) '
    #RPROMPT='%{%F{yellow}%}[%*]'
    RPROMPT='%{%f%b%k%}$(build_rprompt) '
    #echo $(build_prompt)
    #echo $(build_rprompt)
    # TODO 继续完善
elif [[ $theme == "minimal" ]]; then
    PROMPT='%(?,%{${fg_bold[green]}%}☀%{$reset_color%},%{${fg_bold[red]}%}☂%{$reset_color%}) %{$fg_bold[yellow]%}%~%{$reset_color%}$(git_prompt_info) %{${fg_bold[magenta]}%}%(!.➤.⨠)%{${reset_color}%} '
else
    PROMPT="%{${fg_bold[blue]}%}✿ %{${fg_bold[red]}%}%m%(!.☢.❖)%n %{${fg_bold[magenta]}%}:: %{${fg_bold[cyan]}%}%~%{${fg_bold[cyan]}%}$(git_prompt_info) %(?,%{${fg_bold[green]}%}☀%{$reset_color%},%{${fg_bold[red]}%}☂%{$reset_color%}) %{${fg_bold[blue]}%}»%{${reset_color}%} "
    RPROMPT='%{${fg_bold[yellow]}%}< %w %T %! > %{${fg_bold[blue]}%}✿%{${reset_color}%}'
    # TODO 继续完善
fi

#在 Emacs终端 中使用 Zsh 的一些设置
if [[ "$TERM" == "dumb" ]]; then
setopt No_zle
PROMPT='%n@%M %/
>>'
fi
}

# 终端开启256色支持
if [ -e /usr/share/terminfo/x/xterm-256color ]; then
  export TERM='xterm-256color'
fi

#}}}

#清空历史记录
#cat /dev/null > ${HOME}/.zhistory
#cat /dev/null > ${HOME}/.zsh_history

#标题栏、任务栏样式{{{
#case $TERM in (*xterm*|*rxvt*|(dt|k|E)term)
# preexec () { print -Pn "\e]0;%n@%M//%/\ $1\a" }
#;;
#esac
#}}}

#关于历史纪录的配置 {{{
#历史纪录条目数量
export HISTSIZE=100000
#注销后保存的历史纪录条目数量
export SAVEHIST=100000
#历史纪录文件
export HISTFILE=~/.zhistory
#以附加的方式写入历史纪录
setopt INC_APPEND_HISTORY
#如果连续输入的命令相同，历史纪录中只保留一个
setopt HIST_IGNORE_DUPS
#为历史纪录中的命令添加时间戳
setopt EXTENDED_HISTORY

#启用 cd 命令的历史纪录，cd -[TAB]进入历史路径
setopt AUTO_PUSHD
#相同的历史路径只保留一个
setopt PUSHD_IGNORE_DUPS

#在命令前添加空格，不将此命令添加到纪录文件中
setopt HIST_IGNORE_SPACE
#}}}

#每个目录使用独立的历史纪录{{{
#cd() {
#builtin cd "$@"                             # do actual cd
#fc -W                                       # write current history  file
#local HISTDIR="$HOME/.zsh_history$PWD"      # use nested folders for history
#if  [ ! -d "$HISTDIR" ] ; then          # create folder if needed
#mkdir -p "$HISTDIR"
#fi
#export HISTFILE="$HISTDIR/zhistory"     # set new history file
#touch $HISTFILE
#local ohistsize=$HISTSIZE
#HISTSIZE=0                              # Discard previous dir's history
#HISTSIZE=$ohistsize                     # Prepare for new dir's history
#fc -R                                       #read from current histfile
#}
#mkdir -p $HOME/.zsh_history$PWD
#export HISTFILE="$HOME/.zsh_history$PWD/zhistory"
 
#function allhistory { cat $(find $HOME/.zsh_history -name zhistory) }
#function convhistory {
#sort $1 | uniq |
#sed 's/^:\([ 0-9]*\):[0-9]*;\(.*\)/\1::::::\2/' |
#awk -F"::::::" '{ $1=strftime("%Y-%m-%d %T",$1) "|"; print }'
#}
#使用 histall 命令查看全部历史纪录
#function histall { convhistory =(allhistory) |
#sed '/^.\{20\} *cd/i\\' }
#使用 hist 查看当前目录历史纪录
#function hist { convhistory $HISTFILE }
 
#全部历史纪录 top50
#function top50 { allhistory | awk -F':[ 0-9]*:[0-9]*;' '{ $1="" ; print }' | sed 's/ /\n/g' | sed '/^$/d' | sort | uniq -c | sort -nr | head -n 50 }
 
#}}}

#杂项 {{{
#允许在交互模式中使用注释 例如：
#cmd #这是注释
setopt INTERACTIVE_COMMENTS

#启用自动 cd，输入目录名回车进入目录
#稍微有点混乱，不如 cd 补全实用
#setopt AUTO_CD

#禁用自动解析通配符
setopt no_nomatch

#扩展路径
#/v/c/p/p => /var/cache/pacman/pkg
setopt complete_in_word

#禁用 core dumps
#limit coredumpsize 0

#Emacs风格 键绑定
bindkey -e
#设置 [DEL]键 为向后删除
bindkey "\e[3~" delete-char

#以下字符视为单词的一部分
WORDCHARS='*?_-[]~=&;!#$%^(){}<>'
#}}}

#自动补全功能 {{{
setopt AUTO_LIST
setopt AUTO_MENU
#开启此选项，补全时会直接选中菜单项
#setopt MENU_COMPLETE

#自动补全缓存
zstyle ':completion::complete:*' use-cache on
zstyle ':completion::complete:*' cache-path .zcache
zstyle ':completion:*:cd:*' ignore-parents parent pwd

#自动补全选项
zstyle ':completion:*' verbose yes
zstyle ':completion:*' menu select
zstyle ':completion:*:*:default' force-list always
zstyle ':completion:*' select-prompt '%SSelect: lines: %L matches: %M [%p]'

zstyle ':completion:*:match:*' original only
zstyle ':completion::prefix-1:*' completer _complete
zstyle ':completion:predict:*' completer _complete
zstyle ':completion:incremental:*' completer _complete _correct
zstyle ':completion:*' completer _complete _prefix _correct _prefix _match _approximate

#路径补全
zstyle ':completion:*' expand 'yes'
zstyle ':completion:*' squeeze-slashes 'yes'
zstyle ':completion::complete:*' '\\'


#修正大小写
zstyle ':completion:*' matcher-list '' 'm:{a-zA-Z}={A-Za-z}'
#错误校正   坑爹功能
#zstyle ':completion:*' completer _complete _match _approximate
#zstyle ':completion:*:match:*' original only
#zstyle ':completion:*:approximate:*' max-errors 1 numeric

#kill 命令补全
compdef pkill=killall
zstyle ':completion:*:*:kill:*' menu yes select
zstyle ':completion:*:*:*:*:processes' force-list always
zstyle ':completion:*:processes' command 'ps -au$USER'

#补全类型提示分组
zstyle ':completion:*:matches' group 'yes'
zstyle ':completion:*' group-name ''
zstyle ':completion:*:options' description 'yes'
zstyle ':completion:*:options' auto-description '%d'
zstyle ':completion:*:descriptions' format $'\e[01;33m -- %d --\e[0m'
zstyle ':completion:*:messages' format $'\e[01;35m -- %d --\e[0m'
zstyle ':completion:*:warnings' format $'\e[01;31m -- No Matches Found --\e[0m'
zstyle ':completion:*:corrections' format $'\e[01;32m -- %d (errors: %e) --\e[0m'

# cd ~ 补全顺序
zstyle ':completion:*:-tilde-:*' group-order 'named-directories' 'path-directories' 'users' 'expand'
#}}}

##行编辑高亮模式 {{{
# Ctrl+@ 设置标记，标记和光标点之间为 region
zle_highlight=(region:bg=magenta #选中区域
special:bold #特殊字符
isearch:underline) #搜索时使用的关键字
#}}}

##空行(光标在行首)补全 "cd " {{{
user-complete(){
case $BUFFER in
"" ) # 空行填入 "cd "
BUFFER="cd "
zle end-of-line
zle expand-or-complete
;;
"cd " ) # TAB + 空格 替换为 "cd ~"
BUFFER="cd ~"
zle end-of-line
zle expand-or-complete
;;
" " )
BUFFER="!?"
zle end-of-line
;;
"cd --" ) # "cd --" 替换为 "cd +"
BUFFER="cd +"
zle end-of-line
zle expand-or-complete
;;
"cd +-" ) # "cd +-" 替换为 "cd -"
BUFFER="cd -"
zle end-of-line
zle expand-or-complete
;;
* )
zle expand-or-complete
;;
esac
}
zle -N user-complete
bindkey "\t" user-complete

#显示 path-directories ，避免候选项唯一时直接选中
cdpath="/home"
#}}}

##在命令前插入 sudo {{{
#定义功能
sudo-command-line() {
[[ -z $BUFFER ]] && zle up-history
[[ $BUFFER != sudo\ * ]] && BUFFER="sudo $BUFFER"
zle end-of-line #光标移动到行末
}
zle -N sudo-command-line
#定义快捷键为： [Esc] [Esc]
bindkey "\e\e" sudo-command-line
#}}}

checkcommand() {
  if type $1 &>/dev/null; then
    return 0
  else
    return 1
  fi
}

# 默认编辑器
export EDITOR=vim

# 本地二进制路径
export PATH=$cdir/bin:$PATH

#命令别名 {{{
if [ -f ~/.vimrc ];then
  alias vim='vim -u ~/.vimrc'
fi
alias vi='vim'

# if [ $(id -u) -ne 0 ] && type fping -h &>/dev/null && ! fping -h &>/dev/null 
# then
#   sudo chown root:root $(which fping)
#   sudo chmod ug+s $(which fping)
# fi

# type fping &>/dev/null && alias ping='fping -e'

# type pycp &>/dev/null && alias cp='pycp -g'
# type pymv &>/dev/null && alias mv='pymv -g'

type acp &>/dev/null && alias cp='acp -gR'

type amv &>/dev/null && alias mv='amv -g'

type htop &>/dev/null && alias top='htop'

type bat &>/dev/null && alias cat='bat --paging never -p'

if type highlight &>/dev/null;then
  alias highlight="highlight -D $cdir/highlight"
  ! highlight -h &> /dev/null && type highlight.low &>/dev/null && alias highlight="highlight.low -D $cdir/highlight --add-config-dir=$cdir/highlight"
  alias cat='h'
  # alias highlight="highlight -D $cdir/highlight --config-file=$cdir/highlight/filetypes.conf"
fi

alias ls='ls -F --color=auto'
type exa &>/dev/null && exa --help &>/dev/null && alias ls='exa -F --color=auto'
alias ll='ls -l'
alias la='ls -a'
alias l='ls'
alias cls='clear'
alias mkdir='mkdir -p'
alias ps='ps -elf'
alias ss='ss -asntup'
alias open='xdg-open'
alias grep='grep --color=auto'
alias sim='sudo -E vim'
alias glog='git log --oneline --graph --decorate'
alias glg='git log --graph --pretty=format:"%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset" --abbrev-commit --date=relative'
alias rs='rsync -ahr --info=progress2 --no-i-r'
alias gm='gomacro'
alias py='python3'
# docker alias
alias dkc='docker-compose'
alias dps='docker container ls'
alias drm='docker rm'
alias drmi='docker rmi'
alias dimg='docker image ls'
alias dnet='docker network ls'
alias dvol='docker volume ls'
alias drun='docker run'
alias dexe='docker exec'
alias dtmp='docker run -it --rm'
alias dcmd='docker inspect --format="{{json .Config.Cmd}}"'
alias denv='docker inspect --format="{{json .Config.Env}}"'
alias dmnt='docker inspect --format="{{json .Mounts}}"'
alias daddr='docker inspect --format="{{json .NetworkSettings.Networks}}"'
alias dport='docker inspect --format="{{json .NetworkSettings.Ports}}"'
alias dpull='docker pull'
alias dpush='docker push'
#alias md='mkdir'
#alias jobs='jobs -l'
#alias cp='rsync -avP --progress'
#alias mv='rsync -avP --progress --remove-source-files'
#alias ping='fping -e'
#alias sudo='sudo env PATH=$PATH'
#alias history='history -fi'
# rsync alias
alias rs='rsync -ahr --info=progress2 --no-i-r'
alias rcp="rs"
function rmv() {
  rsync -ahr --info=progress2 --no-i-r --remove-source-files "$@" && rm -rf "${@:1:${#}-1}"
}

#[Esc][h] man 当前命令时，显示简短说明
alias run-help >&/dev/null && unalias run-help
autoload run-help

#历史命令 top10
alias top10='print -l ${(o)history%% *} | uniq -c | sort -nr | head -n 10'
#}}}

#路径别名 {{{
#进入相应的路径时只要 cd ~xxx
hash -d HIST="$HISTDIR"
#}}}

#{{{自定义补全
#补全 ping
zstyle ':completion:*:ping:*' hosts g.cn facebook.com

# def pacman-color completion as pacman
# compdef pacman-color=pacman
# }}}

#{{{ F1 计算器
arith-eval-echo() {
LBUFFER="${LBUFFER}echo \$(( "
RBUFFER=" ))$RBUFFER"
}
zle -N arith-eval-echo
bindkey "^[[11~" arith-eval-echo
#}}}

####{{{
#function timeconv { date -d @$1 +"%Y-%m-%d %T" }

# }}}

# pygments 高亮显示输出
function pat() {
    local style="monokai"
    if [ $# -eq 0 ]; then
        pygmentize -P style=$style -P tabsize=4 -f terminal256 -g
    else
        for NAME in $@; do
            pygmentize -P style=$style -P tabsize=4 -f terminal256 -g "$NAME"
        done
    fi
}
# alias cat='pat'

# 输出的指定字符串高亮
function hl() {
	declare -A fg_color_map
	fg_color_map[black]=30
	fg_color_map[red]=31
	fg_color_map[green]=32
	fg_color_map[yellow]=33
	fg_color_map[blue]=34
	fg_color_map[magenta]=35
	fg_color_map[cyan]=36
	 
	fg_c=$(echo -e "\e[1;${fg_color_map[$1]}m")
	c_rs=$'\e[0m'
	sed -u s"/$2/$fg_c\0$c_rs/g"
}

# highlight 高亮显示输出
# alias cat='h'
function h() {
  #local cmds="*ls* *echo* *cat* *ps* *who* *pwd* *arch* *lspci* *lsusb* *time* *date* *cal* *tree* *iconv* *env* *set* *df* *dfc* *tree* *free* *alias* *rpm* *yum* *dnf* *apt* *apt-get* *pacman* *sar* *lsattr* *h* *tar* *sort* *grep* *egrep* *recode* *pvs* *pvdisplay* *pvscan* *vgscan* *lvs* *lvdisplay* *vgs* *vgdisplay* *pvcreate* *pvremove* *lvcreate* *lvremove* *lvextend* *lvresize* *lvreduce* *lvrename* *lvconvert* *lvscan* *lvchange* *vgchange* *vgcreate* *vgreduce* *vgextend**vgrename* *ip* *ping* *ifconfig* *route* *hostname* *iwlist* *host* *npm* *nslookup* *whois* *dig* *dmesg* *rsync* *cp* *mv* *glog* *git* *mtr* *ss* *fd* *ag* *bat* *pycp* *pymv* *highlight* *xargs* *more* *less* *wd* *dict* *exa* *sysctl* *nmap* *convert* *ethtool* *smbclient* *mount* *umount* *tar* *zip* *unzip* *firewall-cmd* *iptables* *gs* *systemctl* *journalctl* *hostnamectl* *dmidecode* *hdparm* *which* *whereis* *locate* *jobs* *uname* *head* *tail* *cut* *tr* *sed* *awk* *file* *lsof* *netstat* *vmstat* *iostat* *docker* *docker-compose* *kubectl* *helm* *redis-cli* *mc* *mysql* *pgsql* *node* *python* *python3* *java* *go* *pip* *pip3* *make* *gcc* *tsc* *perl* *lua* *ruby* *rust* *scala* *julia* "


  if [ ! -t 1 ]
  then
    \cat $@
  else
    #local sn=()
    declare -a sn
    while [ $# -gt 0 ] && [[ "$1" == -* ]]
    do
        sn+=("$1")
        shift
    done
    if [ $# -eq 0 ]
    then
        #CAT="highlight -O xterm256 -t 4 -s bipolar -S sh"
        #highlight -O xterm256 -t 4 -s $style -S $syntax
        if type highlight &>/dev/null;then
            highlight -S sh ${sn[*]} -O xterm256 -t 4 -s bipolar
        elif type bat &>/dev/null;then
            bat ${sn[*]} --paging never -p
        elif type ccat &>/dev/null;then
            ccat $${sn[*]}
        else
            \cat ${sn[*]}
        fi
    else
        if [[ $(file "$1" | grep -o "text" | wc -l ) -lt 1 ]];then
            echo -en "\033[1m"
            \cat ${sn[*]} $@
            echo -en "\033[0m"
        else
        if type highlight &>/dev/null;then
            highlight ${sn[*]} -O xterm256 -t 4 -s bipolar $@ 2> /dev/null || highlight -S sh ${sn[*]} -O xterm256 -t 4 -s bipolar $@ 2> /dev/null || \cat ${sn[*]} $@
        elif type bat &>/dev/null;then
            bat ${sn[*]} --paging never -p $@ 2> /dev/null || \cat ${sn[*]} $@
        elif type ccat &>/dev/null;then
            ccat ${sn[*]} $@ 2> /dev/null || \cat ${sn[*]} $@
        else
            echo -en "\033[1m"
            \cat ${sn[*]} $@
            echo -en "\033[0m"
        fi
    fi
  fi
#        if [[ "$cmds" =~ "*"$1"*" ]]
#         then
#             #$* | highlight -O xterm256 -t 4 -s $style 2> /dev/null || $* | highlight -O xterm256 -t 4 -s $style -S sh
#             $* | highlight -O xterm256 -t 4 -s $style -S $syntax
#         else
#             highlight -O xterm256 -t 4 -s $style $@ 2> /dev/null || highlight -O xterm256 -t 4 -s $style -S $syntax $@
#         fi
  #highlight -O xterm256 -t 4 -s $style $@ 2> /dev/null || highlight -O xterm256 -t 4 -s $style -S $syntax $@ 2> /dev/null || bat -p $@ 2> /dev/null || cat $@
  fi
}

command_not_found_handler () {      #if the command is not found, let bash show the message and advice(zsh could only show "command not found").
        runcnf=1 
        retval=127 
        [ ! -S /var/run/dbus/system_bus_socket ] && runcnf=0 
        [ ! -x /usr/libexec/packagekitd ] && runcnf=0 
        if [ $runcnf -eq 1 ]
        then
                /usr/libexec/pk-command-not-found $@
                retval=$? 
        fi
        return 0                   #if return $retval, both the bash and zsh messages will be shown in the terminal, if return 0, only the bash massage will be shown.
}

####{{{
#function command_not_found_handler() {
#python /usr/lib/command-not-found $1
#return 0
#}
# }}}

#zmodload zsh/mathfunc
#autoload -U zsh-mime-setup
#zsh-mime-setup
#setopt EXTENDED_GLOB
#autoload -U promptinit
#promptinit
#prompt redhat
 
#setopt correctall
#autoload compinstall

## END OF FILE #################################################################
# vim:filetype=zsh foldmethod=marker autoindent expandtab shiftwidth=4

#autojump
if [ $commands[autojump] ]; then # check if autojump is installed
  if [ -f $HOME/.autojump/etc/profile.d/autojump.zsh ]; then # manual user-local installation
    . $HOME/.autojump/etc/profile.d/autojump.zsh
  elif [ -f /usr/share/autojump/autojump.zsh ]; then # debian and ubuntu package
    . /usr/share/autojump/autojump.zsh
  elif [ -f /etc/profile.d/autojump.zsh ]; then # manual installation
    . /etc/profile.d/autojump.zsh
  elif [ -f /etc/profile.d/autojump.sh ]; then # gentoo installation
    . /etc/profile.d/autojump.sh
  elif [ -f /usr/local/share/autojump/autojump.zsh ]; then # freebsd installation
    . /usr/local/share/autojump/autojump.zsh
  elif [ -f /opt/local/etc/profile.d/autojump.zsh ]; then # mac os x with ports
    . /opt/local/etc/profile.d/autojump.zsh
  elif [ $commands[brew] -a -f `brew --prefix`/etc/autojump.zsh ]; then # mac os x with brew
    . `brew --prefix`/etc/autojump.zsh
  fi
fi

# web_search from terminal
function web_search() {
  # get the open command
  local open_cmd
  if [[ "$OSTYPE" = darwin* ]]; then
    open_cmd='open'
  else
    open_cmd='xdg-open'
  fi

  # check whether the search engine is supported
  if [[ ! $1 =~ '(google|bing|baidu|pan115|zhongzilou)' ]];
  then
    echo "Search engine $1 not supported."
    return 1
  fi

  local url="https://www.$1.com"

  # no keyword provided, simply open the search engine homepage
  if [[ $# -le 1 ]]; then
    $open_cmd "$url"
    return
  fi
  local kwd=($@)
  unset kwd[1]
  kwd=$(echo $kwd)
  case $1 in
    "baidu")
        url="${url}/s?wd=$kwd"
        ;;
    "pan115")
        url="${url/https/http}/search?key=$kwd"
        ;;
     "zhongzilou")
        url="${url}/list/$kwd""/1"
        ;;
    *)
        url="${url}/search?q=$kwd"
  esac

  #shift   # shift out $1

  #url="${url%?}" # remove the last '+'
  nohup $open_cmd "$url" &> /dev/null 
}

alias baidu='web_search baidu'
alias bing='web_search bing'
alias google='web_search google'
alias pan='web_search pan115'
alias zzl='web_search zhongzilou'
# add your own !bang searches here

# 打印256色输出效果
format_number() {
  local c=$'\u254F'
  if [ $1 -lt 10 ]; then
    printf "$c %d" $1
  else
    printf "$c%02d" $(($1%100))
  fi
}

somecolors() {
  local from="$1"
  local to="$2"
  local prefix="$3"
  local line

  for line in \
      "\e[2mdim      " \
      "normal   " \
      "\e[1mbold     " \
      "\e[1;2mbold+dim "; do
    echo -ne "$line"
    i=$from
    while [ $i -le $to ]; do
      echo -ne "\e[$prefix${i}m"
      format_number $i
      i=$((i+1))
    done
    echo $'\e[0m\e[K'
  done
}

allcolors() {
  echo "-- 8 standard colors: SGR ${1}0..${1}7 --"
  somecolors 0 7 "$1"
  echo
  echo "-- 8 bright colors: SGR ${2}0..${2}7 --"
  somecolors 0 7 "$2"
  echo
  echo "-- 256 colors: SGR ${1}8;5;0..255 --"
  somecolors 0 15 "${1}8;5;"
  echo
  somecolors  16  51 "${1}8;5;"
  somecolors  52  87 "${1}8;5;"
  somecolors  88 123 "${1}8;5;"
  somecolors 124 159 "${1}8;5;"
  somecolors 160 195 "${1}8;5;"
  somecolors 196 231 "${1}8;5;"
  echo
  somecolors 232 255 "${1}8;5;"
}

function printcolor() {
  allcolors 3 9
  echo
  allcolors 4 10
}

# zshrc重载
function src() {
	local cache="$ZSH_CACHE_DIR"
	autoload -U compinit zrecompile
	compinit -i -d "$cache/zcomp-$HOST"

	for f in ~/.zshrc "$cache/zcomp-$HOST"; do
		zrecompile -p $f && command rm -f $f.zwc.old
	done

	# Use $SHELL if available; remove leading dash if login shell
	[[ -n "$SHELL" ]] && exec ${SHELL#-} || exec zsh
}

#extract自动解压，同样适用于bash
function extract() {
 if [ -z "$1" ]; then
    # display usage if no parameters given
    echo "Usage: extract <path/file_name>.<zip|rar|bz2|gz|tar|tbz2|tgz|Z|7z|xz|ex|tar.bz2|tar.gz|tar.xz>"
    echo "       extract <path/file_name_1.ext> [path/file_name_2.ext] [path/file_name_3.ext]"
 else
    for n in $@
    do
      if [ -f "$n" ] ; then
          case "${n%,}" in
            *.tar.bz2|*.tar.gz|*.tar.xz|*.tbz2|*.tgz|*.txz|*.tar) 
                         tar xvf "$n"       ;;
            *.lzma)      unlzma ./"$n"      ;;
            *.bz2)       bunzip2 ./"$n"     ;;
            *.rar)       unrar x -ad ./"$n" ;;
            *.gz)        gunzip ./"$n"      ;;
            *.zip)       unzip ./"$n"       ;;
            *.z)         uncompress ./"$n"  ;;
            *.7z|*.arj|*.cab|*.chm|*.deb|*.dmg|*.iso|*.lzh|*.msi|*.rpm|*.udf|*.wim|*.xar)
                         7z x ./"$n"        ;;
            *.xz)        unxz ./"$n"        ;;
            *.exe)       cabextract ./"$n"  ;;
            *)
                         echo "extract: '$n' - unknown archive method"
                         return 1
                         ;;
          esac
      else
          echo "'$n' - file does not exist"
          return 1
      fi
    done
fi
}

alias x='extract'

#漂亮又实用的命令高亮界面
setopt extended_glob
 TOKENS_FOLLOWED_BY_COMMANDS=('|' '||' ';' '&' '&&' 'sudo' 'do' 'time' 'strace')
 
recolor-cmd() {
    region_highlight=()
    colorize=true
    start_pos=0
    for arg in ${(z)BUFFER}; do
        ((start_pos+=${#BUFFER[$start_pos+1,-1]}-${#${BUFFER[$start_pos+1,-1]## #}}))
        ((end_pos=$start_pos+${#arg}))
        if $colorize; then
            colorize=false
            res=$(LC_ALL=C builtin type $arg 2>/dev/null)
            case $res in
                *'reserved word'*)   style="fg=magenta,bold";;
                *'alias for'*)       style="fg=cyan,bold";;
                *'shell builtin'*)   style="fg=yellow,bold";;
                *'shell function'*)  style='fg=green,bold';;
                *"$arg is"*)
                    [[ $arg = 'sudo' ]] && style="fg=red,bold" || style="fg=blue,bold";;
                *)                   style="fg=cyan,bold";;
            esac
            region_highlight+=("$start_pos $end_pos $style")
        fi
        [[ ${${TOKENS_FOLLOWED_BY_COMMANDS[(r)${arg//|/\|}]}:+yes} = 'yes' ]] && colorize=true
        start_pos=$end_pos
    done
}
 
#高亮命令和自动补全命令的插件均注册了self-insert和backward-delete-char事件，一次后注册的会覆盖前注册的，所以直接在后注册的函数中调用命令高亮的函数即可
check-cmd-self-insert() { zle .self-insert && recolor-cmd }
check-cmd-backward-delete-char() { zle .backward-delete-char && recolor-cmd }
 
zle -N self-insert check-cmd-self-insert
zle -N backward-delete-char check-cmd-backward-delete-char

#彩色补全菜单
# eval $(dircolors -b)

#文件类型高亮(废弃)
# export LS_COLORS='no=00:fi=00:di=01;34:ln=00;36:pi=40;33:so=01;35:do=01;35:bd=40;33;01:cd=40;33;01:or=41;33;01:ex=00;32:*=02;36:*.cmd=01;32:*.CMD=01;32:*.exe=01;32:*.EXE=01;32:*.com=01;32:*.COM=01;32:*.ocx=01;32:*.OCX=01;32:*.bat=01;32:*.BAT=01;32:*.rsp=01;32:*.RSP=01;32:*.btm=01;32:*.BTM=01;32:*.dll=01;32:*.DLL=01;32:*.tar=00;31:*.TAR=00;31:*.xz=00;31:*.XZ=00;31:*.tbz=00;31:*.TBZ=00;31:*.tgz=00;31:*.TGZ=00;31:*.rpm=00;31:*.RPM=00;31:*.deb=00;31:*.DEB=00;31:*.arj=00;31:*.ARJ=00;31:*.jar=00;31:*.JAR=00;31:*.taz=00;31:*.TAZ=00;31:*.lzh=00;31:*.LZH=00;31:*.lzma=00;31:*.LZMA=00;31:*.rar=00;31:*.RAR=00;31:*.zip=00;31:*.ZIP=00;31:*.iso=00;31:*.ISO=00;31:*.zoo=00;31:*.ZOO=00;31:*.z=00;31:*.Z=00;31:*.gz=00;31:*.GZ=00;31:*.bz2=00;31:*.BZ2=00;31:*.tb2=00;31:*.TB2=00;31:*.tz2=00;31:*.TZ2=00;31:*.tbz2=00;31:*.TBZ2=00;31:*.bin=00;31:*.BIN=00;31:*.bak=00;31:*.BAK=00;31:*.new=00;31:*.NEW=00;31:*.old=00;31:*.OLD=00;31:*.cab=00;31:*.CAB=00;31:*.mdf=01;31:*.MDF=01;31:*.vdi=01;31:*.VDI=01;31:*.box=01;31:*.BOX=01;31:*.vbox=01;31:*.VBOX=01;31:*.bmp=01;33:*.BMP=01;33:*.fli=01;33:*.FLI=01;33:*.gif=01;33:*.GIF=01;33:*.jpg=01;33:*.JPG=01;33:*.jpeg=01;33:*.JPEG=01;33:*.mng=01;33:*.MNG=01;33:*.mov=01;33:*.MOV=01;33:*.mpg=01;33:*.MPG=01;33:*.pcx=01;33:*.PCX=01;33:*.pbm=01;33:*.PBM=01;33:*.pgm=01;33:*.PGM=01;33:*.png=01;33:*.PNG=01;33:*.ico=01;33:*.ICO=01;33:*.ppm=01;33:*.PPM=01;33:*.tga=01;33:*.TGA=01;33:*.tif=01;33:*.TIF=01;33:*.tiff=01;33:*.TIFF=01;33:*.xbm=01;33:*.XBM=01;33:*.xpm=01;33:*.XPM=01;33:*.dl=01;33:*.DL=01;33:*.gl=01;33:*.GL=01;33:*.avi=01;36:*.AVI=01;36:*.wmv=01;36:*.WMV=01;36:*.rm=01;36:*.RM=01;36:*.rmvb=01;36:*.RMVB=01;36:*mp4=01;36:*.MP4=01;36:*.swf=01;36:*.SWF=01;36:*.mkv=01;36:*.MKV=01;36:*.otf=02;33:*.OTF=02;33:*.ttf=02;33:*.TTF=02;33:*.eot=02;33:*.EOT=02;33:*.woff=02;33:*.WOFF=02;33:*.ttc=02;33:*.TTC=02;33:*.woff2=02;33:*.WOFF2=02;33:*.fon=02;33:*.FON=02;33:*.font=02;33:*.FONT=02;33:*.aiff=02;34:*.AIFF=02;34:*.au=02;34:*.AU=02;34:*.mid=02;34:*.MID=02;34:*.mp3=02;34:*.MP3=02;34:*.ogg=02;34:*.OGG=02;34:*.voc=02;34:*.VOC=02;34:*.wav=02;34:*.WAV=02;34:*.tmp=00;33:*.TMP=00;33:*.temp=00;33:*.TEMP=00;33:*.txt=00;33:*.TXT=00;33:*.ini=00;33:*.INI=00;33:*.cf=00;33:*.CF=00;33:*.conf=00;33:*.CONF=00;33:*.cfg=00;33:*.CFG=00;33:*.properties=00;33:*.PROPERTIES=00;33:*.log=00;33:*.LOG=00;33:*.pas=00;33:*.PAS=00;33:*.pass=00;33:*.PASS=00;33:*.map=00;33:*.MAP=00;33:*.out=00;33:*.OUT=00;33:*.res=00;33:*.RES=00;33:*.todo=00;33:*.TODO=00;33:*.crt=00;33:*.CRT=00;33:*.pem=00;33:*.PEM=00;33:*.doc=00;33:*.DOC=00;33:*.docx=00;33:*.DOCX=00;33:*.xls=00;33:*.XLS=00;33:*.xlsx=00;33:*.XLSX=00;33:*.ppt=00;33:*.PPT=00;33:*.pptx=00;33:*.PPTX=00;33:*.pdf=00;33:*.PDF=00;33:*.chm=00;33:*.CHM=00;33:*.csv=00;33:*.CSV=00;33:*.wps=00;33:*.WPS=00;33:*.desktop=02;35:*.DESKTOP=02;35:*.tx=02;35:*.TX=02;35:*.block=02;35:*.BLOCK=02;35:*.idx=02;35:*.IDX=02;35:*.wiz=02;35:*.WIZ=02;35:*.sock=02;35:*.SOCK=02;35:*.pid=02;35:*.PID=02;35:*.et=02;35:*.ET=02;35:*.zsh=00;35:*.ZSH=00;35:*.fish=00;35:*.FISH=00;35:*.bashrc=00;35:*.BASHRC=00;35:*.zshrc=00;35:*.ZSHRC=00;35:*.vimrc=00;35:*.VIMRC=00;35:*.md=00;35:*.MD=00;35:*.yaml=00;35:*.YAML=00;35:*.xml=00;35:*.XML=00;35:*.xaml=00;35:*.XAML=00;35:*.repo=00;35:*.REPO=00;35:*.js=00;35:*.JS=00;35:*.json=00;35:*.JSON=00;35:*.sh=00;35:*.SH=00;35:*.css=00;35:*.CSS=00;35:*.htm=00;35:*.HTM=00;35:*.html=00;35:*.HTML=00;35:*.svg=00;35:*.SVG=00;35:*.h=00;35:*.H=00;35:*.hpp=00;35:*.HPP=00;35:*.hxx=00;35:*.HXX=00;35:*.c=00;35:*.C=00;35:*.obj=00;35:*.OBJ=00;35:*.cpp=00;35:*.CPP=00;35:*.c++=00;35:*.C++=00;35:*.cxx=00;35:*.CXX=00;35:*.cc=00;35:*.CC=00;35:*.cs=00;35:*.CS=00;35:*.vb=00;35:*.VB=00;35:*.java=00;35:*.JAVA=00;35:*.do=00;35:*.DO=00;35:*.asp=00;35:*.ASP=00;35:*.aspx=00;35:*.ASPX=00;35:*.m=00;35:*.M=00;35:*.php=00;35:*.PHP=00;35:*.jsp=00;35:*.JSP=00;35:*.pl=00;35:*.PL=00;35:*.prl=00;35:*.PRL=00;35:*.lua=00;35:*.LUA=00;35:*.d=00;35:*.D=00;35:*.rs=00;35:*.RS=00;35:*.rb=00;35:*.RB=00;35:*.py=00;35:*.PY=00;35:*.pyc=00;35:*.PYC=00;35:*.s=00;35:*.S=00;35:*.f=00;35:*.F=00;35:*.go=00;35:*.GO=00;35:*.scala=00;35:*.SCALA=00;35:*.sol=00;35:*.SOL=00;35:*.tcl=00;35:*.TCL=00;35:*.sql=00;35:*.SQL=00;35:*.test=02;36:'

# ls命令结果高亮（三种方式效果相同）

# 方式１(不推荐)
#cp $(dirname $0)/DIR_COLORS /etc/DIR_COLORS
#cp $(dirname $0)/DIR_COLORS /etc/DIR_COLORS.256color

# 方式２(格式不直观)
#export LS_COLORS='bd=38;5;68:ca=38;5;17:cd=38;5;5;1:di=38;5;4:do=38;5;127:ex=38;5;129;1:pi=38;5;13:fi=38;5;183:ln=target:mh=38;5;222;1:no=38;5;183:or=48;5;1;38;5;232;1:ow=38;5;220;1:sg=48;5;3;38;5;0:su=38;5;220;1;3;100;1:so=38;5;197:st=38;5;86;48;5;234:tw=48;5;235;38;5;139;3:*LS_COLORS=48;5;89;38;5;197;1;3;4;7:*README=38;5;181;1:*README.rst=38;5;181;1:*LICENSE=38;5;181;1:*COPYING=38;5;181;1:*INSTALL=38;5;181;1:*COPYRIGHT=38;5;181;1:*AUTHORS=38;5;181;1:*HISTORY=38;5;181;1:*CONTRIBUTORS=38;5;181;1:*PATENTS=38;5;181;1:*VERSION=38;5;181;1:*NOTICE=38;5;181;1:*CHANGES=38;5;181;1:*.log=38;5;180:*.LOG=38;5;180:*.txt=38;5;153:*.TXT=38;5;153:*.etx=38;5;51:*.ETX=38;5;51:*.info=38;5;51:*.INFO=38;5;51:*.markdown=38;5;51:*.MARKDOWN=38;5;51:*.md=38;5;51:*.MD=38;5;51:*.mkd=38;5;51:*.MKD=38;5;51:*.nfo=38;5;51:*.NFO=38;5;51:*.pod=38;5;51:*.POD=38;5;51:*.rst=38;5;51:*.RST=38;5;51:*.tex=38;5;51:*.TEX=38;5;51:*.textile=38;5;51:*.TEXTILE=38;5;51:*.json=38;5;135:*.JSON=38;5;135:*.msg=38;5;135:*.MSG=38;5;135:*.pgn=38;5;135:*.PGN=38;5;135:*.rss=38;5;135:*.RSS=38;5;135:*.xml=38;5;135:*.XML=38;5;135:*.yaml=38;5;135:*.YAML=38;5;135:*.yml=38;5;135:*.YML=38;5;135:*.rdata=38;5;135:*.RDATA=38;5;135:*.cbr=38;5;141:*.CBR=38;5;141:*.cbz=38;5;141:*.CBZ=38;5;141:*.chm=38;5;141:*.CHM=38;5;141:*.djvu=38;5;141:*.DJVU=38;5;141:*.pdf=38;5;141:*.PDF=38;5;141:*.wps=38;5;117:*.WPS=38;5;117:*.docm=38;5;111;4:*.DOCM=38;5;111;4:*.doc=38;5;111:*.DOC=38;5;111:*.docx=38;5;111:*.DOCX=38;5;111:*.eps=38;5;111:*.EPS=38;5;111:*.ps=38;5;111:*.PS=38;5;111:*.odb=38;5;111:*.ODB=38;5;111:*.odt=38;5;111:*.ODT=38;5;111:*.rtf=38;5;111:*.RTF=38;5;111:*.odp=38;5;122:*.ODP=38;5;122:*.pps=38;5;122:*.PPS=38;5;122:*.ppt=38;5;122:*.PPT=38;5;122:*.pptx=38;5;122:*.PPTX=38;5;122:*.ppts=38;5;122:*.PPTS=38;5;122:*.pptxm=38;5;122;4:*.PPTXM=38;5;122;4:*.pptsm=38;5;122;4:*.PPTSM=38;5;122;4:*.csv=38;5;139:*.CSV=38;5;139:*.ods=38;5;131:*.ODS=38;5;131:*.xla=38;5;133:*.XLA=38;5;133:*.xls=38;5;178:*.XLS=38;5;178:*.xlsx=38;5;178:*.XLSX=38;5;178:*.xlsxm=38;5;178;4:*.XLSXM=38;5;178;4:*.xltm=38;5;40;4:*.XLTM=38;5;40;4:*.xltx=38;5;40:*.XLTX=38;5;40:*cfg=38;5;190:*conf=38;5;190:*rc=38;5;190:*.cf=38;5;190:*.CF=38;5;190:*.cfg=38;5;190:*.CFG=38;5;190:*.conf=38;5;190:*.CONF=38;5;190:*.rc=38;5;190:*.RC=38;5;190:*.ini=38;5;190:*.INI=38;5;190:*.plist=38;5;190:*.PLIST=38;5;190:*.viminfo=38;5;166:*.VIMINFO=38;5;166:*.vimrc=38;5;166:*.VIMRC=38;5;166:*.zshrc=38;5;166:*.ZSHRC=38;5;166:*.bashrc=38;5;166:*.BASHRC=38;5;166:*.xonshrc=38;5;166:*.XONSHRC=38;5;166:*.pcf=38;5;172:*.PCF=38;5;172:*.psf=38;5;70:*.PSF=38;5;70:*.git=38;5;213:*.GIT=38;5;213:*.gitignore=38;5;225:*.GITIGNORE=38;5;225:*.gitattributes=38;5;225:*.GITATTRIBUTES=38;5;225:*.gitmodules=38;5;225:*.GITMODULES=38;5;225:*.awk=38;5;173:*.AWK=38;5;173:*.bash=38;5;173:*.BASH=38;5;173:*.sed=38;5;173:*.SED=38;5;173:*.sh=38;5;173:*.SH=38;5;173:*.csh=38;5;173:*.CSH=38;5;173:*.ksh=38;5;173:*.KSH=38;5;173:*.fish=38;5;173:*.FISH=38;5;173:*.zsh=38;5;173:*.ZSH=38;5;173:*.vim=38;5;173:*.VIM=38;5;173:*.ahk=38;5;195:*.AHK=38;5;195:*.py=38;5;10:*.PY=38;5;10:*.pl=38;5;93:*.PL=38;5;93:*.t=38;5;93:*.T=38;5;93:*.msql=38;5;97:*.MSQL=38;5;97:*.mysql=38;5;97:*.MYSQL=38;5;97:*.pgsql=38;5;97:*.PGSQL=38;5;97:*.sql=38;5;97:*.SQL=38;5;97:*.tcl=38;5;95;1:*.TCL=38;5;95;1:*.r=38;5;63:*.R=38;5;63:*.gs=38;5;89:*.GS=38;5;89:*.repo=38;5;98:*.REPO=38;5;98:*.asm=38;5;100:*.ASM=38;5;100:*.cl=38;5;96:*.CL=38;5;96:*.lisp=38;5;96:*.LISP=38;5;96:*.lua=38;5;103:*.LUA=38;5;103:*.moon=38;5;116:*.MOON=38;5;116:*.c=38;5;142:*.C=38;5;142:*.h=38;5;136:*.H=38;5;136:*.obj=38;5;144:*.OBJ=38;5;144:*.tcc=38;5;136:*.TCC=38;5;136:*.c++=38;5;10:*.C++=38;5;10:*.h++=38;5;137:*.H++=38;5;137:*.hpp=38;5;137:*.HPP=38;5;137:*.hxx=38;5;137:*.HXX=38;5;137:*.ii=38;5;138:*.II=38;5;138:*.m=38;5;201:*.M=38;5;201:*.cc=38;5;184:*.CC=38;5;184:*.cs=38;5;184:*.CS=38;5;184:*.cp=38;5;214:*.CP=38;5;214:*.cpp=38;5;214:*.CPP=38;5;214:*.cxx=38;5;184:*.CXX=38;5;184:*.asp=38;5;211:*.ASP=38;5;211:*.aspx=38;5;212:*.ASPX=38;5;212:*.cr=38;5;159:*.CR=38;5;159:*.go=38;5;112:*.GO=38;5;112:*.scala=38;5;9:*.SCALA=38;5;9:*.sol=38;5;90:*.SOL=38;5;90:*.f=38;5;58:*.F=38;5;58:*.for=38;5;58:*.FOR=38;5;58:*.ftn=38;5;58:*.FTN=38;5;58:*.s=38;5;30:*.S=38;5;30:*.rs=38;5;79:*.RS=38;5;79:*.d=38;5;23:*.D=38;5;23:*.swift=38;5;33:*.SWIFT=38;5;33:*.sx=38;5;81:*.SX=38;5;81:*.hi=38;5;110:*.HI=38;5;110:*.hs=38;5;69:*.HS=38;5;69:*.lhs=38;5;69:*.LHS=38;5;69:*.pyc=38;5;14:*.PYC=38;5;14:*.css=38;5;216;1:*.CSS=38;5;216;1:*.less=38;5;226;1:*.LESS=38;5;226;1:*.sass=38;5;227;1:*.SASS=38;5;227;1:*.scss=38;5;228;1:*.SCSS=38;5;228;1:*.htm=38;5;200;1:*.HTM=38;5;200;1:*.html=38;5;201;1:*.HTML=38;5;201;1:*.jhtm=38;5;165;1:*.JHTM=38;5;165;1:*.mht=38;5;199;1:*.MHT=38;5;199;1:*.eml=38;5;162;1:*.EML=38;5;162;1:*.mustache=38;5;161;1:*.MUSTACHE=38;5;161;1:*.coffee=38;5;205;1:*.COFFEE=38;5;205;1:*.java=38;5;99;1:*.JAVA=38;5;99;1:*.js=38;5;202;1:*.JS=38;5;202;1:*.mjs=38;5;203;1:*.MJS=38;5;203;1:*.jsm=38;5;204;1:*.JSM=38;5;204;1:*.jsm=38;5;205;1:*.JSM=38;5;205;1:*.jsp=38;5;219;1:*.JSP=38;5;219;1:*.php=38;5;25:*.PHP=38;5;25:*.ctp=38;5;24:*.CTP=38;5;24:*.twig=38;5;67:*.TWIG=38;5;67:*.vb=38;5;146:*.VB=38;5;146:*.vba=38;5;146:*.VBA=38;5;146:*.vbs=38;5;146:*.VBS=38;5;146:*Dockerfile=38;5;88:*Makefile=38;5;88:*MANIFEST=38;5;88:*pm_to_blib=38;5;88:*.dockerfile=38;5;88:*.DOCKERFILE=38;5;88:*.dockerignore=38;5;88:*.DOCKERIGNORE=38;5;88:*.makefile=38;5;88:*.MAKEFILE=38;5;88:*.manifest=38;5;88:*.MANIFEST=38;5;88:*.pm_to_blib=38;5;88:*.PM_TO_BLIB=38;5;88:*.am=38;5;56:*.AM=38;5;56:*.in=38;5;56:*.IN=38;5;56:*.hin=38;5;56:*.HIN=38;5;56:*.scan=38;5;56:*.SCAN=38;5;56:*.m4=38;5;56:*.M4=38;5;56:*.old=38;5;56:*.OLD=38;5;56:*.out=38;5;56:*.OUT=38;5;56:*.skip=38;5;237:*.SKIP=38;5;237:*.diff=48;5;23;38;5;24:*.DIFF=48;5;23;38;5;24:*.patch=48;5;29;38;5;30;1:*.PATCH=48;5;29;38;5;30;1:*.bmp=38;5;64:*.BMP=38;5;64:*.tiff=38;5;64:*.TIFF=38;5;64:*.tif=38;5;64:*.TIF=38;5;64:*.tga=38;5;64:*.TGA=38;5;64:*.cdr=38;5;64:*.CDR=38;5;64:*.gif=38;5;64:*.GIF=38;5;64:*.ico=38;5;64:*.ICO=38;5;64:*.jpeg=38;5;64:*.JPEG=38;5;64:*.jpg=38;5;64:*.JPG=38;5;64:*.fli=38;5;64:*.FLI=38;5;64:*.nth=38;5;64:*.NTH=38;5;64:*.png=38;5;64:*.PNG=38;5;64:*.psd=38;5;64:*.PSD=38;5;64:*.pcx=38;5;64:*.PCX=38;5;64:*.mng=38;5;64:*.MNG=38;5;64:*.pbm=38;5;64:*.PBM=38;5;64:*.pgm=38;5;64:*.PGM=38;5;64:*.xpm=38;5;64:*.XPM=38;5;64:*.ppm=38;5;64:*.PPM=38;5;64:*.ai=38;5;123:*.AI=38;5;123:*.eps=38;5;123:*.EPS=38;5;123:*.epsf=38;5;123:*.EPSF=38;5;123:*.drw=38;5;123:*.DRW=38;5;123:*.ps=38;5;123:*.PS=38;5;123:*.svg=38;5;123:*.SVG=38;5;123:*.avi=38;5;75:*.AVI=38;5;75:*.divx=38;5;75:*.DIVX=38;5;75:*.ifo=38;5;75:*.IFO=38;5;75:*.m2v=38;5;75:*.M2V=38;5;75:*.m4v=38;5;75:*.M4V=38;5;75:*.mkv=38;5;75:*.MKV=38;5;75:*.mov=38;5;75:*.MOV=38;5;75:*.mp4=38;5;75:*.MP4=38;5;75:*.mpeg=38;5;75:*.MPEG=38;5;75:*.mpg=38;5;75:*.MPG=38;5;75:*.ogm=38;5;75:*.OGM=38;5;75:*.rmvb=38;5;75:*.RMVB=38;5;75:*.sample=38;5;75:*.SAMPLE=38;5;75:*.wmv=38;5;75:*.WMV=38;5;75:*.3g2=38;5;73:*.3gp=38;5;73:*.gp3=38;5;73:*.GP3=38;5;73:*.webm=38;5;73:*.WEBM=38;5;73:*.gp4=38;5;73:*.GP4=38;5;73:*.asf=38;5;73:*.ASF=38;5;73:*.flv=38;5;73:*.FLV=38;5;73:*.ts=38;5;73:*.TS=38;5;73:*.ogv=38;5;73:*.OGV=38;5;73:*.f4v=38;5;73:*.F4V=38;5;73:*.vob=38;5;85;1:*.VOB=38;5;85;1:*.3ga=38;5;72;1:*.s3m=38;5;72;1:*.S3M=38;5;72;1:*.aac=38;5;72;1:*.AAC=38;5;72;1:*.au=38;5;72;1:*.AU=38;5;72;1:*.dat=38;5;72;1:*.DAT=38;5;72;1:*.dts=38;5;72;1:*.DTS=38;5;72;1:*.fcm=38;5;72;1:*.FCM=38;5;72;1:*.m4a=38;5;72;1:*.M4A=38;5;72;1:*.mid=38;5;72;1:*.MID=38;5;72;1:*.midi=38;5;72;1:*.MIDI=38;5;72;1:*.mod=38;5;72;1:*.MOD=38;5;72;1:*.mp3=38;5;72;1:*.MP3=38;5;72;1:*.mp4a=38;5;72;1:*.MP4A=38;5;72;1:*.oga=38;5;72;1:*.OGA=38;5;72;1:*.ogg=38;5;72;1:*.OGG=38;5;72;1:*.opus=38;5;72;1:*.OPUS=38;5;72;1:*.s3m=38;5;72;1:*.S3M=38;5;72;1:*.sid=38;5;72;1:*.SID=38;5;72;1:*.wma=38;5;72;1:*.WMA=38;5;72;1:*.ape=38;5;80;1:*.APE=38;5;80;1:*.aiff=38;5;80;1:*.AIFF=38;5;80;1:*.cda=38;5;80;1:*.CDA=38;5;80;1:*.flac=38;5;80;1:*.FLAC=38;5;80;1:*.alac=38;5;80;1:*.ALAC=38;5;80;1:*.midi=38;5;80;1:*.MIDI=38;5;80;1:*.pcm=38;5;80;1:*.PCM=38;5;80;1:*.wav=38;5;80;1:*.WAV=38;5;80;1:*.wv=38;5;80;1:*.WV=38;5;80;1:*.wvc=38;5;80;1:*.WVC=38;5;80;1:*.afm=38;5;157:*.AFM=38;5;157:*.fon=38;5;157:*.FON=38;5;157:*.font=38;5;157:*.FONT=38;5;157:*.fnt=38;5;157:*.FNT=38;5;157:*.pfb=38;5;157:*.PFB=38;5;157:*.pfm=38;5;157:*.PFM=38;5;157:*.ttf=38;5;157:*.TTF=38;5;157:*.ttc=38;5;157:*.TTC=38;5;157:*.otf=38;5;157:*.OTF=38;5;157:*.eot=38;5;157:*.EOT=38;5;157:*.aiff=38;5;157:*.AIFF=38;5;157:*.woff=38;5;157:*.WOFF=38;5;157:*.woff2=38;5;157:*.WOFF2=38;5;157:*.pfa=38;5;156:*.PFA=38;5;156:*.7z=38;5;130:*.a=38;5;130:*.A=38;5;130:*.arj=38;5;130:*.ARJ=38;5;130:*.bz2=38;5;130:*.BZ2=38;5;130:*.cpio=38;5;130:*.CPIO=38;5;130:*.gz=38;5;130:*.GZ=38;5;130:*.bz2=38;5;130:*.BZ2=38;5;130:*.tb2=38;5;130:*.TB2=38;5;130:*.tz2=38;5;130:*.TZ2=38;5;130:*.tbz2=38;5;130:*.TBZ2=38;5;130:*.lrz=38;5;130:*.LRZ=38;5;130:*.lz=38;5;130:*.LZ=38;5;130:*.lzh=38;5;130:*.LZH=38;5;130:*.lzma=38;5;130:*.LZMA=38;5;130:*.lzo=38;5;130:*.LZO=38;5;130:*.rar=38;5;130:*.RAR=38;5;130:*.s7z=38;5;130:*.S7Z=38;5;130:*.sz=38;5;130:*.SZ=38;5;130:*.tar=38;5;130:*.TAR=38;5;130:*.taz=38;5;130:*.TAZ=38;5;130:*.tgz=38;5;130:*.TGZ=38;5;130:*.xz=38;5;130:*.XZ=38;5;130:*.z=38;5;130:*.Z=38;5;130:*.zip=38;5;130:*.ZIP=38;5;130:*.zipx=38;5;130:*.ZIPX=38;5;130:*.zoo=38;5;130:*.ZOO=38;5;130:*.zpaq=38;5;130:*.ZPAQ=38;5;130:*.zz=38;5;130:*.ZZ=38;5;130:*.apk=38;5;160:*.APK=38;5;160:*.deb=38;5;160:*.DEB=38;5;160:*.rpm=38;5;160:*.RPM=38;5;160:*.jad=38;5;160:*.JAD=38;5;160:*.jar=38;5;160:*.JAR=38;5;160:*.cab=38;5;160:*.CAB=38;5;160:*.mdf=38;5;160:*.MDF=38;5;160:*.pak=38;5;160:*.PAK=38;5;160:*.pk3=38;5;160:*.PK3=38;5;160:*.vdf=38;5;160:*.VDF=38;5;160:*.vpk=38;5;160:*.VPK=38;5;160:*.bsp=38;5;160:*.BSP=38;5;160:*.dmg=38;5;160:*.DMG=38;5;160:*.exe=38;5;52:*.EXE=38;5;52:*.msi=38;5;52:*.MSI=38;5;52:*.rsp=38;5;52:*.RSP=38;5;52:*.btm=38;5;52:*.BTM=38;5;52:*.dll=38;5;52:*.DLL=38;5;52:*.osx=38;5;52:*.OSX=38;5;52:*.ocx=38;5;52:*.OCX=38;5;52:*.cmd=38;5;52:*.CMD=38;5;52:*.bat=38;5;52:*.BAT=38;5;52:*.reg=38;5;52:*.REG=38;5;52:*.r[0-9]{0,2}=38;5;145:*.R[0-9]{0,2}=38;5;145:*.zx[0-9]{0,2}=38;5;145:*.ZX[0-9]{0,2}=38;5;145:*.z[0-9]{0,2}=38;5;145:*.Z[0-9]{0,2}=38;5;145:*.part=38;5;24:*.PART=38;5;24:*.iso=38;5;94:*.ISO=38;5;94:*.bin=38;5;94:*.BIN=38;5;94:*.nrg=38;5;94:*.NRG=38;5;94:*.qcow=38;5;94:*.QCOW=38;5;94:*.sparseimage=38;5;94:*.SPARSEIMAGE=38;5;94:*.toast=38;5;94:*.TOAST=38;5;94:*.vcd=38;5;94:*.VCD=38;5;94:*.vmdk=38;5;94:*.VMDK=38;5;94:*.vdi=38;5;94:*.VDI=38;5;94:*.vdisk=38;5;94:*.VDISK=38;5;94:*.box=38;5;94:*.BOX=38;5;94:*.img=38;5;94:*.IMG=38;5;94:*.vbox=38;5;94:*.VBOX=38;5;94:*.accdb=38;5;60:*.ACCDB=38;5;60:*.accde=38;5;60:*.ACCDE=38;5;60:*.accdr=38;5;60:*.ACCDR=38;5;60:*.accdt=38;5;60:*.ACCDT=38;5;60:*.db=38;5;60:*.DB=38;5;60:*.fmp12=38;5;60:*.FMP12=38;5;60:*.fp7=38;5;60:*.FP7=38;5;60:*.localstorage=38;5;60:*.LOCALSTORAGE=38;5;60:*.mdb=38;5;60:*.MDB=38;5;60:*.mde=38;5;60:*.MDE=38;5;60:*.sqlite=38;5;60:*.SQLITE=38;5;60:*.typelib=38;5;60:*.TYPELIB=38;5;60:*.nc=38;5;61:*.NC=38;5;61:*.pacnew=38;5;230:*.PACNEW=38;5;230:*.un~=38;5;230:*.UN~=38;5;230:*.orig=38;5;230:*.ORIG=38;5;230:*.bup=38;5;224:*.BUP=38;5;224:*.test=38;5;224:*.TEST=38;5;224:*.bak=38;5;224:*.BAK=38;5;224:*.res=38;5;224:*.RES=38;5;224:*.hd=38;5;224:*.HD=38;5;224:*.new=38;5;224:*.NEW=38;5;224:*.old=38;5;224:*.OLD=38;5;224:*.o=38;5;224:*.O=38;5;224:*.rlib=38;5;224:*.RLIB=38;5;224:*.swp=38;5;229:*.SWP=38;5;229:*.swo=38;5;229:*.SWO=38;5;229:*.tmp=38;5;229:*.TMP=38;5;229:*.temp=38;5;229:*.TEMP=38;5;229:*.sassc=38;5;229:*.SASSC=38;5;229:*.et=38;5;193:*.ET=38;5;193:*.pid=38;5;193:*.PID=38;5;193:*.state=38;5;193:*.STATE=38;5;193:*lockfile=38;5;193:*.err=38;5;187;1:*.ERR=38;5;187;1:*.error=38;5;187;1:*.ERROR=38;5;187;1:*.stderr=38;5;187;1:*.STDERR=38;5;187;1:*.dump=38;5;189:*.DUMP=38;5;189:*.stackdump=38;5;189:*.STACKDUMP=38;5;189:*.zcompdump=38;5;189:*.ZCOMPDUMP=38;5;189:*.zwc=38;5;189:*.ZWC=38;5;189:*.pcap=38;5;102:*.PCAP=38;5;102:*.cap=38;5;102:*.CAP=38;5;102:*.dmp=38;5;102:*.DMP=38;5;102:*.ds_store=38;5;176:*.DS_STORE=38;5;176:*.localized=38;5;176:*.LOCALIZED=38;5;176:*.cfusertextencoding=38;5;176:*.CFUSERTEXTENCODING=38;5;176:*.allow=38;5;153:*.ALLOW=38;5;153:*.deny=38;5;152:*.DENY=38;5;152:*.service=38;5;43:*.SERVICE=38;5;43:*@.service=38;5;43:*.socket=38;5;43:*.SOCKET=38;5;43:*.sock=38;5;43:*.SOCK=38;5;43:*.swap=38;5;43:*.SWAP=38;5;43:*.device=38;5;43:*.DEVICE=38;5;43:*.mount=38;5;43:*.MOUNT=38;5;43:*.automount=38;5;43:*.AUTOMOUNT=38;5;43:*.target=38;5;43:*.TARGET=38;5;43:*.path=38;5;43:*.PATH=38;5;43:*.timer=38;5;43:*.TIMER=38;5;43:*.snapshot=38;5;43:*.SNAPSHOT=38;5;43:*.tx=38;5;60:*.TX=38;5;60:*.block=38;5;60:*.BLOCK=38;5;60:*.idx=38;5;60:*.IDX=38;5;60:*.wiz=38;5;60:*.WIZ=38;5;60:*.application=38;5;155:*.APPLICATION=38;5;155:*.cue=38;5;155:*.CUE=38;5;155:*.description=38;5;155:*.DESCRIPTION=38;5;155:*.directory=38;5;155:*.DIRECTORY=38;5;155:*.m3u=38;5;155:*.M3U=38;5;155:*.m3u8=38;5;155:*.M3U8=38;5;155:*.md5=38;5;155:*.MD5=38;5;155:*.properties=38;5;155:*.PROPERTIES=38;5;155:*.sfv=38;5;155:*.SFV=38;5;155:*.srt=38;5;155:*.SRT=38;5;155:*.theme=38;5;155:*.THEME=38;5;155:*.torrent=38;5;155:*.TORRENT=38;5;155:*.urlview=38;5;155:*.URLVIEW=38;5;155:*.asc=38;5;192;3:*.ASC=38;5;192;3:*.bfe=38;5;192;3:*.BFE=38;5;192;3:*.enc=38;5;192;3:*.ENC=38;5;192;3:*.gpg=38;5;192;3:*.GPG=38;5;192;3:*.signature=38;5;192;3:*.SIGNATURE=38;5;192;3:*.sig=38;5;192;3:*.SIG=38;5;192;3:*.p12=38;5;192;3:*.P12=38;5;192;3:*.pem=38;5;192;3:*.PEM=38;5;192;3:*.crt=38;5;192;3:*.CRT=38;5;192;3:*.pgp=38;5;192;3:*.PGP=38;5;192;3:*.asc=38;5;192;3:*.ASC=38;5;192;3:*.enc=38;5;192;3:*.ENC=38;5;192;3:*.sig=38;5;192;3:*.SIG=38;5;192;3:*.32x=38;5;107:*.cdi=38;5;107:*.CDI=38;5;107:*.fm2=38;5;107:*.FM2=38;5;107:*.rom=38;5;107:*.ROM=38;5;107:*.sav=38;5;107:*.SAV=38;5;107:*.st=38;5;107:*.ST=38;5;107:*.a00=38;5;108:*.A00=38;5;108:*.a52=38;5;108:*.A52=38;5;108:*.a64=38;5;108:*.A64=38;5;108:*.a78=38;5;108:*.A78=38;5;108:*.adf=38;5;108:*.ADF=38;5;108:*.atr=38;5;108:*.ATR=38;5;108:*.gb=38;5;109:*.GB=38;5;109:*.gba=38;5;109:*.GBA=38;5;109:*.gbc=38;5;109:*.GBC=38;5;109:*.gel=38;5;109:*.GEL=38;5;109:*.gg=38;5;109:*.GG=38;5;109:*.ggl=38;5;109:*.GGL=38;5;109:*.ipk=38;5;109:*.IPK=38;5;109:*.j64=38;5;109:*.J64=38;5;109:*.nds=38;5;109:*.NDS=38;5;109:*.nes=38;5;109:*.NES=38;5;109:*.sms=38;5;194:*.SMS=38;5;194:*.pot=38;5;101:*.POT=38;5;101:*.pcb=38;5;101:*.PCB=38;5;101:*.mm=38;5;101:*.MM=38;5;101:*.pod=38;5;101:*.POD=38;5;101:*.gbr=38;5;101:*.GBR=38;5;101:*.spl=38;5;101:*.SPL=38;5;101:*.scm=38;5;101:*.SCM=38;5;101:*.rproj=38;5;11:*.RPROJ=38;5;11:*.sis=38;5;65:*.SIS=38;5;65:*.1p=38;5;65:*.3p=38;5;65:*.cnc=38;5;65:*.CNC=38;5;65:*.def=38;5;65:*.DEF=38;5;65:*.ex=38;5;65:*.EX=38;5;65:*.example=38;5;65:*.EXAMPLE=38;5;65:*.feature=38;5;65:*.FEATURE=38;5;65:*.ger=38;5;65:*.GER=38;5;65:*.map=38;5;65:*.MAP=38;5;65:*.mf=38;5;65:*.MF=38;5;65:*.mfasl=38;5;65:*.MFASL=38;5;65:*.mi=38;5;65:*.MI=38;5;65:*.mtx=38;5;65:*.MTX=38;5;65:*.pc=38;5;65:*.PC=38;5;65:*.pi=38;5;65:*.PI=38;5;65:*.plt=38;5;65:*.PLT=38;5;65:*.pm=38;5;65:*.PM=38;5;65:*.rb=38;5;65:*.RB=38;5;65:*.rdf=38;5;65:*.RDF=38;5;65:*.rst=38;5;65:*.RST=38;5;65:*.ru=38;5;65:*.RU=38;5;65:*.sch=38;5;65:*.SCH=38;5;65:*.sty=38;5;65:*.STY=38;5;65:*.sug=38;5;65:*.SUG=38;5;65:*.t=38;5;65:*.T=38;5;65:*.tdy=38;5;65:*.TDY=38;5;65:*.tfm=38;5;65:*.TFM=38;5;65:*.tfnt=38;5;65:*.TFNT=38;5;65:*.tg=38;5;65:*.TG=38;5;65:*.vcard=38;5;65:*.VCARD=38;5;65:*.vcf=38;5;65:*.VCF=38;5;65:*.xln=38;5;65:*.XLN=38;5;65:*.iml=38;5;12:*.IML=38;5;12:*.xcconfig=38;5;158:*.XCCONFIG=38;5;158:*.entitlements=38;5;158:*.ENTITLEMENTS=38;5;158:*.strings=38;5;158:*.STRINGS=38;5;158:*.storyboard=38;5;158:*.STORYBOARD=38;5;158:*.xcsettings=38;5;158:*.XCSETTINGS=38;5;158:*.xib=38;5;158:*.XIB=38;5;158:'

# 方式３(依赖外部文件)
eval $(dircolors -b $cdir/DIR_COLORS)

export ZLSCOLORS="${LS_COLORS}"
zmodload zsh/complist
zstyle ':completion:*' list-colors ${(s.:.)LS_COLORS}
zstyle ':completion:*:*:kill:*:processes' list-colors '=(#b) #([0-9]#)*=0=01;31'

#lf目录图标
if type lf &>/dev/null; then
  LF_ICONS=$(sed $cdir/DIR_ICONS -e '/^[ \t]*#/d' -e '/^[ \t]*$/d' -e 's/[ \t]\+/=/g' -e 's/$/ /')
  LF_ICONS=${LF_ICONS//$'\n'/:}
  export LF_ICONS
fi

#强制启用颜色高亮
force_color_prompt=yes

#高亮man帮助文档，zsh中启用colored-man可达到一样效果
export LESS_TERMCAP_mb=$'\E[01;31m'
export LESS_TERMCAP_md=$'\E[01;31m'
export LESS_TERMCAP_me=$'\E[0m'
export LESS_TERMCAP_se=$'\E[0m'
export LESS_TERMCAP_so=$'\E[01;44;33m'
export LESS_TERMCAP_ue=$'\E[0m'
export LESS_TERMCAP_us=$'\E[01;32m'

# syntax color definition
ZSH_HIGHLIGHT_HIGHLIGHTERS=(main brackets pattern)
#ZSH_HIGHLIGHT_HIGHLIGHTERS=(main brackets pattern cursor)

typeset -A ZSH_HIGHLIGHT_STYLES

# ZSH_HIGHLIGHT_STYLES[command]=fg=white,bold
# ZSH_HIGHLIGHT_STYLES[alias]='fg=magenta,bold'

ZSH_HIGHLIGHT_STYLES[default]=fg=magenta,bold
ZSH_HIGHLIGHT_STYLES[unknown-token]=fg=009
ZSH_HIGHLIGHT_STYLES[reserved-word]=fg=009,standout
ZSH_HIGHLIGHT_STYLES[alias]=fg=033,bold
ZSH_HIGHLIGHT_STYLES[builtin]=fg=032,bold
ZSH_HIGHLIGHT_STYLES[function]=fg=031,bold
ZSH_HIGHLIGHT_STYLES[command]=fg=036,bold
ZSH_HIGHLIGHT_STYLES[precommand]=fg=034,bold
ZSH_HIGHLIGHT_STYLES[commandseparator]=fg=034,bold
ZSH_HIGHLIGHT_STYLES[hashed-command]=fg=071,bold
ZSH_HIGHLIGHT_STYLES[path]=fg=214,underline
ZSH_HIGHLIGHT_STYLES[path_pathseparator]=fg=211,underline
ZSH_HIGHLIGHT_STYLES[path_prefix]=fg=125,underline
ZSH_HIGHLIGHT_STYLES[path_prefix_pathseparator]=fg=223,underline
ZSH_HIGHLIGHT_STYLES[globbing]=fg=158,bold
ZSH_HIGHLIGHT_STYLES[history-expansion]=fg=214,bold,underline
ZSH_HIGHLIGHT_STYLES[single-hyphen-option]=fg=180,bold
ZSH_HIGHLIGHT_STYLES[double-hyphen-option]=fg=181,bold
ZSH_HIGHLIGHT_STYLES[back-quoted-argument]=fg=093,bold
ZSH_HIGHLIGHT_STYLES[single-quoted-argument]=fg=062,bold
ZSH_HIGHLIGHT_STYLES[double-quoted-argument]=fg=063,bold
ZSH_HIGHLIGHT_STYLES[double-quoted-argument_unclosed]=fg=064,bold
ZSH_HIGHLIGHT_STYLES[dollar-quoted-argument]=fg=131,bold
ZSH_HIGHLIGHT_STYLES[back-double-quoted-argument]=fg=011,bold
ZSH_HIGHLIGHT_STYLES[dollar-double-quoted-argument]=fg=132,bold
ZSH_HIGHLIGHT_STYLES[comment]=fg=101,bold
ZSH_HIGHLIGHT_STYLES[redirection]=fg=125,bold
ZSH_HIGHLIGHT_STYLES[assign]=fg=magenta,bold
# ZSH_HIGHLIGHT_STYLES[assign]=none

ZSH_HIGHLIGHT_REGEXP+=('\bsudo\b' fg=123,bold)

# vi mode
#使用:<<' 注释内容 ' 的形式注释掉Incremental completion插件，删除空行的单引号可以重新启用
# Updates editor information when the keymap changes.
function zle-keymap-select() {
  zle reset-prompt
  zle -R
}

# Ensure that the prompt is redrawn when the terminal size changes.
TRAPWINCH() {
  zle &&  zle -R
}

zle -N zle-keymap-select
zle -N edit-command-line

bindkey -v

# allow v to edit the command line (standard behaviour)
autoload -Uz edit-command-line
bindkey -M vicmd 'v' edit-command-line

# allow ctrl-p, ctrl-n for navigate history (standard behaviour)
bindkey '^P' up-history
bindkey '^N' down-history

# allow ctrl-h, ctrl-w, ctrl-? for char and word deletion (standard behaviour)
bindkey '^?' backward-delete-char
bindkey '^h' backward-delete-char
bindkey '^w' backward-kill-word

# allow ctrl-r to perform backward search in history
bindkey '^r' history-incremental-search-backward

# allow ctrl-a and ctrl-e to move to beginning/end of line
bindkey '^a' beginning-of-line
bindkey '^e' end-of-line

# 将回车绑定到其他键上(如按a代表回车)
# bindkey 'a' accept-line

# if mode indicator wasn't setup by theme, define default
if [[ "$MODE_INDICATOR" == "" ]]; then
  MODE_INDICATOR="%{$fg_bold[red]%}<%{$fg[red]%}<<%{$reset_color%}"
fi

function vi_mode_prompt_info() {
  echo "${${KEYMAP/vicmd/$MODE_INDICATOR}/(main|viins)/}"
}

# define right prompt, if it wasn't defined by a theme
if [[ "$RPS1" == "" && "$RPROMPT" == "" ]]; then
  RPS1='$(vi_mode_prompt_info)'
fi

# 加载额外功能
autoload -U compinit && compinit
autoload -U promptinit && promptinit
autoload -U add-zsh-hook
autoload -U zcalc

#########################################################################
# 端口开放
#########################################################################

function portopen() {
for port in $*;
do
  sudo /sbin/iptables -I INPUT -p tcp --dport $port -j ACCEPT
done
}

function portclose() {
for port in $*;
do
  sudo /sbin/iptables -I INPUT -p tcp --dport $port -j DROP
done
}

#########################################################################
# 计算器
#########################################################################

function calc() {
    zcalc -e "$*"
}

function c() 
{ 
   local in="$(echo " $*" | sed -e 's/\[/(/g' | sed -e 's/\]/)/g')";
   awk "BEGIN {printf $in}"
}

#########################################################################
# 大小写转换
#########################################################################

function upper() {
    echo "$*" | tr '[:lower:]' '[:upper:]'
}

function lower() {
    echo "$*" | tr '[:upper:]' '[:lower:]'
}

function capitalize() {
    echo "$*" | tr '[:upper:]' '[:lower:]' | sed 's/^\w\|\s\w/\U&/g'
}

#########################################################################
# 批量ping主机
#########################################################################
function pping() {
  trap "exit 1" SIGINT SIGQUIT
  for ipaddr in $@; do
    if [ $(echo "$ipaddr" |tr -cd 1 |wc -c) -eq 2 ]; then
      ipd=$ipaddr
      ipr="1-254"
    else
      ipd=${ipaddr%.*}
      ipr=${ipaddr##*.}    
    fi

    ipstart=$ipr
    ipend=$ipr

    if [[ "$ipr" =~ "-" ]]; then
      ipstart=${ipr%-*}
      ipend=${ipr##*-}
    fi

    for ((i = $ipstart; i <= $ipend; i++)); do
      if [ $ipd == $ipr ]; then
        cip=$ipd
      else
        cip=$ipd.$i
      fi
      if [ ${cip##*.} -le 9 ]; then
        sip="$cip  "
      elif [ ${cip##*.} -le 99 ]; then
        sip="$cip "
      else
        sip="$cip"
      fi
      if ping -c 1 -w 1 $cip &>/dev/null;then
        echo -e "\e[1;32m $sip is up \e[0m"
       else
         echo -e "\e[1;31m $sip is down \e[0m"
      fi
    done
  done
}

#########################################################################
# 自动执行sudo命令(Alt+Enter)
#########################################################################

sudo-command-line() {
    [[ -z $BUFFER ]] && zle up-history
    if [[ $BUFFER == sudo\ * ]]; then
        LBUFFER="${LBUFFER#sudo }"
    elif [[ $BUFFER == $EDITOR\ * ]]; then
        LBUFFER="${LBUFFER#$EDITOR }"
        LBUFFER="sudoedit $LBUFFER"
    elif [[ $BUFFER == sudoedit\ * ]]; then
        LBUFFER="${LBUFFER#sudoedit }"
        LBUFFER="$EDITOR $LBUFFER"
    else
        LBUFFER="sudo $LBUFFER"
    fi
    
    zle accept-line
}
zle -N sudo-command-line
# Defined shortcut keys: [Alt] [Enter]
bindkey -M emacs '^[^M' sudo-command-line
bindkey -M vicmd '^[^M' sudo-command-line
bindkey -M viins '^[^M' sudo-command-line

#########################################################################
# 远程主机批量执行任务(类似ansible)
#########################################################################
function batch() {
  local usage='
USAGE: \n
COMMAND: \n
\tbatch -u [user] -s [password] -p [port] <mode> <hosts> <args> \n
MODES: \n
\tping: ping hosts \n
\tshell: run command on hosts \n
\tscript: execute script on hosts \n
\tcopy: copy local files to hosts \n
\tfetch: fetch files on hosts to local \n
\ttemplate: replace string for the file on hosts \n
EXAMPLES: \n
\tbatch ping "127.0.0.1 196.168.0.1" \n
\tbatch ping "${array[*]}" \n
\tbatch -u root -p 22 shell "${array[*]}" "ps" \n
\tbatch -s password script "${array[*]}" "run.sh" \n
\tbatch copy "${array[*]}" "/home/demo/run.sh /root/" \n
\tbatch template "${array[*]}" "/home/demo/run.sh foo bar" \n
  '
  local user="root"
  local pw=""
  local port="22"
  local cfile=""
  local args=()
  while [ $# -gt 0 ] && [[ "$1" != "--" ]]; do
    while getopts "u:s:p:f:" opt; do
      case $opt in
      u)
        user="$OPTARG"
        ;;
      s)
        pw="sshpass -p $OPTARG"
        if ! type sshpass &>/dev/null; then
          echo "sshpass required"
          exit 1
        fi
        ;;
      p)
        port="$OPTARG"
        ;;
      f)
        cfile="$OPTARG"
        ;;
      *)
        echo -e $usage
        exit 0
        ;;
      esac
    done

    shift $((OPTIND - 1))

    while [ $# -gt 0 ] && (! [[ "$1" =~ ^- ]] || [[ "$1" =~ ^-- ]]); do
      args=("${args[@]}" "$1")
      shift
      OPTIND=1 # reset OPTIND
    done
  done

  local mode=${args[1]}
  local hosts
  local cmd
  if [ -f "$cfile" ]; then
    IFS=$'\r\n' read -d '' -r -a hosts < $cfile
    cmd= ${args[2]}
  else
    hosts=${args[2]}
    cmd=${args[3]}
  fi

  setopt shwordsplit   # Compatible with bash word split
  for host in ${hosts[*]}; do
    case $mode in
    ping)
      if /bin/ping -c 1 -w 1 $host &>/dev/null; then
        echo "$host [succ]"
      else
        echo "$host [fail]"
      fi
      ;;
    shell)
      $pw ssh -p $port $user@$host $cmd
      ;;
    script)
      $pw ssh -p $port $user@$host "bash $cmd"
      ;;
    copy)
      local src=`echo $cmd | awk '{print $1}'`
      local dest=`echo $cmd | awk '{print $2}'`
      $pw scp -r -P $port $src $user@$host:$dest
      ;;
    fetch)
      local src=`echo $cmd | awk '{print $1}'`
      local dest=`echo $cmd | awk '{print $2}'`
      $pw scp -r -P $port $user@$host:$src $dest
      ;;
    template)
      local src=`echo $cmd | awk '{print $1}'`
      local old=`echo $cmd | awk '{print $2}'`
      local new=`echo $cmd | awk '{print $3}'`
      $pw ssh -p $port $user@$host "sed -i 's/$old/$new/g' $src"
      ;;
    *)
      echo $usage
      exit 1
      ;;
    esac
  done
  unsetopt shwordsplit
}

#########################################################################
# lf终端文件管理器(类ranger)
#########################################################################
# lfcleanup() {
#   exec 3>&-
#   rm -rf "$FIFO_UEBERZUG"
# }
    
function ff() {
  if [ -n "$1" ]; then
    cd "$1"
  fi
  
  if type lf &>/dev/null &&  type ueberzug &>/dev/null; then
    if [ -n "$SSH_CLIENT" ] || [ -n "$SSH_TTY" ]; then
      lf "$@"
    else
      rm -rf "$HOME/.cache/lf"
      [ ! -d "$HOME/.cache/lf" ] && mkdir --parents "$HOME/.cache/lf"
      export FIFO_UEBERZUG="$HOME/.cache/lf/ueberzug-$RANDOM"
      mkfifo "$FIFO_UEBERZUG" &> /dev/null
      (ueberzug layer -s <"$FIFO_UEBERZUG" -p bash &) &> /dev/null
      exec 3>"$FIFO_UEBERZUG"
      # trap lfcleanup EXIT
      lf "$@" 3>&-
    fi
  fi
}

#########################################################################
# ansible输出json提取
#########################################################################
function ans() {
  local text
  local host=$1
  local offset
  host=$(echo $host|sed 's/\./\\./g')
  #shift
  #read -d '' text
  #IFS=$'\r\n' read -d '' -r -a text < $2
  if [ $# == 1 ]; then
    text=$(cat $2)
  else
    if [ -f $2 ]; then
      text=$(cat $2)
    else
      offset=$2
      text=$(cat $3)
    fi
  fi
  if [ -n "$offset" ] && type jq >/dev/null 2>&1; then
    echo "$text" | \grep -Pzo "\[?$host\]?.* => {\n(.*\n)*?}\n?"|sed "s/^\S.* => {/{/g" | jq -rsM . | jq -rM .[$offset]
  else
    echo "$text" | \grep -Pzo "\[?$host\]?.* => {\n(.*\n)*?}\n?"|sed "s/^\S.* => {/{/g"
  fi
}

#########################################################################
# kubectl自动补全加载较慢, 启用延迟加载
#########################################################################
function kubectl() {
  if ! type __start_kubectl >/dev/null 2>&1; then
    source <(command kubectl completion zsh)
  fi
  command kubectl "$@"
}
alias k="kubectl"

#########################################################################
# nvm加载较慢, 启用延迟加载
#########################################################################
#function nvm() {
#  [ -s "$HOME/.nvm/nvm.sh" ] && . "$HOME/.nvm/nvm.sh"
#  [ -s "$HOME/.nvm/bash_completion" ] && . "$HOME/.nvm/bash_completion"
#  command nvm "$@"
#}

if [ -f /usr/share/nvm/init-nvm.sh ]; then
  source /usr/share/nvm/init-nvm.sh --no-use
fi

#########################################################################
#目录跳转后自动显示目录内容
#########################################################################
_listpwd() {
  emulate -L zsh
  \ls -F --color=auto
}

chpwd_functions+=(_listpwd)
# add-zsh-hook -Uz chpwd listpwd #(){ listpwd }
# add-zsh-hook -Uz precmd (){ echo "AAA" }
# Hook Functions (http://zsh.sourceforge.net/Doc/Release/Functions.html)
# chpwd: 当前目录改变时触发
# periodic: 设置PERIOD环境变量后, 每隔PERIOD秒触发一次(输出在提示符前)
# precmd: 每次绘制提示符前触发
# preexec: 每次执行命令前触发
# zshaddhistory: 写入历史记录前执行
# zshexit: zsh退出前触发

#########################################################################
# fzf配置
#########################################################################
export FZF_DEFAULT_OPTS='--bind=ctrl-t:top,change:top --bind ctrl-e:down,ctrl-u:up'
#export FZF_DEFAULT_OPTS='--bind ctrl-e:down,ctrl-u:up --preview "[[ $(file --mime {}) =~ binary ]] && echo {} is a binary file || (ccat --color=always {} || highlight -O ansi -l {} || cat {}) 2> /dev/null | head -500"'
#export FZF_DEFAULT_COMMAND='ag --hidden --ignore .git -g ""'
export FZF_DEFAULT_COMMAND='fd'
export FZF_COMPLETION_TRIGGER='\'
export FZF_TMUX=1
export FZF_TMUX_HEIGHT='80%'
export fzf_preview_cmd='[[ $(file --mime {}) =~ binary ]] && echo {} is a binary file || highlight -O ansi -l {} || (ccat --color=always {} || cat {}) 2> /dev/null | head -500'

_fzf_fpath=${0:h}/fzf
fpath+=$_fzf_fpath
autoload -U $_fzf_fpath/*(.:t)
unset _fzf_fpath

_fzf_compgen_path() {
  fd --hidden --follow --exclude ".git" . "$1"
}

# Use fd to generate the list for directory completion
_fzf_compgen_dir() {
  fd --type d --hidden --follow --exclude ".git" . "$1"
}

# (EXPERIMENTAL) Advanced customization of fzf options via _fzf_comprun function
# - The first argument to the function is the name of the command.
# - You should make sure to pass the rest of the arguments to fzf.
_fzf_comprun() {
  local command=$1
  shift

  case "$command" in
    cd)           fzf "$@" --preview 'tree -C {} | head -200' ;;
    export|unset) fzf "$@" --preview "eval 'echo \$'{}" ;;
    ssh)          fzf "$@" --preview 'dig {}' ;;
    *)            fzf "$@" ;;
  esac
}

fzf-redraw-prompt() {
	local precmd
	for precmd in $precmd_functions; do
		$precmd
	done
	zle reset-prompt
}
zle -N fzf-redraw-prompt

zle -N fzf-find-widget
bindkey '^p' fzf-find-widget

fzf-cd-widget() {
	local tokens=(${(z)LBUFFER})
	if (( $#tokens <= 1 )); then
		zle fzf-find-widget 'only_dir'
		if [[ -d $LBUFFER ]]; then
			cd $LBUFFER
			local ret=$?
			LBUFFER=
			zle fzf-redraw-prompt
			return $ret
		fi
	fi
}
zle -N fzf-cd-widget
bindkey '^t' fzf-cd-widget

fzf-history-widget() {
	local num=$(fhistory $LBUFFER)
	local ret=$?
	if [[ -n $num ]]; then
		zle vi-fetch-history -n $num
	fi
	zle reset-prompt
	return $ret
}
zle -N fzf-history-widget
bindkey '^R' fzf-history-widget

fif() {
  if [ ! "$#" -gt 0 ]; then echo "Need a string to search for!"; return 1; fi
  rg --files-with-matches --no-messages "$1" | fzf --preview "highlight -O ansi -l {} 2> /dev/null | rg --colors 'match:bg:yellow' --ignore-case --pretty --context 10 '$1' || rg --ignore-case --pretty --context 10 '$1' {}"
}

find-in-file() {
	grep --line-buffered --color=never -r "" * | fzf
}
zle -N find-in-file
bindkey '^f' find-in-file

#########################################################################
# 显示命令执行时间
#########################################################################
# If command execution time above min. time, plugins will not output time.
ZSH_COMMAND_TIME_MIN_SECONDS=5

# Message to display (set to "" for disable).
ZSH_COMMAND_TIME_MSG="Execution time: %s sec"

# Message color.
ZSH_COMMAND_TIME_COLOR="cyan"

_command_time_preexec() {
  timer=${timer:-$SECONDS}
  ZSH_COMMAND_TIME_MSG=${ZSH_COMMAND_TIME_MSG-"Time: %s"}
  ZSH_COMMAND_TIME_COLOR=${ZSH_COMMAND_TIME_COLOR-"white"}
  export ZSH_COMMAND_TIME=""
}

_command_time_precmd() {
  if [ $timer ]; then
    timer_show=$(($SECONDS - $timer))
    if [ -n "$TTY" ] && [ $timer_show -ge ${ZSH_COMMAND_TIME_MIN_SECONDS:-3} ]; then
      export ZSH_COMMAND_TIME="$timer_show"
      if [ ! -z ${ZSH_COMMAND_TIME_MSG} ]; then
        zsh_command_time
      fi
    fi
    unset timer
  fi
}

zsh_command_time() {
  if [ -n "$ZSH_COMMAND_TIME" ]; then
    timer_show=$(printf '%dh:%02dm:%02ds\n' $(($ZSH_COMMAND_TIME/3600)) $(($ZSH_COMMAND_TIME%3600/60)) $(($ZSH_COMMAND_TIME%60)))
    print -P "%F{$ZSH_COMMAND_TIME_COLOR}$(printf "${ZSH_COMMAND_TIME_MSG}\n" "$timer_show")%f"
  fi
}

precmd_functions+=(_command_time_precmd)
preexec_functions+=(_command_time_preexec)

#使用:<<' 注释内容 ' 的形式注释掉Incremental completion插件，删除空行的单引号可以重新启用
:<<'
# Incremental completion for zsh
# by y.fujii <y-fujii at mimosa-pudica.net>, public domain

autoload -U compinit
zle -N self-insert self-insert-incr
zle -N vi-cmd-mode-incr
zle -N vi-backward-delete-char-incr
zle -N backward-delete-char-incr
zle -N expand-or-complete-prefix-incr
compinit

bindkey -M viins '^[' vi-cmd-mode-incr
bindkey -M viins '^h' vi-backward-delete-char-incr
bindkey -M viins '^?' vi-backward-delete-char-incr
bindkey -M viins '^i' expand-or-complete-prefix-incr
bindkey -M emacs '^h' backward-delete-char-incr
bindkey -M emacs '^?' backward-delete-char-incr
bindkey -M emacs '^i' expand-or-complete-prefix-incr

unsetopt automenu
compdef -d scp
compdef -d tar
compdef -d make
compdef -d java
compdef -d svn
compdef -d cvs

# TODO:
#     cp dir/

now_predict=0

function limit-completion
{
	if ((compstate[nmatches] <= 1)); then
		zle -M ""
#     elif ((compstate[list_lines] > 6)); then
#		 compstate[list]=""
#        zle -M "too many matches."
	fi
}

function correct-prediction
{
	if ((now_predict == 1)); then
		if [[ "$BUFFER" != "$buffer_prd" ]] || ((CURSOR != cursor_org)); then
			now_predict=0
		fi
	fi
	recolor-cmd
}

function remove-prediction
{
	if ((now_predict == 1)); then
		BUFFER="$buffer_org"
		now_predict=0
	fi
	recolor-cmd
}

function show-prediction
{
	#assert(now_predict == 0)
	if
		((PENDING == 0)) &&
		((CURSOR > 1)) &&
		[[ "$PREBUFFER" == "" ]] &&
		[[ "$BUFFER[CURSOR]" != " " ]]
	then
		cursor_org="$CURSOR"
		buffer_org="$BUFFER"
		#comppostfuncs=(limit-completion)
		zle complete-word
		cursor_prd="$CURSOR"
		buffer_prd="$BUFFER"
		if [[ "$buffer_org[1,cursor_org]" == "$buffer_prd[1,cursor_org]" ]]; then
			CURSOR="$cursor_org"
			if [[ "$buffer_org" != "$buffer_prd" ]] || ((cursor_org != cursor_prd)); then
				now_predict=1
			fi
		else
			BUFFER="$buffer_org"
			CURSOR="$cursor_org"
		fi
		echo -n "\e[32m"
	else
		zle -M ""
	fi
}

function preexec
{
	echo -n "\e[39m"
}

function vi-cmd-mode-incr
{
	correct-prediction
	remove-prediction
	zle vi-cmd-mode
}

function self-insert-incr
{
	correct-prediction
	remove-prediction
	if zle .self-insert; then
		show-prediction
	fi
	recolor-cmd
}

function vi-backward-delete-char-incr
{
	correct-prediction
	remove-prediction
	if zle vi-backward-delete-char; then
		show-prediction
	fi
}

function backward-delete-char-incr
{
	correct-prediction
	remove-prediction
	if zle backward-delete-char; then
		show-prediction
	fi
	recolor-cmd
}

function expand-or-complete-prefix-incr
{
	correct-prediction
	if ((now_predict == 1)); then
		CURSOR="$cursor_prd"
		now_predict=0
		comppostfuncs=(limit-completion)
		zle list-choices
	else
		remove-prediction
		zle expand-or-complete-prefix
	fi
	recolor-cmd
}
'
