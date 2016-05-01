(function() {

  var textInput;
  var request = null;
  var classificationLabel;

  function classify() {
    classificationLabel.style.display = 'block';
    classificationLabel.innerText = 'Loading...';
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
    var obj = JSON.parse(classification);

    classificationLabel.innerText = 'Classification: ' + obj.lang;
    classificationLabel.style.display = 'block';
  }

  window.addEventListener('load', function() {
    textInput = document.getElementById('text-input');
    classificationLabel = document.getElementById('classification');
    var classifyButton = document.getElementById('classify-button');
    classifyButton.addEventListener('click', classify);
  });

})();
