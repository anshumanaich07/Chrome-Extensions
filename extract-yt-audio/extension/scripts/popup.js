var convert = $("#convert")

convert.on("click", function() {
  var videoURL = $("#ytURL").val()
  console.log("received after button click: ", videoURL)
  var extractAudioDomain = "http://localhost:8000/extract-audio";
  var downloadAudioDomain = "http://localhost:8000/download-audio"; 
  // var es = new EventSource("messages", {payload: {"ytURL":videoURL}});
  fetch(extractAudioDomain, {
    method: "POST",
    body: JSON.stringify({"videoURL": videoURL})
  }).then(function(res) {
    console.log(res)
    var es = new EventSource();
    es.onmessage = function(e) {
      console.log(e.Data)
    }
    // return res.json();
  }).then(function(res) {
    console.log("hi")
    if (res.Msg === "download") {
      // var es = new EventSource(downloadAudioDomain);
      // es.onmessage = function(e) {
      //   console.dir(e)
      //   console.log(e.Progress)
      // }
      window.open(downloadAudioDomain);
    };
  });
})
