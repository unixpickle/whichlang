(function() {

  var textInput;
  var request = null;

  function classify() {
    if (request !== null) {
      request.abort();
    }
    request = new XMLHttpRequest();
    request.onreadystatechange = function() {
      if (request.readyState === 4) {
        handleClassification(request.responseText);
        request = null;
      }
    };
    var time = new Date().getTime();
    request.open('POST', '/classify?time=' + time, true);
    request.send(textInput.value);
  }

  function handleClassification(classification) {
    var label = document.getElementById('classification');
    label.innerText = 'Classification: ' + classification;
    label.style.display = 'block';
  }

  window.addEventListener('load', function() {
    textInput = document.getElementById('text-input');
    var classifyButton = document.getElementById('classify-button');
    classifyButton.addEventListener('click', classify);
  });

})();
