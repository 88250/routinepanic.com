const port = 6868
const server = server = require('webserver').create()
const page = require('webpage').create()

server.listen(port, function (request, response) {
  page.open('https://translate.google.cn', function () {
    if (status === 'success') {
      page.evaluate(function () {
        document.getElementById('result_box').value = request.post
        document.getElementById('gt-submit').click()
        const intervalTime = setInterval(function () {
          if (document.getElementById('result_box').textContent !== '') {
            clearInterval(intervalTime)
            response.write(document.getElementById('result_box').textContent)
            response.close()
          }
        }, 1000)
      })
    } else {
      response.write(status + 1)
      response.close()
    }
  })
})

console.log('Translation Server is running at port: ' + port)