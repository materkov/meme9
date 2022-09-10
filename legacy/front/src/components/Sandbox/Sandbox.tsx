import React from "react";
import * as schema from "../../api/renderer";
import {UniversalRenderer} from "../UniversalRenderer/UniversalRenderer";
import styles from "./Sandbox.module.css";

const defaultRenderer = `{
  "postRenderer": {
    "post": {
      "id": "88",
      "url": "/posts/88",
      "authorId": "1",
      "authorAvatar": "https://sun6-21.userapi.com/s/v1/ig2/9ixFw8PuBakcmO9IerClB3O1u-iFsX6Xc5kHjK-FCTtxCrIy912u5JNajr2vMrkbRE264vNM0pSECAakY2aQ4Cfa.jpg?size=200x0&quality=96&crop=122,105,561,561&ava=1",
      "authorName": "Макс Матерков",
      "authorUrl": "/users/1",
      "dateDisplay": "2 Jun 2021 14:35",
      "text": "dsf",
      "canLike": true
    },
    "composer": {
      "postId": "88",
      "placeholder": "Напишите здесь свой комментарий..."
    }
  }
}`;

export function Sandbox(props: { data: schema.SandboxRenderer }) {
    const [text, setText] = React.useState(defaultRenderer);
    const [renderer, setRenderer] = React.useState({} as schema.UniversalRenderer);
    const [changeTimeout, setChangeTimeout] = React.useState(null as any);

    React.useEffect(() => {
        updateRenderer(text);
    }, []);

    const onChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        setText(e.target.value);
        if (changeTimeout) {
            clearTimeout(changeTimeout);
        }

        setChangeTimeout(setTimeout(() => {
            updateRenderer(e.target.value);
        }, 200));
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
