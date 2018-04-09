var port = 6868;
var server = server = require('webserver').create();
var page = require('webpage').create();

var service = server.listen(port, function (request, response) {
  console.log(JSON.stringify(request, null, 4));
  var text = JSON.parse(request.post).text;
  console.log(text);

  page.open('https://translate.google.cn/#en/zh-CN/' + encodeURIComponent(text), function(status) {
    console.log(status);  
    var result = page.evaluate(function() {
      return document.getElementById('result_box').textContent;
    });

    console.log(result);

    response.write(result);
    response.close();
  });
});


console.log("Translation engine is running at port: " + port);
