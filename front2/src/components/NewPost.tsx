import React from "react";
import {NewPostRenderer} from "../types";

export const NewPost = (props: { data: NewPostRenderer }) => {
    const [text, setText] = React.useState<string>("");

    function sendPost() {
        fetch("http://localhost:8000/add_post", {
            method: 'POST',
            body: JSON.stringify({text: text})
        }).then(r => r.json()).then(r => {
            window.location.href = r.postUrl;
        });
    }

    return <div>
        <textarea value={text} onChange={event => setText(event.target.value)}/>
        <button onClick={sendPost}>{props.data.sendLabel}</button>
    </div>
}
