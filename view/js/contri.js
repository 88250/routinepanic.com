var editor = CodeMirror.fromTextArea(document.getElementById('content'), {
  mode: 'text/html',
  autoCloseTags: true,
  lineNumbers: true,
  lineWrapping: true,
  foldGutter: true,
  gutters: ['CodeMirror-linenumbers', 'CodeMirror-foldgutter'],
})

// editor.foldCode(CodeMirror.Pos(0, 0));
CodeMirror.commands['selectAll'](editor)
editor.autoFormatRange(editor.getCursor(true), editor.getCursor(false))
CodeMirror.commands['goDocStart'](editor)