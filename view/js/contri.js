let editor

const preview = () => {
  const preview = document.getElementById('contriPreview')
  if (preview.className.indexOf('contri__preview--active') > -1) {
    preview.className = 'contri__preview content-reset'
    preview.innerHTML = ''
  } else {
    $('#dictionary').html('')
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
    if ($.trim(cm.doc.getSelection()) === '') {
      $('#dictionary').html('')
      return
    }
    $.ajax({
      url: '/words/' + cm.doc.getSelection(),
      success: function (result) {
        if (result.code !== 0) {
          $('#dictionary').html('')
          return
        }
        let dicHTML = ''
        if (result.data) {
          dicHTML = `<div class="contri__dic"><div class="fn-flex-center"><span>${result.data.name}&nbsp;&nbsp;</span>`
          if (result.data.phAm) {
            dicHTML += `<span class="ft-fade ft-12">[${result.data.phAm}]&nbsp;&nbsp;</span>`
          }
          if (result.data.phAmMP3) {
            dicHTML += `<svg class="fn-pointer ft-gray contri__vice"><use xlink:href="#iconVice"></use></svg>
  <audio>
    <source = src="${result.data.phAmMP3}" type="audio/mp3">
  </audio>`
          }
          dicHTML += `</div><div class="ft-12">${result.data.means.replace(/\n/g, '<br>')}</div></div>`
        }
        $('#dictionary').html(dicHTML)
      },
    })
  })
  $('#dictionary').on('mouseover', 'svg', function () {
    $('#dictionary audio')[0].play()
  })
}

initEditor()


