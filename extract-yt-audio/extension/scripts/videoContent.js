function getVideoInfo() {
  var videoURL = window.location.href;
  var extractAudioDomain = "http://localhost:8000/extract-audio";
  var downloadAudioDomain = "http://localhost:8000/download-audio"; 

  fetch(extractAudioDomain, {
    method: "POST",
    body: JSON.stringify({"videoURL": videoURL})
  }).then(function(res) {
    return res.json();
  }).then(function(res) {
    if (res.Msg === "download") {
      console.log("here")
      window.open(downloadAudioDomain);
    };
  });
};

getVideoInfo();
