let editor

const _getPreview = (preview) => {
  $.ajax({
    url: config.server + '/html',
    type: 'POST',
    data: JSON.stringify({html: editor.edit.getValue()}),
    success: function (result) {
      preview.innerHTML = result.data
      $('pre code').each(function (i, block) {
        hljs.highlightBlock(block)
      })
    },
  })
}

const preview = () => {
  const preview = document.getElementById('contriPreview')
  if (preview.className.indexOf('contri__preview--active') > -1) {
    preview.className = 'contri__preview content-reset'
    preview.innerHTML = ''
  } else {
    $('#dictionary').html('')
    preview.className = 'contri__preview--active contri__preview content-reset'
    _getPreview(preview)
  }
}

const initEditor = (type) => {
  // init content
  editor = CodeMirror.MergeView(document.getElementById('content'), {
    autoCloseTags: true,
    lineNumbers: true,
    lineWrapping: true,
    foldGutter: true,
    mode: 'text/html',
    gutters: ['CodeMirror-linenumbers', 'CodeMirror-foldgutter'],
    value: html_beautify($('#contentZh').val()),
    origLeft: html_beautify($('#contentEn').val()),
    showDifferences: type === 'review',
    collapseIdentical: type === 'review',
  })

  // resize
  $(window).resize(function () {
    let otherHeight = (document.querySelector('.contri__diff') ? 268 : 223)
    if (type === 'review') {
      otherHeight = (document.querySelector('.contri__diff') ? 318 : 273)
    }
    const contentHeight = $(window).height() - otherHeight
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
      _getPreview(preview)
    }
  })

  // dictionary
  editor.leftOriginal().on('dblclick', function (cm) {
    if ($.trim(cm.doc.getSelection()) === '') {
      $('#dictionary').html('')
      return
    }
    $.ajax({
      url: config.server + '/words/' + cm.doc.getSelection(),
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
          dicHTML += `</div><div class="ft-12">${result.data.means.replace(
            /\n/g, '<br>')}</div></div>`
        }
        $('#dictionary').html(dicHTML)
      },
    })
  })
  $('#dictionary').on('mouseover', 'svg', function () {
    $('#dictionary audio')[0].play()
  })
}

const review = function (passed, id) {
  $.ajax({
    url: config.server + '/reviews/review/' + id,
    type: 'POST',
    data: JSON.stringify({
      passed,
      memo: $('#memo').val(),
      content: editor.edit.getValue(),
      title: $('#title').val()
    }),
    success: function (result) {
      window.location.href = result.data
    },
  })
}