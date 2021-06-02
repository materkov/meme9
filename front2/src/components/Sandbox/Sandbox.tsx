import React from "react";
import * as schema from "../../api/renderer";
import {UniversalRenderer} from "../UniversalRenderer/UniversalRenderer";
import styles from "./Sandbox.module.css";

export function Sandbox(props: { data: schema.SandboxRenderer }) {
    const [text, setText] = React.useState('{}');
    const [renderer, setRenderer] = React.useState({} as schema.UniversalRenderer);
    const [changeTimeout, setChangeTimeout] = React.useState(null as any);

    const onChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        setText(e.target.value);
        if (changeTimeout) {
            clearTimeout(changeTimeout);
        }

        setChangeTimeout(setTimeout(() => {
            updateRenderer(e.target.value);
        }, 500));
    }

    const updateRenderer = (text: string) => {
        try {
            const renderer = JSON.parse(text);
            setRenderer(schema.UniversalRenderer.fromJSON(renderer));
        } catch (e) {
            return;
        }
    }

    return (
        <>
            <textarea className={styles.Props} value={text} onChange={onChange}/>
            <UniversalRenderer data={renderer}/>
        </>
    );
}
