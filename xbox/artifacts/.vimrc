" debug vim
" vim --noplugin -D
" disable all pulgins
" vim -U NONE -u

" 设置.vim文件路径
" vim -u ~/.vimrc  " 设置vim启动时使用的配置文件
set runtimepath=~/.vim,$VIM/vimfiles,$VIMRUNTIME,$VIM/vimfiles/after,~/.vim/after
if version >=800
  set packpath=~/.vim,$VIM/vimfiles,$VIMRUNTIME,$VIM/vimfiles/after,~/.vim/after
endif

set encoding=utf-8
scriptencoding utf-8
" 识别文件编码
set fileencodings=utf-8,ucs-bom,shift-jis,gb18030,gbk,gb2312,cp936
set termencoding=utf-8
set fileencodings=ucs-bom,utf-8,cp936
set fileencoding=utf-8
" 显示中文帮助
set helplang=cn
" 语言设置
set langmenu=zh_CN.UTF-8
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" vim主题站点：http://bytefluent.com/vivify/
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" 显示相关
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
set t_Co=256
set nocompatible  "去掉讨厌的有关vi一致性模式，避免以前版本的一些bug和局限
" 自定义高亮当前行
set cursorline
"hi CursorLine  cterm=NONE   ctermbg=darkred ctermfg=white
"hi CursorColumn cterm=NONE ctermbg=darkred ctermfg=white
"set cuc      " 高亮列
set shortmess=atI   " 启动的时候不显示那个援助乌干达儿童的提示
"set go=             " 不要图形按钮
"color desert     " 设置背景主题
"color maroloccio     " 设置背景主题
"color torte     " 设置背景主题
"colorscheme murphy
"if has('gui_running')
"    set background=light
"else
"    set background=dark
"endif
set guifont=DejaVu\ Sans\ mono\ 11   " 设置字体
"autocmd InsertLeave * se nocul  " 用浅色高亮当前行
"autocmd InsertEnter * se cul    " 用浅色高亮当前行
"set whichwrap+=<,>,h,l   " 允许backspace和光标键跨越行边界(不建议)
set scrolloff=2     " 光标移动到buffer的顶部和底部时保持2行距离
set laststatus=2    " 启动显示状态行(1),总是显示状态行(2)
"set foldenable      " 允许折叠
"set foldmethod=manual   " 手动折叠
" 侦测文件类型
filetype on
" 载入文件类型插件
filetype plugin on
" 为特定文件类型载入相关缩进文件
filetype indent on
" 开启语法检查
syntax enable
syntax on
" 自动缩进
set autoindent
set cindent
" Tab键的宽度
set tabstop=2
" 统一缩进为2
set softtabstop=2
set shiftwidth=2
" 用空格代替制表符
set expandtab
" 在行和段开始处使用制表符
set smarttab
" 显示行号
set number
" set nu
" 显示相对行号
" set relativenumber
" 历史记录数
set history=10000
" 搜索逐字符高亮
set hlsearch
set incsearch
" 命令行高度
set cmdheight=1
" 代码补全
set completeopt=preview,menu
"set completeopt=longest,menu
"set completeopt=menu
" 允许不保存buffer时切换
set hidden
" 共享剪贴板
set clipboard+=unnamed
" 自动保存
set autowrite
"set ruler                   " 打开状态栏标尺
"set cursorline              " 突出显示当前行
set magic                    " 设置魔术
"set guioptions-=T           " 隐藏工具栏
"set guioptions-=m           " 隐藏菜单栏
"set foldcolumn=0
set foldmethod=indent
set nofoldenable
"set foldlevel=3
" 去掉输入错误的提示声音
set noerrorbells
" 去掉按Enter键或其他命令以继续提示
set shortmess=atI
" 在处理未保存或只读文件的时候，弹出确认
set confirm
" 禁止生成临时文件
set nobackup
set noswapfile
set noundofile
" 搜索忽略大小写
set ignorecase
set smartcase
set linespace=0
" 增强模式中的命令行自动完成操作
set wildmenu
" 使回格键（backspace）正常处理indent, eol, start等
set backspace=2
" 允许backspace和光标键跨越行边界
set whichwrap+=<,>,h,l
" 通过使用: commands命令，告诉我们文件的哪一行被改变过
set report=0
" 在被分割的窗口间显示空白，便于阅读
set fillchars=vert:\ ,stl:\ ,stlnc:\
" 高亮显示匹配的括号
set showmatch
" 匹配括号高亮的时间（单位是十分之一秒）
set matchtime=1
" 设置映射延迟和按键延迟
set timeoutlen=500
set ttimeoutlen=10
" 保存全局变量
set viminfo+=!
" 带有如下符号的单词不要被换行分割
set iskeyword+=_,$,@,%,#,-
" 字符间插入的像素行数目
set ruler           " 显示标尺
set showcmd         " 输入的命令显示出来，看的清楚些
set showmode        " 显示当前模式
" 将tab替换为空格
"nmap tt :%s/\t/    /g<CR>
" 仅需要时重绘
"set lazyredraw

""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
"键盘命令
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" 自定义老板键(默认是\)
let mapleader="\\"
" 自定义esc键(默认是<esc>和<C-[>)
" inoremap jo <esc>

" 可以在buffer的任何地方使用鼠标（类似office中在工作区双击鼠标定位）
" 支持滚轮
set mousemodel=popup
set mouse-=a
set selection=exclusive
set selectmode=mouse,key
" 输入法
:let g:vimim_map='c-/'
":let g:vimim_cloud='sougou' " QQ云输入
:let g:vimim_punctuation=0  " 不用中文标点
:let g:vimim_cloud=-1

" 编辑模式下映射方向键使hjkl能移动光标
inoremap <silent><C-H> <Left>
inoremap <silent><C-L> <Right>
inoremap <silent><C-J> <Down>
inoremap <silent><C-K> <Up>

" 大于小于号缩进
xnoremap < <gv
xnoremap > >gv

" 复制粘贴映射
map <C-y> "+y
map <C-p> "+gp

" 禁用箭头按键
"nnoremap <left> <Nop>
"nnoremap <right> <Nop>
"nnoremap <up> <Nop>
"nnoremap <down> <Nop>
"inoremap <left> <Nop>
"inoremap <right> <Nop>
"inoremap <up> <Nop>
"inoremap <down> <Nop>
"vnoremap <left> <Nop>
"vnoremap <right> <Nop>
"vnoremap <up> <Nop>
"vnoremap <down> <Nop>

" 外部粘贴无需手动切换paste模式
"nnoremap :set invpaste paste? imap :set invpaste paste? set pastetoggle=
"set pastetoggle=
"==========
"nnoremap <F2> :set invpaste paste?<CR>
"imap <F2> <C-O>:set invpaste paste?<CR>
"set pastetoggle=<F2>
"==========
let &t_SI .= "\<Esc>[?2004h"
let &t_EI .= "\<Esc>[?2004l"

inoremap <special> <expr> <Esc>[200~ XTermPasteBegin()

function! XTermPasteBegin()
  set pastetoggle=<Esc>[201~
  set paste
  return ""
endfunction

" 自定义statusline函数(ale插件实现)
function! LinterStatus() abort
    let l:counts = ale#statusline#Count(bufnr(''))

    let l:all_errors = l:counts.error + l:counts.style_error
    let l:all_non_errors = l:counts.total - l:all_errors

    return l:counts.total == 0 ? 'OK' : printf(
    \   '%dW %dE',
    \   all_non_errors,
    \   all_errors
    \)
endfunction

" 自定义状态栏
"set statusline=%F%m%r%h%w\ [FORMAT=%{&ff}]\ [TYPE=%Y]\ [ASCII=/%03.3b]\ [POS=%l,%v]\ [LEN=%L]\ [%p%%]\ %{strftime(\"%d/%m/%y\ -\ %H:%M\")}\ %{LinterStatus()}   "状态行显示的内容
set statusline=%F%m%r%h%w\ [FORMAT=%{&ff}]\ [TYPE=%Y]\ [ASCII=/%03.3b]\ [POS=%l,%v]\ [LEN=%L]\ [%p%%]\ %{strftime(\"%d/%m/%y\ -\ %H:%M\")}   "状态行显示的内容


" 设置当文件被改动时自动载入
set autoread
augroup filetype_cpp
    autocmd!
    autocmd FileType c,cpp map <buffer> <leader><space> :w<cr>:make<cr>
augroup END
augroup filetype_dict
    autocmd!
    autocmd FileType php set dict+=~/.vim/dict/php_funclist.dict
    autocmd FileType css set dict+=~/.vim/dict/css.dict
    autocmd FileType c set dict+=~/.vim/dict/c.dict
    autocmd FileType cpp set dict+=~/.vim/dict/cpp.dict
    autocmd FileType scale set dict+=~/.vim/dict/scale.dict
    autocmd FileType javascript set dict+=~/.vim/dict/javascript.dict
    autocmd FileType html set dict+=~/.vim/dict/javascript.dict
    autocmd FileType html set dict+=~/.vim/dict/css.dict
augroup END

"修改文件为unix格式
map <C-!> :call FormatClean()<CR>
func! FormatClean()
    exec 'set fileformat=unix'
    exec 'w'
endfunc

" 统一缩进
"nnoremap <F2> gg=G<CR>
" 设置代码折叠
"map <F2> :call CodeFold()<CR>
"func! CodeFold()
    "exec 'set foldmethod=indent'
    "zC
    "zO
"endfunc

" 去空行
nnoremap <F7> :g/^\s*$/d<CR>

"taglist
":nmap <silent> <F4> <ESC>:Tlist<RETURN>

"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
"""""按F5编译运行
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
map <F5> :call CompileRun()<CR>
func! CompileRun()
    exec 'w'
    if search('mpi\.h') != 0
        exec '!mpirun -np 4 ./%<'
    elseif &filetype ==? 'c'
        let l:compilecmd='!gcc '
        let l:compileflag='-o %< '
        if search('mpi\.h') != 0
            let l:compilecmd = '!mpicc '
        endif
        if search('glut\.h') != 0
            let l:compileflag .= ' -lglut -lGLU -lGL '
        endif
        if search('cv\.h') != 0
            let l:compileflag .= ' -lcv -lhighgui -lcvaux '
        endif
        if search('omp\.h') != 0
            let l:compileflag .= ' -fopenmp '
        endif
        if search('math\.h') != 0
            let l:compileflag .= ' -lm '
        endif
        exec l:compilecmd.' % '.l:compileflag
        :!time ./%<
        :!rm -f ./%<
    elseif &filetype ==? 'cpp'
        let l:compilecmd='!g++ '
        let l:compileflag='-o %< '
        if search('mpi\.h') != 0
            let l:compilecmd = '!mpic++ '
        endif
        if search('glut\.h') != 0
            let l:compileflag .= ' -lglut -lGLU -lGL '
        endif
        if search('cv\.h') != 0
            let l:compileflag .= ' -lcv -lhighgui -lcvaux '
        endif
        if search('omp\.h') != 0
            let l:compileflag .= ' -fopenmp '
        endif
        if search('math\.h') != 0
            let l:compileflag .= ' -lm '
        endif
        exec l:compilecmd.' % '.l:compileflag
        :!time ./%<
        :!rm -f ./%<
    elseif &filetype ==? 'java'
        :!javac %
        :!time java %<
    elseif &filetype ==? 'sh'
        :!time bash %
    elseif &filetype ==? 'python'
        :!time python3 %
    elseif &filetype ==? 'javascript'
        :!time node %
    elseif &filetype ==? 'typescript'
        :!time tsc %
    elseif &filetype ==? 'html'
        :!firefox % &
    elseif &filetype ==? 'php'
        :!time php -f %
    elseif &filetype ==? 'lua'
        :!time lua %
    elseif &filetype ==? 'go'
        "":!go build %<
        :!time go run %
    elseif &filetype ==? 'julia'
        :!time julia %
    elseif &filetype ==? 'ruby'
        :!time ruby %
    elseif &filetype ==? 'rust'
        :!time rust %
    elseif &filetype ==? 'mkd'
        :!~/.vim/markdown.pl % > %.html &
        :!firefox %.html &
    endif
endfunc

""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
"""""""C,C++的调试
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
map <F6> :call Rungdb()<CR>
func! Rungdb()
    exec 'w'
    exec '!g++ % -g -o %<'
    exec '!gdb ./%<'
endfunc

""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
"""""""代码格式化
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
map <F4> :call FormatSrc()<CR>

" 定义FormartSrc()
func! FormatSrc()
    exec 'w'
    if &filetype ==? 'c'
        exec '!astyle --style=ansi -a --suffix=none %'
    elseif &filetype ==? 'cpp' || &filetype ==? 'hpp'
        exec 'r !astyle --style=ansi --one-line=keep-statements -a --suffix=none %> /dev/null 2>&1'
    elseif &filetype ==? 'perl'
        exec '!astyle --style=gnu --suffix=none %'
    elseif &filetype ==? 'py'||&filetype ==? 'python'
        exec 'r !autopep8 -i --aggressive %'
    elseif &filetype ==? 'java'
        exec '!astyle --style=java --suffix=none %'
    elseif &filetype ==? 'jsp'
        exec '!astyle --style=gnu --suffix=none %'
    elseif &filetype ==? 'xml'
        exec '!astyle --style=gnu --suffix=none %'
    else
        exec 'normal gg=G'
        return
    endif
    exec 'e! %'
endfunc

""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
""实用设置
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
if has('autocmd')
    augroup ac
        autocmd!
        autocmd BufReadPost * if line("'\"") > 0 && line("'\"") <= line("$") | exe "normal g`\"" | endif
    augroup END
endif

""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
""切换行号显示
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
nmap <F3> :call ToggleLineNumber()<CR>
func! ToggleLineNumber()
  if !exists("b:togglenum")
    let b:togglenum=1
  endif
  if b:togglenum==0
    execute "set number"
    execute "set norelativenumber"
    let b:togglenum=1
  else
    if b:togglenum==1
      execute "set number"
      execute "set relativenumber"
      let b:togglenum=2
    else
      execute "set nonumber"
      execute "set norelativenumber"
      let b:togglenum=0
    endif
  endif
endfunc

""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
""解决输入法问题
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
let g:input_toggle = 0          " 设置默认为英文输入法
function! Fcitx2en()
   let s:input_status = system("fcitx5-remote")
   if s:input_status == 2
      let g:input_toggle = 1
      let l:a = system("fcitx5-remote -c")
   endif
endfunction

function! Fcitx2zh()
   let s:input_status = system("fcitx5-remote")
   if s:input_status != 2 && g:input_toggle == 1
      let l:a = system("fcitx5-remote -o")
      let g:input_toggle = 0
   endif
endfunction

set timeoutlen=150
autocmd InsertLeave * call Fcitx2en()    " 退出Insert模式时关闭中文输入法
autocmd InsertEnter * call Fcitx2zh()   " 进入Insert模式时开启中文输入法

" 为C程序提供自动缩进
" 自动补全
"":inoremap ( ()<ESC>i
"":inoremap ) <c-r>=ClosePair(')')<CR>
":inoremap { {<CR>}<ESC>O
":inoremap } <c-r>=ClosePair('}')<CR>
"":inoremap [ []<ESC>i
"":inoremap ] <c-r>=ClosePair(']')<CR>
"":inoremap " ""<ESC>i
"":inoremap ' ''<ESC>i
""function! ClosePair(char)
""  if getline('.')[col('.') - 1] == a:char
""      return "\<Right>"
""  else
""      return a:char
""  endif
""endfunction


" 能够漂亮地显示.NFO文件
function! SetFileEncodings(encodings)
    let b:myfileencodingsbak=&fileencodings
    let &fileencodings=a:encodings
endfunction
function! RestoreFileEncodings()
    let &fileencodings=b:myfileencodingsbak
    unlet b:myfileencodingsbak
endfunction
augroup auc
    autocmd!
    autocmd BufReadPre *.nfo call SetFileEncodings('cp437')|setambiwidth=single autocmd BufReadPost *.nfo callRestoreFileEncodings()
augroup END

nmap <silent>di :!curl dict://dict.org/d:<cword><CR><CR>

" execute a command and show its output in a split window
command! -nargs=* -complete=shellcmd Rsplit execute "new | r! <args>"

" execute a command and show its output in a new tab
command! -nargs=* -complete=shellcmd Rtab execute "tabnew | r! <args>"

" :new | 0read ! <command>
