let editor

const preview = () => {
  const preview = document.getElementById('contriPreview')
  if (preview.className.indexOf('contri__preview--active') > -1) {
    preview.className = 'contri__preview content-reset'
    preview.innerHTML = ''
  } else {
    preview.className = 'contri__preview--active contri__preview content-reset'
    setTimeout(() => {
      preview.innerHTML = editor.edit.getValue()
    }, 150)
  }
}

const initEditor = () => {
  // format zh
  const editorZh = CodeMirror.fromTextArea(document.getElementById('contentZh'),
    {
      mode: 'text/html',
    })
  CodeMirror.commands['selectAll'](editorZh)
  editorZh.autoFormatRange(editorZh.getCursor(true), editorZh.getCursor(false))

  // format en
  const editorEn = CodeMirror.fromTextArea(document.getElementById('contentEn'),
    {
      mode: 'text/html',
    })
  CodeMirror.commands['selectAll'](editorEn)
  editorEn.autoFormatRange(editorEn.getCursor(true), editorEn.getCursor(false))

  // init content
  const contentHeight = window.innerHeight -
    (document.querySelector('.contri__diff') ? 313 : 268)
  editor = CodeMirror.MergeView(document.getElementById('content'), {
    autoCloseTags: true,
    lineNumbers: true,
    lineWrapping: true,
    foldGutter: true,
    mode: 'text/html',
    gutters: ['CodeMirror-linenumbers', 'CodeMirror-foldgutter'],
    value: editorZh.getValue(),
    origLeft: editorEn.getValue(),
    connect: 'align',
    showDifferences: false,
    height: contentHeight,
  })

  // resize
  editor.edit.setSize('100%', contentHeight)
  editor.leftOriginal().setSize('100%', contentHeight)
  document.querySelector('.CodeMirror-merge').style.height = contentHeight +
    'px'

  // preview
  editor.edit.on('change', function (cm) {
    document.getElementById('contentValue').value = cm.getValue()
    const preview = document.getElementById('contriPreview')
    if (preview.className.indexOf('contri__preview--active') > -1) {
      preview.innerHTML = editor.edit.getValue()
    }
  })

  // dictionary
  editor.leftOriginal().on('dblclick', function (cm) {

    console.log(cm, cm.doc.getSelection())
  })
}

initEditor()


