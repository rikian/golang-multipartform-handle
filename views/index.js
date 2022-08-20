function qs(elm) {
    return document.querySelector(elm)
}

function ce(elm) {
    return document.createElement(elm)
}

var preview = new FileReader
var upload = new FileReader
var imgPreview = qs("#imgPrev")
var fileImage = qs("#file-image")

fileImage.addEventListener("change", function() {
    var prev = this.files[0]
    preview.readAsDataURL(prev)
    preview.onload = function() {
        imgPreview.src = this.result
    }
})

upload.addEventListener("load", function() {
    console.log(this)
})