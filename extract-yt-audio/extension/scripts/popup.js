var convert = $("#convert")

// domains
var getURLDomain  = "http://localhost:8000/get-url";
var extractAudioDomain = "http://localhost:8000/extract-audio";
var downloadAudioDomain = "http://localhost:8000/download-audio"; 

async function sendURL(domain, videoURL) {
  var res = fetch(domain, {method: "POST", body:  JSON.stringify({"videoURL": videoURL})});
  return res;
}

convert.on("click", function() {
  var videoURL = $("#ytURL").val()
  console.log("received after button click: ", videoURL);

  var res;
  res = sendURL(getURLDomain, videoURL);
  res.then(function(res) {
    return res.json();
  }).then(function(res) {
    if (res.Msg == "download") {
      var source = new EventSource(extractAudioDomain);
      source.onmessage = function (event) {
        console.log('data received from backend: ', event.data);
        if (event.data == "100%") { 
          source.close(); 
          window.open(downloadAudioDomain)
        }
      };
    }
  });
})