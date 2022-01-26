var convert = $("#convert")
var extractAudioDomain = "http://localhost:8000/extract-audio";
var downloadAudioDomain = "http://localhost:8000/download-audio"; 

async function extractAudio(domain, videoURL) {
  var res = await fetch(domain, {method: "POST", body:  JSON.stringify({"videoURL": videoURL})});
  var reader = res.body.getReader()
  var res = await reader.read()
  return res;
}

convert.on("click", function() {
  var videoURL = $("#ytURL").val()
  console.log("received after button click: ", videoURL);

  var res;
  res = extractAudio(extractAudioDomain, videoURL);
  res.then(function(res) {
    var textDecoder = new TextDecoder();
    console.log(textDecoder.decode(res.value))
  })  
})