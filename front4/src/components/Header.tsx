import styles from "./Header.module.css";
import React, {useEffect} from "react";
import {QueryParams} from "../types";
import {api} from "../api";
import {getByID, getByType, storeOnChanged, storeSubscribe, storeUnsubscribe} from "../store/store";
import {getViewer, getVkAuthUrl} from "../store/viewer";

export function Header() {
    const [authURL, setAuthURL] = React.useState('');
    const [viewerName, setViewerName] = React.useState("");
    const [viewerId, setViewerId] = React.useState("");

    useEffect(() => {
        const cb = () => {
            const viewer = getByType("Query");
            if (viewer && viewer.type == "Query" && viewerId != viewer.viewer) {
                setViewerId(viewer.viewer || "");
                getViewer();

                const q: QueryParams = {
                    viewer: {
                        inner: {
                            name: {}
                        }
                    }
                }
                api(q).then(result => {
                    // @ts-ignore
                    setViewerName(result.viewer.name);
                })
            }

            const vkAuthUrl = getByType("VkAuthURL");
            if (vkAuthUrl && vkAuthUrl.type == "VkAuthURL") {
                setAuthURL(vkAuthUrl.url || "");
            }

        }
        storeSubscribe(cb);

        getVkAuthUrl();

        return () => storeUnsubscribe(cb);
    }, []);

    const onExit = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        const token = getByType("CurrentRoute");
        if (token && token.type == "CurrentRoute") {
            setViewerName("");
            setViewerId("");

            localStorage.removeItem("authToken");

            token.url = "/";
            storeOnChanged();
        }
    };

    return (
        <div className={styles.header}>
            <a href={"/"} className={styles.headerLink}>meme9</a>

            {viewerId && <div>{viewerName} | <a href={"#"} onClick={onExit}>Выйти</a></div>}
            {!viewerId && <a href={authURL}>Войти</a>}
        </div>
    )
}
