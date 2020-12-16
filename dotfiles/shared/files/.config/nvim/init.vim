call plug#begin(stdpath('data') . '/plugged')

Plug 'morhetz/gruvbox'                              " Gruvbox theme
Plug 'neoclide/coc.nvim', {'branch': 'release'}     " coc.nvim, lsp client
Plug 'junegunn/fzf', { 'do': { -> fzf#install() } }
Plug 'junegunn/fzf.vim'
Plug 'fatih/vim-go'
Plug 'jackguo380/vim-lsp-cxx-highlight'             " semantic highlighting for c/c++

call plug#end()

"
" Indentation
"

set smartindent                 " Makes indenting 'smart'
set autoindent                  " Auto indent
set shiftwidth=4		        " indentation is 4 spaces wide
set tabstop=4			        " tabs are 4 spaces wide
set softtabstop=4		        " tabs are 4 spaces wide
set expandtab			        " tabs are expanded into spaces

"
" Appearence
"

set background=dark
let g:gruvbox_contrast_dark='hard'
set termguicolors
colorscheme gruvbox
syntax enable
set ruler                       " show location information
set cmdheight=2                 " more space for displaying messages
set signcolumn=yes              " always show signcolumn so it doesn't keep popping in and out
set number
set scrolloff=5                 " Keep 5 lines at top/bottom of screen


"
" KEYS
"

" use space as leader key
let mapleader=" "

" use jj to exit insert mode
inoremap jj <esc>

" Use alt + hjkl to resize windows
nnoremap <M-j> :resize -2<CR>
nnoremap <M-k> :resize +2<CR>
nnoremap <M-h> :vertical resize -2<CR>
nnoremap <M-l> :vertical resize +2<CR>

" tab completion for coc.nvim
inoremap <silent><expr> <TAB>
      \ pumvisible() ? "\<C-n>" : "\<TAB>"
      "\ <SID>check_back_space() ? "\<TAB>" :
      "\ coc#refresh()
inoremap <expr><S-TAB> pumvisible() ? "\<C-p>" : "\<C-h>"
" use ctrl-k and ctrl-j to move through completions
inoremap <expr><C-k> pumvisible() ? "\<C-p>" : "\<C-h>"
inoremap <expr><C-j> pumvisible() ? "\<C-n>" : "\<C-h>"

function! s:check_back_space() abort
  let col = col('.') - 1
  return !col || getline('.')[col - 1]  =~# '\s'
endfunction

" use <cr> to confirm completion
if exists('*complete_info')
  inoremap <expr> <cr> complete_info()["selected"] != "-1" ? "\<C-y>" : "\<C-g>u\<CR>"
else
  inoremap <expr> <cr> pumvisible() ? "\<C-y>" : "\<C-g>u\<CR>"
endif

" Use `[g` and `]g` to navigate diagnostics
nmap <silent> [g <Plug>(coc-diagnostic-prev)
nmap <silent> ]g <Plug>(coc-diagnostic-next)

" GoTo code navigation.
nmap <silent> gd <Plug>(coc-definition)
nmap <silent> gy <Plug>(coc-type-definition)
nmap <silent> gi <Plug>(coc-implementation)
nmap <silent> gr <Plug>(coc-references)

" Use K to show documentation in preview window.
nnoremap <silent> K :call <SID>show_documentation()<CR>

function! s:show_documentation()
    if (index(['vim','help'], &filetype) >= 0)
        execute 'h '.expand('<cword>')
    else
        call CocAction('doHover')
    endif
endfunction

" Highlight the symbol and its references when holding the cursor.
autocmd CursorHold * silent call CocActionAsync('highlight')

" Symbol renaming.
nmap <leader>rn <Plug>(coc-rename)
" Formatting selected code.
xmap <leader>f  <Plug>(coc-format-selected)
nmap <leader>f  <Plug>(coc-format-selected)

augroup mygroup
  autocmd!
  " Setup formatexpr specified filetype(s).
  autocmd FileType typescript,json setl formatexpr=CocAction('formatSelected')
  " Update signature help on jump placeholder.
  autocmd User CocJumpPlaceholder call CocActionAsync('showSignatureHelp')
augroup end

" Applying codeAction to the selected region.
xmap <leader>a  <Plug>(coc-codeaction-selected)
nmap <leader>a  <Plug>(coc-codeaction-selected)

" Remap keys for applying codeAction to the current line.
nmap <leader>ac  <Plug>(coc-codeaction)
" Apply AutoFix to problem on the current line.
nmap <leader>qf  <Plug>(coc-fix-current)

" Map function and class text objects
xmap if <Plug>(coc-funcobj-i)
omap if <Plug>(coc-funcobj-i)
xmap af <Plug>(coc-funcobj-a)
omap af <Plug>(coc-funcobj-a)
xmap ic <Plug>(coc-classobj-i)
omap ic <Plug>(coc-classobj-i)
xmap ac <Plug>(coc-classobj-a)
omap ac <Plug>(coc-classobj-a)

" Use CTRL-S for selections ranges.

" Requires 'textDocument/selectionRange' support of LS, ex: coc-tsserver
nmap <silent> <C-s> <Plug>(coc-range-select)
xmap <silent> <C-s> <Plug>(coc-range-select)

" Add `:Format` command to format current buffer.
command! -nargs=0 Format :call CocAction('format')

" Add `:Fold` command to fold current buffer.
command! -nargs=? Fold :call     CocAction('fold', <f-args>)

" Add `:OR` command for organize imports of the current buffer.
command! -nargs=0 OR   :call     CocAction('runCommand', 'editor.action.organizeImport')

" Add (Neo)Vim's native statusline support.
" NOTE: Please see `:h coc-status` for integrations with external plugins that
" provide custom statusline: lightline.vim, vim-airline.
set statusline^=%{coc#status()}%{get(b:,'coc_current_function','')}

" Mappings using CoCList:
"" Show all diagnostics.
"nnoremap <silent> <space>a  :<C-u>CocList diagnostics<cr>
"" Manage extensions.
"nnoremap <silent> <space>e  :<C-u>CocList extensions<cr>
"" Show commands.
"nnoremap <silent> <space>c  :<C-u>CocList commands<cr>
"" Find symbol of current document.
"nnoremap <silent> <space>o  :<C-u>CocList outline<cr>
"" Search workspace symbols.
"nnoremap <silent> <space>s  :<C-u>CocList -I symbols<cr>
"" Do default action for next item.
"nnoremap <silent> <space>j  :<C-u>CocNext<CR>
"" Do default action for previous item.
"nnoremap <silent> <space>k  :<C-u>CocPrev<CR>
"" Resume latest coc list.
"nnoremap <silent> <space>p  :<C-u>CocListResume<CR>

let g:coc_global_extensions = ['coc-clangd', 'coc-css', 'coc-go', 'coc-html', 'coc-json', 'coc-phpls', 'coc-tsserver']

" FZF
let g:fzf_action = {'ctrl-t': 'tab split', 'ctrl-x': 'split', 'ctrl-v': 'vsplit' }
let g:fzf_layout = {'window': {'width': 0.8, 'height': 0.8}}
nnoremap <leader>p :Files<CR>
nnoremap <leader>g :Rg 

"
" Misc
"

set mouse=a
set hidden                      " allow edited buffers to be hidden
set updatetime=300
set shortmess+=c                " don't pass messages to |ins-completion-menu|
set clipboard=unnamedplus       " use system clipboard

"
" Vim-go
"

" Disable go def in vim go as this will be handled by coc
let g:go_def_mapping_enabled = 0
" Use goimports to format so imports are taken care of
let g:go_fmt_command = "goimports"
