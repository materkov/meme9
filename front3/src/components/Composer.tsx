import React from "react";
import {ComposerRenderer} from "./types";
import {StateContext} from "./StoreConnector";

export const Composer = (props: { data: ComposerRenderer }) => {
    const state = React.useContext(StateContext);
    const [text, setText] = React.useState<string>("");

    const onSend = () => {
        state.do({
            createPost: {
                text: text,
            }
        })
    };

    const onTextChanged = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        setText(e.target.value);
    };

    return <>
        <h2>{props.data.placeholder}</h2>
        <textarea value={text} onChange={onTextChanged}/>
        <br/>
        <button onClick={onSend}>Отправить</button>
    </>
}
