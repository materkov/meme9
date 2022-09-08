import React, {useEffect} from "react";
import styles from "./Header.module.css";
import {Link} from "./Link";
import {apiHost, User} from "../store/types";
import {emitCustomEvent, useCustomEventListener} from "react-custom-events";

function getViewer(): Promise<User> {
    return new Promise<User>((resolve, reject) => {

        fetch(apiHost + "/browse?url=/viewer", {
            headers: {
                'authorization': 'Bearer ' + localStorage.getItem('authToken'),
            }
        })
            .then(r => r.json())
            .then(r => {
                resolve(r.componentData[0] as User)
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
        getViewer().then(setViewer);
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
                    <Link href={viewer.href}>{viewer.name}</Link> | <a onClick={onLogout} href={"#"}>Выход</a>
                    </>
                }
            </div>
        </div>
    )
}