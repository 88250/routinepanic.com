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
  })

  // resize
  $(window).resize(function () {
    const contentHeight = $(window).height() -
      (document.querySelector('.contri__diff') ? 313 : 268)
    editor.edit.setSize('100%', contentHeight)
    editor.leftOriginal().setSize('100%', contentHeight)
    $('.CodeMirror-merge').height(contentHeight)
    $('#contriPreview').outerHeight(contentHeight + 2)
  })
  $(window).resize()

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
    $.ajax({
      url: '/words/' + cm.doc.getSelection(),
      success: function (result) {
        if (result.code !== 0) {
          alert(result.msg)
          return
        }
        console.log(result)
        let dicHTML = ''
        if (result.data) {
          dicHTML = `<div class="fn-flex-center">
  <span>${result.data.name}&nbsp;&nbsp;</span>
  <span class="ft-fade ft-12">[${result.data.phAm}]&nbsp;&nbsp;</span>
  <svg class="fn-pointer ft-gray contri__vice"><use xlink:href="#iconVice"></use></svg>
</div>
<div class="ft-12 ft-fade">${result.data.means}</div>`
        }
        $('#dictionary').html(dicHTML)
      },
    })
  })
}

initEditor()


