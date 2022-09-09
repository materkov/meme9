import React, {useEffect} from "react";
import styles from "./Header.module.css";
import {Link} from "./Link";
import {api, User} from "../store/types";
import {emitCustomEvent, useCustomEventListener} from "react-custom-events";

function getViewer(): Promise<User | null> {
    return new Promise<User | null>((resolve, reject) => {
        api("/viewer").then(r => {
            resolve(r[0]);
        })
    })

}

const vkURL = "https://oauth.vk.com/authorize?client_id=7260220&response_type=code&redirect_uri=" + location.origin + "/vk-callback"

export function Header() {
    const [viewer, setViewer] = React.useState<User | undefined | null>();
    useEffect(() => {
        refreshUser();
    }, [])

    const onLogout = (e: MouseEvent) => {
        e.preventDefault();
        localStorage.removeItem('authToken');

        setViewer(null);
        emitCustomEvent('onAuthorized');
    }

    const refreshUser = () => {
        api("/viewer").then(r => {
            setViewer(r[0]);
        })
    }

    useCustomEventListener('onAuthorized', () => {
        refreshUser();
    });

    return (
        <div className={styles.header}>
            <Link href="/" className={styles.logo}>meme</Link>

            <div className={styles.userName}>
                {viewer === null && <a href={vkURL}>Авторизация</a>}
                {viewer !== null && viewer !== undefined &&
                    <>
                        <Link href={"/users/" + viewer.id}>{viewer.name}</Link> | <a onClick={onLogout}
                                                                                     href={"#"}>Выход</a>
                    </>
                }
            </div>
        </div>
    )
}