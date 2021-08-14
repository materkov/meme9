package controller

type NewPostRenderer struct {

}

func (n *NewPostRenderer) Render() string {
	return `
<script>
function sendPost() {
	fetch("/add_post", {
	  method: 'POST',
	  body: JSON.stringify({"text": document.querySelector("#text").value})
	}).then(r => r.json()).then(r => {
		window.location.href = r.postUrl;
	});
}
</script>
<textarea id="text"></textarea>
<button onclick="sendPost()">Отправить</button>
`
}